package di

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ferriyusra/movie-service/internal/api"
	"github.com/ferriyusra/movie-service/internal/api/handler"
	"github.com/ferriyusra/movie-service/internal/platform"

	genreRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/genre"
	movieRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/movie"
	refreshTokenRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/refresh_token"
	reservationRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/reservation"
	reservationSeatRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/reservation_seat"
	seatRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/seat"
	showtimeRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/showtime"
	theaterRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/theater"
	userRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/user"

	genreSvc "github.com/ferriyusra/movie-service/internal/service/genre"
	healthSvc "github.com/ferriyusra/movie-service/internal/service/health"
	movieSvc "github.com/ferriyusra/movie-service/internal/service/movie"
	reservationSvc "github.com/ferriyusra/movie-service/internal/service/reservation"
	showtimeSvc "github.com/ferriyusra/movie-service/internal/service/showtime"
	theaterSvc "github.com/ferriyusra/movie-service/internal/service/theater"
	tokenSvc "github.com/ferriyusra/movie-service/internal/service/token"
	userSvc "github.com/ferriyusra/movie-service/internal/service/user"
)

// Container holds all application dependencies
type Container struct {
	Config   *platform.Config
	Router   *gin.Engine
	Services *Services
}

// Services holds all service layer dependencies
type Services struct {
	Health   healthSvc.HealthService
	User     userSvc.UserService
	Token    tokenSvc.TokenService
	Genre    genreSvc.GenreService
	Movie    movieSvc.MovieService
	Theater  theaterSvc.TheaterService
	Showtime    showtimeSvc.ShowtimeService
	Reservation reservationSvc.ReservationService
}

// Handlers holds all HTTP handler dependencies
type Handlers struct {
	Health      *handler.HealthHandler
	User        *handler.UserHandler
	Genre       *handler.GenreHandler
	Movie       *handler.MovieHandler
	Theater     *handler.TheaterHandler
	Showtime    *handler.ShowtimeHandler
	Reservation *handler.ReservationHandler
}

// NewContainer creates and initializes a new dependency container
func NewContainer(cfg *platform.Config) (*Container, error) {
	if !cfg.Auth.DevMode {
		if cfg.Auth.JWTAccessSecret == "" || cfg.Auth.JWTRefreshSecret == "" {
			return nil, fmt.Errorf("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be set in production")
		}
	}

	accessSecret := cfg.Auth.JWTAccessSecret
	if accessSecret == "" {
		accessSecret = "dev-access-secret-DO-NOT-USE-IN-PRODUCTION"
	}
	refreshSecret := cfg.Auth.JWTRefreshSecret
	if refreshSecret == "" {
		refreshSecret = "dev-refresh-secret-DO-NOT-USE-IN-PRODUCTION"
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Auth.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	db := cfg.Database.Gorm
	if db == nil {
		var err error
		db, err = platform.InitializeDatabase(cfg)
		if err != nil {
			return nil, fmt.Errorf("initializing database: %w", err)
		}
		cfg.Database.Gorm = db
	}

	// Seed admin account on first startup
	if err := platform.SeedAdminAccount(db); err != nil {
		return nil, fmt.Errorf("seeding admin account: %w", err)
	}

	// Initialize repositories
	userRepository, err := userRepo.NewGORMUserRepository(db)
	if err != nil {
		return nil, fmt.Errorf("initializing user repository: %w", err)
	}
	refreshTokenRepository, err := refreshTokenRepo.NewGORMRefreshTokenRepository(db)
	if err != nil {
		return nil, fmt.Errorf("initializing refresh token repository: %w", err)
	}
	genreRepository := genreRepo.NewGORMGenreRepository(db)
	movieRepository := movieRepo.NewGORMMovieRepository(db)
	theaterRepository := theaterRepo.NewGORMTheaterRepository(db)
	seatRepository := seatRepo.NewGORMSeatRepository(db)
	showtimeRepository := showtimeRepo.NewGORMShowtimeRepository(db)
	reservationRepository := reservationRepo.NewGORMReservationRepository(db)
	reservationSeatRepository := reservationSeatRepo.NewGORMReservationSeatRepository(db)

	// Initialize token service
	tokenConfig := tokenSvc.TokenConfig{
		AccessTokenSecret:  accessSecret,
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenSecret: refreshSecret,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}
	tokenService := tokenSvc.NewTokenService(tokenConfig)

	// Initialize services
	services := &Services{
		Health: healthSvc.NewHealthService(),
		User:   userSvc.NewUserService(userRepository, refreshTokenRepository, tokenService),
		Token:  tokenService,
		Genre:    genreSvc.NewGenreService(genreRepository),
		Movie:    movieSvc.NewMovieService(movieRepository, genreRepository),
		Theater:  theaterSvc.NewTheaterService(theaterRepository, seatRepository),
		Showtime:    showtimeSvc.NewShowtimeService(showtimeRepository, movieRepository, theaterRepository),
		Reservation: reservationSvc.NewReservationService(reservationRepository, reservationSeatRepository, showtimeRepository, seatRepository),
	}

	// Initialize handlers
	handlers := &Handlers{
		Health:      handler.NewHealthHandler(services.Health),
		User:        handler.NewUserHandler(services.User),
		Genre:       handler.NewGenreHandler(services.Genre),
		Movie:       handler.NewMovieHandler(services.Movie),
		Theater:     handler.NewTheaterHandler(services.Theater),
		Showtime:    handler.NewShowtimeHandler(services.Showtime),
		Reservation: handler.NewReservationHandler(services.Reservation),
	}

	// Setup routes
	api.SetupRoutes(r, handlers.User, handlers.Genre, handlers.Movie, handlers.Theater, handlers.Showtime, handlers.Reservation, services.Token)
	api.SetupHealthRoutes(r, handlers.Health)

	return &Container{
		Config:   cfg,
		Router:   r,
		Services: services,
	}, nil
}

package di

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ferriyusra/movie-service/internal/api"
	"github.com/ferriyusra/movie-service/internal/api/handler"
	"github.com/ferriyusra/movie-service/internal/platform"

	refreshTokenRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/refresh_token"
	userRepo "github.com/ferriyusra/movie-service/internal/repository/implementations/user"

	healthSvc "github.com/ferriyusra/movie-service/internal/service/health"
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
	Health healthSvc.HealthService
	User   userSvc.UserService
	Token  tokenSvc.TokenService
}

// Handlers holds all HTTP handler dependencies
type Handlers struct {
	Health *handler.HealthHandler
	User   *handler.UserHandler
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

	// Initialize repositories
	userRepository, err := userRepo.NewGORMUserRepository(db)
	if err != nil {
		return nil, fmt.Errorf("initializing user repository: %w", err)
	}
	refreshTokenRepository, err := refreshTokenRepo.NewGORMRefreshTokenRepository(db)
	if err != nil {
		return nil, fmt.Errorf("initializing refresh token repository: %w", err)
	}

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
	}

	// Initialize handlers
	handlers := &Handlers{
		Health: handler.NewHealthHandler(services.Health),
		User:   handler.NewUserHandler(services.User),
	}

	// Setup routes
	api.SetupRoutes(r, handlers.User, services.Token)
	api.SetupHealthRoutes(r, handlers.Health)

	return &Container{
		Config:   cfg,
		Router:   r,
		Services: services,
	}, nil
}

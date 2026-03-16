package di

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ferriyusra/clean-arch-go-gin/internal/api"
	"github.com/ferriyusra/clean-arch-go-gin/internal/api/handler"
	"github.com/ferriyusra/clean-arch-go-gin/internal/platform"

	counterRepo "github.com/ferriyusra/clean-arch-go-gin/internal/repository/implementations/counter"
	messageRepo "github.com/ferriyusra/clean-arch-go-gin/internal/repository/implementations/message"
	refreshTokenRepo "github.com/ferriyusra/clean-arch-go-gin/internal/repository/implementations/refresh_token"
	userRepo "github.com/ferriyusra/clean-arch-go-gin/internal/repository/implementations/user"

	counterSvc "github.com/ferriyusra/clean-arch-go-gin/internal/service/counter"
	csrfSvc "github.com/ferriyusra/clean-arch-go-gin/internal/service/csrf"
	healthSvc "github.com/ferriyusra/clean-arch-go-gin/internal/service/health"
	messageSvc "github.com/ferriyusra/clean-arch-go-gin/internal/service/message"
	tokenSvc "github.com/ferriyusra/clean-arch-go-gin/internal/service/token"
	userSvc "github.com/ferriyusra/clean-arch-go-gin/internal/service/user"
)

// Container holds all application dependencies
type Container struct {
	Config   *platform.Config
	Router   *gin.Engine
	Services *Services
}

// Services holds all service layer dependencies
type Services struct {
	Message messageSvc.MessageService
	Health  healthSvc.HealthService
	Counter counterSvc.CounterService
	User    userSvc.UserService
	Token   tokenSvc.TokenService
	CSRF    csrfSvc.CSRFService
}

// Handlers holds all HTTP handler dependencies
type Handlers struct {
	Message *handler.MessageHandler
	Health  *handler.HealthHandler
	Counter *handler.CounterHandler
	User    *handler.UserHandler
}

// NewContainer creates and initializes a new dependency container
func NewContainer(cfg *platform.Config) (*Container, error) {
	// Validate secrets in production
	if !cfg.Auth.DevMode {
		if cfg.Auth.JWTAccessSecret == "" || cfg.Auth.JWTRefreshSecret == "" {
			return nil, fmt.Errorf("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be set in production")
		}
		if cfg.Auth.CSRFSecret == "" {
			return nil, fmt.Errorf("CSRF_SECRET must be set in production")
		}
	}

	// Use defaults in dev mode if not set
	accessSecret := cfg.Auth.JWTAccessSecret
	if accessSecret == "" {
		accessSecret = "dev-access-secret-DO-NOT-USE-IN-PRODUCTION"
	}
	refreshSecret := cfg.Auth.JWTRefreshSecret
	if refreshSecret == "" {
		refreshSecret = "dev-refresh-secret-DO-NOT-USE-IN-PRODUCTION"
	}
	csrfSecret := cfg.Auth.CSRFSecret
	if csrfSecret == "" {
		csrfSecret = "dev-csrf-secret-DO-NOT-USE-IN-PRODUCTION"
	}

	// Initialize Gin
	r := gin.Default()

	// Setup CORS middleware before routes
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Auth.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// Initialize database
	db := cfg.Database.Gorm
	if db == nil {
		var err error
		db, err = platform.InitializeDatabase(cfg)
		if err != nil {
			return nil, fmt.Errorf("initializing database: %w", err)
		}
		cfg.Database.Gorm = db
	}

	// Initialize repositories (handle errors)
	counterRepository, err := counterRepo.NewGORMCounterRepository(db)
	if err != nil {
		return nil, fmt.Errorf("initializing counter repository: %w", err)
	}
	messageRepository, err := messageRepo.NewGORMMessageRepository(db)
	if err != nil {
		return nil, fmt.Errorf("initializing message repository: %w", err)
	}
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
		Message: messageSvc.NewMessageService(messageRepository),
		Health:  healthSvc.NewHealthService(),
		Counter: counterSvc.NewCounterService(counterRepository),
		User:    userSvc.NewUserService(userRepository, refreshTokenRepository, tokenService),
		Token:   tokenService,
		CSRF:    csrfSvc.NewCSRFService(csrfSecret),
	}

	// Initialize handlers
	handlers := &Handlers{
		Message: handler.NewMessageHandler(services.Message),
		Health:  handler.NewHealthHandler(services.Health),
		Counter: handler.NewCounterHandler(services.Counter),
		User:    handler.NewUserHandler(services.User, services.Token, services.CSRF, cfg.Auth.DevMode),
	}

	// Setup routes with dependencies
	api.SetupRoutes(r, handlers.Message, handlers.Counter, handlers.User, services.Token, services.CSRF)
	api.SetupHealthRoutes(r, handlers.Health)

	return &Container{
		Config:   cfg,
		Router:   r,
		Services: services,
	}, nil
}

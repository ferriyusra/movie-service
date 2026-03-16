package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ferriyusra/clean-arch-go-gin/internal/api/handler"
	"github.com/ferriyusra/clean-arch-go-gin/internal/api/middleware"
	"github.com/ferriyusra/clean-arch-go-gin/internal/service/csrf"
	"github.com/ferriyusra/clean-arch-go-gin/internal/service/token"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	r *gin.Engine,
	messageHandler *handler.MessageHandler,
	counterHandler *handler.CounterHandler,
	userHandler *handler.UserHandler,
	tokenService token.TokenService,
	csrfService csrf.CSRFService,
) {
	api := r.Group("/api")

	// Public routes (no authentication required)
	api.GET("/message", messageHandler.GetMessage)

	// Auth routes (public)
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/refresh", middleware.CSRFMiddleware(csrfService), userHandler.Refresh)
	api.GET("/csrf", userHandler.GetCSRFToken)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(tokenService))

	// Auth protected endpoints
	protected.GET("/auth/me", userHandler.GetMe)

	// Logout requires auth + CSRF protection (it's a POST request)
	protected.POST("/auth/logout", middleware.CSRFMiddleware(csrfService), userHandler.Logout)

	// Counter endpoints (protected)
	protected.GET("/counter", counterHandler.GetCounter)

	// Counter POST requires auth + CSRF protection
	protected.POST("/counter", middleware.CSRFMiddleware(csrfService), counterHandler.IncrementCounter)
}

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(r *gin.Engine, healthHandler *handler.HealthHandler) {
	api := r.Group("/api")
	api.GET("/health", healthHandler.Check)
}

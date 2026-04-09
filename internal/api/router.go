package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ferriyusra/movie-service/internal/api/handler"
	"github.com/ferriyusra/movie-service/internal/api/middleware"
	"github.com/ferriyusra/movie-service/internal/service/token"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	tokenService token.TokenService,
) {
	api := r.Group("/api")

	// Auth routes (public)
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/refresh", userHandler.Refresh)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(tokenService))

	// Auth protected endpoints
	protected.GET("/auth/me", userHandler.GetMe)
	protected.POST("/auth/logout", userHandler.Logout)
}

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(r *gin.Engine, healthHandler *handler.HealthHandler) {
	api := r.Group("/api")
	api.GET("/health", healthHandler.Check)
}

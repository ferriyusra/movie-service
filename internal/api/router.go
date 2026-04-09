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
	genreHandler *handler.GenreHandler,
	movieHandler *handler.MovieHandler,
	theaterHandler *handler.TheaterHandler,
	tokenService token.TokenService,
) {
	api := r.Group("/api")

	// Auth routes (public)
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/refresh", userHandler.Refresh)

	// Public routes
	api.GET("/genres", genreHandler.List)
	api.GET("/movies", movieHandler.List)
	api.GET("/movies/:id", movieHandler.Get)
	api.GET("/theaters", theaterHandler.List)
	api.GET("/theaters/:id", theaterHandler.Get)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(tokenService))

	// Auth protected endpoints
	protected.GET("/auth/me", userHandler.GetMe)
	protected.POST("/auth/logout", userHandler.Logout)

	// Admin routes (require authentication + admin role)
	admin := protected.Group("")
	admin.Use(middleware.AdminOnly())

	admin.POST("/genres", genreHandler.Create)
	admin.DELETE("/genres/:id", genreHandler.Delete)
	admin.POST("/movies", movieHandler.Create)
	admin.PUT("/movies/:id", movieHandler.Update)
	admin.DELETE("/movies/:id", movieHandler.Delete)
	admin.POST("/theaters", theaterHandler.Create)
	admin.PUT("/theaters/:id", theaterHandler.Update)
}

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(r *gin.Engine, healthHandler *handler.HealthHandler) {
	api := r.Group("/api")
	api.GET("/health", healthHandler.Check)
}

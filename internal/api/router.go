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
	showtimeHandler *handler.ShowtimeHandler,
	reservationHandler *handler.ReservationHandler,
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
	api.GET("/showtimes", showtimeHandler.ListByDate)
	api.GET("/showtimes/:id", showtimeHandler.Get)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(tokenService))

	// Auth protected endpoints
	protected.GET("/auth/me", userHandler.GetMe)
	protected.POST("/auth/logout", userHandler.Logout)

	// Reservation routes (authenticated users)
	protected.POST("/reservations", reservationHandler.Create)
	protected.GET("/reservations", reservationHandler.List)
	protected.GET("/reservations/:id", reservationHandler.Get)
	protected.DELETE("/reservations/:id", reservationHandler.Cancel)

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
	admin.POST("/showtimes", showtimeHandler.Create)
	admin.PUT("/showtimes/:id", showtimeHandler.Update)
	admin.DELETE("/showtimes/:id", showtimeHandler.Delete)
}

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(r *gin.Engine, healthHandler *handler.HealthHandler) {
	api := r.Group("/api")
	api.GET("/health", healthHandler.Check)
}

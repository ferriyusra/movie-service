package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/api/middleware"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
	adminSvc "github.com/ferriyusra/movie-service/internal/service/admin"
	showtimeSvc "github.com/ferriyusra/movie-service/internal/service/showtime"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	adminService    adminSvc.AdminService
	showtimeService showtimeSvc.ShowtimeService
}

// NewAdminHandler creates a new instance of AdminHandler
func NewAdminHandler(adminService adminSvc.AdminService, showtimeService showtimeSvc.ShowtimeService) *AdminHandler {
	return &AdminHandler{adminService: adminService, showtimeService: showtimeService}
}

// ListReservations handles GET /api/admin/reservations
func (h *AdminHandler) ListReservations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filters := interfaces.ReservationFilters{
		Status:   c.Query("status"),
		DateFrom: c.Query("dateFrom"),
		DateTo:   c.Query("dateTo"),
		Page:     page,
		Limit:    limit,
	}

	if userIDStr := c.Query("userId"); userIDStr != "" {
		id, err := uuid.Parse(userIDStr)
		if err == nil {
			filters.UserID = &id
		}
	}
	if showtimeIDStr := c.Query("showtimeId"); showtimeIDStr != "" {
		id, err := uuid.Parse(showtimeIDStr)
		if err == nil {
			filters.ShowtimeID = &id
		}
	}

	reservations, total, err := h.adminService.ListAllReservations(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch reservations"))
		return
	}

	c.JSON(http.StatusOK, response.OKWithMeta("Reservations retrieved", reservations, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: int(total),
	}))
}

// CancelReservation handles DELETE /api/admin/reservations/:id
func (h *AdminHandler) CancelReservation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid reservation ID"))
		return
	}

	if err := h.adminService.CancelAnyReservation(c.Request.Context(), id); err != nil {
		if err.Error() == "reservation not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Reservation cancelled", nil))
}

// PromoteUser handles PATCH /api/admin/users/:id/promote
func (h *AdminHandler) PromoteUser(c *gin.Context) {
	adminID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Err("Unauthorized"))
		return
	}

	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid user ID"))
		return
	}

	if err := h.adminService.PromoteUser(c.Request.Context(), adminID, userID); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("User promoted to admin", nil))
}

// ListUsers handles GET /api/admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.adminService.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch users"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Users retrieved", users))
}

// ListShowtimes handles GET /api/admin/showtimes
func (h *AdminHandler) ListShowtimes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filters := interfaces.ShowtimeFilters{
		MovieTitle: c.Query("movieTitle"),
		DateFrom:   c.Query("dateFrom"),
		DateTo:     c.Query("dateTo"),
		Page:       page,
		Limit:      limit,
	}

	showtimes, total, err := h.showtimeService.ListAllShowtimes(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch showtimes"))
		return
	}

	c.JSON(http.StatusOK, response.OKWithMeta("Showtimes retrieved", showtimes, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: int(total),
	}))
}

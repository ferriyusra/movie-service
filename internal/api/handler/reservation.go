package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/api/middleware"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	reservationSvc "github.com/ferriyusra/movie-service/internal/service/reservation"
)

// ReservationHandler handles reservation-related HTTP requests
type ReservationHandler struct {
	reservationService reservationSvc.ReservationService
}

// NewReservationHandler creates a new instance of ReservationHandler
func NewReservationHandler(reservationService reservationSvc.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

// Create handles POST /api/reservations (authenticated user)
func (h *ReservationHandler) Create(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Err("Unauthorized"))
		return
	}

	req := &request.CreateReservationRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

	resp, err := h.reservationService.CreateReservation(c.Request.Context(), userID, req)
	if err != nil {
		if err.Error() == "showtime not found" || err.Error() == "cannot reserve seats for a past showtime" {
			c.JSON(http.StatusBadRequest, response.Err(err.Error()))
			return
		}
		if len(err.Error()) > 20 && err.Error()[:20] == "seats already reserv" {
			c.JSON(http.StatusConflict, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.OK("Reservation created", resp))
}

// Cancel handles DELETE /api/reservations/:id (authenticated user)
func (h *ReservationHandler) Cancel(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Err("Unauthorized"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid reservation ID"))
		return
	}

	if err := h.reservationService.CancelReservation(c.Request.Context(), userID, id); err != nil {
		switch err.Error() {
		case "reservation not found":
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
		case "you can only cancel your own reservations":
			c.JSON(http.StatusForbidden, response.Err(err.Error()))
		case "reservation is already cancelled":
			c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		case "cancellation window has passed (must be at least 1 hour before showtime)":
			c.JSON(http.StatusUnprocessableEntity, response.Err(err.Error()))
		default:
			c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.OK("Reservation cancelled", nil))
}

// Get handles GET /api/reservations/:id (authenticated user)
func (h *ReservationHandler) Get(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Err("Unauthorized"))
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid reservation ID"))
		return
	}

	resp, err := h.reservationService.GetReservation(c.Request.Context(), userID, id)
	if err != nil {
		if err.Error() == "reservation not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		if err.Error() == "you can only view your own reservations" {
			c.JSON(http.StatusForbidden, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch reservation"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Reservation retrieved", resp))
}

// List handles GET /api/reservations (authenticated user)
func (h *ReservationHandler) List(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Err("Unauthorized"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	status := c.Query("status")

	reservations, total, err := h.reservationService.ListMyReservations(c.Request.Context(), userID, status, page, limit)
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

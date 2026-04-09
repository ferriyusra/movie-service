package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	showtimeSvc "github.com/ferriyusra/movie-service/internal/service/showtime"
)

// ShowtimeHandler handles showtime-related HTTP requests
type ShowtimeHandler struct {
	showtimeService showtimeSvc.ShowtimeService
}

// NewShowtimeHandler creates a new instance of ShowtimeHandler
func NewShowtimeHandler(showtimeService showtimeSvc.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{showtimeService: showtimeService}
}

// Create handles POST /api/showtimes (admin only)
func (h *ShowtimeHandler) Create(c *gin.Context) {
	req := &request.CreateShowtimeRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

	resp, err := h.showtimeService.CreateShowtime(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "showtime conflicts with an existing showtime in this theater" {
			c.JSON(http.StatusConflict, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.OK("Showtime created", resp))
}

// Update handles PUT /api/showtimes/:id (admin only)
func (h *ShowtimeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid showtime ID"))
		return
	}

	req := &request.UpdateShowtimeRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

	resp, err := h.showtimeService.UpdateShowtime(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "showtime not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		if err.Error() == "showtime conflicts with an existing showtime in this theater" {
			c.JSON(http.StatusConflict, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Showtime updated", resp))
}

// Delete handles DELETE /api/showtimes/:id (admin only)
func (h *ShowtimeHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid showtime ID"))
		return
	}

	if err := h.showtimeService.DeleteShowtime(c.Request.Context(), id); err != nil {
		if err.Error() == "showtime not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Showtime deleted", nil))
}

// Get handles GET /api/showtimes/:id (public)
func (h *ShowtimeHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid showtime ID"))
		return
	}

	resp, err := h.showtimeService.GetShowtime(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "showtime not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch showtime"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Showtime retrieved", resp))
}

// ListByDate handles GET /api/showtimes?date=YYYY-MM-DD (public)
func (h *ShowtimeHandler) ListByDate(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid date format, use YYYY-MM-DD"))
		return
	}

	showtimes, err := h.showtimeService.ListShowtimesByDate(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch showtimes"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Showtimes retrieved", showtimes))
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	theaterSvc "github.com/ferriyusra/movie-service/internal/service/theater"
)

// TheaterHandler handles theater-related HTTP requests
type TheaterHandler struct {
	theaterService theaterSvc.TheaterService
}

// NewTheaterHandler creates a new instance of TheaterHandler
func NewTheaterHandler(theaterService theaterSvc.TheaterService) *TheaterHandler {
	return &TheaterHandler{theaterService: theaterService}
}

// Create handles POST /api/theaters (admin only)
func (h *TheaterHandler) Create(c *gin.Context) {
	req := &request.CreateTheaterRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

	resp, err := h.theaterService.CreateTheater(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.OK("Theater created", resp))
}

// Update handles PUT /api/theaters/:id (admin only)
func (h *TheaterHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid theater ID"))
		return
	}

	req := &request.UpdateTheaterRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

	resp, err := h.theaterService.UpdateTheater(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "theater not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Theater updated", resp))
}

// Get handles GET /api/theaters/:id (public)
func (h *TheaterHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid theater ID"))
		return
	}

	resp, err := h.theaterService.GetTheater(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "theater not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch theater"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Theater retrieved", resp))
}

// List handles GET /api/theaters (public)
func (h *TheaterHandler) List(c *gin.Context) {
	theaters, err := h.theaterService.ListTheaters(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch theaters"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Theaters retrieved", theaters))
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	genreSvc "github.com/ferriyusra/movie-service/internal/service/genre"
)

// GenreHandler handles genre-related HTTP requests
type GenreHandler struct {
	genreService genreSvc.GenreService
}

// NewGenreHandler creates a new instance of GenreHandler
func NewGenreHandler(genreService genreSvc.GenreService) *GenreHandler {
	return &GenreHandler{
		genreService: genreService,
	}
}

// Create handles POST /api/genres (admin only)
func (h *GenreHandler) Create(c *gin.Context) {
	req := &request.CreateGenreRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

	resp, err := h.genreService.CreateGenre(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.OK("Genre created", resp))
}

// Delete handles DELETE /api/genres/:id (admin only)
func (h *GenreHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid genre ID"))
		return
	}

	if err := h.genreService.DeleteGenre(c.Request.Context(), id); err != nil {
		if err.Error() == "genre not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Genre deleted", nil))
}

// List handles GET /api/genres (public)
func (h *GenreHandler) List(c *gin.Context) {
	genres, err := h.genreService.ListGenres(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch genres"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Genres retrieved", genres))
}

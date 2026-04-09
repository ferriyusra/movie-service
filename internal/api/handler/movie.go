package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
	movieSvc "github.com/ferriyusra/movie-service/internal/service/movie"
)

// MovieHandler handles movie-related HTTP requests
type MovieHandler struct {
	movieService movieSvc.MovieService
}

// NewMovieHandler creates a new instance of MovieHandler
func NewMovieHandler(movieService movieSvc.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

// Create handles POST /api/movies (admin only)
func (h *MovieHandler) Create(c *gin.Context) {
	req := &request.CreateMovieRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

	resp, err := h.movieService.CreateMovie(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.OK("Movie created", resp))
}

// Update handles PUT /api/movies/:id (admin only)
func (h *MovieHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid movie ID"))
		return
	}

	req := &request.UpdateMovieRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

	resp, err := h.movieService.UpdateMovie(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Movie updated", resp))
}

// Delete handles DELETE /api/movies/:id (admin only)
func (h *MovieHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid movie ID"))
		return
	}

	if err := h.movieService.DeleteMovie(c.Request.Context(), id); err != nil {
		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Movie deleted", nil))
}

// Get handles GET /api/movies/:id (public)
func (h *MovieHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid movie ID"))
		return
	}

	resp, err := h.movieService.GetMovie(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, response.Err(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch movie"))
		return
	}

	c.JSON(http.StatusOK, response.OK("Movie retrieved", resp))
}

// List handles GET /api/movies (public)
func (h *MovieHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filters := interfaces.MovieFilters{
		Title: c.Query("title"),
		Page:  page,
		Limit: limit,
	}

	if genreIDStr := c.Query("genre_id"); genreIDStr != "" {
		genreID, err := uuid.Parse(genreIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Err("Invalid genre_id"))
			return
		}
		filters.GenreID = &genreID
	}

	movies, total, err := h.movieService.ListMovies(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to fetch movies"))
		return
	}

	c.JSON(http.StatusOK, response.OKWithMeta("Movies retrieved", movies, &response.Meta{
		Page:  page,
		Limit: limit,
		Total: int(total),
	}))
}

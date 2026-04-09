package movie

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

// MovieService defines the interface for movie operations
type MovieService interface {
	CreateMovie(ctx context.Context, req *request.CreateMovieRequest) (*response.MovieResponse, error)
	UpdateMovie(ctx context.Context, id uuid.UUID, req *request.UpdateMovieRequest) (*response.MovieResponse, error)
	DeleteMovie(ctx context.Context, id uuid.UUID) error
	GetMovie(ctx context.Context, id uuid.UUID) (*response.MovieResponse, error)
	ListMovies(ctx context.Context, filters interfaces.MovieFilters) ([]response.MovieResponse, int64, error)
}

// movieService is the concrete implementation of MovieService
type movieService struct {
	movieRepository interfaces.MovieRepository
	genreRepository interfaces.GenreRepository
}

// NewMovieService creates a new instance of MovieService
func NewMovieService(
	movieRepository interfaces.MovieRepository,
	genreRepository interfaces.GenreRepository,
) MovieService {
	return &movieService{
		movieRepository: movieRepository,
		genreRepository: genreRepository,
	}
}

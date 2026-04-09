package genre

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

// GenreService defines the interface for genre operations
type GenreService interface {
	CreateGenre(ctx context.Context, req *request.CreateGenreRequest) (*response.GenreResponse, error)
	DeleteGenre(ctx context.Context, id uuid.UUID) error
	ListGenres(ctx context.Context) ([]response.GenreResponse, error)
}

// genreService is the concrete implementation of GenreService
type genreService struct {
	genreRepository interfaces.GenreRepository
}

// NewGenreService creates a new instance of GenreService
func NewGenreService(genreRepository interfaces.GenreRepository) GenreService {
	return &genreService{
		genreRepository: genreRepository,
	}
}

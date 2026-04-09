package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// MovieRepository defines the interface for movie data access
type MovieRepository interface {
	Create(ctx context.Context, movie entity.MovieEntity) (*uuid.UUID, error)
	Update(ctx context.Context, movie entity.MovieEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.MovieEntity, error)
	FindAll(ctx context.Context, filters MovieFilters) ([]entity.MovieEntity, int64, error)
	ReplaceGenres(ctx context.Context, movieID uuid.UUID, genres []entity.GenreEntity) error
}

// MovieFilters holds filter/pagination parameters for listing movies
type MovieFilters struct {
	Title   string
	GenreID *uuid.UUID
	Page    int
	Limit   int
}

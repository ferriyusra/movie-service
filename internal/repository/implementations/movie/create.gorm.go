package movie

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// Create creates a new movie
func (r *GORMMovieRepository) Create(ctx context.Context, movie entity.MovieEntity) (*uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := r.db.WithContext(ctx).Create(&movie).Error; err != nil {
		return nil, err
	}

	return &movie.ID, nil
}

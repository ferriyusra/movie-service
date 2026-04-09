package movie

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// Delete soft-deletes a movie by ID
func (r *GORMMovieRepository) Delete(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return r.db.WithContext(ctx).Delete(&entity.MovieEntity{}, "id = ?", id).Error
}

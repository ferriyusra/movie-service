package movie

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// Update updates an existing movie
func (r *GORMMovieRepository) Update(ctx context.Context, movie entity.MovieEntity) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Omit Genres to prevent GORM from upserting the many2many association
	return r.db.WithContext(ctx).Omit("Genres").Save(&movie).Error
}

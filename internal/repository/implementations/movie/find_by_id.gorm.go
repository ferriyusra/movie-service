package movie

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// FindByID returns a movie by ID with its genres preloaded, or nil if not found
func (r *GORMMovieRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.MovieEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var movie entity.MovieEntity
	err := r.db.WithContext(ctx).Preload("Genres").Where("id = ?", id).First(&movie).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &movie, nil
}

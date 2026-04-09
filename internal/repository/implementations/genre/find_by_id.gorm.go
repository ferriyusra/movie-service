package genre

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// FindByID returns a genre by ID, or nil if not found
func (r *GORMGenreRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.GenreEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var genre entity.GenreEntity
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&genre).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &genre, nil
}

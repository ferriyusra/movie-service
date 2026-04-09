package genre

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// Create creates a new genre
func (r *GORMGenreRepository) Create(ctx context.Context, genre entity.GenreEntity) (*uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := r.db.WithContext(ctx).Create(&genre).Error; err != nil {
		return nil, err
	}

	return &genre.ID, nil
}

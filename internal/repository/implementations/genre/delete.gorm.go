package genre

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// Delete removes a genre by ID
func (r *GORMGenreRepository) Delete(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return r.db.WithContext(ctx).Delete(&entity.GenreEntity{}, "id = ?", id).Error
}

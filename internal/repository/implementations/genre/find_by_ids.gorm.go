package genre

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// FindByIDs returns genres matching the given IDs
func (r *GORMGenreRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.GenreEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var genres []entity.GenreEntity
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&genres).Error; err != nil {
		return nil, err
	}

	return genres, nil
}

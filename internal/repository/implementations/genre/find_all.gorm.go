package genre

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// FindAll returns all genres
func (r *GORMGenreRepository) FindAll(ctx context.Context) ([]entity.GenreEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var genres []entity.GenreEntity
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&genres).Error; err != nil {
		return nil, err
	}

	return genres, nil
}

package theater

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMTheaterRepository) FindAll(ctx context.Context) ([]entity.TheaterEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var theaters []entity.TheaterEntity
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&theaters).Error; err != nil {
		return nil, err
	}
	return theaters, nil
}

package theater

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMTheaterRepository) Update(ctx context.Context, theater entity.TheaterEntity) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return r.db.WithContext(ctx).Save(&theater).Error
}

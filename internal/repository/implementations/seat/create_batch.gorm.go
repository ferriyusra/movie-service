package seat

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMSeatRepository) CreateBatch(ctx context.Context, seats []entity.SeatEntity) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return r.db.WithContext(ctx).Create(&seats).Error
}

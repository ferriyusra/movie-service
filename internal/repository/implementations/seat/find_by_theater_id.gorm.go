package seat

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMSeatRepository) FindByTheaterID(ctx context.Context, theaterID uuid.UUID) ([]entity.SeatEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var seats []entity.SeatEntity
	if err := r.db.WithContext(ctx).
		Where("theater_id = ?", theaterID).
		Order("row ASC, number ASC").
		Find(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}

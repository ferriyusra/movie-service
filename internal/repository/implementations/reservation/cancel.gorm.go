package reservation

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMReservationRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.ReservationEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "cancelled",
			"cancelled_at": &now,
		}).Error
}

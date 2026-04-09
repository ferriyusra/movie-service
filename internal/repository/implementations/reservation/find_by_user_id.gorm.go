package reservation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMReservationRepository) FindByUserID(ctx context.Context, userID uuid.UUID, status string, page, limit int) ([]entity.ReservationEntity, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.ReservationEntity{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var reservations []entity.ReservationEntity
	if err := query.
		Preload("Showtime").Preload("Showtime.Movie").Preload("Showtime.Theater").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&reservations).Error; err != nil {
		return nil, 0, err
	}

	return reservations, total, nil
}

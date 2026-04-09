package reservation

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

func (r *GORMReservationRepository) FindAll(ctx context.Context, filters interfaces.ReservationFilters) ([]entity.ReservationEntity, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.ReservationEntity{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.ShowtimeID != nil {
		query = query.Where("showtime_id = ?", *filters.ShowtimeID)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.DateFrom != "" {
		query = query.Where("created_at >= ?", filters.DateFrom)
	}
	if filters.DateTo != "" {
		query = query.Where("created_at <= ?", filters.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.Limit
	var reservations []entity.ReservationEntity
	if err := query.
		Preload("Showtime").Preload("Showtime.Movie").Preload("Showtime.Theater").Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(filters.Limit).
		Find(&reservations).Error; err != nil {
		return nil, 0, err
	}

	return reservations, total, nil
}

package reservation_seat

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMReservationSeatRepository) FindByShowtimeID(ctx context.Context, showtimeID uuid.UUID) ([]entity.ReservationSeatEntity, error) {
	var seats []entity.ReservationSeatEntity
	if err := r.db.WithContext(ctx).
		Preload("Seat").
		Joins("JOIN reservations ON reservations.id = reservation_seats.reservation_id").
		Where("reservation_seats.showtime_id = ? AND reservations.status = 'confirmed'", showtimeID).
		Find(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}

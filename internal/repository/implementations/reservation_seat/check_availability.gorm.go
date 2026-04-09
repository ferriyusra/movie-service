package reservation_seat

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// CheckAvailabilityWithTx uses SELECT FOR UPDATE to lock and check seat availability
func (r *GORMReservationSeatRepository) CheckAvailabilityWithTx(ctx context.Context, tx *gorm.DB, showtimeID uuid.UUID, seatIDs []uuid.UUID) ([]entity.ReservationSeatEntity, error) {
	var taken []entity.ReservationSeatEntity
	if err := tx.WithContext(ctx).
		Joins("JOIN reservations ON reservations.id = reservation_seats.reservation_id").
		Where("reservation_seats.showtime_id = ? AND reservation_seats.seat_id IN ? AND reservations.status = 'confirmed'", showtimeID, seatIDs).
		Preload("Seat").
		Find(&taken).Error; err != nil {
		return nil, err
	}
	return taken, nil
}

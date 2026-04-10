package reservation_seat

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// CheckAvailabilityWithTx locks existing reservation_seat rows for the given
// showtime using SELECT ... FOR UPDATE, then checks if any of the requested
// seats are already taken. The row-level lock prevents two concurrent
// transactions from both reading "0 taken" and then both inserting.
func (r *GORMReservationSeatRepository) CheckAvailabilityWithTx(ctx context.Context, tx *gorm.DB, showtimeID uuid.UUID, seatIDs []uuid.UUID) ([]entity.ReservationSeatEntity, error) {
	var taken []entity.ReservationSeatEntity
	if err := tx.WithContext(ctx).
		Set("gorm:query_option", "FOR UPDATE").
		Joins("JOIN reservations ON reservations.id = reservation_seats.reservation_id").
		Where("reservation_seats.showtime_id = ? AND reservation_seats.seat_id IN ? AND reservations.status = 'confirmed'", showtimeID, seatIDs).
		Preload("Seat").
		Find(&taken).Error; err != nil {
		return nil, err
	}
	return taken, nil
}

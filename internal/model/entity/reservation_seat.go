package entity

import "github.com/google/uuid"

// ReservationSeatEntity links a reservation to specific seats
type ReservationSeatEntity struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	ReservationID uuid.UUID `gorm:"not null;index"`
	SeatID        uuid.UUID `gorm:"not null;uniqueIndex:idx_seat_showtime"`
	ShowtimeID    uuid.UUID `gorm:"not null;uniqueIndex:idx_seat_showtime"`
	Reservation   ReservationEntity `gorm:"foreignKey:ReservationID"`
	Seat          SeatEntity        `gorm:"foreignKey:SeatID"`
}

func (ReservationSeatEntity) TableName() string {
	return "reservation_seats"
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

// ReservationEntity represents a seat reservation
type ReservationEntity struct {
	ID               uuid.UUID  `gorm:"primaryKey"`
	BookingReference string     `gorm:"uniqueIndex;type:varchar(30);not null"`
	UserID           uuid.UUID  `gorm:"not null;index"`
	ShowtimeID       uuid.UUID  `gorm:"not null;index"`
	Status           string     `gorm:"type:varchar(20);default:'confirmed'"`
	TotalAmount      float64    `gorm:"type:decimal(10,2);not null"`
	CancelledAt      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	User             UserEntity     `gorm:"foreignKey:UserID"`
	Showtime         ShowtimeEntity `gorm:"foreignKey:ShowtimeID"`
}

func (ReservationEntity) TableName() string {
	return "reservations"
}

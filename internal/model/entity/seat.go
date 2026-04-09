package entity

import "github.com/google/uuid"

// SeatEntity represents a seat in a theater
type SeatEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	TheaterID uuid.UUID `gorm:"not null;index"`
	Row       string    `gorm:"type:char(1);not null"`
	Number    int       `gorm:"not null"`
	Label     string    `gorm:"type:varchar(10);not null"`
	Type      string    `gorm:"type:varchar(20);default:'standard'"`
}

func (SeatEntity) TableName() string {
	return "seats"
}

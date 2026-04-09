package entity

import (
	"time"

	"github.com/google/uuid"
)

// TheaterEntity represents a cinema theater
type TheaterEntity struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Location    string    `gorm:"type:varchar(255)"`
	TotalSeats  int       `gorm:"not null"`
	Rows        int       `gorm:"not null"`
	SeatsPerRow int       `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TheaterEntity) TableName() string {
	return "theaters"
}

package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShowtimeEntity represents a movie showtime in a theater
type ShowtimeEntity struct {
	ID        uuid.UUID      `gorm:"primaryKey"`
	MovieID   uuid.UUID      `gorm:"not null;index"`
	TheaterID uuid.UUID      `gorm:"not null;index"`
	StartTime time.Time      `gorm:"not null"`
	EndTime   time.Time      `gorm:"not null"`
	Price     float64        `gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Movie     MovieEntity    `gorm:"foreignKey:MovieID"`
	Theater   TheaterEntity  `gorm:"foreignKey:TheaterID"`
}

func (ShowtimeEntity) TableName() string {
	return "showtimes"
}

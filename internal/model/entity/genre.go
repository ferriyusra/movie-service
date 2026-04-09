package entity

import (
	"time"

	"github.com/google/uuid"
)

// GenreEntity represents a movie genre
type GenreEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex;type:varchar(100);not null"`
	CreatedAt time.Time
}

func (GenreEntity) TableName() string {
	return "genres"
}

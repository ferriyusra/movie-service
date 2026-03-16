package entity

import (
	"time"

	"github.com/google/uuid"
)

// CounterEntity represents the counter table schema
type CounterEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Value     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (CounterEntity) TableName() string {
	return "counters"
}

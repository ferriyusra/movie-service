package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserEntity represents a user in the system
type UserEntity struct {
	ID        uuid.UUID      `gorm:"primaryKey"`
	Email     string         `gorm:"uniqueIndex"`
	Password  []byte
	Name      string
	Role      string         `gorm:"type:varchar(20);default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserEntity) TableName() string {
	return "users"
}

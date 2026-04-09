package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MovieEntity represents a movie in the system
type MovieEntity struct {
	ID          uuid.UUID      `gorm:"primaryKey"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	PosterURL   string         `gorm:"type:varchar(500)"`
	DurationMin int            `gorm:"not null"`
	Language    string         `gorm:"type:varchar(50);default:'English'"`
	ReleaseDate *time.Time     `gorm:"type:date"`
	Rating      float64        `gorm:"type:decimal(3,1);default:0.0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Genres      []GenreEntity  `gorm:"many2many:movie_genres;"`
}

func (MovieEntity) TableName() string {
	return "movies"
}

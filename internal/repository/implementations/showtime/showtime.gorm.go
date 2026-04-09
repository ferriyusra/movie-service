package showtime

import "gorm.io/gorm"

// GORMShowtimeRepository is a GORM implementation of ShowtimeRepository
type GORMShowtimeRepository struct {
	db *gorm.DB
}

// NewGORMShowtimeRepository creates a new GORM showtime repository
func NewGORMShowtimeRepository(db *gorm.DB) *GORMShowtimeRepository {
	return &GORMShowtimeRepository{db: db}
}

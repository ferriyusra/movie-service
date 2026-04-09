package theater

import "gorm.io/gorm"

// GORMTheaterRepository is a GORM implementation of TheaterRepository
type GORMTheaterRepository struct {
	db *gorm.DB
}

// NewGORMTheaterRepository creates a new GORM theater repository
func NewGORMTheaterRepository(db *gorm.DB) *GORMTheaterRepository {
	return &GORMTheaterRepository{db: db}
}

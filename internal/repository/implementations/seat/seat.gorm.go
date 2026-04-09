package seat

import "gorm.io/gorm"

// GORMSeatRepository is a GORM implementation of SeatRepository
type GORMSeatRepository struct {
	db *gorm.DB
}

// NewGORMSeatRepository creates a new GORM seat repository
func NewGORMSeatRepository(db *gorm.DB) *GORMSeatRepository {
	return &GORMSeatRepository{db: db}
}

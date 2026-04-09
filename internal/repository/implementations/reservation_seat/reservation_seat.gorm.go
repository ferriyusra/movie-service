package reservation_seat

import "gorm.io/gorm"

type GORMReservationSeatRepository struct {
	db *gorm.DB
}

func NewGORMReservationSeatRepository(db *gorm.DB) *GORMReservationSeatRepository {
	return &GORMReservationSeatRepository{db: db}
}

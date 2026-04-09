package reservation

import (
	"context"

	"gorm.io/gorm"
)

type GORMReservationRepository struct {
	db *gorm.DB
}

func NewGORMReservationRepository(db *gorm.DB) *GORMReservationRepository {
	return &GORMReservationRepository{db: db}
}

func (r *GORMReservationRepository) BeginTx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

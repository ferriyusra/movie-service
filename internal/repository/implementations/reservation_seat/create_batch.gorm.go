package reservation_seat

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

func (r *GORMReservationSeatRepository) CreateBatchWithTx(ctx context.Context, tx *gorm.DB, seats []entity.ReservationSeatEntity) error {
	return tx.WithContext(ctx).Create(&seats).Error
}

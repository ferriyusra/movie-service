package reservation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

func (r *GORMReservationRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, reservation entity.ReservationEntity) (*uuid.UUID, error) {
	if err := tx.WithContext(ctx).Create(&reservation).Error; err != nil {
		return nil, err
	}
	return &reservation.ID, nil
}

package reservation

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

func (r *GORMReservationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.ReservationEntity, error) {
	var reservation entity.ReservationEntity
	err := r.db.WithContext(ctx).
		Preload("Showtime").Preload("Showtime.Movie").Preload("Showtime.Theater").
		Where("id = ?", id).First(&reservation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reservation, nil
}

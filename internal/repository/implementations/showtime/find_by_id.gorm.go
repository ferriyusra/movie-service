package showtime

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

func (r *GORMShowtimeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.ShowtimeEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var showtime entity.ShowtimeEntity
	err := r.db.WithContext(ctx).
		Preload("Movie").Preload("Movie.Genres").Preload("Theater").
		Where("id = ?", id).First(&showtime).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &showtime, nil
}

package showtime

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMShowtimeRepository) Create(ctx context.Context, showtime entity.ShowtimeEntity) (*uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if err := r.db.WithContext(ctx).Create(&showtime).Error; err != nil {
		return nil, err
	}
	return &showtime.ID, nil
}

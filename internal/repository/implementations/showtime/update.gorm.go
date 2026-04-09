package showtime

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMShowtimeRepository) Update(ctx context.Context, showtime entity.ShowtimeEntity) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return r.db.WithContext(ctx).Save(&showtime).Error
}

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
	// Use Omit to prevent GORM from upserting preloaded associations (Movie, Theater)
	return r.db.WithContext(ctx).Omit("Movie", "Theater").Save(&showtime).Error
}

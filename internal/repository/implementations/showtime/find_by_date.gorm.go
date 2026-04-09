package showtime

import (
	"context"
	"time"

	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMShowtimeRepository) FindByDate(ctx context.Context, date time.Time) ([]entity.ShowtimeEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var showtimes []entity.ShowtimeEntity
	if err := r.db.WithContext(ctx).
		Preload("Movie").Preload("Movie.Genres").Preload("Theater").
		Where("start_time >= ? AND start_time < ?", startOfDay, endOfDay).
		Order("start_time ASC").
		Find(&showtimes).Error; err != nil {
		return nil, err
	}
	return showtimes, nil
}

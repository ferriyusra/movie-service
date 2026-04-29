package showtime

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

func (r *GORMShowtimeRepository) FindAll(ctx context.Context, filters interfaces.ShowtimeFilters) ([]entity.ShowtimeEntity, int64, error) {
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	query := r.db.WithContext(ctx).Model(&entity.ShowtimeEntity{}).
		Joins("JOIN movies ON movies.id = showtimes.movie_id AND movies.deleted_at IS NULL")

	if filters.MovieTitle != "" {
		query = query.Where("movies.title LIKE ?", "%"+filters.MovieTitle+"%")
	}
	if filters.DateFrom != "" {
		query = query.Where("showtimes.start_time >= ?", filters.DateFrom)
	}
	if filters.DateTo != "" {
		query = query.Where("showtimes.start_time <= ?", filters.DateTo+" 23:59:59")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.Limit
	var showtimes []entity.ShowtimeEntity
	if err := query.
		Preload("Movie").Preload("Movie.Genres").Preload("Theater").
		Order("start_time ASC").
		Offset(offset).Limit(filters.Limit).
		Find(&showtimes).Error; err != nil {
		return nil, 0, err
	}

	return showtimes, total, nil
}

package movie

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

// FindAll returns movies matching the given filters with pagination
func (r *GORMMovieRepository) FindAll(ctx context.Context, filters interfaces.MovieFilters) ([]entity.MovieEntity, int64, error) {
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	query := r.db.WithContext(ctx).Model(&entity.MovieEntity{})

	if filters.Title != "" {
		query = query.Where("title ILIKE ?", "%"+filters.Title+"%")
	}

	if filters.GenreID != nil {
		query = query.Joins("JOIN movie_genres ON movie_genres.movie_entity_id = movies.id").
			Where("movie_genres.genre_entity_id = ?", *filters.GenreID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.Limit
	var movies []entity.MovieEntity
	if err := query.Preload("Genres").
		Order("created_at DESC").
		Offset(offset).
		Limit(filters.Limit).
		Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

package movie

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// ReplaceGenres replaces all genres for a movie
func (r *GORMMovieRepository) ReplaceGenres(ctx context.Context, movieID uuid.UUID, genres []entity.GenreEntity) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remove existing associations
		if err := tx.Exec("DELETE FROM movie_genres WHERE movie_entity_id = ?", movieID).Error; err != nil {
			return err
		}

		// Insert new associations
		for _, g := range genres {
			if err := tx.Exec("INSERT INTO movie_genres (movie_entity_id, genre_entity_id) VALUES (?, ?)", movieID, g.ID).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

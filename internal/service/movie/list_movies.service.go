package movie

import (
	"context"
	"fmt"

	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

// ListMovies returns movies matching the given filters with pagination
func (s *movieService) ListMovies(ctx context.Context, filters interfaces.MovieFilters) ([]response.MovieResponse, int64, error) {
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	movies, total, err := s.movieRepository.FindAll(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("listing movies: %w", err)
	}

	result := make([]response.MovieResponse, len(movies))
	for i, m := range movies {
		genreResponses := make([]response.GenreResponse, len(m.Genres))
		for j, g := range m.Genres {
			genreResponses[j] = response.GenreResponse{ID: g.ID, Name: g.Name}
		}
		result[i] = response.MovieResponse{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			PosterURL:   m.PosterURL,
			DurationMin: m.DurationMin,
			Language:    m.Language,
			ReleaseDate: m.ReleaseDate,
			Rating:      m.Rating,
			Genres:      genreResponses,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		}
	}

	return result, total, nil
}

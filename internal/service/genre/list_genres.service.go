package genre

import (
	"context"
	"fmt"

	"github.com/ferriyusra/movie-service/internal/model/response"
)

// ListGenres returns all genres
func (s *genreService) ListGenres(ctx context.Context) ([]response.GenreResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	genres, err := s.genreRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing genres: %w", err)
	}

	result := make([]response.GenreResponse, len(genres))
	for i, g := range genres {
		result[i] = response.GenreResponse{
			ID:   g.ID,
			Name: g.Name,
		}
	}

	return result, nil
}

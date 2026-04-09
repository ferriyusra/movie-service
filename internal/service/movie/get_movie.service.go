package movie

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// GetMovie retrieves a movie by ID
func (s *movieService) GetMovie(ctx context.Context, id uuid.UUID) (*response.MovieResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	movie, err := s.movieRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding movie: %w", err)
	}
	if movie == nil {
		return nil, errors.New("movie not found")
	}

	genreResponses := make([]response.GenreResponse, len(movie.Genres))
	for i, g := range movie.Genres {
		genreResponses[i] = response.GenreResponse{ID: g.ID, Name: g.Name}
	}

	return &response.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		PosterURL:   movie.PosterURL,
		DurationMin: movie.DurationMin,
		Language:    movie.Language,
		ReleaseDate: movie.ReleaseDate,
		Rating:      movie.Rating,
		Genres:      genreResponses,
		CreatedAt:   movie.CreatedAt,
		UpdatedAt:   movie.UpdatedAt,
	}, nil
}

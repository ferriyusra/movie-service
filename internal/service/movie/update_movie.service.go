package movie

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// UpdateMovie updates an existing movie and reassigns genres
func (s *movieService) UpdateMovie(ctx context.Context, id uuid.UUID, req *request.UpdateMovieRequest) (*response.MovieResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	existing, err := s.movieRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding movie: %w", err)
	}
	if existing == nil {
		return nil, errors.New("movie not found")
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("movie title is required")
	}
	if req.DurationMin <= 0 {
		return nil, errors.New("duration must be greater than zero")
	}

	// Validate genre IDs
	var genres []entity.GenreEntity
	if len(req.GenreIDs) > 0 {
		genres, err = s.genreRepository.FindByIDs(ctx, req.GenreIDs)
		if err != nil {
			return nil, fmt.Errorf("finding genres: %w", err)
		}
		if len(genres) != len(req.GenreIDs) {
			return nil, errors.New("one or more genre IDs are invalid")
		}
	}

	var releaseDate *time.Time
	if req.ReleaseDate != "" {
		t, err := time.Parse("2006-01-02", req.ReleaseDate)
		if err != nil {
			return nil, errors.New("invalid release_date format, use YYYY-MM-DD")
		}
		releaseDate = &t
	}

	language := req.Language
	if language == "" {
		language = "English"
	}

	existing.Title = title
	existing.Description = req.Description
	existing.PosterURL = req.PosterURL
	existing.DurationMin = req.DurationMin
	existing.Language = language
	existing.ReleaseDate = releaseDate

	if err := s.movieRepository.Update(ctx, *existing); err != nil {
		return nil, fmt.Errorf("updating movie: %w", err)
	}

	if err := s.movieRepository.ReplaceGenres(ctx, id, genres); err != nil {
		return nil, fmt.Errorf("replacing genres: %w", err)
	}

	genreResponses := make([]response.GenreResponse, len(genres))
	for i, g := range genres {
		genreResponses[i] = response.GenreResponse{ID: g.ID, Name: g.Name}
	}

	return &response.MovieResponse{
		ID:          existing.ID,
		Title:       existing.Title,
		Description: existing.Description,
		PosterURL:   existing.PosterURL,
		DurationMin: existing.DurationMin,
		Language:    existing.Language,
		ReleaseDate: existing.ReleaseDate,
		Rating:      existing.Rating,
		Genres:      genreResponses,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   existing.UpdatedAt,
	}, nil
}

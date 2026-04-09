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

// CreateMovie creates a new movie and attaches genres
func (s *movieService) CreateMovie(ctx context.Context, req *request.CreateMovieRequest) (*response.MovieResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
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
		var err error
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

	movie := entity.MovieEntity{
		ID:          uuid.New(),
		Title:       title,
		Description: req.Description,
		PosterURL:   req.PosterURL,
		DurationMin: req.DurationMin,
		Language:    language,
		ReleaseDate: releaseDate,
	}

	id, err := s.movieRepository.Create(ctx, movie)
	if err != nil {
		return nil, fmt.Errorf("creating movie: %w", err)
	}

	if len(genres) > 0 {
		if err := s.movieRepository.ReplaceGenres(ctx, *id, genres); err != nil {
			return nil, fmt.Errorf("attaching genres: %w", err)
		}
	}

	genreResponses := make([]response.GenreResponse, len(genres))
	for i, g := range genres {
		genreResponses[i] = response.GenreResponse{ID: g.ID, Name: g.Name}
	}

	return &response.MovieResponse{
		ID:          *id,
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

package showtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// GetShowtime returns a showtime with movie and theater details
func (s *showtimeService) GetShowtime(ctx context.Context, id uuid.UUID) (*response.ShowtimeDetailResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	showtime, err := s.showtimeRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding showtime: %w", err)
	}
	if showtime == nil {
		return nil, errors.New("showtime not found")
	}

	genreResponses := make([]response.GenreResponse, len(showtime.Movie.Genres))
	for i, g := range showtime.Movie.Genres {
		genreResponses[i] = response.GenreResponse{ID: g.ID, Name: g.Name}
	}

	return &response.ShowtimeDetailResponse{
		ShowtimeResponse: response.ShowtimeResponse{
			ID:        showtime.ID,
			MovieID:   showtime.MovieID,
			TheaterID: showtime.TheaterID,
			StartTime: showtime.StartTime,
			EndTime:   showtime.EndTime,
			Price:     showtime.Price,
		},
		Movie: response.MovieResponse{
			ID:          showtime.Movie.ID,
			Title:       showtime.Movie.Title,
			Description: showtime.Movie.Description,
			DurationMin: showtime.Movie.DurationMin,
			Language:    showtime.Movie.Language,
			Genres:      genreResponses,
		},
		Theater: response.TheaterResponse{
			ID:          showtime.Theater.ID,
			Name:        showtime.Theater.Name,
			Location:    showtime.Theater.Location,
			TotalSeats:  showtime.Theater.TotalSeats,
			Rows:        showtime.Theater.Rows,
			SeatsPerRow: showtime.Theater.SeatsPerRow,
		},
	}, nil
}

package showtime

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// CreateShowtime creates a new showtime with conflict detection
func (s *showtimeService) CreateShowtime(ctx context.Context, req *request.CreateShowtimeRequest) (*response.ShowtimeResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if req.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, errors.New("invalid start_time format, use RFC3339")
	}

	movie, err := s.movieRepository.FindByID(ctx, req.MovieID)
	if err != nil {
		return nil, fmt.Errorf("finding movie: %w", err)
	}
	if movie == nil {
		return nil, errors.New("movie not found")
	}

	theater, err := s.theaterRepository.FindByID(ctx, req.TheaterID)
	if err != nil {
		return nil, fmt.Errorf("finding theater: %w", err)
	}
	if theater == nil {
		return nil, errors.New("theater not found")
	}

	// Auto-calculate end_time = start_time + duration + 15min cleaning buffer
	endTime := startTime.Add(time.Duration(movie.DurationMin)*time.Minute + 15*time.Minute)

	// Check for conflicts
	conflict, err := s.showtimeRepository.CheckConflict(ctx, req.TheaterID, startTime, endTime, nil)
	if err != nil {
		return nil, fmt.Errorf("checking conflict: %w", err)
	}
	if conflict != nil {
		return nil, errors.New("showtime conflicts with an existing showtime in this theater")
	}

	showtime := entity.ShowtimeEntity{
		ID:        uuid.New(),
		MovieID:   req.MovieID,
		TheaterID: req.TheaterID,
		StartTime: startTime,
		EndTime:   endTime,
		Price:     req.Price,
	}

	id, err := s.showtimeRepository.Create(ctx, showtime)
	if err != nil {
		return nil, fmt.Errorf("creating showtime: %w", err)
	}

	return &response.ShowtimeResponse{
		ID:        *id,
		MovieID:   showtime.MovieID,
		TheaterID: showtime.TheaterID,
		StartTime: showtime.StartTime,
		EndTime:   showtime.EndTime,
		Price:     showtime.Price,
	}, nil
}

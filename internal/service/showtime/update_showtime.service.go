package showtime

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// UpdateShowtime updates an existing showtime
func (s *showtimeService) UpdateShowtime(ctx context.Context, id uuid.UUID, req *request.UpdateShowtimeRequest) (*response.ShowtimeResponse, error) {
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

	existing, err := s.showtimeRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding showtime: %w", err)
	}
	if existing == nil {
		return nil, errors.New("showtime not found")
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

	endTime := startTime.Add(time.Duration(movie.DurationMin)*time.Minute + 15*time.Minute)

	conflict, err := s.showtimeRepository.CheckConflict(ctx, req.TheaterID, startTime, endTime, &id)
	if err != nil {
		return nil, fmt.Errorf("checking conflict: %w", err)
	}
	if conflict != nil {
		return nil, errors.New("showtime conflicts with an existing showtime in this theater")
	}

	existing.MovieID = req.MovieID
	existing.TheaterID = req.TheaterID
	existing.StartTime = startTime
	existing.EndTime = endTime
	existing.Price = req.Price

	if err := s.showtimeRepository.Update(ctx, *existing); err != nil {
		return nil, fmt.Errorf("updating showtime: %w", err)
	}

	return &response.ShowtimeResponse{
		ID:          existing.ID,
		MovieID:     existing.MovieID,
		MovieTitle:  movie.Title,
		TheaterID:   existing.TheaterID,
		TheaterName: theater.Name,
		StartTime:   existing.StartTime,
		EndTime:     existing.EndTime,
		Price:       existing.Price,
	}, nil
}

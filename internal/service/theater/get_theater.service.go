package theater

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// GetTheater returns a theater with its seat layout
func (s *theaterService) GetTheater(ctx context.Context, id uuid.UUID) (*response.TheaterDetailResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	theater, err := s.theaterRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding theater: %w", err)
	}
	if theater == nil {
		return nil, errors.New("theater not found")
	}

	seats, err := s.seatRepository.FindByTheaterID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding seats: %w", err)
	}

	seatResponses := make([]response.SeatResponse, len(seats))
	for i, seat := range seats {
		seatResponses[i] = response.SeatResponse{
			ID:     seat.ID,
			Row:    seat.Row,
			Number: seat.Number,
			Label:  seat.Label,
			Type:   seat.Type,
		}
	}

	return &response.TheaterDetailResponse{
		TheaterResponse: response.TheaterResponse{
			ID:          theater.ID,
			Name:        theater.Name,
			Location:    theater.Location,
			TotalSeats:  theater.TotalSeats,
			Rows:        theater.Rows,
			SeatsPerRow: theater.SeatsPerRow,
		},
		Seats: seatResponses,
	}, nil
}

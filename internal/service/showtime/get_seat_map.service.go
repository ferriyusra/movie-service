package showtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// GetSeatMap returns all seats for a showtime with their availability status
func (s *showtimeService) GetSeatMap(ctx context.Context, showtimeID uuid.UUID) ([]response.SeatMapResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	showtime, err := s.showtimeRepository.FindByID(ctx, showtimeID)
	if err != nil {
		return nil, fmt.Errorf("finding showtime: %w", err)
	}
	if showtime == nil {
		return nil, errors.New("showtime not found")
	}

	// Get all seats for the theater
	seats, err := s.seatRepository.FindByTheaterID(ctx, showtime.TheaterID)
	if err != nil {
		return nil, fmt.Errorf("finding seats: %w", err)
	}

	// Get reserved seats for this showtime
	reservedSeats, err := s.reservationSeatRepository.FindByShowtimeID(ctx, showtimeID)
	if err != nil {
		return nil, fmt.Errorf("finding reserved seats: %w", err)
	}

	reservedSet := make(map[uuid.UUID]bool, len(reservedSeats))
	for _, rs := range reservedSeats {
		reservedSet[rs.SeatID] = true
	}

	result := make([]response.SeatMapResponse, len(seats))
	for i, seat := range seats {
		status := "available"
		if reservedSet[seat.ID] {
			status = "reserved"
		}
		result[i] = response.SeatMapResponse{
			SeatResponse: response.SeatResponse{
				ID:     seat.ID,
				Row:    seat.Row,
				Number: seat.Number,
				Label:  seat.Label,
				Type:   seat.Type,
			},
			Status: status,
		}
	}

	return result, nil
}

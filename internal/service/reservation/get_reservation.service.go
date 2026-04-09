package reservation

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// GetReservation returns a reservation by ID with ownership check
func (s *reservationService) GetReservation(ctx context.Context, userID uuid.UUID, reservationID uuid.UUID) (*response.ReservationDetailResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	reservation, err := s.reservationRepository.FindByID(ctx, reservationID)
	if err != nil {
		return nil, fmt.Errorf("finding reservation: %w", err)
	}
	if reservation == nil {
		return nil, errors.New("reservation not found")
	}

	if reservation.UserID != userID {
		return nil, errors.New("you can only view your own reservations")
	}

	// Fetch reservation seats
	resSeats, err := s.reservationSeatRepository.FindByShowtimeID(ctx, reservation.ShowtimeID)
	if err != nil {
		return nil, fmt.Errorf("finding reservation seats: %w", err)
	}

	// Filter to only this reservation's seats
	var seatResponses []response.SeatResponse
	for _, rs := range resSeats {
		if rs.ReservationID == reservationID {
			seatResponses = append(seatResponses, response.SeatResponse{
				ID:     rs.Seat.ID,
				Row:    rs.Seat.Row,
				Number: rs.Seat.Number,
				Label:  rs.Seat.Label,
				Type:   rs.Seat.Type,
			})
		}
	}

	genreResponses := make([]response.GenreResponse, len(reservation.Showtime.Movie.Genres))
	for i, g := range reservation.Showtime.Movie.Genres {
		genreResponses[i] = response.GenreResponse{ID: g.ID, Name: g.Name}
	}

	return &response.ReservationDetailResponse{
		ReservationResponse: response.ReservationResponse{
			ID:               reservation.ID,
			BookingReference: reservation.BookingReference,
			ShowtimeID:       reservation.ShowtimeID,
			Status:           reservation.Status,
			TotalAmount:      reservation.TotalAmount,
			CancelledAt:      reservation.CancelledAt,
			CreatedAt:        reservation.CreatedAt,
			Seats:            seatResponses,
		},
		Movie: response.MovieResponse{
			ID:          reservation.Showtime.Movie.ID,
			Title:       reservation.Showtime.Movie.Title,
			DurationMin: reservation.Showtime.Movie.DurationMin,
			Genres:      genreResponses,
		},
		Showtime: response.ShowtimeResponse{
			ID:        reservation.Showtime.ID,
			StartTime: reservation.Showtime.StartTime,
			EndTime:   reservation.Showtime.EndTime,
			Price:     reservation.Showtime.Price,
		},
		Theater: response.TheaterResponse{
			ID:   reservation.Showtime.Theater.ID,
			Name: reservation.Showtime.Theater.Name,
		},
	}, nil
}

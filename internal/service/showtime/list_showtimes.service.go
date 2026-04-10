package showtime

import (
	"context"
	"fmt"
	"time"

	"github.com/ferriyusra/movie-service/internal/model/response"
)

// ListShowtimesByDate returns all showtimes for a given date
func (s *showtimeService) ListShowtimesByDate(ctx context.Context, date time.Time) ([]response.ShowtimeResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	showtimes, err := s.showtimeRepository.FindByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("listing showtimes: %w", err)
	}

	result := make([]response.ShowtimeResponse, len(showtimes))
	for i, st := range showtimes {
		// Calculate available seats
		totalSeats := 0
		seats, err := s.seatRepository.FindByTheaterID(ctx, st.TheaterID)
		if err == nil {
			totalSeats = len(seats)
		}
		reservedSeats, err := s.reservationSeatRepository.FindByShowtimeID(ctx, st.ID)
		reservedCount := 0
		if err == nil {
			reservedCount = len(reservedSeats)
		}

		result[i] = response.ShowtimeResponse{
			ID:             st.ID,
			MovieID:        st.MovieID,
			MovieTitle:     st.Movie.Title,
			TheaterID:      st.TheaterID,
			TheaterName:    st.Theater.Name,
			StartTime:      st.StartTime,
			EndTime:        st.EndTime,
			Price:          st.Price,
			AvailableSeats: totalSeats - reservedCount,
		}
	}

	return result, nil
}

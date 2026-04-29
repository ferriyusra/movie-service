package showtime

import (
	"context"
	"fmt"

	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

// ListAllShowtimes returns all showtimes with optional filters (admin use)
func (s *showtimeService) ListAllShowtimes(ctx context.Context, filters interfaces.ShowtimeFilters) ([]response.ShowtimeResponse, int64, error) {
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	showtimes, total, err := s.showtimeRepository.FindAll(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("listing all showtimes: %w", err)
	}

	result := make([]response.ShowtimeResponse, len(showtimes))
	for i, st := range showtimes {
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

	return result, total, nil
}

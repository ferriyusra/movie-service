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
		result[i] = response.ShowtimeResponse{
			ID:        st.ID,
			MovieID:   st.MovieID,
			TheaterID: st.TheaterID,
			StartTime: st.StartTime,
			EndTime:   st.EndTime,
			Price:     st.Price,
		}
	}

	return result, nil
}

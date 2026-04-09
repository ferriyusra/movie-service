package theater

import (
	"context"
	"fmt"

	"github.com/ferriyusra/movie-service/internal/model/response"
)

// ListTheaters returns all theaters
func (s *theaterService) ListTheaters(ctx context.Context) ([]response.TheaterResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	theaters, err := s.theaterRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing theaters: %w", err)
	}

	result := make([]response.TheaterResponse, len(theaters))
	for i, t := range theaters {
		result[i] = response.TheaterResponse{
			ID:          t.ID,
			Name:        t.Name,
			Location:    t.Location,
			TotalSeats:  t.TotalSeats,
			Rows:        t.Rows,
			SeatsPerRow: t.SeatsPerRow,
		}
	}

	return result, nil
}

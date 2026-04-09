package theater

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// UpdateTheater updates a theater's name and location
func (s *theaterService) UpdateTheater(ctx context.Context, id uuid.UUID, req *request.UpdateTheaterRequest) (*response.TheaterResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	existing, err := s.theaterRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding theater: %w", err)
	}
	if existing == nil {
		return nil, errors.New("theater not found")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("theater name is required")
	}

	existing.Name = name
	existing.Location = req.Location

	if err := s.theaterRepository.Update(ctx, *existing); err != nil {
		return nil, fmt.Errorf("updating theater: %w", err)
	}

	return &response.TheaterResponse{
		ID:          existing.ID,
		Name:        existing.Name,
		Location:    existing.Location,
		TotalSeats:  existing.TotalSeats,
		Rows:        existing.Rows,
		SeatsPerRow: existing.SeatsPerRow,
	}, nil
}

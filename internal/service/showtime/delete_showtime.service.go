package showtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// DeleteShowtime soft-deletes a showtime
func (s *showtimeService) DeleteShowtime(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	existing, err := s.showtimeRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("finding showtime: %w", err)
	}
	if existing == nil {
		return errors.New("showtime not found")
	}

	if err := s.showtimeRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting showtime: %w", err)
	}

	return nil
}

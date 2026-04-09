package movie

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// DeleteMovie soft-deletes a movie by ID
func (s *movieService) DeleteMovie(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	existing, err := s.movieRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("finding movie: %w", err)
	}
	if existing == nil {
		return errors.New("movie not found")
	}

	if err := s.movieRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting movie: %w", err)
	}

	return nil
}

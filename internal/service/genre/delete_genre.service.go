package genre

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// DeleteGenre deletes a genre by ID after verifying it exists
func (s *genreService) DeleteGenre(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	existing, err := s.genreRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("finding genre: %w", err)
	}
	if existing == nil {
		return errors.New("genre not found")
	}

	if err := s.genreRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting genre: %w", err)
	}

	return nil
}

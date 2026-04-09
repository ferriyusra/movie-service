package reservation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CancelReservation cancels a user's reservation with ownership and time checks
func (s *reservationService) CancelReservation(ctx context.Context, userID uuid.UUID, reservationID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	reservation, err := s.reservationRepository.FindByID(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("finding reservation: %w", err)
	}
	if reservation == nil {
		return errors.New("reservation not found")
	}

	if reservation.UserID != userID {
		return errors.New("you can only cancel your own reservations")
	}

	if reservation.Status == "cancelled" {
		return errors.New("reservation is already cancelled")
	}

	// Check 1-hour cancellation window
	if reservation.Showtime.StartTime.Before(time.Now().Add(1 * time.Hour)) {
		return errors.New("cancellation window has passed (must be at least 1 hour before showtime)")
	}

	if err := s.reservationRepository.Cancel(ctx, reservationID); err != nil {
		return fmt.Errorf("cancelling reservation: %w", err)
	}

	return nil
}

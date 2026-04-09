package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// CancelAnyReservation cancels any reservation without time restriction (admin only)
func (s *adminService) CancelAnyReservation(ctx context.Context, reservationID uuid.UUID) error {
	reservation, err := s.reservationRepository.FindByID(ctx, reservationID)
	if err != nil {
		return fmt.Errorf("finding reservation: %w", err)
	}
	if reservation == nil {
		return errors.New("reservation not found")
	}
	if reservation.Status == "cancelled" {
		return errors.New("reservation is already cancelled")
	}

	return s.reservationRepository.Cancel(ctx, reservationID)
}

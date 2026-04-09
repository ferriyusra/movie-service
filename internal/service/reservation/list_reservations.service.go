package reservation

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// ListMyReservations returns a paginated list of the user's reservations
func (s *reservationService) ListMyReservations(ctx context.Context, userID uuid.UUID, status string, page, limit int) ([]response.ReservationResponse, int64, error) {
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	reservations, total, err := s.reservationRepository.FindByUserID(ctx, userID, status, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("listing reservations: %w", err)
	}

	result := make([]response.ReservationResponse, len(reservations))
	for i, r := range reservations {
		result[i] = response.ReservationResponse{
			ID:               r.ID,
			BookingReference: r.BookingReference,
			ShowtimeID:       r.ShowtimeID,
			Status:           r.Status,
			TotalAmount:      r.TotalAmount,
			CancelledAt:      r.CancelledAt,
			CreatedAt:        r.CreatedAt,
		}
	}

	return result, total, nil
}

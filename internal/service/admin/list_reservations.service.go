package admin

import (
	"context"
	"fmt"

	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

func (s *adminService) ListAllReservations(ctx context.Context, filters interfaces.ReservationFilters) ([]response.ReservationResponse, int64, error) {
	reservations, total, err := s.reservationRepository.FindAll(ctx, filters)
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

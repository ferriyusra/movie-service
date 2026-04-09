package reservation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// CreateReservation creates a reservation within a single DB transaction
func (s *reservationService) CreateReservation(ctx context.Context, userID uuid.UUID, req *request.CreateReservationRequest) (*response.ReservationResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if len(req.SeatIDs) < 1 || len(req.SeatIDs) > 10 {
		return nil, errors.New("at least 1 and at most 10 seats required")
	}

	showtime, err := s.showtimeRepository.FindByID(ctx, req.ShowtimeID)
	if err != nil {
		return nil, fmt.Errorf("finding showtime: %w", err)
	}
	if showtime == nil {
		return nil, errors.New("showtime not found")
	}

	if showtime.StartTime.Before(time.Now()) {
		return nil, errors.New("cannot reserve seats for a past showtime")
	}

	// Begin transaction
	tx := s.reservationRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check seat availability within transaction
	taken, err := s.reservationSeatRepository.CheckAvailabilityWithTx(ctx, tx, req.ShowtimeID, req.SeatIDs)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("checking availability: %w", err)
	}
	if len(taken) > 0 {
		tx.Rollback()
		labels := make([]string, len(taken))
		for i, t := range taken {
			labels[i] = t.Seat.Label
		}
		return nil, fmt.Errorf("seats already reserved: %v", labels)
	}

	// Calculate total amount
	totalAmount := float64(len(req.SeatIDs)) * showtime.Price

	// Generate booking reference: CINE-YYYYMMDD-XXXXX
	bookingRef := fmt.Sprintf("CINE-%s-%s",
		time.Now().Format("20060102"),
		uuid.New().String()[:5],
	)

	reservation := entity.ReservationEntity{
		ID:               uuid.New(),
		BookingReference: bookingRef,
		UserID:           userID,
		ShowtimeID:       req.ShowtimeID,
		Status:           "confirmed",
		TotalAmount:      totalAmount,
	}

	id, err := s.reservationRepository.CreateWithTx(ctx, tx, reservation)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("creating reservation: %w", err)
	}

	// Create reservation seats
	reservationSeats := make([]entity.ReservationSeatEntity, len(req.SeatIDs))
	for i, seatID := range req.SeatIDs {
		reservationSeats[i] = entity.ReservationSeatEntity{
			ID:            uuid.New(),
			ReservationID: *id,
			SeatID:        seatID,
			ShowtimeID:    req.ShowtimeID,
		}
	}

	if err := s.reservationSeatRepository.CreateBatchWithTx(ctx, tx, reservationSeats); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("creating reservation seats: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	// Fetch seat details for response
	seats, err := s.seatRepository.FindByTheaterID(ctx, showtime.TheaterID)
	if err != nil {
		return nil, fmt.Errorf("fetching seats: %w", err)
	}

	seatMap := make(map[uuid.UUID]entity.SeatEntity)
	for _, seat := range seats {
		seatMap[seat.ID] = seat
	}

	seatResponses := make([]response.SeatResponse, 0, len(req.SeatIDs))
	for _, seatID := range req.SeatIDs {
		if seat, ok := seatMap[seatID]; ok {
			seatResponses = append(seatResponses, response.SeatResponse{
				ID:     seat.ID,
				Row:    seat.Row,
				Number: seat.Number,
				Label:  seat.Label,
				Type:   seat.Type,
			})
		}
	}

	return &response.ReservationResponse{
		ID:               *id,
		BookingReference: bookingRef,
		ShowtimeID:       req.ShowtimeID,
		Status:           "confirmed",
		TotalAmount:      totalAmount,
		CreatedAt:        reservation.CreatedAt,
		Seats:            seatResponses,
	}, nil
}

package reservation

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

// ReservationService defines the interface for reservation operations
type ReservationService interface {
	CreateReservation(ctx context.Context, userID uuid.UUID, req *request.CreateReservationRequest) (*response.ReservationResponse, error)
	CancelReservation(ctx context.Context, userID uuid.UUID, reservationID uuid.UUID) error
	GetReservation(ctx context.Context, userID uuid.UUID, reservationID uuid.UUID) (*response.ReservationDetailResponse, error)
	ListMyReservations(ctx context.Context, userID uuid.UUID, status string, page, limit int) ([]response.ReservationResponse, int64, error)
}

// reservationService is the concrete implementation of ReservationService
type reservationService struct {
	reservationRepository     interfaces.ReservationRepository
	reservationSeatRepository interfaces.ReservationSeatRepository
	showtimeRepository        interfaces.ShowtimeRepository
	seatRepository            interfaces.SeatRepository
}

// NewReservationService creates a new instance of ReservationService
func NewReservationService(
	reservationRepository interfaces.ReservationRepository,
	reservationSeatRepository interfaces.ReservationSeatRepository,
	showtimeRepository interfaces.ShowtimeRepository,
	seatRepository interfaces.SeatRepository,
) ReservationService {
	return &reservationService{
		reservationRepository:     reservationRepository,
		reservationSeatRepository: reservationSeatRepository,
		showtimeRepository:        showtimeRepository,
		seatRepository:            seatRepository,
	}
}

package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

// AdminService defines the interface for admin operations
type AdminService interface {
	ListAllReservations(ctx context.Context, filters interfaces.ReservationFilters) ([]response.ReservationResponse, int64, error)
	CancelAnyReservation(ctx context.Context, reservationID uuid.UUID) error
	PromoteUser(ctx context.Context, adminID uuid.UUID, userID uuid.UUID) error
	ListUsers(ctx context.Context) ([]response.GetUser, error)
}

// adminService is the concrete implementation of AdminService
type adminService struct {
	reservationRepository interfaces.ReservationRepository
	userRepository        interfaces.UserRepository
}

// NewAdminService creates a new instance of AdminService
func NewAdminService(
	reservationRepository interfaces.ReservationRepository,
	userRepository interfaces.UserRepository,
) AdminService {
	return &adminService{
		reservationRepository: reservationRepository,
		userRepository:        userRepository,
	}
}

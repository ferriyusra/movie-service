package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// ReservationRepository defines the interface for reservation data access
type ReservationRepository interface {
	CreateWithTx(ctx context.Context, tx *gorm.DB, reservation entity.ReservationEntity) (*uuid.UUID, error)
	Cancel(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ReservationEntity, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, status string, page, limit int) ([]entity.ReservationEntity, int64, error)
	FindAll(ctx context.Context, filters ReservationFilters) ([]entity.ReservationEntity, int64, error)
	BeginTx(ctx context.Context) *gorm.DB
}

// ReservationFilters holds filter parameters for admin reservation listing
type ReservationFilters struct {
	UserID     *uuid.UUID
	MovieID    *uuid.UUID
	ShowtimeID *uuid.UUID
	Status     string
	DateFrom   string
	DateTo     string
	Page       int
	Limit      int
}

package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// ReservationSeatRepository defines the interface for reservation seat data access
type ReservationSeatRepository interface {
	CreateBatchWithTx(ctx context.Context, tx *gorm.DB, seats []entity.ReservationSeatEntity) error
	FindByShowtimeID(ctx context.Context, showtimeID uuid.UUID) ([]entity.ReservationSeatEntity, error)
	CheckAvailabilityWithTx(ctx context.Context, tx *gorm.DB, showtimeID uuid.UUID, seatIDs []uuid.UUID) ([]entity.ReservationSeatEntity, error)
}

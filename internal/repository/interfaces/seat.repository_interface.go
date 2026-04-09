package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// SeatRepository defines the interface for seat data access
type SeatRepository interface {
	CreateBatch(ctx context.Context, seats []entity.SeatEntity) error
	FindByTheaterID(ctx context.Context, theaterID uuid.UUID) ([]entity.SeatEntity, error)
}

package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// TheaterRepository defines the interface for theater data access
type TheaterRepository interface {
	Create(ctx context.Context, theater entity.TheaterEntity) (*uuid.UUID, error)
	Update(ctx context.Context, theater entity.TheaterEntity) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.TheaterEntity, error)
	FindAll(ctx context.Context) ([]entity.TheaterEntity, error)
}

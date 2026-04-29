package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// ShowtimeFilters holds filter parameters for admin showtime listing
type ShowtimeFilters struct {
	MovieTitle string
	DateFrom   string
	DateTo     string
	Page       int
	Limit      int
}

// ShowtimeRepository defines the interface for showtime data access
type ShowtimeRepository interface {
	Create(ctx context.Context, showtime entity.ShowtimeEntity) (*uuid.UUID, error)
	Update(ctx context.Context, showtime entity.ShowtimeEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ShowtimeEntity, error)
	FindByDate(ctx context.Context, date time.Time) ([]entity.ShowtimeEntity, error)
	FindAll(ctx context.Context, filters ShowtimeFilters) ([]entity.ShowtimeEntity, int64, error)
	CheckConflict(ctx context.Context, theaterID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (*entity.ShowtimeEntity, error)
}

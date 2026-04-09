package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// GenreRepository defines the interface for genre data access
type GenreRepository interface {
	Create(ctx context.Context, genre entity.GenreEntity) (*uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context) ([]entity.GenreEntity, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.GenreEntity, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.GenreEntity, error)
}

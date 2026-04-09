package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// RefreshTokenRepository defines the interface for refresh token data access
type RefreshTokenRepository interface {
	Create(ctx context.Context, token entity.RefreshTokenEntity) error
	FindByToken(ctx context.Context, tokenString string) (*entity.RefreshTokenEntity, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteByToken(ctx context.Context, tokenString string) error
}

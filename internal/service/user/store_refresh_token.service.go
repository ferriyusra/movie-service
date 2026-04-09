package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// storeRefreshToken persists a refresh token in the database.
func (s *userService) storeRefreshToken(ctx context.Context, userID uuid.UUID, tokenStr string) error {
	token := entity.RefreshTokenEntity{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(s.tokenService.RefreshTokenExpiry()),
	}
	if err := s.refreshTokenRepository.Create(ctx, token); err != nil {
		return fmt.Errorf("storing refresh token: %w", err)
	}
	return nil
}

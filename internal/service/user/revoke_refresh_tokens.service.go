package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RevokeRefreshTokens deletes all refresh tokens for a user
func (s *userService) RevokeRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := s.refreshTokenRepository.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("revoking refresh tokens: %w", err)
	}
	return nil
}

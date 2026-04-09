package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// Refresh validates a refresh token, rotates it, and returns a new access token and refresh token.
func (s *userService) Refresh(ctx context.Context, refreshToken string) (*response.RefreshResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	storedToken, err := s.refreshTokenRepository.FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("verifying refresh token: %w", err)
	}
	if storedToken == nil {
		return nil, errors.New("refresh token has been revoked")
	}
	if storedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	accessToken, err := s.tokenService.GenerateAccessToken(claims.UserID, claims.Email, claims.Name)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	// Rotate refresh token — revoke the old one and issue a new one
	if err := s.refreshTokenRepository.DeleteByToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("revoking old refresh token: %w", err)
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	if err := s.storeRefreshToken(ctx, claims.UserID, newRefreshToken); err != nil {
		return nil, err
	}

	return &response.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// GetUser retrieves user information by ID
func (s *userService) GetUser(ctx context.Context, userID string) (*response.GetUser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &response.GetUser{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

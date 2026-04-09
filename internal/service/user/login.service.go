package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"golang.org/x/crypto/bcrypt"
)

// Login authenticates a user, generates auth tokens, and returns the full auth response.
func (s *userService) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	user, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("finding user by email: %w", err)
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	if err := s.storeRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		User: response.GetUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

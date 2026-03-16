package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/ferriyusra/clean-arch-go-gin/internal/model/request"
	"github.com/ferriyusra/clean-arch-go-gin/internal/model/response"
	"golang.org/x/crypto/bcrypt"
)

// Login authenticates a user and returns user info (tokens set via cookies in handler)
func (s *userService) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Find user by email
	user, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("finding user by email: %w", err)
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return &response.LoginResponse{
		User: response.GetUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}

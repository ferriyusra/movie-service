package admin

import (
	"context"
	"fmt"

	"github.com/ferriyusra/movie-service/internal/model/response"
)

// ListUsers returns all users
func (s *adminService) ListUsers(ctx context.Context) ([]response.GetUser, error) {
	users, err := s.userRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}

	result := make([]response.GetUser, len(users))
	for i, u := range users {
		result[i] = response.GetUser{
			ID:    u.ID,
			Email: u.Email,
			Name:  u.Name,
		}
	}

	return result, nil
}

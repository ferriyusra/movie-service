package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// PromoteUser promotes a user to admin, preventing self-demotion
func (s *adminService) PromoteUser(ctx context.Context, adminID uuid.UUID, userID uuid.UUID) error {
	if adminID == userID {
		return errors.New("cannot promote yourself")
	}

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}
	if user.Role == "admin" {
		return errors.New("user is already an admin")
	}

	return s.userRepository.Update(ctx, userID, entity.UserEntity{Role: "admin"})
}

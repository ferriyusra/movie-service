package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// FindByID finds a user by ID in GORM
func (r *GORMUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.UserEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var user entity.UserEntity
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("finding user by id: %w", err)
	}

	return &user, nil
}

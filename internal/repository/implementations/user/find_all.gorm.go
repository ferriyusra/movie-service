package user

import (
	"context"

	"github.com/ferriyusra/movie-service/internal/model/entity"
)

// FindAll returns all users
func (r *GORMUserRepository) FindAll(ctx context.Context) ([]entity.UserEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var users []entity.UserEntity
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

package refresh_token

import (
	"context"
	"fmt"

	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// FindByToken finds a refresh token by its token string
func (r *GORMRefreshTokenRepository) FindByToken(ctx context.Context, tokenString string) (*entity.RefreshTokenEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var token entity.RefreshTokenEntity
	if err := r.db.WithContext(ctx).Where("token = ?", tokenString).First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("finding refresh token: %w", err)
	}
	return &token, nil
}

package refresh_token

import (
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// GORMRefreshTokenRepository is a GORM implementation of RefreshTokenRepository
type GORMRefreshTokenRepository struct {
	db *gorm.DB
}

// RefreshTokenModel represents the refresh_token_entities table schema
type RefreshTokenModel = entity.RefreshTokenEntity

// NewGORMRefreshTokenRepository creates a new GORM refresh token repository
func NewGORMRefreshTokenRepository(db *gorm.DB) (*GORMRefreshTokenRepository, error) {
	if err := db.AutoMigrate(&RefreshTokenModel{}); err != nil {
		return nil, err
	}
	return &GORMRefreshTokenRepository{db: db}, nil
}

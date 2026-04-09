package movie

import (
	"gorm.io/gorm"
)

// GORMMovieRepository is a GORM implementation of MovieRepository
type GORMMovieRepository struct {
	db *gorm.DB
}

// NewGORMMovieRepository creates a new GORM movie repository
func NewGORMMovieRepository(db *gorm.DB) *GORMMovieRepository {
	return &GORMMovieRepository{db: db}
}

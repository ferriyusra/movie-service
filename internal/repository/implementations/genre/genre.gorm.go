package genre

import (
	"gorm.io/gorm"
)

// GORMGenreRepository is a GORM implementation of GenreRepository
type GORMGenreRepository struct {
	db *gorm.DB
}

// NewGORMGenreRepository creates a new GORM genre repository
func NewGORMGenreRepository(db *gorm.DB) *GORMGenreRepository {
	return &GORMGenreRepository{db: db}
}

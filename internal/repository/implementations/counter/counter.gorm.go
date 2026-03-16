package counter

import (
	"github.com/google/uuid"
	"github.com/ferriyusra/clean-arch-go-gin/internal/model/entity"
	"gorm.io/gorm"
)

// GORMCounterRepository is a GORM implementation of CounterRepository
type GORMCounterRepository struct {
	db *gorm.DB
}

type CounterModel = entity.CounterEntity

// NewGORMCounterRepository creates a new GORM counter repository
func NewGORMCounterRepository(db *gorm.DB) (*GORMCounterRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&CounterModel{}); err != nil {
		return nil, err
	}

	// Initialize counter if it doesn't exist
	var count int64
	db.Model(&CounterModel{}).Count(&count)
	if count == 0 {
		db.Create(&CounterModel{ID: uuid.New(), Value: 0})
	}

	return &GORMCounterRepository{
		db: db,
	}, nil
}

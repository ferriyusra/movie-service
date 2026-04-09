package platform

import (
	"fmt"

	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDatabase(cfg *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	dbConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	if cfg.Database.Type == "postgres" {
		dialector = postgres.Open(cfg.Database.DSN)
	} else {
		dialector = sqlite.Open(cfg.Database.DSN)
	}
	db, err := gorm.Open(dialector, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	sqlDB, err := db.DB() // Get the underlying generic *sql.DB
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings from your config
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Auto-migrate all entity schemas
	if err := db.AutoMigrate(
		&entity.UserEntity{},
		&entity.RefreshTokenEntity{},
		&entity.GenreEntity{},
		&entity.MovieEntity{},
		&entity.TheaterEntity{},
		&entity.SeatEntity{},
		&entity.ShowtimeEntity{},
		&entity.ReservationEntity{},
		&entity.ReservationSeatEntity{},
	); err != nil {
		return nil, fmt.Errorf("auto-migrating database: %w", err)
	}

	return db, nil
}

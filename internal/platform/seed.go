package platform

import (
	"log"

	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedAdminAccount creates the default admin account if it doesn't already exist.
func SeedAdminAccount(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.UserEntity{}).Where("email = ?", "admin@cinema.com").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // admin already exists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := entity.UserEntity{
		ID:       uuid.New(),
		Email:    "admin@cinema.com",
		Password: hashedPassword,
		Name:     "Admin",
		Role:     "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("Seeded admin account: admin@cinema.com")
	return nil
}

package platform

import (
	"testing"

	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&entity.UserEntity{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestSeedAdminAccount(t *testing.T) {
	tests := []struct {
		name           string
		preExisting    bool
		expectedCount  int64
	}{
		{
			name:          "should create admin when not exists",
			preExisting:   false,
			expectedCount: 1,
		},
		{
			name:          "should not duplicate admin when already exists",
			preExisting:   true,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			if tt.preExisting {
				// Seed once first
				if err := SeedAdminAccount(db); err != nil {
					t.Fatalf("first seed failed: %v", err)
				}
			}

			// Run seed
			err := SeedAdminAccount(db)
			if err != nil {
				t.Fatalf("SeedAdminAccount error: %v", err)
			}

			var count int64
			db.Model(&entity.UserEntity{}).Where("email = ?", "admin@cinema.com").Count(&count)
			if count != tt.expectedCount {
				t.Errorf("expected %d admin(s), got %d", tt.expectedCount, count)
			}

			// Verify admin fields
			var admin entity.UserEntity
			db.Where("email = ?", "admin@cinema.com").First(&admin)

			if admin.Role != "admin" {
				t.Errorf("expected role 'admin', got '%s'", admin.Role)
			}
			if admin.Name != "Admin" {
				t.Errorf("expected name 'Admin', got '%s'", admin.Name)
			}
			if err := bcrypt.CompareHashAndPassword(admin.Password, []byte("Admin@123")); err != nil {
				t.Errorf("password hash doesn't match 'Admin@123': %v", err)
			}
		})
	}
}

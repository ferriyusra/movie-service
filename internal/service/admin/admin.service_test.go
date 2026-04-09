package admin

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestCancelAnyReservation(t *testing.T) {
	reservationID := uuid.New()

	tests := []struct {
		name            string
		mockReservation *entity.ReservationEntity
		mockCancelErr   error
		expectedError   bool
		expectedMsg     string
	}{
		{
			name: "should cancel reservation successfully",
			mockReservation: &entity.ReservationEntity{
				ID:     reservationID,
				Status: "confirmed",
			},
			expectedError: false,
		},
		{
			name:            "should return error when not found",
			mockReservation: nil,
			expectedError:   true,
			expectedMsg:     "reservation not found",
		},
		{
			name: "should return error when already cancelled",
			mockReservation: &entity.ReservationEntity{
				ID:     reservationID,
				Status: "cancelled",
			},
			expectedError: true,
			expectedMsg:   "reservation is already cancelled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReservationRepo := mock.NewMockReservationRepository(ctrl)
			mockUserRepo := mock.NewMockUserRepository(ctrl)

			mockReservationRepo.EXPECT().
				FindByID(gomock.Any(), reservationID).
				Return(tt.mockReservation, nil).
				Times(1)

			if tt.mockReservation != nil && tt.mockReservation.Status != "cancelled" {
				mockReservationRepo.EXPECT().
					Cancel(gomock.Any(), reservationID).
					Return(tt.mockCancelErr).
					Times(1)
			}

			svc := NewAdminService(mockReservationRepo, mockUserRepo)
			err := svc.CancelAnyReservation(context.Background(), reservationID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.expectedMsg != "" && err.Error() != tt.expectedMsg {
					t.Errorf("expected error '%s', got '%s'", tt.expectedMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPromoteUser(t *testing.T) {
	adminID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name          string
		adminID       uuid.UUID
		userID        uuid.UUID
		mockUser      *entity.UserEntity
		expectedError bool
		expectedMsg   string
	}{
		{
			name:    "should promote user successfully",
			adminID: adminID,
			userID:  userID,
			mockUser: &entity.UserEntity{
				ID:   userID,
				Role: "user",
			},
			expectedError: false,
		},
		{
			name:          "should return error on self-promotion",
			adminID:       adminID,
			userID:        adminID,
			expectedError: true,
			expectedMsg:   "cannot promote yourself",
		},
		{
			name:          "should return error when user not found",
			adminID:       adminID,
			userID:        userID,
			mockUser:      nil,
			expectedError: true,
			expectedMsg:   "user not found",
		},
		{
			name:    "should return error when already admin",
			adminID: adminID,
			userID:  userID,
			mockUser: &entity.UserEntity{
				ID:   userID,
				Role: "admin",
			},
			expectedError: true,
			expectedMsg:   "user is already an admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReservationRepo := mock.NewMockReservationRepository(ctrl)
			mockUserRepo := mock.NewMockUserRepository(ctrl)

			if tt.adminID != tt.userID {
				mockUserRepo.EXPECT().
					FindByID(gomock.Any(), tt.userID).
					Return(tt.mockUser, nil).
					Times(1)

				if tt.mockUser != nil && tt.mockUser.Role != "admin" {
					mockUserRepo.EXPECT().
						Update(gomock.Any(), tt.userID, gomock.Any()).
						Return(nil).
						Times(1)
				}
			}

			svc := NewAdminService(mockReservationRepo, mockUserRepo)
			err := svc.PromoteUser(context.Background(), tt.adminID, tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.expectedMsg != "" && err.Error() != tt.expectedMsg {
					t.Errorf("expected error '%s', got '%s'", tt.expectedMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// Suppress unused import warning
var _ = time.Now

package reservation

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestCancelReservation(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	reservationID := uuid.New()

	tests := []struct {
		name            string
		userID          uuid.UUID
		reservationID   uuid.UUID
		mockReservation *entity.ReservationEntity
		mockCancelErr   error
		expectedError   bool
		expectedMsg     string
	}{
		{
			name:          "should cancel reservation successfully",
			userID:        userID,
			reservationID: reservationID,
			mockReservation: &entity.ReservationEntity{
				ID:     reservationID,
				UserID: userID,
				Status: "confirmed",
				Showtime: entity.ShowtimeEntity{
					StartTime: time.Now().Add(2 * time.Hour),
				},
			},
			expectedError: false,
		},
		{
			name:            "should return error when reservation not found",
			userID:          userID,
			reservationID:   reservationID,
			mockReservation: nil,
			expectedError:   true,
			expectedMsg:     "reservation not found",
		},
		{
			name:          "should return error when not owner",
			userID:        otherUserID,
			reservationID: reservationID,
			mockReservation: &entity.ReservationEntity{
				ID:     reservationID,
				UserID: userID,
				Status: "confirmed",
				Showtime: entity.ShowtimeEntity{
					StartTime: time.Now().Add(2 * time.Hour),
				},
			},
			expectedError: true,
			expectedMsg:   "you can only cancel your own reservations",
		},
		{
			name:          "should return error when already cancelled",
			userID:        userID,
			reservationID: reservationID,
			mockReservation: &entity.ReservationEntity{
				ID:     reservationID,
				UserID: userID,
				Status: "cancelled",
				Showtime: entity.ShowtimeEntity{
					StartTime: time.Now().Add(2 * time.Hour),
				},
			},
			expectedError: true,
			expectedMsg:   "reservation is already cancelled",
		},
		{
			name:          "should return error when within 1 hour of showtime",
			userID:        userID,
			reservationID: reservationID,
			mockReservation: &entity.ReservationEntity{
				ID:     reservationID,
				UserID: userID,
				Status: "confirmed",
				Showtime: entity.ShowtimeEntity{
					StartTime: time.Now().Add(30 * time.Minute),
				},
			},
			expectedError: true,
			expectedMsg:   "cancellation window has passed (must be at least 1 hour before showtime)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReservationRepo := mock.NewMockReservationRepository(ctrl)
			mockReservationSeatRepo := mock.NewMockReservationSeatRepository(ctrl)
			mockShowtimeRepo := mock.NewMockShowtimeRepository(ctrl)
			mockSeatRepo := mock.NewMockSeatRepository(ctrl)

			mockReservationRepo.EXPECT().
				FindByID(gomock.Any(), tt.reservationID).
				Return(tt.mockReservation, nil).
				Times(1)

			if tt.mockReservation != nil &&
				tt.mockReservation.UserID == tt.userID &&
				tt.mockReservation.Status == "confirmed" &&
				tt.mockReservation.Showtime.StartTime.After(time.Now().Add(1*time.Hour)) {
				mockReservationRepo.EXPECT().
					Cancel(gomock.Any(), tt.reservationID).
					Return(tt.mockCancelErr).
					Times(1)
			}

			svc := NewReservationService(mockReservationRepo, mockReservationSeatRepo, mockShowtimeRepo, mockSeatRepo)
			err := svc.CancelReservation(context.Background(), tt.userID, tt.reservationID)

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

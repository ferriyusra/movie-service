package theater

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestGetTheater(t *testing.T) {
	testID := uuid.New()

	tests := []struct {
		name           string
		id             uuid.UUID
		mockTheater    *entity.TheaterEntity
		mockTheaterErr error
		mockSeats      []entity.SeatEntity
		mockSeatErr    error
		expectedError  bool
		expectedMsg    string
	}{
		{
			name: "should get theater with seats",
			id:   testID,
			mockTheater: &entity.TheaterEntity{
				ID: testID, Name: "Theater 1", TotalSeats: 2, Rows: 1, SeatsPerRow: 2,
			},
			mockSeats: []entity.SeatEntity{
				{ID: uuid.New(), TheaterID: testID, Row: "A", Number: 1, Label: "A1", Type: "standard"},
				{ID: uuid.New(), TheaterID: testID, Row: "A", Number: 2, Label: "A2", Type: "standard"},
			},
			expectedError: false,
		},
		{
			name:          "should return error when theater not found",
			id:            testID,
			mockTheater:   nil,
			expectedError: true,
			expectedMsg:   "theater not found",
		},
		{
			name:           "should return error when find fails",
			id:             testID,
			mockTheater:    nil,
			mockTheaterErr: errors.New("db error"),
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTheaterRepo := mock.NewMockTheaterRepository(ctrl)
			mockSeatRepo := mock.NewMockSeatRepository(ctrl)

			mockTheaterRepo.EXPECT().
				FindByID(gomock.Any(), tt.id).
				Return(tt.mockTheater, tt.mockTheaterErr).
				Times(1)

			if tt.mockTheater != nil && tt.mockTheaterErr == nil {
				mockSeatRepo.EXPECT().
					FindByTheaterID(gomock.Any(), tt.id).
					Return(tt.mockSeats, tt.mockSeatErr).
					Times(1)
			}

			svc := NewTheaterService(mockTheaterRepo, mockSeatRepo)
			result, err := svc.GetTheater(context.Background(), tt.id)

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
				if result == nil {
					t.Fatal("expected non-nil result")
				}
				if len(result.Seats) != len(tt.mockSeats) {
					t.Errorf("expected %d seats, got %d", len(tt.mockSeats), len(result.Seats))
				}
			}
		})
	}
}

package theater

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestCreateTheater(t *testing.T) {
	tests := []struct {
		name            string
		request         *request.CreateTheaterRequest
		mockCreateID    *uuid.UUID
		mockCreateErr   error
		mockSeatErr     error
		expectedError   bool
		expectedMsg     string
		expectedSeats   int
	}{
		{
			name: "should create theater with auto-generated seats",
			request: &request.CreateTheaterRequest{
				Name:        "Theater 1",
				Location:    "Floor 1",
				Rows:        3,
				SeatsPerRow: 5,
			},
			mockCreateID: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
			expectedError: false,
			expectedSeats: 15, // 3 rows * 5 seats
		},
		{
			name: "should return error when name is empty",
			request: &request.CreateTheaterRequest{
				Name:        "",
				Rows:        3,
				SeatsPerRow: 5,
			},
			expectedError: true,
			expectedMsg:   "theater name is required",
		},
		{
			name: "should return error when rows is zero",
			request: &request.CreateTheaterRequest{
				Name:        "Theater 1",
				Rows:        0,
				SeatsPerRow: 5,
			},
			expectedError: true,
			expectedMsg:   "rows must be between 1 and 26",
		},
		{
			name: "should return error when rows exceeds 26",
			request: &request.CreateTheaterRequest{
				Name:        "Theater 1",
				Rows:        27,
				SeatsPerRow: 5,
			},
			expectedError: true,
			expectedMsg:   "rows must be between 1 and 26",
		},
		{
			name: "should return error when seats_per_row is zero",
			request: &request.CreateTheaterRequest{
				Name:        "Theater 1",
				Rows:        3,
				SeatsPerRow: 0,
			},
			expectedError: true,
			expectedMsg:   "seats_per_row must be between 1 and 30",
		},
		{
			name: "should return error when theater create fails",
			request: &request.CreateTheaterRequest{
				Name:        "Theater 1",
				Location:    "Floor 1",
				Rows:        3,
				SeatsPerRow: 5,
			},
			mockCreateID:  nil,
			mockCreateErr: errors.New("db error"),
			expectedError: true,
		},
		{
			name: "should return error when seat batch create fails",
			request: &request.CreateTheaterRequest{
				Name:        "Theater 1",
				Location:    "Floor 1",
				Rows:        3,
				SeatsPerRow: 5,
			},
			mockCreateID: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
			mockSeatErr:   errors.New("seat db error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTheaterRepo := mock.NewMockTheaterRepository(ctrl)
			mockSeatRepo := mock.NewMockSeatRepository(ctrl)

			passesValidation := tt.request.Name != "" &&
				tt.request.Rows >= 1 && tt.request.Rows <= 26 &&
				tt.request.SeatsPerRow >= 1 && tt.request.SeatsPerRow <= 30

			if passesValidation {
				mockTheaterRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(entity.TheaterEntity{})).
					Return(tt.mockCreateID, tt.mockCreateErr).
					Times(1)

				if tt.mockCreateErr == nil && tt.mockCreateID != nil {
					mockSeatRepo.EXPECT().
						CreateBatch(gomock.Any(), gomock.Any()).
						Return(tt.mockSeatErr).
						Times(1)
				}
			}

			svc := NewTheaterService(mockTheaterRepo, mockSeatRepo)
			result, err := svc.CreateTheater(context.Background(), tt.request)

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
				if result.Name != tt.request.Name {
					t.Errorf("expected name '%s', got '%s'", tt.request.Name, result.Name)
				}
				if len(result.Seats) != tt.expectedSeats {
					t.Errorf("expected %d seats, got %d", tt.expectedSeats, len(result.Seats))
				}
				if result.TotalSeats != tt.expectedSeats {
					t.Errorf("expected total_seats %d, got %d", tt.expectedSeats, result.TotalSeats)
				}
			}
		})
	}
}

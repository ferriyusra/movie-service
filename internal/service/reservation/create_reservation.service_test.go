package reservation

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestCreateReservation(t *testing.T) {
	userID := uuid.New()
	showtimeID := uuid.New()
	seatID1 := uuid.New()
	seatID2 := uuid.New()

	tests := []struct {
		name          string
		request       *request.CreateReservationRequest
		mockShowtime  *entity.ShowtimeEntity
		expectedError bool
		expectedMsg   string
	}{
		{
			name: "should return error when no seats provided",
			request: &request.CreateReservationRequest{
				ShowtimeID: showtimeID,
				SeatIDs:    []uuid.UUID{},
			},
			expectedError: true,
			expectedMsg:   "at least 1 and at most 10 seats required",
		},
		{
			name: "should return error when more than 10 seats",
			request: &request.CreateReservationRequest{
				ShowtimeID: showtimeID,
				SeatIDs:    make([]uuid.UUID, 11),
			},
			expectedError: true,
			expectedMsg:   "at least 1 and at most 10 seats required",
		},
		{
			name: "should return error when showtime not found",
			request: &request.CreateReservationRequest{
				ShowtimeID: showtimeID,
				SeatIDs:    []uuid.UUID{seatID1, seatID2},
			},
			mockShowtime:  nil,
			expectedError: true,
			expectedMsg:   "showtime not found",
		},
		{
			name: "should return error when showtime is in the past",
			request: &request.CreateReservationRequest{
				ShowtimeID: showtimeID,
				SeatIDs:    []uuid.UUID{seatID1},
			},
			mockShowtime: &entity.ShowtimeEntity{
				ID:        showtimeID,
				StartTime: time.Now().Add(-1 * time.Hour),
				Price:     50000,
				Movie:     entity.MovieEntity{Title: "Test"},
				Theater:   entity.TheaterEntity{Name: "Theater 1"},
			},
			expectedError: true,
			expectedMsg:   "cannot reserve seats for a past showtime",
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

			if len(tt.request.SeatIDs) >= 1 && len(tt.request.SeatIDs) <= 10 {
				mockShowtimeRepo.EXPECT().
					FindByID(gomock.Any(), tt.request.ShowtimeID).
					Return(tt.mockShowtime, nil).
					Times(1)
			}

			svc := NewReservationService(mockReservationRepo, mockReservationSeatRepo, mockShowtimeRepo, mockSeatRepo)
			result, err := svc.CreateReservation(context.Background(), userID, tt.request)

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
			}
		})
	}
}

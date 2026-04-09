package showtime

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestCreateShowtime(t *testing.T) {
	movieID := uuid.New()
	theaterID := uuid.New()
	startTime := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)

	tests := []struct {
		name             string
		request          *request.CreateShowtimeRequest
		mockMovie        *entity.MovieEntity
		mockMovieErr     error
		mockTheater      *entity.TheaterEntity
		mockTheaterErr   error
		mockConflict     *entity.ShowtimeEntity
		mockConflictErr  error
		mockCreateID     *uuid.UUID
		mockCreateErr    error
		expectedError    bool
		expectedMsg      string
	}{
		{
			name: "should create showtime successfully",
			request: &request.CreateShowtimeRequest{
				MovieID:   movieID,
				TheaterID: theaterID,
				StartTime: startTime,
				Price:     50000,
			},
			mockMovie:   &entity.MovieEntity{ID: movieID, Title: "Inception", DurationMin: 148},
			mockTheater: &entity.TheaterEntity{ID: theaterID, Name: "Theater 1"},
			mockConflict: nil,
			mockCreateID: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
			expectedError: false,
		},
		{
			name: "should return error when movie not found",
			request: &request.CreateShowtimeRequest{
				MovieID:   movieID,
				TheaterID: theaterID,
				StartTime: startTime,
				Price:     50000,
			},
			mockMovie:     nil,
			expectedError: true,
			expectedMsg:   "movie not found",
		},
		{
			name: "should return error when theater not found",
			request: &request.CreateShowtimeRequest{
				MovieID:   movieID,
				TheaterID: theaterID,
				StartTime: startTime,
				Price:     50000,
			},
			mockMovie:     &entity.MovieEntity{ID: movieID, DurationMin: 148},
			mockTheater:   nil,
			expectedError: true,
			expectedMsg:   "theater not found",
		},
		{
			name: "should return error on conflict",
			request: &request.CreateShowtimeRequest{
				MovieID:   movieID,
				TheaterID: theaterID,
				StartTime: startTime,
				Price:     50000,
			},
			mockMovie:    &entity.MovieEntity{ID: movieID, DurationMin: 148},
			mockTheater:  &entity.TheaterEntity{ID: theaterID},
			mockConflict: &entity.ShowtimeEntity{ID: uuid.New()},
			expectedError: true,
			expectedMsg:   "showtime conflicts with an existing showtime in this theater",
		},
		{
			name: "should return error when price is zero",
			request: &request.CreateShowtimeRequest{
				MovieID:   movieID,
				TheaterID: theaterID,
				StartTime: startTime,
				Price:     0,
			},
			expectedError: true,
			expectedMsg:   "price must be greater than zero",
		},
		{
			name: "should return error when start_time is invalid",
			request: &request.CreateShowtimeRequest{
				MovieID:   movieID,
				TheaterID: theaterID,
				StartTime: "not-a-date",
				Price:     50000,
			},
			expectedError: true,
			expectedMsg:   "invalid start_time format, use RFC3339",
		},
		{
			name: "should return error when create fails",
			request: &request.CreateShowtimeRequest{
				MovieID:   movieID,
				TheaterID: theaterID,
				StartTime: startTime,
				Price:     50000,
			},
			mockMovie:    &entity.MovieEntity{ID: movieID, DurationMin: 148},
			mockTheater:  &entity.TheaterEntity{ID: theaterID},
			mockConflict: nil,
			mockCreateErr: errors.New("db error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockShowtimeRepo := mock.NewMockShowtimeRepository(ctrl)
			mockMovieRepo := mock.NewMockMovieRepository(ctrl)
			mockTheaterRepo := mock.NewMockTheaterRepository(ctrl)

			// Price and start_time validation happen first
			if tt.request.Price <= 0 || tt.expectedMsg == "invalid start_time format, use RFC3339" {
				// No repo calls expected
			} else {
				mockMovieRepo.EXPECT().
					FindByID(gomock.Any(), tt.request.MovieID).
					Return(tt.mockMovie, tt.mockMovieErr).
					Times(1)

				if tt.mockMovie != nil {
					mockTheaterRepo.EXPECT().
						FindByID(gomock.Any(), tt.request.TheaterID).
						Return(tt.mockTheater, tt.mockTheaterErr).
						Times(1)

					if tt.mockTheater != nil {
						mockShowtimeRepo.EXPECT().
							CheckConflict(gomock.Any(), tt.request.TheaterID, gomock.Any(), gomock.Any(), gomock.Nil()).
							Return(tt.mockConflict, tt.mockConflictErr).
							Times(1)

						if tt.mockConflict == nil && tt.mockConflictErr == nil {
							mockShowtimeRepo.EXPECT().
								Create(gomock.Any(), gomock.Any()).
								Return(tt.mockCreateID, tt.mockCreateErr).
								Times(1)
						}
					}
				}
			}

			mockSeatRepo := mock.NewMockSeatRepository(ctrl)
			mockResSeatRepo := mock.NewMockReservationSeatRepository(ctrl)
			svc := NewShowtimeService(mockShowtimeRepo, mockMovieRepo, mockTheaterRepo, mockSeatRepo, mockResSeatRepo)
			result, err := svc.CreateShowtime(context.Background(), tt.request)

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
				if result.Price != tt.request.Price {
					t.Errorf("expected price %f, got %f", tt.request.Price, result.Price)
				}
			}
		})
	}
}

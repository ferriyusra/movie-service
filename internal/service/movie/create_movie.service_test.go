package movie

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

func TestCreateMovie(t *testing.T) {
	genreID1 := uuid.New()
	genreID2 := uuid.New()

	tests := []struct {
		name            string
		request         *request.CreateMovieRequest
		mockGenres      []entity.GenreEntity
		mockGenresErr   error
		mockCreateID    *uuid.UUID
		mockCreateErr   error
		mockReplaceErr  error
		expectedError   bool
		expectedMsg     string
	}{
		{
			name: "should create movie successfully",
			request: &request.CreateMovieRequest{
				Title:       "Inception",
				Description: "A mind-bending thriller",
				DurationMin: 148,
				Language:    "English",
				ReleaseDate: "2010-07-16",
				GenreIDs:    []uuid.UUID{genreID1, genreID2},
			},
			mockGenres: []entity.GenreEntity{
				{ID: genreID1, Name: "Sci-Fi"},
				{ID: genreID2, Name: "Action"},
			},
			mockCreateID: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
			expectedError: false,
		},
		{
			name: "should return error when title is empty",
			request: &request.CreateMovieRequest{
				Title:       "",
				DurationMin: 120,
			},
			expectedError: true,
			expectedMsg:   "movie title is required",
		},
		{
			name: "should return error when duration is zero",
			request: &request.CreateMovieRequest{
				Title:       "Inception",
				DurationMin: 0,
			},
			expectedError: true,
			expectedMsg:   "duration must be greater than zero",
		},
		{
			name: "should return error when genre not found",
			request: &request.CreateMovieRequest{
				Title:       "Inception",
				DurationMin: 148,
				GenreIDs:    []uuid.UUID{genreID1, genreID2},
			},
			mockGenres:    []entity.GenreEntity{{ID: genreID1, Name: "Sci-Fi"}},
			mockGenresErr: nil,
			expectedError: true,
			expectedMsg:   "one or more genre IDs are invalid",
		},
		{
			name: "should return error when repository create fails",
			request: &request.CreateMovieRequest{
				Title:       "Inception",
				DurationMin: 148,
				GenreIDs:    []uuid.UUID{genreID1},
			},
			mockGenres:    []entity.GenreEntity{{ID: genreID1, Name: "Sci-Fi"}},
			mockCreateID:  nil,
			mockCreateErr: errors.New("db error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMovieRepo := mock.NewMockMovieRepository(ctrl)
			mockGenreRepo := mock.NewMockGenreRepository(ctrl)

			// Only expect genre lookup when we pass validation
			if tt.request.Title != "" && tt.request.DurationMin > 0 && len(tt.request.GenreIDs) > 0 {
				mockGenreRepo.EXPECT().
					FindByIDs(gomock.Any(), tt.request.GenreIDs).
					Return(tt.mockGenres, tt.mockGenresErr).
					Times(1)
			}

			// Only expect create when genres are valid
			if tt.mockGenres != nil && len(tt.mockGenres) == len(tt.request.GenreIDs) && tt.mockGenresErr == nil {
				mockMovieRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(tt.mockCreateID, tt.mockCreateErr).
					Times(1)

				if tt.mockCreateErr == nil && tt.mockCreateID != nil {
					mockMovieRepo.EXPECT().
						ReplaceGenres(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(tt.mockReplaceErr).
						Times(1)
				}
			}

			svc := NewMovieService(mockMovieRepo, mockGenreRepo)
			result, err := svc.CreateMovie(context.Background(), tt.request)

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
				if result.Title != tt.request.Title {
					t.Errorf("expected title '%s', got '%s'", tt.request.Title, result.Title)
				}
				if len(result.Genres) != len(tt.request.GenreIDs) {
					t.Errorf("expected %d genres, got %d", len(tt.request.GenreIDs), len(result.Genres))
				}
			}
		})
	}
}

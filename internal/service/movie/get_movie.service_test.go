package movie

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestGetMovie(t *testing.T) {
	testID := uuid.New()

	tests := []struct {
		name          string
		id            uuid.UUID
		mockReturn    *entity.MovieEntity
		mockError     error
		expectedError bool
		expectedMsg   string
	}{
		{
			name: "should get movie successfully",
			id:   testID,
			mockReturn: &entity.MovieEntity{
				ID:          testID,
				Title:       "Inception",
				DurationMin: 148,
				Genres:      []entity.GenreEntity{{ID: uuid.New(), Name: "Sci-Fi"}},
			},
			expectedError: false,
		},
		{
			name:          "should return error when movie not found",
			id:            testID,
			mockReturn:    nil,
			mockError:     nil,
			expectedError: true,
			expectedMsg:   "movie not found",
		},
		{
			name:          "should return error when repository fails",
			id:            testID,
			mockReturn:    nil,
			mockError:     errors.New("database error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMovieRepo := mock.NewMockMovieRepository(ctrl)
			mockGenreRepo := mock.NewMockGenreRepository(ctrl)

			mockMovieRepo.EXPECT().
				FindByID(gomock.Any(), tt.id).
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			svc := NewMovieService(mockMovieRepo, mockGenreRepo)
			result, err := svc.GetMovie(context.Background(), tt.id)

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
				if result.Title != "Inception" {
					t.Errorf("expected title 'Inception', got '%s'", result.Title)
				}
			}
		})
	}
}

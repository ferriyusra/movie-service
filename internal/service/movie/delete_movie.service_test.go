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

func TestDeleteMovie(t *testing.T) {
	testID := uuid.New()

	tests := []struct {
		name            string
		id              uuid.UUID
		mockFindReturn  *entity.MovieEntity
		mockFindError   error
		mockDeleteError error
		expectedError   bool
		expectedMsg     string
	}{
		{
			name:           "should delete movie successfully",
			id:             testID,
			mockFindReturn: &entity.MovieEntity{ID: testID, Title: "Inception"},
			expectedError:  false,
		},
		{
			name:          "should return error when movie not found",
			id:            testID,
			mockFindReturn: nil,
			expectedError: true,
			expectedMsg:   "movie not found",
		},
		{
			name:           "should return error when delete fails",
			id:             testID,
			mockFindReturn: &entity.MovieEntity{ID: testID, Title: "Inception"},
			mockDeleteError: errors.New("delete error"),
			expectedError:  true,
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
				Return(tt.mockFindReturn, tt.mockFindError).
				Times(1)

			if tt.mockFindReturn != nil && tt.mockFindError == nil {
				mockMovieRepo.EXPECT().
					Delete(gomock.Any(), tt.id).
					Return(tt.mockDeleteError).
					Times(1)
			}

			svc := NewMovieService(mockMovieRepo, mockGenreRepo)
			err := svc.DeleteMovie(context.Background(), tt.id)

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

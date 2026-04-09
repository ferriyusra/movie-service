package movie

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestListMovies(t *testing.T) {
	tests := []struct {
		name          string
		filters       interfaces.MovieFilters
		mockReturn    []entity.MovieEntity
		mockTotal     int64
		mockError     error
		expectedLen   int
		expectedTotal int64
		expectedError bool
	}{
		{
			name:    "should list movies successfully",
			filters: interfaces.MovieFilters{Page: 1, Limit: 20},
			mockReturn: []entity.MovieEntity{
				{ID: uuid.New(), Title: "Inception", DurationMin: 148},
				{ID: uuid.New(), Title: "The Matrix", DurationMin: 136},
			},
			mockTotal:     2,
			expectedLen:   2,
			expectedTotal: 2,
			expectedError: false,
		},
		{
			name:          "should return empty list",
			filters:       interfaces.MovieFilters{Page: 1, Limit: 20},
			mockReturn:    []entity.MovieEntity{},
			mockTotal:     0,
			expectedLen:   0,
			expectedTotal: 0,
			expectedError: false,
		},
		{
			name:          "should return error when repository fails",
			filters:       interfaces.MovieFilters{Page: 1, Limit: 20},
			mockReturn:    nil,
			mockTotal:     0,
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
				FindAll(gomock.Any(), tt.filters).
				Return(tt.mockReturn, tt.mockTotal, tt.mockError).
				Times(1)

			svc := NewMovieService(mockMovieRepo, mockGenreRepo)
			result, total, err := svc.ListMovies(context.Background(), tt.filters)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) != tt.expectedLen {
					t.Errorf("expected %d movies, got %d", tt.expectedLen, len(result))
				}
				if total != tt.expectedTotal {
					t.Errorf("expected total %d, got %d", tt.expectedTotal, total)
				}
			}
		})
	}
}

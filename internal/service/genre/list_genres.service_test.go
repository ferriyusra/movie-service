package genre

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestListGenres(t *testing.T) {
	tests := []struct {
		name          string
		mockReturn    []entity.GenreEntity
		mockError     error
		expectedLen   int
		expectedError bool
	}{
		{
			name: "should list genres successfully",
			mockReturn: []entity.GenreEntity{
				{ID: uuid.New(), Name: "Action"},
				{ID: uuid.New(), Name: "Comedy"},
			},
			mockError:     nil,
			expectedLen:   2,
			expectedError: false,
		},
		{
			name:          "should return empty list",
			mockReturn:    []entity.GenreEntity{},
			mockError:     nil,
			expectedLen:   0,
			expectedError: false,
		},
		{
			name:          "should return error when repository fails",
			mockReturn:    nil,
			mockError:     errors.New("database error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockGenreRepository(ctrl)
			mockRepo.EXPECT().
				FindAll(gomock.Any()).
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			svc := NewGenreService(mockRepo)
			result, err := svc.ListGenres(context.Background())

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) != tt.expectedLen {
					t.Errorf("expected %d genres, got %d", tt.expectedLen, len(result))
				}
			}
		})
	}
}

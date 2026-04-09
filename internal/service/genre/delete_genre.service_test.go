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

func TestDeleteGenre(t *testing.T) {
	testID := uuid.New()

	tests := []struct {
		name            string
		id              uuid.UUID
		mockFindReturn  *entity.GenreEntity
		mockFindError   error
		mockDeleteError error
		expectedError   bool
		expectedMsg     string
	}{
		{
			name:            "should delete genre successfully",
			id:              testID,
			mockFindReturn:  &entity.GenreEntity{ID: testID, Name: "Action"},
			mockFindError:   nil,
			mockDeleteError: nil,
			expectedError:   false,
		},
		{
			name:           "should return error when genre not found",
			id:             testID,
			mockFindReturn: nil,
			mockFindError:  nil,
			expectedError:  true,
			expectedMsg:    "genre not found",
		},
		{
			name:           "should return error when find fails",
			id:             testID,
			mockFindReturn: nil,
			mockFindError:  errors.New("database error"),
			expectedError:  true,
		},
		{
			name:            "should return error when delete fails",
			id:              testID,
			mockFindReturn:  &entity.GenreEntity{ID: testID, Name: "Action"},
			mockFindError:   nil,
			mockDeleteError: errors.New("delete error"),
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockGenreRepository(ctrl)

			mockRepo.EXPECT().
				FindByID(gomock.Any(), tt.id).
				Return(tt.mockFindReturn, tt.mockFindError).
				Times(1)

			if tt.mockFindReturn != nil && tt.mockFindError == nil {
				mockRepo.EXPECT().
					Delete(gomock.Any(), tt.id).
					Return(tt.mockDeleteError).
					Times(1)
			}

			svc := NewGenreService(mockRepo)
			err := svc.DeleteGenre(context.Background(), tt.id)

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

package genre

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

func TestCreateGenre(t *testing.T) {
	tests := []struct {
		name          string
		request       *request.CreateGenreRequest
		mockReturn    *uuid.UUID
		mockError     error
		expectedError bool
		expectedMsg   string
	}{
		{
			name:    "should create genre successfully",
			request: &request.CreateGenreRequest{Name: "Action"},
			mockReturn: func() *uuid.UUID {
				id := uuid.New()
				return &id
			}(),
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "should return error when name is empty",
			request:       &request.CreateGenreRequest{Name: ""},
			expectedError: true,
			expectedMsg:   "genre name is required",
		},
		{
			name:          "should return error when repository fails",
			request:       &request.CreateGenreRequest{Name: "Action"},
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

			if tt.name != "should return error when name is empty" {
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.AssignableToTypeOf(entity.GenreEntity{})).
					Return(tt.mockReturn, tt.mockError).
					Times(1)
			}

			svc := NewGenreService(mockRepo)
			result, err := svc.CreateGenre(context.Background(), tt.request)

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
			}
		})
	}
}

func TestCreateGenreContextCancellation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockGenreRepository(ctrl)
	svc := NewGenreService(mockRepo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.CreateGenre(ctx, &request.CreateGenreRequest{Name: "Action"})
	if err == nil {
		t.Error("expected context.Canceled error, got nil")
	}
}

package theater

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/repository/mock"
)

func TestListTheaters(t *testing.T) {
	tests := []struct {
		name          string
		mockReturn    []entity.TheaterEntity
		mockError     error
		expectedLen   int
		expectedError bool
	}{
		{
			name: "should list theaters successfully",
			mockReturn: []entity.TheaterEntity{
				{ID: uuid.New(), Name: "Theater 1", TotalSeats: 100},
				{ID: uuid.New(), Name: "Theater 2", TotalSeats: 200},
			},
			expectedLen:   2,
			expectedError: false,
		},
		{
			name:          "should return error when repository fails",
			mockReturn:    nil,
			mockError:     errors.New("db error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTheaterRepo := mock.NewMockTheaterRepository(ctrl)
			mockSeatRepo := mock.NewMockSeatRepository(ctrl)

			mockTheaterRepo.EXPECT().
				FindAll(gomock.Any()).
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			svc := NewTheaterService(mockTheaterRepo, mockSeatRepo)
			result, err := svc.ListTheaters(context.Background())

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) != tt.expectedLen {
					t.Errorf("expected %d theaters, got %d", tt.expectedLen, len(result))
				}
			}
		})
	}
}

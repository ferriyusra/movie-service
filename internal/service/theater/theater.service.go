package theater

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

// TheaterService defines the interface for theater operations
type TheaterService interface {
	CreateTheater(ctx context.Context, req *request.CreateTheaterRequest) (*response.TheaterDetailResponse, error)
	UpdateTheater(ctx context.Context, id uuid.UUID, req *request.UpdateTheaterRequest) (*response.TheaterResponse, error)
	GetTheater(ctx context.Context, id uuid.UUID) (*response.TheaterDetailResponse, error)
	ListTheaters(ctx context.Context) ([]response.TheaterResponse, error)
}

// theaterService is the concrete implementation of TheaterService
type theaterService struct {
	theaterRepository interfaces.TheaterRepository
	seatRepository    interfaces.SeatRepository
}

// NewTheaterService creates a new instance of TheaterService
func NewTheaterService(
	theaterRepository interfaces.TheaterRepository,
	seatRepository interfaces.SeatRepository,
) TheaterService {
	return &theaterService{
		theaterRepository: theaterRepository,
		seatRepository:    seatRepository,
	}
}

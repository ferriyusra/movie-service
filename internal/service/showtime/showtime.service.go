package showtime

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/repository/interfaces"
)

// ShowtimeService defines the interface for showtime operations
type ShowtimeService interface {
	CreateShowtime(ctx context.Context, req *request.CreateShowtimeRequest) (*response.ShowtimeResponse, error)
	UpdateShowtime(ctx context.Context, id uuid.UUID, req *request.UpdateShowtimeRequest) (*response.ShowtimeResponse, error)
	DeleteShowtime(ctx context.Context, id uuid.UUID) error
	GetShowtime(ctx context.Context, id uuid.UUID) (*response.ShowtimeDetailResponse, error)
	ListShowtimesByDate(ctx context.Context, date time.Time) ([]response.ShowtimeResponse, error)
}

// showtimeService is the concrete implementation of ShowtimeService
type showtimeService struct {
	showtimeRepository interfaces.ShowtimeRepository
	movieRepository    interfaces.MovieRepository
	theaterRepository  interfaces.TheaterRepository
}

// NewShowtimeService creates a new instance of ShowtimeService
func NewShowtimeService(
	showtimeRepository interfaces.ShowtimeRepository,
	movieRepository interfaces.MovieRepository,
	theaterRepository interfaces.TheaterRepository,
) ShowtimeService {
	return &showtimeService{
		showtimeRepository: showtimeRepository,
		movieRepository:    movieRepository,
		theaterRepository:  theaterRepository,
	}
}

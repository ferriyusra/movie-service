package response

import (
	"time"

	"github.com/google/uuid"
)

// ShowtimeResponse represents a showtime in API responses
type ShowtimeResponse struct {
	ID             uuid.UUID `json:"id"`
	MovieID        uuid.UUID `json:"movie_id"`
	TheaterID      uuid.UUID `json:"theater_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Price          float64   `json:"price"`
	AvailableSeats int       `json:"available_seats"`
}

// ShowtimeDetailResponse includes movie and theater info
type ShowtimeDetailResponse struct {
	ShowtimeResponse
	Movie   MovieResponse   `json:"movie"`
	Theater TheaterResponse `json:"theater"`
}

// SeatMapResponse represents a seat with availability status
type SeatMapResponse struct {
	SeatResponse
	Status string `json:"status"` // available, reserved
}

package response

import (
	"time"

	"github.com/google/uuid"
)

// ShowtimeResponse represents a showtime in API responses
type ShowtimeResponse struct {
	ID             uuid.UUID `json:"id"`
	MovieID        uuid.UUID `json:"movieId"`
	TheaterID      uuid.UUID `json:"theaterId"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	Price          float64   `json:"price"`
	AvailableSeats int       `json:"availableSeats"`
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

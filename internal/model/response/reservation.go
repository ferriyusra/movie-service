package response

import (
	"time"

	"github.com/google/uuid"
)

// ReservationResponse represents a reservation in API responses
type ReservationResponse struct {
	ID               uuid.UUID      `json:"id"`
	BookingReference string         `json:"booking_reference"`
	ShowtimeID       uuid.UUID      `json:"showtime_id"`
	Status           string         `json:"status"`
	TotalAmount      float64        `json:"total_amount"`
	CancelledAt      *time.Time     `json:"cancelled_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	Seats            []SeatResponse `json:"seats"`
}

// ReservationDetailResponse includes movie, showtime, and theater info
type ReservationDetailResponse struct {
	ReservationResponse
	Movie    MovieResponse    `json:"movie"`
	Showtime ShowtimeResponse `json:"showtime"`
	Theater  TheaterResponse  `json:"theater"`
}

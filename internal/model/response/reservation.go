package response

import (
	"time"

	"github.com/google/uuid"
)

// ReservationResponse represents a reservation in API responses
type ReservationResponse struct {
	ID               uuid.UUID      `json:"id"`
	BookingReference string         `json:"bookingReference"`
	ShowtimeID       uuid.UUID      `json:"showtimeId"`
	Status           string         `json:"status"`
	TotalAmount      float64        `json:"totalAmount"`
	CancelledAt      *time.Time     `json:"cancelledAt,omitempty"`
	CreatedAt        time.Time      `json:"createdAt"`
	Seats            []SeatResponse `json:"seats"`
}

// ReservationDetailResponse includes movie, showtime, and theater info
type ReservationDetailResponse struct {
	ReservationResponse
	Movie    MovieResponse    `json:"movie"`
	Showtime ShowtimeResponse `json:"showtime"`
	Theater  TheaterResponse  `json:"theater"`
}

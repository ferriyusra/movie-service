package request

import "github.com/google/uuid"

// CreateShowtimeRequest is the input for creating a showtime
type CreateShowtimeRequest struct {
	MovieID   uuid.UUID `json:"movieId"`
	TheaterID uuid.UUID `json:"theaterId"`
	StartTime string    `json:"startTime"`
	Price     float64   `json:"price"`
}

// UpdateShowtimeRequest is the input for updating a showtime
type UpdateShowtimeRequest struct {
	MovieID   uuid.UUID `json:"movieId"`
	TheaterID uuid.UUID `json:"theaterId"`
	StartTime string    `json:"startTime"`
	Price     float64   `json:"price"`
}

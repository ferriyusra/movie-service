package request

import "github.com/google/uuid"

// CreateShowtimeRequest is the input for creating a showtime
type CreateShowtimeRequest struct {
	MovieID   uuid.UUID `json:"movie_id"`
	TheaterID uuid.UUID `json:"theater_id"`
	StartTime string    `json:"start_time"`
	Price     float64   `json:"price"`
}

// UpdateShowtimeRequest is the input for updating a showtime
type UpdateShowtimeRequest struct {
	MovieID   uuid.UUID `json:"movie_id"`
	TheaterID uuid.UUID `json:"theater_id"`
	StartTime string    `json:"start_time"`
	Price     float64   `json:"price"`
}

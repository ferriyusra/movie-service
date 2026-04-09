package request

import "github.com/google/uuid"

// CreateReservationRequest is the input for creating a reservation
type CreateReservationRequest struct {
	ShowtimeID uuid.UUID   `json:"showtime_id"`
	SeatIDs    []uuid.UUID `json:"seat_ids"`
}

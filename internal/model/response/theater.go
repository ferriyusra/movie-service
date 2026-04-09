package response

import "github.com/google/uuid"

// TheaterResponse represents a theater in API responses
type TheaterResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	TotalSeats  int       `json:"totalSeats"`
	Rows        int       `json:"rows"`
	SeatsPerRow int       `json:"seatsPerRow"`
}

// SeatResponse represents a seat in API responses
type SeatResponse struct {
	ID     uuid.UUID `json:"id"`
	Row    string    `json:"row"`
	Number int       `json:"number"`
	Label  string    `json:"label"`
	Type   string    `json:"type"`
}

// TheaterDetailResponse includes seat layout
type TheaterDetailResponse struct {
	TheaterResponse
	Seats []SeatResponse `json:"seats"`
}

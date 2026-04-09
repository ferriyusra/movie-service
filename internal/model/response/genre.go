package response

import "github.com/google/uuid"

// GenreResponse represents a genre in API responses
type GenreResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

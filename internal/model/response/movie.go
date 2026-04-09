package response

import (
	"time"

	"github.com/google/uuid"
)

// MovieResponse represents a movie in API responses
type MovieResponse struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	PosterURL   string          `json:"posterUrl"`
	DurationMin int             `json:"durationMin"`
	Language    string          `json:"language"`
	ReleaseDate *time.Time      `json:"releaseDate"`
	Rating      float64         `json:"rating"`
	Genres      []GenreResponse `json:"genres"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

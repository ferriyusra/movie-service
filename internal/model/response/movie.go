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
	PosterURL   string          `json:"poster_url"`
	DurationMin int             `json:"duration_min"`
	Language    string          `json:"language"`
	ReleaseDate *time.Time      `json:"release_date"`
	Rating      float64         `json:"rating"`
	Genres      []GenreResponse `json:"genres"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

package request

import "github.com/google/uuid"

// CreateMovieRequest is the input for creating a movie
type CreateMovieRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PosterURL   string    `json:"poster_url"`
	DurationMin int       `json:"duration_min"`
	Language    string    `json:"language"`
	ReleaseDate string    `json:"release_date"`
	GenreIDs    []uuid.UUID `json:"genre_ids"`
}

// UpdateMovieRequest is the input for updating a movie
type UpdateMovieRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PosterURL   string    `json:"poster_url"`
	DurationMin int       `json:"duration_min"`
	Language    string    `json:"language"`
	ReleaseDate string    `json:"release_date"`
	GenreIDs    []uuid.UUID `json:"genre_ids"`
}

package request

import "github.com/google/uuid"

// CreateMovieRequest is the input for creating a movie
type CreateMovieRequest struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	PosterURL   string      `json:"posterUrl"`
	DurationMin int         `json:"durationMin"`
	Language    string      `json:"language"`
	ReleaseDate string      `json:"releaseDate"`
	GenreIDs    []uuid.UUID `json:"genreIds"`
}

// UpdateMovieRequest is the input for updating a movie
type UpdateMovieRequest struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	PosterURL   string      `json:"posterUrl"`
	DurationMin int         `json:"durationMin"`
	Language    string      `json:"language"`
	ReleaseDate string      `json:"releaseDate"`
	GenreIDs    []uuid.UUID `json:"genreIds"`
}

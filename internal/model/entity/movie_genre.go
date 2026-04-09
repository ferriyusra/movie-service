package entity

import "github.com/google/uuid"

// MovieGenreEntity represents the join table between movies and genres
type MovieGenreEntity struct {
	MovieEntityID uuid.UUID `gorm:"primaryKey;column:movie_entity_id"`
	GenreEntityID uuid.UUID `gorm:"primaryKey;column:genre_entity_id"`
}

func (MovieGenreEntity) TableName() string {
	return "movie_genres"
}

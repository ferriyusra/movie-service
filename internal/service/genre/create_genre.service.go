package genre

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
)

// CreateGenre creates a new genre
func (s *genreService) CreateGenre(ctx context.Context, req *request.CreateGenreRequest) (*response.GenreResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("genre name is required")
	}

	genre := entity.GenreEntity{
		ID:   uuid.New(),
		Name: name,
	}

	id, err := s.genreRepository.Create(ctx, genre)
	if err != nil {
		return nil, fmt.Errorf("creating genre: %w", err)
	}

	return &response.GenreResponse{
		ID:   *id,
		Name: name,
	}, nil
}

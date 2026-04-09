package theater

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMTheaterRepository) Create(ctx context.Context, theater entity.TheaterEntity) (*uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := r.db.WithContext(ctx).Create(&theater).Error; err != nil {
		return nil, err
	}
	return &theater.ID, nil
}

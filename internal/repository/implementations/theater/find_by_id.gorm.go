package theater

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

func (r *GORMTheaterRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.TheaterEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var theater entity.TheaterEntity
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&theater).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &theater, nil
}

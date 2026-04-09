package showtime

import (
	"context"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
)

func (r *GORMShowtimeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return r.db.WithContext(ctx).Delete(&entity.ShowtimeEntity{}, "id = ?", id).Error
}

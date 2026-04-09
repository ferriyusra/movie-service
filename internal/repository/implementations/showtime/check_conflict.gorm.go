package showtime

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/entity"
	"gorm.io/gorm"
)

// CheckConflict checks if a showtime overlaps with an existing one in the same theater.
// Returns the conflicting showtime if found, nil otherwise.
func (r *GORMShowtimeRepository) CheckConflict(ctx context.Context, theaterID uuid.UUID, startTime, endTime time.Time, excludeID *uuid.UUID) (*entity.ShowtimeEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	query := r.db.WithContext(ctx).
		Where("theater_id = ?", theaterID).
		Where("start_time < ? AND end_time > ?", endTime, startTime)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var conflict entity.ShowtimeEntity
	err := query.First(&conflict).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &conflict, nil
}

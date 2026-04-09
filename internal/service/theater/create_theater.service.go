package theater

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

// CreateTheater creates a theater and auto-generates seats
func (s *theaterService) CreateTheater(ctx context.Context, req *request.CreateTheaterRequest) (*response.TheaterDetailResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("theater name is required")
	}
	if req.Rows < 1 || req.Rows > 26 {
		return nil, errors.New("rows must be between 1 and 26")
	}
	if req.SeatsPerRow < 1 || req.SeatsPerRow > 30 {
		return nil, errors.New("seats_per_row must be between 1 and 30")
	}

	totalSeats := req.Rows * req.SeatsPerRow
	theater := entity.TheaterEntity{
		ID:          uuid.New(),
		Name:        name,
		Location:    req.Location,
		TotalSeats:  totalSeats,
		Rows:        req.Rows,
		SeatsPerRow: req.SeatsPerRow,
	}

	id, err := s.theaterRepository.Create(ctx, theater)
	if err != nil {
		return nil, fmt.Errorf("creating theater: %w", err)
	}

	// Auto-generate seats
	seats := make([]entity.SeatEntity, 0, totalSeats)
	for r := 0; r < req.Rows; r++ {
		row := string(rune('A' + r))
		for n := 1; n <= req.SeatsPerRow; n++ {
			seats = append(seats, entity.SeatEntity{
				ID:        uuid.New(),
				TheaterID: *id,
				Row:       row,
				Number:    n,
				Label:     fmt.Sprintf("%s%d", row, n),
				Type:      "standard",
			})
		}
	}

	if err := s.seatRepository.CreateBatch(ctx, seats); err != nil {
		return nil, fmt.Errorf("creating seats: %w", err)
	}

	seatResponses := make([]response.SeatResponse, len(seats))
	for i, seat := range seats {
		seatResponses[i] = response.SeatResponse{
			ID:     seat.ID,
			Row:    seat.Row,
			Number: seat.Number,
			Label:  seat.Label,
			Type:   seat.Type,
		}
	}

	return &response.TheaterDetailResponse{
		TheaterResponse: response.TheaterResponse{
			ID:          *id,
			Name:        theater.Name,
			Location:    theater.Location,
			TotalSeats:  totalSeats,
			Rows:        theater.Rows,
			SeatsPerRow: theater.SeatsPerRow,
		},
		Seats: seatResponses,
	}, nil
}

package counter

import (
	"context"

	"github.com/ferriyusra/clean-arch-go-gin/internal/model/response"
)

// IncrementCounter increments the counter and returns the new value as a response
func (s *counterService) IncrementCounter(ctx context.Context) (*response.GetCounter, error) {
	value, err := s.repo.IncrementCounter(ctx)
	if err != nil {
		return nil, err
	}
	return &response.GetCounter{
		Value: value,
	}, nil
}

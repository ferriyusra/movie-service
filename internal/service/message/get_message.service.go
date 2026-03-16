package message

import (
	"context"

	"github.com/ferriyusra/clean-arch-go-gin/internal/model/response"
)

// GetMessage returns a greeting message using the response model
func (s *messageService) GetMessage(ctx context.Context) (*response.GetMessage, error) {
	message, err := s.repo.GetMessage(ctx, "default")
	if err != nil {
		return nil, err
	}
	if message == nil {
		return nil, nil
	}
	return &response.GetMessage{
		Content: *message,
	}, nil
}

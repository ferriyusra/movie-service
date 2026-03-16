package message

import (
	"context"

	"github.com/ferriyusra/clean-arch-go-gin/internal/model/response"
	"github.com/ferriyusra/clean-arch-go-gin/internal/repository/interfaces"
)

// MessageService defines the interface for message operations
type MessageService interface {
	GetMessage(ctx context.Context) (*response.GetMessage, error)
}

// messageService is the concrete implementation of MessageService
type messageService struct {
	repo interfaces.MessageRepository
}

// NewMessageService creates a new instance of MessageService
func NewMessageService(repo interfaces.MessageRepository) MessageService {
	return &messageService{
		repo: repo,
	}
}

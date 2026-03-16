package response

import "github.com/google/uuid"

// GetUser represents a user in the system
type GetUser struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

type LoginResponse struct {
	User GetUser `json:"user"`
}

type RegisterResponse struct {
	User GetUser `json:"user"`
}

type RefreshResponse struct {
	Message string `json:"message"`
}

type CSRFTokenResponse struct {
	Token string `json:"token"`
}

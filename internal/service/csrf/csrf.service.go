package csrf

// CSRFService handles CSRF token generation and validation
type CSRFService interface {
	GenerateToken() (string, error)
	ValidateToken(token string) bool
}

// csrfService implements CSRFService
type csrfService struct {
	secret []byte
}

// NewCSRFService creates a new CSRF service with HMAC-based token validation
func NewCSRFService(secret string) CSRFService {
	return &csrfService{secret: []byte(secret)}
}

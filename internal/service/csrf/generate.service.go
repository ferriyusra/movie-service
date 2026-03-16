package csrf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateToken generates a HMAC-signed CSRF token (nonce.signature)
func (s *csrfService) GenerateToken() (string, error) {
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("generating CSRF nonce: %w", err)
	}

	mac := hmac.New(sha256.New, s.secret)
	mac.Write(nonce)
	sig := mac.Sum(nil)

	return hex.EncodeToString(nonce) + "." + hex.EncodeToString(sig), nil
}

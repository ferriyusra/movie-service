package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// ValidateToken validates a CSRF token by verifying the HMAC signature
func (s *csrfService) ValidateToken(token string) bool {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return false
	}

	nonce, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}

	sig, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, s.secret)
	mac.Write(nonce)
	expected := mac.Sum(nil)

	return hmac.Equal(sig, expected)
}

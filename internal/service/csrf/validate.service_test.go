package csrf

import (
	"testing"
)

func TestValidateTokenWithGeneratedToken(t *testing.T) {
	service := NewCSRFService("test-secret")

	token, err := service.GenerateToken()
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if !service.ValidateToken(token) {
		t.Errorf("expected generated token to be valid")
	}
}

func TestValidateTokenConsistency(t *testing.T) {
	service := NewCSRFService("test-secret")

	token, _ := service.GenerateToken()

	r1 := service.ValidateToken(token)
	r2 := service.ValidateToken(token)
	r3 := service.ValidateToken(token)

	if !r1 || !r2 || !r3 {
		t.Errorf("validation should be consistently true for a valid token")
	}
}

func TestValidateTokenRejectsEmpty(t *testing.T) {
	service := NewCSRFService("test-secret")

	if service.ValidateToken("") {
		t.Errorf("expected empty token to be invalid")
	}
}

func TestValidateTokenRejectsTampered(t *testing.T) {
	service := NewCSRFService("test-secret")

	token, _ := service.GenerateToken()

	if service.ValidateToken(token + "x") {
		t.Errorf("expected tampered token to be invalid")
	}
}

func TestValidateTokenRejectsWrongSecret(t *testing.T) {
	service1 := NewCSRFService("secret-1")
	service2 := NewCSRFService("secret-2")

	token, _ := service1.GenerateToken()

	if service2.ValidateToken(token) {
		t.Errorf("expected token from different secret to be invalid")
	}
}

func TestValidateTokenRejectsInvalidFormat(t *testing.T) {
	service := NewCSRFService("test-secret")

	invalidTokens := []string{
		"no-dot-separator",
		".",
		".only-signature",
		"only-nonce.",
		"not-hex.not-hex",
		"abc123",
		"token with spaces",
		"!@#$%^&*()",
	}

	for _, token := range invalidTokens {
		if service.ValidateToken(token) {
			t.Errorf("expected token %q to be invalid", token)
		}
	}
}

func TestValidateTokenRejectsWrongSignature(t *testing.T) {
	service := NewCSRFService("test-secret")

	token, _ := service.GenerateToken()

	// Replace signature with a different valid hex string
	parts := splitToken(token)
	if len(parts) == 2 {
		fakeToken := parts[0] + "." + "00000000000000000000000000000000" + "00000000000000000000000000000000"
		if service.ValidateToken(fakeToken) {
			t.Errorf("expected token with wrong signature to be invalid")
		}
	}
}

func splitToken(token string) []string {
	for i, ch := range token {
		if ch == '.' {
			return []string{token[:i], token[i+1:]}
		}
	}
	return []string{token}
}

package csrf

import (
	"strings"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	service := NewCSRFService("test-secret")
	token, err := service.GenerateToken()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Errorf("expected non-empty token, got empty string")
	}

	// Verify token format: nonce.signature
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		t.Fatalf("expected token format 'nonce.signature', got %q", token)
	}

	// nonce = 32 bytes = 64 hex chars
	if len(parts[0]) != 64 {
		t.Errorf("expected nonce length 64, got %d", len(parts[0]))
	}

	// signature = SHA256 = 32 bytes = 64 hex chars
	if len(parts[1]) != 64 {
		t.Errorf("expected signature length 64, got %d", len(parts[1]))
	}
}

func TestGenerateTokenUniqueness(t *testing.T) {
	service := NewCSRFService("test-secret")

	token1, _ := service.GenerateToken()
	token2, _ := service.GenerateToken()

	if token1 == token2 {
		t.Errorf("expected different tokens, but got same token: %s", token1)
	}
}

func TestGenerateTokenRandomness(t *testing.T) {
	service := NewCSRFService("test-secret")
	tokens := make(map[string]bool)

	for i := 0; i < 100; i++ {
		token, err := service.GenerateToken()
		if err != nil {
			t.Fatalf("token generation failed on iteration %d: %v", i, err)
		}
		if tokens[token] {
			t.Errorf("generated duplicate token on iteration %d: %s", i, token)
		}
		tokens[token] = true
	}

	if len(tokens) != 100 {
		t.Errorf("expected 100 unique tokens, got %d", len(tokens))
	}
}

package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateRefreshToken validates and parses a refresh token
func (s *tokenService) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	claims := &TokenClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.RefreshTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

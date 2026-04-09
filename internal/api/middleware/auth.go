package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/ferriyusra/movie-service/internal/service/token"
)

const (
	UserIDCtxKey    = "user_id"
	UserEmailCtxKey = "user_email"
	ClaimsCtxKey    = "claims"
)

// AuthMiddleware validates JWT token from Authorization: Bearer header
func AuthMiddleware(tokenService token.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractBearerToken(c)
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, response.Err("Missing authentication token"))
			c.Abort()
			return
		}

		claims, err := tokenService.ValidateAccessToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Err("Invalid or expired token"))
			c.Abort()
			return
		}

		c.Set(UserIDCtxKey, claims.UserID.String())
		c.Set(UserEmailCtxKey, claims.Email)
		c.Set(ClaimsCtxKey, claims)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT from Authorization header but doesn't require it
func OptionalAuthMiddleware(tokenService token.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractBearerToken(c)
		if tokenStr == "" {
			c.Next()
			return
		}

		claims, err := tokenService.ValidateAccessToken(tokenStr)
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserIDCtxKey, claims.UserID.String())
		c.Set(UserEmailCtxKey, claims.Email)
		c.Set(ClaimsCtxKey, claims)

		c.Next()
	}
}

// extractBearerToken parses the Authorization header and returns the raw token string.
func extractBearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(UserIDCtxKey)
	if !exists {
		return uuid.UUID{}, fmt.Errorf("user not found in context")
	}

	id, err := uuid.Parse(userID.(string))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid user id")
	}

	return id, nil
}

// GetEmailFromContext extracts email from context
func GetEmailFromContext(c *gin.Context) string {
	email, exists := c.Get(UserEmailCtxKey)
	if !exists {
		return ""
	}
	return email.(string)
}

// GetClaimsFromContext extracts token claims from context
func GetClaimsFromContext(c *gin.Context) *token.TokenClaims {
	claims, exists := c.Get(ClaimsCtxKey)
	if !exists {
		return nil
	}
	return claims.(*token.TokenClaims)
}

package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/clean-arch-go-gin/internal/model/response"
	"github.com/ferriyusra/clean-arch-go-gin/internal/service/csrf"
	"github.com/ferriyusra/clean-arch-go-gin/internal/service/token"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
	UserIDCtxKey       = "user_id"
	UserEmailCtxKey    = "user_email"
	ClaimsCtxKey       = "claims"
)

// AuthMiddleware validates JWT token from HTTP-only cookie
func AuthMiddleware(tokenService token.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from HTTP-only cookie
		tokenStr, err := c.Cookie(AccessTokenCookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Err("Missing authentication token"))
			c.Abort()
			return
		}

		// Validate token
		claims, err := tokenService.ValidateAccessToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Err("Invalid or expired token"))
			c.Abort()
			return
		}

		// Store user info in context
		c.Set(UserIDCtxKey, claims.UserID.String())
		c.Set(UserEmailCtxKey, claims.Email)
		c.Set(ClaimsCtxKey, claims)

		c.Next()
	}
}

// OptionalAuthMiddleware validates JWT but doesn't require it
func OptionalAuthMiddleware(tokenService token.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from HTTP-only cookie
		tokenStr, err := c.Cookie(AccessTokenCookie)
		if err != nil {
			c.Next()
			return
		}

		// Validate token
		claims, err := tokenService.ValidateAccessToken(tokenStr)
		if err != nil {
			c.Next()
			return
		}

		// Store user info in context
		c.Set(UserIDCtxKey, claims.UserID.String())
		c.Set(UserEmailCtxKey, claims.Email)
		c.Set(ClaimsCtxKey, claims)

		c.Next()
	}
}

// CSRFMiddleware validates CSRF tokens for state-changing operations
func CSRFMiddleware(csrfService csrf.CSRFService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only validate for state-changing operations
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			csrfToken := c.GetHeader("X-CSRF-Token")
			if csrfToken == "" {
				c.JSON(http.StatusForbidden, response.Err("Missing CSRF token"))
				c.Abort()
				return
			}
			if !csrfService.ValidateToken(csrfToken) {
				c.JSON(http.StatusForbidden, response.Err("Invalid CSRF token"))
				c.Abort()
				return
			}
		}

		c.Next()
	}
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

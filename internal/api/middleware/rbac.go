package middleware

import (
	"net/http"

	"github.com/ferriyusra/movie-service/internal/model/response"
	"github.com/gin-gonic/gin"
)

// AdminOnly middleware restricts access to users with role "admin" in JWT claims.
// Must be used after AuthMiddleware so that claims are available in context.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaimsFromContext(c)
		if claims == nil {
			c.JSON(http.StatusUnauthorized, response.Err("Missing authentication"))
			c.Abort()
			return
		}

		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, response.Err("Admin access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

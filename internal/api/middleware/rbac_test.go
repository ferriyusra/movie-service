package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/service/token"
)

func TestAdminOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupContext   func(c *gin.Context)
		expectedStatus int
	}{
		{
			name: "should allow admin access",
			setupContext: func(c *gin.Context) {
				c.Set(ClaimsCtxKey, &token.TokenClaims{
					UserID: uuid.New(),
					Email:  "admin@cinema.com",
					Name:   "Admin",
					Role:   "admin",
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "should deny regular user access",
			setupContext: func(c *gin.Context) {
				c.Set(ClaimsCtxKey, &token.TokenClaims{
					UserID: uuid.New(),
					Email:  "user@example.com",
					Name:   "User",
					Role:   "user",
				})
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "should deny when no claims in context",
			setupContext:   func(c *gin.Context) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "should deny empty role",
			setupContext: func(c *gin.Context) {
				c.Set(ClaimsCtxKey, &token.TokenClaims{
					UserID: uuid.New(),
					Email:  "user@example.com",
					Name:   "User",
					Role:   "",
				})
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.Use(func(c *gin.Context) {
				tt.setupContext(c)
				c.Next()
			})
			r.Use(AdminOnly())
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, c.Request)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

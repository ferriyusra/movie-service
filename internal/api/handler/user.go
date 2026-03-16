package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/clean-arch-go-gin/internal/api/middleware"
	"github.com/ferriyusra/clean-arch-go-gin/internal/model/request"
	"github.com/ferriyusra/clean-arch-go-gin/internal/model/response"
	"github.com/ferriyusra/clean-arch-go-gin/internal/service/csrf"
	"github.com/ferriyusra/clean-arch-go-gin/internal/service/token"
	userSvc "github.com/ferriyusra/clean-arch-go-gin/internal/service/user"
)

const (
	accessTokenMaxAge  = 15 * 60          // 15 minutes
	refreshTokenMaxAge = 7 * 24 * 60 * 60 // 7 days
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService  userSvc.UserService
	tokenService token.TokenService
	csrfService  csrf.CSRFService
	devMode      bool
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService userSvc.UserService, tokenService token.TokenService, csrfService csrf.CSRFService, devMode bool) *UserHandler {
	return &UserHandler{
		userService:  userService,
		tokenService: tokenService,
		csrfService:  csrfService,
		devMode:      devMode,
	}
}

func (h *UserHandler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	secure := !h.devMode
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.AccessTokenCookie, accessToken, accessTokenMaxAge, "/", "", secure, true)
	c.SetCookie(middleware.RefreshTokenCookie, refreshToken, refreshTokenMaxAge, "/", "", secure, true)
}

func (h *UserHandler) clearAuthCookies(c *gin.Context) {
	secure := !h.devMode
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.AccessTokenCookie, "", -1, "/", "", secure, true)
	c.SetCookie(middleware.RefreshTokenCookie, "", -1, "/", "", secure, true)
}

// Register handles POST /api/auth/register requests
func (h *UserHandler) Register(c *gin.Context) {
	req := &request.RegisterUserRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

	// Validate input
	errs := make(map[string]string)
	if req.Email == "" {
		errs["email"] = "Email is required"
	}
	if req.Password == "" {
		errs["password"] = "Password is required"
	}
	if req.Name == "" {
		errs["name"] = "Name is required"
	}
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, response.ValidationErr("Validation failed", errs))
		return
	}

	// Register user
	resp, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	// Generate tokens
	accessToken, err := h.tokenService.GenerateAccessToken(resp.User.ID, resp.User.Email, resp.User.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to generate access token"))
		return
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(resp.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to generate refresh token"))
		return
	}

	// Store refresh token
	expiresAt := time.Now().Add(time.Duration(refreshTokenMaxAge) * time.Second)
	if err := h.userService.StoreRefreshToken(c.Request.Context(), resp.User.ID, refreshToken, expiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to store refresh token"))
		return
	}

	h.setAuthCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusCreated, response.OK("Registration successful", resp))
}

// Login handles POST /api/auth/login requests
func (h *UserHandler) Login(c *gin.Context) {
	req := &request.LoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

	// Validate input
	errs := make(map[string]string)
	if req.Email == "" {
		errs["email"] = "Email is required"
	}
	if req.Password == "" {
		errs["password"] = "Password is required"
	}
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, response.ValidationErr("Validation failed", errs))
		return
	}

	// Login user
	resp, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Err(err.Error()))
		return
	}

	// Generate tokens
	accessToken, err := h.tokenService.GenerateAccessToken(resp.User.ID, resp.User.Email, resp.User.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to generate access token"))
		return
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(resp.User.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to generate refresh token"))
		return
	}

	// Store refresh token
	expiresAt := time.Now().Add(time.Duration(refreshTokenMaxAge) * time.Second)
	if err := h.userService.StoreRefreshToken(c.Request.Context(), resp.User.ID, refreshToken, expiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to store refresh token"))
		return
	}

	h.setAuthCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusOK, response.OK("Login successful", resp))
}

// Refresh handles POST /api/auth/refresh requests
func (h *UserHandler) Refresh(c *gin.Context) {
	tokenStr, err := c.Cookie(middleware.RefreshTokenCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Err("Missing refresh token"))
		return
	}

	resp, err := h.userService.Refresh(c.Request.Context(), tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Err(err.Error()))
		return
	}

	// resp.Message contains the new access token
	secure := !h.devMode
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.AccessTokenCookie, resp.Message, accessTokenMaxAge, "/", "", secure, true)

	c.JSON(http.StatusOK, response.OK("Token refreshed", nil))
}

// Logout handles POST /api/auth/logout requests
func (h *UserHandler) Logout(c *gin.Context) {
	// Revoke refresh tokens from database
	userID, exists := c.Get(middleware.UserIDCtxKey)
	if exists {
		if id, err := uuid.Parse(userID.(string)); err == nil {
			_ = h.userService.RevokeRefreshTokens(c.Request.Context(), id)
		}
	}

	h.clearAuthCookies(c)

	c.JSON(http.StatusOK, response.OK("Logged out successfully", nil))
}

// GetMe handles GET /api/auth/me requests (protected)
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDCtxKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Err("Unauthorized"))
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("User retrieved", user))
}

// GetCSRFToken handles GET /api/csrf requests
func (h *UserHandler) GetCSRFToken(c *gin.Context) {
	token, err := h.csrfService.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Failed to generate CSRF token"))
		return
	}

	c.JSON(http.StatusOK, response.OK("CSRF token generated", map[string]string{"token": token}))
}

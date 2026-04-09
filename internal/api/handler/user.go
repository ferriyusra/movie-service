package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ferriyusra/movie-service/internal/api/middleware"
	"github.com/ferriyusra/movie-service/internal/model/request"
	"github.com/ferriyusra/movie-service/internal/model/response"
	userSvc "github.com/ferriyusra/movie-service/internal/service/user"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService userSvc.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService userSvc.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles POST /api/auth/register requests
func (h *UserHandler) Register(c *gin.Context) {
	req := &request.RegisterUserRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

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

	resp, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.OK("Registration successful", resp))
}

// Login handles POST /api/auth/login requests
func (h *UserHandler) Login(c *gin.Context) {
	req := &request.LoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Err("Invalid request body"))
		return
	}

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

	resp, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Login successful", resp))
}

// Refresh handles POST /api/auth/refresh requests
func (h *UserHandler) Refresh(c *gin.Context) {
	req := &request.RefreshTokenRequest{}
	if err := c.ShouldBindJSON(req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, response.Err("refresh_token is required"))
		return
	}

	resp, err := h.userService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Token refreshed", resp))
}

// Logout handles POST /api/auth/logout requests
func (h *UserHandler) Logout(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDCtxKey)
	if exists {
		if id, err := uuid.Parse(userID.(string)); err == nil {
			_ = h.userService.RevokeRefreshTokens(c.Request.Context(), id)
		}
	}

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

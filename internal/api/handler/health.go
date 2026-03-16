package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ferriyusra/clean-arch-go-gin/internal/model/response"
	"github.com/ferriyusra/clean-arch-go-gin/internal/service/health"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	service health.HealthService
}

// NewHealthHandler creates a new instance of HealthHandler
func NewHealthHandler(svc health.HealthService) *HealthHandler {
	return &HealthHandler{
		service: svc,
	}
}

// Check handles GET /api/health requests
func (h *HealthHandler) Check(c *gin.Context) {
	status, err := h.service.Check(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK("Service is healthy", status))
}

package handler

import (
	service "honda-leasing-api/internal/services"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service service.DashboardService
}

func NewDashboardHandler(s service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: s}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {

	result, err := h.service.GetSummary(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": result})
}

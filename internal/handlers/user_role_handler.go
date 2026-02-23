package handler

import (
	"strconv"

	service "honda-leasing-api/internal/services"

	"github.com/gin-gonic/gin"
)

type UserRoleHandler struct {
	service service.UserRoleService
}

func NewUserRoleHandler(s service.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{service: s}
}

func (h *UserRoleHandler) AssignRole(c *gin.Context) {

	userIDStr := c.Param("id")
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	var req struct {
		RoleID     int64 `json:"role_id"`
		AssignedBy int64 `json:"assigned_by"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AssignRole(
		c.Request.Context(),
		userID,
		req.RoleID,
		req.AssignedBy,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "role assigned"})
}

func (h *UserRoleHandler) GetUserRoles(c *gin.Context) {

	userIDStr := c.Param("id")
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	roles, err := h.service.GetUserRoles(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": roles})
}

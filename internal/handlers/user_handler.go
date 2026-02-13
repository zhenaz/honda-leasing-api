package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"honda-leasing-api/internal/domain/model"
	service "honda-leasing-api/internal/services"

	"github.com/gin-gonic/gin"
)

var ErrUserNotFound = errors.New("user not found") // sementara, nanti idealnya dari service/repo

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		// Kalau repo kamu masih balikin gorm.ErrRecordNotFound, cek itu juga di sini / di service.
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// contoh: jangan bocorin password/pin kalau json tagnya masih kebuka
	_ = context.Background()
	_ = model.User{} // biar import kepake kalau IDE kamu masih warning

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	type reqBody struct {
		Username    *string `json:"username"`
		PhoneNumber string  `json:"phone_number" binding:"required"`
		Email       *string `json:"email"`
		FullName    string  `json:"full_name" binding:"required"`
		Password    *string `json:"password"`
	}

	var req reqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := &model.User{
		Username:    req.Username,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		FullName:    req.FullName,
		Password:    req.Password,
	}

	if err := h.service.CreateUser(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// response aman (nggak balikin password/pin)
	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"user_id":       u.UserID,
			"username":      u.Username,
			"phone_number":  u.PhoneNumber,
			"email":         u.Email,
			"full_name":     u.FullName,
			"is_active":     u.IsActive,
			"created_at":    u.CreatedAt,
			"updated_at":    u.UpdatedAt,
		},
	})
}
package handler

import (
	"net/http"

	"honda-leasing-api/internal/domain/query"
	"honda-leasing-api/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	q := query.Use(h.db)

	// 1️⃣ Cari user
	user, err := q.User.WithContext(c.Request.Context()).
		Where(q.User.PhoneNumber.Eq(req.PhoneNumber)).
		First()

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	// 2️⃣ Ambil user_roles
	userRoles, err := q.UserRole.WithContext(c.Request.Context()).
		Where(q.UserRole.UserID.Eq(user.UserID)).
		Find()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed get user roles",
		})
		return
	}

	// 3️⃣ Ambil role_name dari roles table
	var roles []string

	for _, ur := range userRoles {

		role, err := q.Role.WithContext(c.Request.Context()).
			Where(q.Role.RoleID.Eq(ur.RoleID)).
			First()

		if err == nil {
			roles = append(roles, role.RoleName)
		}
	}

	// 4️⃣ Generate JWT
	token, err := utils.GenerateToken(user.UserID, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

package middleware

import (
	"net/http"
	"strings"

	"honda-leasing-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// inject user_id ke context
		c.Set("user_id", claims.UserID)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

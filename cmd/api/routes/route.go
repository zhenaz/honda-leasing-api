package routes

import (
	handler "honda-leasing-api/internal/handlers"
	repository "honda-leasing-api/internal/repositories"
	service "honda-leasing-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// init repo
	userRepo := repository.NewUserRepository(db)

	// init service
	userService := service.NewUserService(userRepo)

	// init handler
	userHandler := handler.NewUserHandler(userService)

	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "success", "message": "ok"})
		})

		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetByID)
			users.POST("", userHandler.CreateUser)
		}
	}
}

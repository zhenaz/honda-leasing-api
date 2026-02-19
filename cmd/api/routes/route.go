package routes

import (
	handler "honda-leasing-api/internal/handlers"
	repository "honda-leasing-api/internal/repositories"
	service "honda-leasing-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {

	// =====================
	// DASHBOARD MODULE
	// =====================

	dashboardService := service.NewDashboardService(db)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	// =====================
	// USER MODULE
	// =====================

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// =====================
	// LEASING MODULE
	// =====================

	leasingService := service.NewLeasingService(db)
	leasingHandler := handler.NewLeasingHandler(leasingService)

	// =====================
	// PAYMENT MODULE
	// =====================

	paymentService := service.NewPaymentService(db)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// =====================
	// ROUTES
	// =====================

	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "success", "message": "ok"})
		})

		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/summary", dashboardHandler.GetSummary)
		}

		// USER ROUTES
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetByID)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// LEASING ROUTES
		leasing := api.Group("/leasing")
		{
			leasing.POST("/contracts", leasingHandler.CreateContract)
			leasing.GET("/contracts/:id", leasingHandler.GetContractDetail)

		}

		// PAYMENT ROUTES
		payments := api.Group("/payments")
		{
			payments.POST("", paymentHandler.PayInstallment)
		}
	}
}

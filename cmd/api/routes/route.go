package routes

import (
	handler "honda-leasing-api/internal/handlers"
	"honda-leasing-api/internal/middleware"
	repository "honda-leasing-api/internal/repositories"
	service "honda-leasing-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {

	// =====================
	// INIT SERVICES & HANDLERS
	// =====================

	// Dashboard
	dashboardService := service.NewDashboardService(db)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	// User
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// User Role
	userRoleRepo := repository.NewUserRoleRepository(db)
	userRoleService := service.NewUserRoleService(userRoleRepo)
	userRoleHandler := handler.NewUserRoleHandler(userRoleService)

	// Leasing
	leasingService := service.NewLeasingService(db)
	leasingHandler := handler.NewLeasingHandler(leasingService)

	// Payment
	paymentService := service.NewPaymentService(db)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Auth
	authHandler := handler.NewAuthHandler(db)

	// =====================
	// API GROUP
	// =====================

	api := router.Group("/api/v1")

	// ---------------------
	// PUBLIC ROUTES
	// ---------------------
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "ok",
		})
	})

	api.POST("/login", authHandler.Login)

	// ---------------------
	// PROTECTED ROUTES (JWT)
	// ---------------------
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		// =====================
		// DASHBOARD
		// =====================
		dashboard := protected.Group("/dashboard")
		{
			dashboard.GET("/summary", dashboardHandler.GetSummary)
		}

		// =====================
		// USER
		// =====================
		users := protected.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetByID)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)

			// Role Management
			users.POST("/:id/roles", userRoleHandler.AssignRole)
			users.GET("/:id/roles", userRoleHandler.GetUserRoles)
		}

		// =====================
		// LEASING
		// =====================
		leasing := protected.Group("/leasing")
		{
			leasing.POST("/contracts", leasingHandler.CreateContract)
			leasing.GET("/contracts/:id", leasingHandler.GetContractDetail)
			leasing.PUT("/contracts/:id/approve", leasingHandler.ApproveContract)
			leasing.GET("/contracts", leasingHandler.ListContracts)
			leasing.PUT("/tasks/:id/complete", leasingHandler.CompleteTask)
		}

		// =====================
		// PAYMENT
		// =====================
		payments := protected.Group("/payments")
		{
			payments.POST("", paymentHandler.PayInstallment)
			payments.POST("/check-overdue", paymentHandler.CheckOverdue)
		}
	}
}

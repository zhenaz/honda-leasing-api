package main

import (
	"log"

	config "honda-leasing-api/internal/configs"
	// handler "honda-leasing-api/internal/handlers"
	// repository "honda-leasing-api/internal/repositories"
	// service "honda-leasing-api/internal/services"
	route "honda-leasing-api/cmd/api/routes"
	"honda-leasing-api/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Configuration
	cfg := config.Load() // Pastikan fungsi ini ada di package configs

	// 2. Initialize Database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	// 3. Dependency Injection
	// Kita ambil field .DB karena repository butuh *gorm.DB, bukan *database.Database
	// userRepo := repository.NewUserRepository(db.DB)
	// userService := service.NewUserService(userRepo)
	// userHandler := handler.NewUserHandler(userService)

	// 4. Setup Gin Router
	r := gin.Default()
	route.SetupRoutes(r, db.DB)

	// User Routes
	// v1 := r.Group("/api/v1")
	// {
	// 	v1.GET("/health", func(c *gin.Context) {
	// 		c.JSON(200, gin.H{"status": "success", "message": "Server is up and running"})
	// 	})

	// 	v1.GET("/users", userHandler.GetUsers)
	// 	v1.GET("/users/:id", userHandler.GetByID)
	// }

	
	// Run Server
	// Pastikan cfg.Server.Address berisi format ":8080" atau "localhost:8080"
	log.Printf("🚀 Server running at %s", cfg.Server.Address)

	if err := r.Run(cfg.Server.Address); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

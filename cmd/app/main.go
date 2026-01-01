package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"

	"restaurant/internal/session/handler"
	"restaurant/internal/session/repository"
	"restaurant/internal/session/service"
)

func main() {
	// Database connection
	connStr := "postgres://username:password@localhost/restaurant?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize repositories
	menuRepo := repository.NewMenuRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	sessionRepo := repository.NewPostgresRepository(db)

	// Initialize services with proper dependency injection
	menuService := service.NewMenuService(menuRepo)
	orderService := service.NewOrderService(orderRepo, menuService) // Inject menuService for validation
	sessionService := service.NewService(sessionRepo)

	// Initialize handlers
	menuHandler := handler.NewMenuHandler(menuService)
	orderHandler := handler.NewOrderHandler(orderService)
	sessionHandler := handler.NewHandler(sessionService)

	// Setup Gin router
	router := gin.Default()

	// Register routes from each handler
	menuHandler.RegisterRoutes(router)
	orderHandler.RegisterRoutes(router)
	sessionHandler.RegisterRoutes(router)

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed:", err)
	}
}

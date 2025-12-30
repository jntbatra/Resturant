package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

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

	// Setup routes
	mux := http.NewServeMux()

	// Menu routes
	mux.HandleFunc("/menu", menuHandler.ListMenuItems)
	mux.HandleFunc("/menu/create", menuHandler.CreateMenuItem)
	mux.HandleFunc("/menu/item", menuHandler.GetMenuItem)

	// Order routes
	mux.HandleFunc("/orders", orderHandler.ListOrders)
	mux.HandleFunc("/orders/create", orderHandler.CreateOrder)
	mux.HandleFunc("/orders/item", orderHandler.GetOrder)

	// Session routes
	mux.HandleFunc("/sessions", sessionHandler.ListSessions)
	mux.HandleFunc("/sessions/create", sessionHandler.CreateSession)

	// Start server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed:", err)
	}
}

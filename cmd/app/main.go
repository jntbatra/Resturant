package main

// @title Restaurant Management API
// @version 1.0
// @description Restaurant management system with sessions, menus, and orders
// @termsOfService http://example.com/terms
// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "restaurant/docs"

	"restaurant/internal/middleware"
	"restaurant/internal/pool"
	"restaurant/internal/session/handler"
	"restaurant/internal/session/repository"
	"restaurant/internal/session/service"
	"restaurant/internal/shutdown"
)

func main() {
	// Initialize shutdown manager
	shutdownMgr := shutdown.NewManager(30 * time.Second)

	// Database connection string from environment variables
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "restaurant"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "restaurant_password"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "restaurant"
	}
	dbSSLMode := os.Getenv("DB_SSLMODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	log.Printf("Connecting to database: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Apply connection pool configuration for optimal performance
	poolConfig := pool.DefaultPoolConfig()
	pool.ApplyPoolConfig(db, poolConfig)
	log.Printf("Database pool configured: MaxOpenConns=%d, MaxIdleConns=%d",
		poolConfig.MaxOpenConns, poolConfig.MaxIdleConns)

	// Register database shutdown hook (will execute last in reverse order)
	shutdownMgr.RegisterHook(func(ctx context.Context) error {
		log.Println("Closing database connection...")
		return db.Close()
	})

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

	// Initialize and register middlewares in OPTIMAL order
	// Middleware order matters! Each layer wraps the next, so order affects what gets caught/processed

	// 1. ERROR HANDLER - OUTERMOST layer, catches panics from all downstream middleware
	router.Use(middleware.ErrorHandler())

	// 2. CORS - Handle cross-origin requests (browser requirement)
	router.Use(middleware.CORSMiddleware())

	// 3. REQUEST ID - Generate unique ID early for tracking/logging/debugging
	router.Use(middleware.RequestIDMiddleware())

	// 4. LOGGING - Log all requests (uses RequestID from step 3)
	router.Use(middleware.LoggingMiddleware())

	// 5. REQUEST SIZE LIMIT - Prevent DOS attacks (check size before processing)
	router.Use(middleware.RequestSizeLimitMiddleware(1024 * 1024)) // 1MB limit

	// 6. RATE LIMITING - Final security check before handlers (prevent abuse)
	middleware.InitRateLimiter(middleware.RateLimitConfig{
		RequestsPerSecond: 100,
		BurstSize:         200,
	})
	router.Use(middleware.RateLimitMiddleware())

	// API Documentation endpoint - Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Redirect /docs to Swagger UI
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Register routes from each handler
	menuHandler.RegisterRoutes(router)
	orderHandler.RegisterRoutes(router)
	sessionHandler.RegisterRoutes(router)

	// Create HTTP server with graceful shutdown support
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Register HTTP server shutdown hook (will execute first in reverse order)
	shutdownMgr.RegisterHook(func(ctx context.Context) error {
		log.Println("Shutting down HTTP server...")
		return server.Shutdown(ctx)
	})

	// Start listening for shutdown signals in a goroutine
	go shutdownMgr.Wait()

	// Start server
	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed:", err)
	}

	// Wait for shutdown to complete
	<-shutdownMgr.Done()
	log.Println("Application shutdown complete")
}

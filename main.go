package main

import (
	"log"
	"os"

	"transaction-service/internal/adapters/database"
	"transaction-service/internal/adapters/handlers"
	"transaction-service/internal/application/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found, using default environment variables")
	}

	// Initialize the database
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed successfully")

	// Initialize repositories
	userRepo := database.NewUserRepository(db)
	transactionRepo := database.NewTransactionRepository(db)

	// Initialize services
	transactionService := services.NewTransactionService(userRepo, transactionRepo)

	// Initialize the HTTP handlers
	httpHandler := handlers.NewHandler(transactionService)

	// Set up Gin HTTP router
	router := gin.Default()

	// Add middleware for error handling and logging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Set up routes
	httpHandler.SetupRoutes(router)

	// Get port from environment variables or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

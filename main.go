package main

import (
	"context"
	"expense-tracker-api/repository"
	"expense-tracker-api/services"
	"expense-tracker-api/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	if err := repository.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repository.CloseDB()

	// Migrate existing schema (remove user_id columns if they exist)
	if err := repository.MigrateSchema(); err != nil {
		log.Printf("Warning: Migration failed (may not be needed): %v", err)
	}

	// Create database schema
	if err := repository.CreateSchema(); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// Initialize repositories
	categoryRepo := repository.NewCategoryRepository()
	expenseRepo := repository.NewExpenseRepository()
	budgetRepo := repository.NewBudgetRepository()

	// Initialize services
	categoryService := services.NewCategoryService(categoryRepo)
	expenseService := services.NewExpenseService(expenseRepo, categoryRepo)
	budgetService := services.NewBudgetService(budgetRepo, expenseRepo)

	// Setup router
	router := transport.SetupRouter(categoryService, expenseService, budgetService)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

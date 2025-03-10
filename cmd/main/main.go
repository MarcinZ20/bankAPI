package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MarcinZ20/bankAPI/api/routes"
	"github.com/MarcinZ20/bankAPI/internal/app"
	"github.com/MarcinZ20/bankAPI/internal/database"
	"github.com/MarcinZ20/bankAPI/internal/importer"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create a base context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Println("Received shutdown signal. Initiating graceful shutdown...")
		cancel()
	}()

	// Initialize database
	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Disconnect(ctx)

	// Import data from spreadsheet
	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	if spreadsheetID == "" {
		log.Fatal("SPREADSHEET_ID environment variable is not set")
	}

	log.Println("Starting data import...")
	if err := importer.ImportSpreadsheetData(ctx, spreadsheetID); err != nil {
		log.Fatalf("Failed to import data: %v", err)
	}
	log.Println("Data import completed successfully")

	// Initialize and configure API server
	appConfig := app.Initialize()
	routes.BankRoutes(appConfig.Server)

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Starting server on port %s\n", os.Getenv("API_SERVER_PORT"))
		if err := appConfig.Server.Listen(os.Getenv("API_SERVER_PORT")); err != nil {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for shutdown or error
	select {
	case err := <-serverErrors:
		log.Printf("Server error: %v\n", err)
	case <-ctx.Done():
		log.Println("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := appConfig.Server.ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("Error during server shutdown: %v\n", err)
		}
	}
}

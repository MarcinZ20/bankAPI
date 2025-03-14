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
	"github.com/MarcinZ20/bankAPI/internal/services"
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

	log.Println("Successfully connected to database")

	// Initialize services
	serviceManager := services.NewServiceManager(db)
	if !serviceManager.IsInitialized() {
		log.Fatal("Failed to initialize services")
	}
	log.Println("Services initialized successfully")

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

	serverErrors := make(chan error, 1)
	go func() {
		port := os.Getenv("API_SERVER_PORT")
		if port == "" {
			port = ":8080"
		}
		log.Printf("Starting server on port %s\n", port)
		if err := appConfig.Server.Listen(port); err != nil {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

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

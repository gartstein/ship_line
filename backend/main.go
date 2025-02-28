package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ship_line/db/bolt"
	"ship_line/handlers"
	"ship_line/services"
	"syscall"
	"time"
)

func main() {
	// Open (or create) the Bolt database.
	// TODO: move to config
	storage, err := bolt.NewBoltStorage("pack_sizes.db")
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer storage.Close()

	// Load pack sizes from DB.
	packSizes, err := storage.GetPackSizes()
	if err != nil {
		log.Fatalf("failed to load pack sizes: %v", err)
	}
	// If no pack sizes stored, use default values and store them.
	// TODO: move to config
	if len(packSizes) == 0 {
		packSizes = []int{5000, 2000, 1000, 500, 250}
		if err := storage.SetPackSizes(packSizes); err != nil {
			log.Fatalf("failed to store default pack sizes: %v", err)
		}
	}

	// Create PackService with the DB repository.
	packService := services.NewPackService(storage)

	// Pass the service to the handler.
	router := handlers.SetupRouter(packService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the server in a goroutine.
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server running on port 8080")

	// Listen for the interrupt signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Gracefully shut down, waiting up to 5 seconds for current operations to finish.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

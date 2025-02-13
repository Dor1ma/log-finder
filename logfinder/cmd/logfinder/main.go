package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dor1ma/log-finder/logfinder/internal/config"
	"github.com/Dor1ma/log-finder/logfinder/internal/server/handlers"
	"github.com/Dor1ma/log-finder/logfinder/internal/server/routers"
	"github.com/Dor1ma/log-finder/logfinder/internal/service"
	"github.com/Dor1ma/log-finder/logfinder/internal/storage/repository"
)

func main() {
	cfg := config.Load()

	repo, err := repository.NewLogRepository(
		cfg.LogDir,
		cfg.MaxOpenFiles,
		cfg.FileCacheTTL,
	)

	if err != nil {
		log.Fatalf("Error occured during repo creating: %v", err)
	}

	useCase := service.NewLogService(repo, cfg.CacheTTL)
	handler := handlers.NewLogHandler(useCase)

	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: routers.NewRouter(handler, cfg.RateLimit),
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.ServerPort)

	<-done
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

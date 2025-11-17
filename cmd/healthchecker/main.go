package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"distributed-URL-health-checker/internal/checker"
	"distributed-URL-health-checker/internal/server"
	"distributed-URL-health-checker/pkg/config"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hc := checker.New(cfg.CheckInterval, cfg.Timeout)
	for _, url := range cfg.URLs {
		hc.AddURL(url)
	}

	hc.Start(ctx)
	srv := server.New(hc, cfg.Port)

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down gracefully...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}

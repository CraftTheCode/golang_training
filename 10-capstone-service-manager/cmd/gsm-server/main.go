package main

import (
	"context"
	"errors"
	"fmt"
	"golang_training/10-capstone-service-manager/internal/api"
	"golang_training/10-capstone-service-manager/internal/config"
	"golang_training/10-capstone-service-manager/internal/storage"
	"golang_training/10-capstone-service-manager/internal/worker"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg := config.Load()
	store := storage.NewMemoryStore()

	// Seed some default sample services for demonstration
	_, _ = store.Create("Cloudflare DNS", "https://1.1.1.1")
	_, _ = store.Create("Google Public API", "https://dns.google")

	monitor := worker.NewHealthMonitor(store, cfg, logger)

	// Context for background workers
	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()

	// Start background health checking monitor
	go monitor.Start(workerCtx)

	serverMux := api.NewServer(store, monitor)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           serverMux.Handler(),
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("GSM Server listening", "port", cfg.Port)
		serverErrors <- httpServer.ListenAndServe()
	}()

	// Graceful shutdown listener
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Fatal server error", "err", err)
			os.Exit(1)
		}
	case sig := <-shutdown:
		logger.Info("Shutdown signal received", "signal", sig.String())

		// Stop background workers first
		cancelWorkers()

		// Drain in-flight HTTP connections
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("Forced shutdown", "err", err)
			_ = httpServer.Close()
		}
		logger.Info("GSM Server exited cleanly")
	}
}

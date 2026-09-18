package main

import (
	"context"
	"errors"
	"golang_training/06-networking-and-http/project-rest-api/handler"
	"golang_training/06-networking-and-http/project-rest-api/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	taskHandler := handler.NewTaskHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", taskHandler.Health)
	mux.HandleFunc("GET /api/v1/tasks", taskHandler.ListTasks)
	mux.HandleFunc("POST /api/v1/tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /api/v1/tasks/{id}", taskHandler.GetTask)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", taskHandler.DeleteTask)

	// Build middleware pipeline: Recoverer -> RequestID -> Logger -> Mux
	pipeline := middleware.Recoverer(
		middleware.RequestID(
			middleware.Logger(logger)(mux),
		),
	)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           pipeline,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Server run loop
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("Starting HTTP Server", "addr", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Signal interception for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server fatal error", "err", err)
			os.Exit(1)
		}
	case sig := <-shutdown:
		logger.Info("Shutdown signal received", "signal", sig.String())

		// Allow up to 10 seconds for in-flight requests to complete
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Graceful shutdown failed, forcing close", "err", err)
			_ = server.Close()
			os.Exit(1)
		}
		logger.Info("Server stopped gracefully")
	}
}

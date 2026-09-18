package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

type Middleware func(http.Handler) http.Handler

// Chain combines multiple middlewares from left to right
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// RequestIDMiddleware injects a simulated unique correlation ID into request context
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := fmt.Sprintf("req-%d", time.Now().UnixNano()%10000)
		ctx := context.WithValue(r.Context(), "requestID", reqID)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware logs request metadata and duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := r.Context().Value("requestID")

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		fmt.Printf("[LOG] ID=%v | %s %s | Duration: %v\n",
			reqID, r.Method, r.URL.Path, duration)
	})
}

// RecoveryMiddleware prevents server crashes from handler panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				reqID := r.Context().Value("requestID")
				fmt.Printf("[PANIC RECOVERED] ID=%v | Error: %v\n", reqID, rec)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!\n"))
	})

	mux.HandleFunc("GET /panic-test", func(w http.ResponseWriter, r *http.Request) {
		panic("critical unexpected database driver crash!")
	})

	// Wrap entire router with middlewares
	pipeline := Chain(mux, RecoveryMiddleware, RequestIDMiddleware, LoggingMiddleware)

	ts := httptest.NewServer(pipeline)
	defer ts.Close()

	fmt.Println("=== 1. Normal Request through Middleware Chain ===")
	res, _ := http.Get(ts.URL + "/hello")
	fmt.Println("Status:", res.Status, "| Header X-Request-ID:", res.Header.Get("X-Request-ID"))

	fmt.Println("\n=== 2. Panicking Request Handled by Recovery Middleware ===")
	resPanic, _ := http.Get(ts.URL + "/panic-test")
	fmt.Println("Status:", resPanic.Status)
}

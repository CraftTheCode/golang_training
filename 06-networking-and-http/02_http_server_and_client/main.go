package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

type StatusResponse struct {
	Status    string    `json:"status"`
	Uptime    string    `json:"uptime"`
	Timestamp time.Time `json:"timestamp"`
}

func statusHandler(startTime time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := StatusResponse{
			Status:    "HEALTHY",
			Uptime:    time.Since(startTime).Round(time.Second).String(),
			Timestamp: time.Now().UTC(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func main() {
	startTime := time.Now()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", statusHandler(startTime))

	// In-memory test server for reliable demonstration
	ts := httptest.NewServer(mux)
	defer ts.Close()

	fmt.Println("=== 1. Production HTTP Server Configuration ===")
	// Note: When deploying standalone, configure timeouts explicitly:
	_ = &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,  // Protects against Slowloris attacks
		ReadTimeout:       10 * time.Second, // Max duration reading entire request
		WriteTimeout:      15 * time.Second, // Max duration writing response
		IdleTimeout:       60 * time.Second, // Max keep-alive connection idle duration
	}
	fmt.Println("Configured http.Server timeouts: ReadHeaderTimeout, ReadTimeout, WriteTimeout, IdleTimeout.")

	fmt.Println("\n=== 2. Making Requests with Robust http.Client ===")
	// NEVER use http.DefaultClient in production (it has no timeout!)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, ts.URL+"/health", nil)
	req.Header.Set("User-Agent", "Go-Learning-Client/1.0")

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request failed: %v\n", err)
		return
	}
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)
	fmt.Printf("Response Status: %s\n", res.Status)
	fmt.Printf("Response Body: %s\n", string(bodyBytes))
}

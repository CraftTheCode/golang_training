package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang_training/10-capstone-service-manager/internal/model"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var defaultBaseURL = "http://localhost:8080"

func printUsage() {
	fmt.Fprintf(os.Stderr, `Go Service Manager (GSM) CLI Client

Usage:
  gsm-cli <command> [arguments]

Commands:
  list                          List all registered services and statuses
  register <name> <url>         Register a new endpoint to monitor
  check <id>                    Trigger an immediate manual health probe
  delete <id>                   Delete a registered service
  health                        Check if the GSM server daemon is responsive
`)
}

func getBaseURL() string {
	if val := os.Getenv("GSM_SERVER_URL"); val != "" {
		return strings.TrimRight(val, "/")
	}
	return defaultBaseURL
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	cmd := strings.ToLower(os.Args[1])
	baseURL := getBaseURL()

	switch cmd {
	case "health":
		resp, err := client.Get(baseURL + "/health")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to GSM server: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("GSM Server Status: %s\n%s\n", resp.Status, string(body))

	case "list":
		resp, err := client.Get(baseURL + "/api/v1/services")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching services: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		var services []model.MonitoredService
		if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode response: %v\n", err)
			os.Exit(1)
		}

		if len(services) == 0 {
			fmt.Println("No services registered. Use 'gsm-cli register <name> <url>' to add one.")
			return
		}

		fmt.Println("-----------------------------------------------------------------------------------------")
		fmt.Printf("%-4s | %-18s | %-10s | %-10s | %s\n", "ID", "NAME", "STATUS", "LATENCY", "URL")
		fmt.Println("-----------------------------------------------------------------------------------------")
		for _, s := range services {
			latencyStr := fmt.Sprintf("%dms", s.LatencyMs)
			if s.LatencyMs == 0 {
				latencyStr = "N/A"
			}
			fmt.Printf("%-4d | %-18s | %-10s | %-10s | %s\n",
				s.ID, s.Name, s.Status, latencyStr, s.URL)
		}
		fmt.Println("-----------------------------------------------------------------------------------------")

	case "register":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: gsm-cli register <name> <url>")
			os.Exit(1)
		}
		name := os.Args[2]
		url := os.Args[3]

		payload, _ := json.Marshal(model.CreateServiceRequest{Name: name, URL: url})
		resp, err := client.Post(baseURL+"/api/v1/services", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to register service: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			fmt.Fprintf(os.Stderr, "Registration failed [%s]: %s\n", resp.Status, string(body))
			os.Exit(1)
		}

		var created model.MonitoredService
		_ = json.NewDecoder(resp.Body).Decode(&created)
		fmt.Printf("✓ Registered service #%d: %s (%s)\n", created.ID, created.Name, created.URL)

	case "check":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: gsm-cli check <id>")
			os.Exit(1)
		}
		id := os.Args[2]
		resp, err := client.Post(fmt.Sprintf("%s/api/v1/services/%s/check", baseURL, id), "application/json", nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Probe request failed: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		var result model.HealthCheckResult
		_ = json.NewDecoder(resp.Body).Decode(&result)
		fmt.Printf("✓ Probe complete for Service #%s -> Status: %s (Latency: %dms)\n",
			id, result.Status, result.LatencyMs)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: gsm-cli delete <id>")
			os.Exit(1)
		}
		id := os.Args[2]
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/services/%s", baseURL, id), nil)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Delete request failed: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			fmt.Printf("✓ Service #%s deleted successfully.\n", id)
		} else {
			fmt.Printf("Failed to delete service #%s (Status: %s)\n", id, resp.Status)
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

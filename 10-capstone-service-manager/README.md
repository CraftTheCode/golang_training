# Module 10: Capstone Project — Go Service Manager (`GSM`)

A production-grade, distributed service monitoring and management application that brings together all the skills mastered across the curriculum:

```text
                                  +-----------------------+
                                  |     gsm-cli (CLI)     |
                                  +-----------------------+
                                              |
                                              v (HTTP REST API)
+-----------------------------------------------------------------------------------+
|                              gsm-server (Daemon)                                  |
|                                                                                   |
|  +---------------------------+             +-----------------------------------+  |
|  |     HTTP REST Server      |             |     Background Health Monitor     |  |
|  |  - Service CRUD           |             |  - Bounded Goroutine Worker Pool  |  |
|  |  - Live Status Endpoints  |             |  - Periodic HTTP Probing          |  |
|  |  - Graceful Shutdown      |             |  - Latency & State Tracking       |  |
|  +---------------------------+             +-----------------------------------+  |
|               \                                   /                               |
|                v                                 v                                |
|             +---------------------------------------+                             |
|             |          Storage Engine Layer         |                             |
|             |  - Thread-safe in-memory & SQL ready  |                             |
|             +---------------------------------------+                             |
+-----------------------------------------------------------------------------------+
```

---

## Key Features

1. **Dual Architecture (CLI + HTTP Server)**:
   - `gsm-server`: Runs HTTP REST API and manages background health check workers.
   - `gsm-cli`: Lightweight CLI tool for operators to register, list, and probe services.
2. **Background Concurrency (Worker Pools & Timers)**:
   - Evaluates health of all registered services periodically using bounded worker goroutines.
   - Measures response latencies and detects service outages automatically.
3. **Graceful Lifecycle Management**:
   - Clean shutdown intercepting `SIGINT`/`SIGTERM` with context timeouts.
4. **Comprehensive Automated Tests**:
   - Unit tests, handler tests, and race-safe concurrency verification.

---

## Quickstart & Usage

### 1. Build Both Binaries
```bash
go build -o ./bin/gsm-server ./10-capstone-service-manager/cmd/gsm-server/main.go
go build -o ./bin/gsm-cli ./10-capstone-service-manager/cmd/gsm-cli/main.go
```

### 2. Start the Server
```bash
./bin/gsm-server
# Or run directly:
go run ./10-capstone-service-manager/cmd/gsm-server/main.go
```

### 3. Manage Services via CLI
```bash
# Register services to monitor
go run ./10-capstone-service-manager/cmd/gsm-cli/main.go register "Auth Service" "https://httpbin.org/status/200"
go run ./10-capstone-service-manager/cmd/gsm-cli/main.go register "Payment Gateway" "https://httpbin.org/delay/1"

# List all services and their live health status
go run ./10-capstone-service-manager/cmd/gsm-cli/main.go list

# Trigger an immediate manual health check
go run ./10-capstone-service-manager/cmd/gsm-cli/main.go check 1

# Delete a registered service
go run ./10-capstone-service-manager/cmd/gsm-cli/main.go delete 1
```

### 4. Running the Tests
```bash
go test -v -race ./10-capstone-service-manager/...
```

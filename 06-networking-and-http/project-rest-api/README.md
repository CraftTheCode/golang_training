# Project: Production Go REST API

A clean REST API built entirely using Go's modern standard library (`net/http` with Go 1.22+ routing enhancements).

## Features
- **Standard Library Routing**: Uses `r.PathValue("id")` and HTTP method prefixes (`GET /api/v1/tasks/{id}`).
- **Graceful Shutdown**: Intercepts `os.Interrupt` and `SIGTERM`, allowing 10 seconds for active requests to finish cleanly.
- **Middleware Pipeline**: Request ID, Structured Logging, Panic Recovery.
- **Structured Error Responses**: Clean JSON errors with proper HTTP status codes.
- **Automated Tests**: Tested using `net/http/httptest`.

## Running the Server
```bash
go run ./06-networking-and-http/project-rest-api/main.go
```

## Running the Automated Tests
```bash
go test -v ./06-networking-and-http/project-rest-api/...
```

## Sample Curl Commands
```bash
# Check health
curl -i http://localhost:8080/health

# Create a task
curl -i -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Design microservice boundary"}'

# List tasks
curl -i http://localhost:8080/api/v1/tasks

# Get task by ID
curl -i http://localhost:8080/api/v1/tasks/1

# Delete task by ID
curl -i -X DELETE http://localhost:8080/api/v1/tasks/1
```

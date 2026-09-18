# Module 06: Networking & HTTP Servers

## 1. Conceptual Foundation

Go's standard library `net/http` package is one of the most battle-tested HTTP server implementations in modern computing, serving billions of requests per second at Google, Cloudflare, and Netflix without needing external web frameworks.

Key architectural concepts:
- **Each incoming HTTP request is handled in its own dedicated goroutine**.
- **Context is passed into every `http.Request`** (`r.Context()`), ensuring that if a client disconnects or times out, downstream database queries and computations can abort immediately.
- **Go 1.22+ Enhanced Routing**: Native method matching (`GET /items`) and path wildcards (`/items/{id}`) directly in `http.ServeMux` without third-party routers like gorilla/mux or chi.

---

## 2. Architecture: HTTP Request Lifecycle & Middleware Chaining

```text
Client Request
      |
      v
+-------------------------------------------------------------+
|                        http.Server                          |
|  - ReadHeaderTimeout, IdleTimeout, WriteTimeout             |
+-------------------------------------------------------------+
      |
      v
+-------------------------------------------------------------+
|                    Middleware Pipeline                      |
|  - RecoveryMiddleware (recovers panics, returns 500)        |
|  - RequestIDMiddleware (generates unique trace ID)          |
|  - LoggingMiddleware (logs method, path, status, latency)   |
+-------------------------------------------------------------+
      |
      v
+-------------------------------------------------------------+
|                      http.ServeMux                          |
|  - Matches: "GET /api/v1/tasks" or "POST /api/v1/tasks"     |
+-------------------------------------------------------------+
      |
      v
+-------------------------------------------------------------+
|                      Resource Handler                       |
|  - JSON decode, domain validation, service response         |
+-------------------------------------------------------------+
```

---

## 3. Design Discussion: Graceful Shutdown

Servers must not be abruptly terminated with `SIGKILL` or `os.Exit(1)`. When updating or deploying:
1. Intercept `SIGINT` (Ctrl+C) and `SIGTERM`.
2. Stop listening for new incoming connections.
3. Call `server.Shutdown(ctx)` with a grace period timeout (e.g. 10 seconds).
4. Allow in-flight requests and background transactions to complete cleanly.
5. Exit safely with code `0`.

---

## 4. Module Directory Structure

```text
06-networking-and-http/
├── 01_tcp_echo_server/main.go          # Low-level net.Listen and net.Conn
├── 02_http_server_and_client/main.go   # http.Server, custom timeouts, http.Client
├── 03_middleware_patterns/main.go      # Middleware chaining mechanics
└── project-rest-api/                   # Production REST API with Go 1.22 routing
    ├── handler/task_handler.go         # CRUD endpoints with JSON validation
    ├── middleware/middleware.go        # Logger, Recovery, RequestID
    ├── main.go                         # Server setup with Graceful Shutdown
    ├── api_test.go                     # httptest end-to-end handler tests
    └── README.md                       # API documentation & curl examples
```

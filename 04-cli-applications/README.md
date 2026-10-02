# Module 04: CLI Application Development & Architecture

## 1. Conceptual Foundation

Building CLI tools in Go is one of the language's greatest strengths (Docker, Kubernetes/kubectl, Terraform, and Hugo are all written in Go). Go compiles to a single, static machine-code binary with zero external runtime dependencies.

A professional CLI must exhibit:
1. **Predictable Exit Codes**: `0` for success, non-zero (`1`, `2`, etc.) for errors.
2. **Stream Separation**: Normal output to `os.Stdout`; diagnostics, progress, and errors to `os.Stderr`.
3. **Structured Subcommands**: Subcommand routing (e.g. `git commit`, `kubectl get`).
4. **Layered Architecture**: Decoupling the CLI transport from business logic and persistence.

---

## 2. Architecture: Layered CLI Pattern

```text
+-------------------------------------------------------------+
|                      CLI Entrypoint                         |
|  - Parses flags, arguments, subcommands (cmd/gomanager)     |
+-------------------------------------------------------------+
                              |
                              v
+-------------------------------------------------------------+
|                     Application Service                     |
|  - Business rules, validation, orchestrates operations      |
+-------------------------------------------------------------+
                              |
                              v
+-------------------------------------------------------------+
|                    Repository Interface                     |
|  - Abstract contract for persistence (TaskRepository)       |
+-------------------------------------------------------------+
                              |
                              v
+-------------------------------------------------------------+
|                   Storage Implementation                    |
|  - JSON file persistence (tasks.json)                       |
+-------------------------------------------------------------+
```

---

## 3. Standard Library Tools & Core Concepts

### A. The `flag` Package: Pointer Semantics
```go
// Returns a POINTER (*string, *int, *bool) because the flags are defined before
// os.Args are parsed. flag.Parse() writes into the memory addresses of these pointers.
port := flag.Int("port", 8080, "HTTP server listening port")
debug := flag.Bool("debug", false, "Enable verbose debug logging")

flag.Parse() // Must be called after definitions, before reading *port or *debug

fmt.Printf("Starting on port %d (debug=%t)\n", *port, *debug)
```
- `flag.Args()` returns remaining positional arguments (e.g. `mycli -v arg1 arg2` -> `["arg1", "arg2"]`).
- `flag.NArg()` returns the count of non-flag arguments.

### B. Stream Separation: `os.Stdout` vs `os.Stderr`
Unix philosophy dictates that **only primary program output goes to stdout**.
- Diagnostics, logs, errors, and `--help` output must go to `os.Stderr`:
```go
fmt.Fprintf(os.Stderr, "Error: missing required argument\n")
os.Exit(1) // Non-zero indicates failure to calling shell/CI
```
This ensures downstream commands can safely pipe output: `mycli -format=json | jq .` without log lines corrupting the JSON stream.

### C. Detecting Piped Input (Terminal vs Pipe)
```go
stat, _ := os.Stdin.Stat()
if (stat.Mode() & os.ModeCharDevice) == 0 {
    // Input is being piped via stdin: `cat data.txt | mycli`
    data, _ := io.ReadAll(os.Stdin)
} else {
    // Input is interactive: user launched `mycli` directly in terminal
}
```

### D. JSON: Streaming (`Encoder`/`Decoder`) vs In-Memory (`Marshal`/`Unmarshal`)

| Operation | Best Used For | Memory Behavior |
| :--- | :--- | :--- |
| `json.NewEncoder(w).Encode(v)` | Files, HTTP responses, network sockets | Streams directly to `io.Writer` without allocating intermediate `[]byte` |
| `json.NewDecoder(r).Decode(&v)` | Reading from files, HTTP request bodies | Reads chunks directly from `io.Reader` |
| `json.Marshal(v)` | Small objects, caching in Redis/memory | Allocates entire encoded JSON into a single `[]byte` slice |
| `json.Unmarshal(b, &v)` | Parsing an existing `[]byte` in memory | Requires complete payload pre-buffered in memory |

---


## 4. Module Directory Structure

```text
04-cli-applications/
├── 01_flags_and_io/main.go             # Flags, standard I/O, exit codes
├── 02_json_and_files/main.go           # File operations and JSON streaming
└── project-gomanager/                  # Complete Layered Task Management CLI
    ├── cmd/gomanager/main.go           # Command routing (init, config, task)
    ├── internal/
    │   ├── model/task.go               # Domain entity
    │   ├── repository/json_storage.go  # File-based JSON persistence
    │   └── service/task_service.go     # Core business logic & validation
    └── README.md                       # CLI usage and commands
```

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

## 3. Standard Library Tools
- `flag`: Standard flag parsing package, supports `flag.NewFlagSet` for subcommands.
- `os`: Args, environment, exit codes (`os.Exit`), file descriptors.
- `bufio`: Buffered reading (`bufio.NewScanner`) and writing.
- `encoding/json`: Serialization and streaming (`json.Encoder` / `json.Decoder`).
- `log/slog`: Structured logging built into Go 1.21+.

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

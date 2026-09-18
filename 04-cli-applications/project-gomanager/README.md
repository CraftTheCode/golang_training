# Project: `gomanager` CLI Task Manager

A command-line application demonstrating layered Clean Architecture in Go:

```text
CLI (cmd/gomanager)
       |
       v
Application Service (internal/service)
       |
       v
Repository Interface (internal/repository)
       |
       v
JSON File Storage (tasks.json / config.json)
```

## Available Commands

```bash
# Initialize storage directory and default config
go run ./04-cli-applications/project-gomanager/cmd/gomanager init

# View current configuration
go run ./04-cli-applications/project-gomanager/cmd/gomanager config show

# Create new tasks
go run ./04-cli-applications/project-gomanager/cmd/gomanager task create "Implement database migrations"
go run ./04-cli-applications/project-gomanager/cmd/gomanager task create "Configure Docker build"

# List all tasks
go run ./04-cli-applications/project-gomanager/cmd/gomanager task list

# Delete a task by ID
go run ./04-cli-applications/project-gomanager/cmd/gomanager task delete 1
```

## Running the Unit Tests
```bash
go test -v ./04-cli-applications/project-gomanager/internal/service/...
```

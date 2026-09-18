# Module 09: Docker, Containerization & Linux Deployment

## 1. Conceptual Foundation

One of Go's greatest superpowers is its ability to produce a **single, fully static machine binary** with zero dynamic shared library dependencies (`libc`, `glibc`).
By setting `CGO_ENABLED=0`, the compiler bundles everything needed into the binary, allowing it to run inside an empty **`scratch`** Docker image (0 MB base size) or a minimal **`alpine`** image (~5 MB).

---

## 2. Cross-Compilation Superpowers

You can compile a Linux binary directly from Windows PowerShell in a single command:

```powershell
# Compile a static 64-bit Linux binary from Windows
$env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="amd64"
go build -ldflags="-s -w" -o ./bin/server-linux ./06-networking-and-http/project-rest-api/main.go
```

**What does `-ldflags="-s -w"` do?**
- `-s`: Strips the symbol table (disables symbol names).
- `-w`: Strips DWARF debug information.
- **Result**: Reduces binary size by 40% to 60% without affecting runtime performance.

---

## 3. Production Docker Architecture: Multi-Stage Builds

```text
STAGE 1: Builder (golang:1.22-alpine)
  - Copies go.mod, downloads cached dependencies
  - Compiles binary with CGO_ENABLED=0 and -ldflags="-s -w"
  - Generates nonroot user and TLS certificates

STAGE 2: Final Runtime Image (scratch or alpine:3.20)
  - Copies ONLY the compiled binary and ca-certificates from Stage 1
  - Drops root privileges (runs as nonroot:10001)
  - Final image size: ~15 MB!
```

---

## 4. Module Directory Structure

```text
09-docker-and-deployment/
├── README.md               # Guide to cross-compilation and containerization
├── Dockerfile              # Multi-stage production build recipe
├── docker-compose.yml      # Orchestrates REST API + PostgreSQL + Health Checks
└── linux_systemd_guide.md  # Production systemd service unit guide
```

# ==============================================================================
# STAGE 1: Build Environment
# ==============================================================================
FROM golang:1.22-alpine AS builder

# Install SSL certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create an unprivileged non-root user for maximum container security
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 10001 \
    appuser

WORKDIR /build

# Leverage Docker layer caching: download dependencies first
COPY go.mod ./
RUN go mod download

# Copy application source code
COPY . .

# Compile fully static binary:
# CGO_ENABLED=0 disables C bindings (pure Go static binary)
# -ldflags="-s -w" strips symbol tables and debug metadata
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /build/api-server \
    ./06-networking-and-http/project-rest-api/main.go

# ==============================================================================
# STAGE 2: Minimal Production Runtime
# ==============================================================================
FROM alpine:3.20

# Import certificates, timezone info, and non-root user from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy compiled static binary
COPY --from=builder /build/api-server /app/api-server

# Drop root privileges! Run as non-root user
USER 10001:10001

EXPOSE 8080

# Health check probe
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/api-server"]

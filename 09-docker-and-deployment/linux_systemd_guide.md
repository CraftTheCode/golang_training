# Linux Systemd Deployment Guide for Go Applications

While containers are popular, running a statically compiled Go binary natively under Linux with **systemd** provides minimal latency, zero virtualization overhead, and automated process supervision.

---

## 1. Deploying the Binary to Linux

1. Cross-compile your binary for Linux:
   ```bash
   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./api-server ./06-networking-and-http/project-rest-api/main.go
   ```

2. Copy binary to the target Linux host:
   ```bash
   scp api-server user@your-server-ip:/usr/local/bin/api-server
   ```

3. Create dedicated system service user:
   ```bash
   sudo useradd -rs /bin/false apprunner
   sudo chmod +x /usr/local/bin/api-server
   ```

---

## 2. Creating the Systemd Service Unit

Create `/etc/systemd/system/goservice.service`:

```ini
[Unit]
Description=Go Production REST API Service
After=network.target postgresql.service

[Service]
Type=simple
User=apprunner
Group=apprunner
WorkingDirectory=/usr/local/bin
ExecStart=/usr/local/bin/api-server
Restart=always
RestartSec=5s

# Security Hardening
NoNewPrivileges=true
ProtectSystem=full
ProtectHome=true

# Environment Variables
Environment=PORT=8080
Environment=ENVIRONMENT=production
Environment=DATABASE_URL=postgres://appuser:secret@localhost:5432/gotraining?sslmode=require

# Standard output and error to systemd journal
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

---

## 3. Managing the Service

```bash
# Reload systemd daemon to pick up new service
sudo systemctl daemon-reload

# Enable service to start on system boot
sudo systemctl enable goservice.service

# Start service now
sudo systemctl start goservice.service

# Check live service status
sudo systemctl status goservice.service

# Stream live application logs
sudo journalctl -u goservice.service -f
```

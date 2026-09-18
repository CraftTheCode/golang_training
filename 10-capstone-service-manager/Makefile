.PHONY: all build test test-race run-server clean

all: test build

build:
	@echo "Building gsm-server and gsm-cli..."
	go build -o ./bin/gsm-server ./10-capstone-service-manager/cmd/gsm-server/main.go
	go build -o ./bin/gsm-cli ./10-capstone-service-manager/cmd/gsm-cli/main.go

test:
	@echo "Running all tests..."
	go test -v ./10-capstone-service-manager/...

test-race:
	@echo "Running tests with race detector..."
	go test -v -race ./10-capstone-service-manager/...

run-server:
	go run ./10-capstone-service-manager/cmd/gsm-server/main.go

clean:
	@rm -rf ./bin

# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Seed the database with initial Super Admin
seed:
	@echo "Seeding database..."
	@go run cmd/seed/main.go

# Seed the database with sample users
seed-users:
	@echo "Adding sample users to database..."
	@go run cmd/seed-users/main.go

# Seed the database with 30 new users
seed-30-users:
	@echo "Adding 30 new users to database..."
	@go run cmd/seed-30-users/main.go

# Seed the database with institutions
seed-institutions:
	@echo "Seeding institutions..."
	@go run cmd/seed-institutions/main.go

# Migrate existing users to use institution IDs
migrate-institutions:
	@echo "Migrating users to institution IDs..."
	@go run cmd/migrate-institutions/main.go

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run docker-down itest seed seed-users seed-30-users seed-institutions migrate-institutions

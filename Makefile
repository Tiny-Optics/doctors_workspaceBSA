# Makefile for BLOODSA Doctor's Workspace

.PHONY: help dev prod build up down logs clean restart

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)BLOODSA Doctor's Workspace - Docker Commands$(NC)"
	@echo ""
	@echo "$(GREEN)Usage:$(NC)"
	@echo "  make [target]"
	@echo ""
	@echo "$(GREEN)Development Targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

# Development Commands
dev: ## Start all services in development mode
	@echo "$(BLUE)Starting development environment...$(NC)"
	docker compose up

dev-build: ## Build and start development services
	@echo "$(BLUE)Building and starting development environment...$(NC)"
	docker compose up --build

dev-detach: ## Start development services in background
	@echo "$(BLUE)Starting development environment in background...$(NC)"
	docker compose up -d

# Production Commands
prod: ## Start all services in production mode
	@echo "$(BLUE)Starting production environment...$(NC)"
	docker compose -f docker-compose.prod.yml up -d

prod-build: ## Build and start production services
	@echo "$(BLUE)Building and starting production environment...$(NC)"
	docker compose -f docker-compose.prod.yml up -d --build

prod-down: ## Stop production services
	@echo "$(BLUE)Stopping production environment...$(NC)"
	docker compose -f docker-compose.prod.yml down

# General Commands
build: ## Build all development images
	@echo "$(BLUE)Building development images...$(NC)"
	docker compose build

up: dev-detach ## Alias for dev-detach

down: ## Stop all development services
	@echo "$(BLUE)Stopping development environment...$(NC)"
	docker compose down

restart: ## Restart all development services
	@echo "$(BLUE)Restarting development environment...$(NC)"
	docker compose restart

# Logs
logs: ## View logs from all services
	docker compose logs -f

logs-backend: ## View backend logs
	docker compose logs -f backend

logs-frontend: ## View frontend logs
	docker compose logs -f frontend

logs-db: ## View database logs
	docker compose logs -f mongodb

# Individual Service Control
backend-restart: ## Restart backend service
	docker compose restart backend

frontend-restart: ## Restart frontend service
	docker compose restart frontend

db-restart: ## Restart database service
	docker compose restart mongodb

# Shell Access
shell-backend: ## Access backend container shell
	docker compose exec backend sh

shell-frontend: ## Access frontend container shell
	docker compose exec frontend sh

shell-db: ## Access MongoDB shell
	docker compose exec mongodb mongosh -u $${BLUEPRINT_DB_USERNAME:-melkey} -p $${BLUEPRINT_DB_ROOT_PASSWORD:-password1234}

# Testing
test-backend: ## Run backend tests in container
	docker compose exec backend go test ./...

test-frontend: ## Run frontend tests in container
	docker compose exec frontend npm test

# Cleanup
clean: ## Remove containers, networks, and volumes
	@echo "$(YELLOW)Warning: This will remove all containers, networks, and volumes!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(BLUE)Cleaning up...$(NC)"; \
		docker compose down -v; \
	fi

clean-force: ## Force remove everything without confirmation
	@echo "$(BLUE)Force cleaning...$(NC)"
	docker compose down -v

prune: ## Remove all unused Docker resources
	@echo "$(YELLOW)Warning: This will remove all unused Docker resources!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(BLUE)Pruning Docker system...$(NC)"; \
		docker system prune -a; \
	fi

# Status
ps: ## Show running containers
	docker compose ps

status: ps ## Alias for ps

# Environment Setup
env: ## Create .env file from .env.example
	@if [ ! -f .env ]; then \
		echo "$(BLUE)Creating .env file from .env.example...$(NC)"; \
		cp .env.example .env; \
		echo "$(GREEN)Created .env file. Please edit it with your credentials.$(NC)"; \
	else \
		echo "$(YELLOW).env file already exists!$(NC)"; \
	fi

# Health Checks
health: ## Check health of all services
	@echo "$(BLUE)Checking service health...$(NC)"
	@echo ""
	@echo "$(GREEN)Backend:$(NC)"
	@curl -s http://localhost:8080/health 2>/dev/null || echo "Backend not responding"
	@echo ""
	@echo "$(GREEN)Frontend:$(NC)"
	@curl -s http://localhost:5173 -o /dev/null && echo "Frontend is running" || echo "Frontend not responding"
	@echo ""
	@echo "$(GREEN)MongoDB:$(NC)"
	@docker compose exec -T mongodb mongosh --eval "db.adminCommand('ping')" --quiet 2>/dev/null || echo "MongoDB not responding"



# Makefile for BLOODSA Doctor's Workspace

.PHONY: help dev prod build start stop restart logs clean backup

# Default target
help:
	@echo "BLOODSA Doctor's Workspace - Make Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Start development environment"
	@echo "  make dev-logs     - View development logs"
	@echo "  make dev-stop     - Stop development environment"
	@echo ""
	@echo "Production:"
	@echo "  make prod         - Deploy production environment"
	@echo "  make prod-build   - Build production images"
	@echo "  make prod-start   - Start production containers"
	@echo "  make prod-stop    - Stop production containers"
	@echo "  make prod-restart - Restart production containers"
	@echo "  make prod-logs    - View production logs"
	@echo ""
	@echo "Maintenance:"
	@echo "  make backup       - Backup database and uploads"
	@echo "  make clean        - Remove all containers and volumes (WARNING: deletes data!)"
	@echo "  make update       - Pull latest code and rebuild"
	@echo ""

# Development commands
dev:
	docker compose up -d
	@echo "Development environment started"
	@echo "Frontend: http://localhost:5173"
	@echo "Backend:  http://localhost:8080"

dev-logs:
	docker compose logs -f

dev-stop:
	docker compose down

# Production commands
prod:
	@echo "Deploying production environment..."
	@chmod +x deploy.sh
	@./deploy.sh

prod-build:
	docker compose -f docker-compose.prod.yml build --no-cache

prod-start:
	docker compose -f docker-compose.prod.yml up -d

prod-stop:
	docker compose -f docker-compose.prod.yml down

prod-restart:
	docker compose -f docker-compose.prod.yml restart

prod-logs:
	docker compose -f docker-compose.prod.yml logs -f

prod-ps:
	docker compose -f docker-compose.prod.yml ps

# Update application
update:
	@echo "Pulling latest changes..."
	git pull origin main
	@echo "Rebuilding containers..."
	docker compose -f docker-compose.prod.yml build
	@echo "Restarting services..."
	docker compose -f docker-compose.prod.yml up -d
	@echo "Update complete!"

# Backup
backup:
	@echo "Creating backup..."
	@mkdir -p backups
	@docker exec bloodsa_mongodb_prod mongodump --out=/dump
	@docker cp bloodsa_mongodb_prod:/dump ./backups/mongo_$(shell date +%Y%m%d_%H%M%S)
	@docker run --rm -v sop_uploads:/data -v $(PWD)/backups:/backup alpine tar czf /backup/uploads_$(shell date +%Y%m%d_%H%M%S).tar.gz /data
	@echo "Backup created in ./backups/"

# Clean (WARNING: Deletes all data!)
clean:
	@echo "WARNING: This will delete all containers and data!"
	@read -p "Are you sure? (yes/no): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		docker compose down -v; \
		docker compose -f docker-compose.prod.yml down -v; \
		docker system prune -af; \
		echo "Cleanup complete"; \
	else \
		echo "Cancelled"; \
	fi

# Database shell
db-shell:
	docker exec -it bloodsa_mongodb_prod mongosh -u ${BLUEPRINT_DB_USERNAME} -p ${BLUEPRINT_DB_ROOT_PASSWORD}

# Backend shell
backend-shell:
	docker exec -it bloodsa_backend_prod sh

# Frontend shell
frontend-shell:
	docker exec -it bloodsa_frontend_prod sh

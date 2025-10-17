# Makefile for BLOODSA Doctor's Workspace

.PHONY: help dev cloudways build start stop restart logs clean backup

# Default target
help:
	@echo "BLOODSA Doctor's Workspace - Make Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Start development environment"
	@echo "  make dev-logs     - View development logs"
	@echo "  make dev-stop     - Stop development environment"
	@echo ""
	@echo "Production (Cloudways Shared Hosting):"
	@echo "  make cloudways        - Deploy on Cloudways shared server"
	@echo "  make cloudways-logs   - View deployment logs"
	@echo "  make cloudways-stop   - Stop deployment"
	@echo "  make cloudways-restart- Restart containers"
	@echo "  make cloudways-ps     - Check container status"
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

# Cloudways commands (Production)
cloudways:
	@echo "Deploying to Cloudways shared server..."
	@chmod +x deploy-cloudways.sh
	@./deploy-cloudways.sh

cloudways-logs:
	docker compose -f docker-compose.cloudways.yml --env-file .env.cloudways logs -f

cloudways-stop:
	docker compose -f docker-compose.cloudways.yml down

cloudways-restart:
	docker compose -f docker-compose.cloudways.yml --env-file .env.cloudways restart

cloudways-ps:
	docker compose -f docker-compose.cloudways.yml ps

# Update application
update:
	@echo "Pulling latest changes..."
	git pull origin main
	@echo "Rebuilding containers..."
	docker compose -f docker-compose.cloudways.yml --env-file .env.cloudways build
	@echo "Restarting services..."
	docker compose -f docker-compose.cloudways.yml --env-file .env.cloudways up -d
	@echo "Update complete!"

# Backup
backup:
	@echo "Creating backup..."
	@mkdir -p backups
	@docker exec bloodsa_doctors_mongodb mongodump --out=/dump
	@docker cp bloodsa_doctors_mongodb:/dump ./backups/mongo_$(shell date +%Y%m%d_%H%M%S)
	@docker run --rm -v bloodsa_doctors_sop_uploads:/data -v $(PWD)/backups:/backup alpine tar czf /backup/uploads_$(shell date +%Y%m%d_%H%M%S).tar.gz /data
	@echo "Backup created in ./backups/"

# Clean (WARNING: Deletes all data!)
clean:
	@echo "WARNING: This will delete all containers and data!"
	@read -p "Are you sure? (yes/no): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		docker compose down -v; \
		docker compose -f docker-compose.cloudways.yml down -v; \
		docker system prune -af; \
		echo "Cleanup complete"; \
	else \
		echo "Cancelled"; \
	fi

# Database shell
db-shell:
	docker exec -it bloodsa_doctors_mongodb mongosh -u ${BLUEPRINT_DB_USERNAME} -p ${BLUEPRINT_DB_ROOT_PASSWORD}

# Backend shell
backend-shell:
	docker exec -it bloodsa_doctors_backend sh

# Frontend shell
frontend-shell:
	docker exec -it bloodsa_doctors_frontend sh

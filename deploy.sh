#!/bin/bash

# Production Deployment Script for BLOODSA Doctor's Workspace
# This script builds and deploys the application in production mode

set -e  # Exit on any error

echo "========================================="
echo "BLOODSA Doctor's Workspace - Production Deployment"
echo "========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env.production exists
if [ ! -f .env.production ]; then
    echo -e "${RED}ERROR: .env.production file not found!${NC}"
    echo "Please create it from .env.production.example and configure your values."
    exit 1
fi

echo -e "${GREEN}✓${NC} Environment file found"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}ERROR: Docker is not running${NC}"
    exit 1
fi

echo -e "${GREEN}✓${NC} Docker is running"

# Load production environment
export $(cat .env.production | grep -v '^#' | xargs)

echo ""
echo "Building production images..."
echo "========================================="

# Build images
docker compose -f docker-compose.prod.yml build --no-cache

echo ""
echo -e "${GREEN}✓${NC} Images built successfully"
echo ""
echo "Starting services..."
echo "========================================="

# Stop existing containers
docker compose -f docker-compose.prod.yml down

# Start services
docker compose -f docker-compose.prod.yml up -d

echo ""
echo "Waiting for services to be healthy..."
sleep 10

# Check health
BACKEND_HEALTHY=$(docker inspect --format='{{.State.Health.Status}}' bloodsa_backend_prod 2>/dev/null || echo "unknown")
MONGODB_HEALTHY=$(docker inspect --format='{{.State.Health.Status}}' bloodsa_mongodb_prod 2>/dev/null || echo "unknown")

echo ""
echo "Service Status:"
echo "========================================="
echo "MongoDB: $MONGODB_HEALTHY"
echo "Backend: $BACKEND_HEALTHY"
echo "Frontend: $(docker inspect --format='{{.State.Health.Status}}' bloodsa_frontend_prod 2>/dev/null || echo 'unknown')"
echo "Nginx: $(docker inspect --format='{{.State.Health.Status}}' bloodsa_nginx_prod 2>/dev/null || echo 'unknown')"
echo ""

# Show running containers
docker compose -f docker-compose.prod.yml ps

echo ""
echo "========================================="
echo -e "${GREEN}Deployment Complete!${NC}"
echo "========================================="
echo ""
echo "Access your application at:"
echo "  HTTP:  http://$(hostname -I | awk '{print $1}')"
echo "  Local: http://localhost"
echo ""
echo "Next steps:"
echo "  1. Configure SSL certificate (see docs/PRODUCTION_DEPLOYMENT.md)"
echo "  2. Test the application"
echo "  3. Set up automated backups"
echo "  4. Monitor logs: docker compose -f docker-compose.prod.yml logs -f"
echo ""


#!/bin/bash

# Cloudways Deployment Script for BLOODSA Doctor's Workspace
# Optimized for shared hosting environment with 40+ other apps

set -e

echo "========================================="
echo "BLOODSA Doctor's Workspace - Cloudways Deployment"
echo "========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check environment file
if [ ! -f .env.cloudways ]; then
    echo -e "${RED}ERROR: .env.cloudways file not found!${NC}"
    echo "Please create it from .env.cloudways.example"
    exit 1
fi

echo -e "${GREEN}✓${NC} Environment file found"

# Load environment
export $(cat .env.cloudways | grep -v '^#' | xargs)

# Check for port conflicts
echo ""
echo "Checking for port conflicts..."
if sudo netstat -tulpn | grep -q ":${BLOODSA_BACKEND_PORT:-8880}"; then
    echo -e "${YELLOW}WARNING: Port ${BLOODSA_BACKEND_PORT:-8880} is already in use!${NC}"
    echo "Please change BLOODSA_BACKEND_PORT in .env.cloudways"
    exit 1
fi

if sudo netstat -tulpn | grep -q ":${BLOODSA_FRONTEND_PORT:-8881}"; then
    echo -e "${YELLOW}WARNING: Port ${BLOODSA_FRONTEND_PORT:-8881} is already in use!${NC}"
    echo "Please change BLOODSA_FRONTEND_PORT in .env.cloudways"
    exit 1
fi

echo -e "${GREEN}✓${NC} Ports ${BLOODSA_BACKEND_PORT:-8880} and ${BLOODSA_FRONTEND_PORT:-8881} are available"

# Check Docker
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}ERROR: Docker is not running${NC}"
    exit 1
fi

echo -e "${GREEN}✓${NC} Docker is running"

# Check available resources
AVAILABLE_MEM=$(free -m | awk 'NR==2{print $7}')
if [ "$AVAILABLE_MEM" -lt 1200 ]; then
    echo -e "${YELLOW}WARNING: Low available memory (${AVAILABLE_MEM}MB). Recommended: 1200MB+${NC}"
    echo "Your app needs ~900MB. Continue anyway? (yes/no)"
    read -r CONTINUE
    if [ "$CONTINUE" != "yes" ]; then
        echo "Deployment cancelled"
        exit 1
    fi
fi

echo ""
echo "Building images (this may take a few minutes)..."
echo "========================================="

# Build with no cache for clean production build
docker compose -f docker-compose.cloudways.yml --env-file .env.cloudways build --no-cache

echo ""
echo -e "${GREEN}✓${NC} Images built successfully"
echo ""
echo "Starting services..."
echo "========================================="

# Stop existing containers if any
docker compose -f docker-compose.cloudways.yml down 2>/dev/null || true

# Start services
docker compose -f docker-compose.cloudways.yml --env-file .env.cloudways up -d

echo ""
echo "Waiting for services to be healthy (may take 30-60 seconds)..."
sleep 40

# Check health
echo ""
echo "Service Status:"
echo "========================================="
docker compose -f docker-compose.cloudways.yml ps

echo ""
echo "Testing backend health..."
HEALTH_CHECK=$(curl -s http://localhost:${BLOODSA_BACKEND_PORT:-8880}/health || echo "failed")

if echo "$HEALTH_CHECK" | grep -q "healthy"; then
    echo -e "${GREEN}✓${NC} Backend is healthy"
else
    echo -e "${YELLOW}⚠${NC} Backend health check: $HEALTH_CHECK"
fi

# Show resource usage
echo ""
echo "Current Resource Usage:"
echo "========================================="
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}" | grep bloodsa

echo ""
echo "========================================="
echo -e "${GREEN}Deployment Complete!${NC}"
echo "========================================="
echo ""
echo "Your application is running on:"
echo "  Backend:  http://localhost:${BLOODSA_BACKEND_PORT:-8880}"
echo "  Frontend: http://localhost:${BLOODSA_FRONTEND_PORT:-8881}"
echo ""
echo "⚠️  NEXT STEPS FOR CLOUDWAYS:"
echo "  1. Configure reverse proxy in Cloudways panel"
echo "  2. Point your domain to these ports:"
echo "     - API: localhost:${BLOODSA_BACKEND_PORT:-8880}"
echo "     - Web: localhost:${BLOODSA_FRONTEND_PORT:-8881}"
echo "  3. Enable SSL in Cloudways (Let's Encrypt)"
echo "  4. Test your domain URL"
echo ""
echo "📚 Full guide: docs/CLOUDWAYS_DEPLOYMENT.md"
echo ""
echo "Monitor logs: docker compose -f docker-compose.cloudways.yml logs -f"
echo ""


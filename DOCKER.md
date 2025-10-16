# Docker Setup Guide

This document provides instructions for running the BLOODSA Doctor's Workspace using Docker.

## 📋 Prerequisites

- Docker Engine 20.10+ or Docker Desktop
- Docker Compose V2 (comes with Docker Desktop)
- At least 2GB of free RAM
- Ports 5173, 8080, and 27017 available

## 🚀 Quick Start

### Development Mode

1. **Copy environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your credentials if needed
   ```

2. **Start all services:**
   ```bash
   docker compose up
   ```

3. **Access the application:**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - MongoDB: localhost:27017

### Production Mode

1. **Build and start production containers:**
   ```bash
   docker compose -f docker-compose.prod.yml up -d --build
   ```

2. **Access the application:**
   - Frontend: http://localhost
   - Backend API: http://localhost:8080

## 🛠️ Available Commands

### Development

```bash
# Start all services
docker compose up

# Start services in background
docker compose up -d

# Stop all services
docker compose down

# View logs
docker compose logs -f

# View logs for specific service
docker compose logs -f backend
docker compose logs -f frontend

# Rebuild containers
docker compose up --build

# Restart a specific service
docker compose restart backend
```

### Production

```bash
# Start production services
docker compose -f docker-compose.prod.yml up -d

# Stop production services
docker compose -f docker-compose.prod.yml down

# View production logs
docker compose -f docker-compose.prod.yml logs -f
```

### Cleanup

```bash
# Stop and remove containers, networks
docker compose down

# Remove containers, networks, and volumes
docker compose down -v

# Remove all unused Docker resources
docker system prune -a
```

## 🏗️ Architecture

The Docker setup consists of three main services:

### 1. MongoDB (`mongodb`)
- **Image:** mongo:latest
- **Port:** 27017
- **Volume:** Persistent data storage
- **Purpose:** Primary database for the application

### 2. Backend (`backend`)
- **Base Image:** golang:1.24.4-alpine
- **Port:** 8080
- **Hot-Reload:** Air (development mode)
- **Purpose:** Go API server with Gin framework

### 3. Frontend (`frontend`)
- **Base Image:** node:22.12.0-alpine
- **Port:** 5173 (dev), 80 (prod)
- **Hot-Reload:** Vite HMR (development mode)
- **Purpose:** Vue.js 3 + TypeScript SPA

## 🔄 Hot-Reloading

### Backend (Air)
The backend uses [Air](https://github.com/air-verse/air) for automatic reloading when Go files change.

- Configuration: `backend/.air.toml`
- Any changes to `.go` files will trigger automatic rebuild
- Build logs available in `backend/tmp/build-errors.log`

### Frontend (Vite)
The frontend uses Vite's built-in HMR (Hot Module Replacement).

- Any changes to `.vue`, `.ts`, or `.js` files trigger instant updates
- No page refresh needed for most changes

## 📁 Volume Mounts

### Development Mode
- **Backend:** `./backend` → `/app` (with tmp/ excluded)
- **Frontend:** `./frontend` → `/app` (with node_modules excluded)
- **MongoDB:** Named volume `mongo_data`

This allows you to edit code on your host machine and see changes reflected immediately in the containers.

## 🔐 Environment Variables

Create a `.env` file in the root directory based on `.env.example`:

```env
# Backend Configuration
PORT=8080
APP_ENV=local

# Database Configuration
BLUEPRINT_DB_HOST=localhost  # Use 'mongodb' for docker-compose
BLUEPRINT_DB_PORT=27017
BLUEPRINT_DB_USERNAME=melkey
BLUEPRINT_DB_ROOT_PASSWORD=password1234

# Frontend Configuration
VITE_API_URL=http://localhost:8080
```

**Note:** When running in Docker Compose, the backend uses `BLUEPRINT_DB_HOST=mongodb` (the service name) instead of `localhost`.

## 🐛 Troubleshooting

### Port Already in Use
```bash
# Find what's using the port
lsof -i :8080
# or
netstat -ano | findstr :8080  # Windows

# Stop the process or change the port in .env
PORT=8081
```

### MongoDB Connection Issues
```bash
# Check MongoDB container status
docker compose ps mongodb

# View MongoDB logs
docker compose logs mongodb

# Restart MongoDB
docker compose restart mongodb
```

### Hot-Reload Not Working
```bash
# For backend - rebuild the container
docker compose up --build backend

# For frontend - ensure volumes are mounted correctly
docker compose down
docker compose up frontend
```

### Permission Issues (Linux/Mac)
```bash
# If you encounter permission errors with volumes
sudo chown -R $USER:$USER ./backend/tmp
sudo chown -R $USER:$USER ./frontend/node_modules
```

### Container Won't Start
```bash
# View container logs
docker compose logs [service-name]

# Check container status
docker compose ps

# Recreate containers
docker compose down
docker compose up --force-recreate
```

## 📊 Health Checks

### MongoDB Health Check
The MongoDB service includes a health check that verifies the database is accepting connections before the backend starts.

### Manual Health Checks
```bash
# Check backend health
curl http://localhost:8080/health

# Check MongoDB connection
docker compose exec mongodb mongosh -u melkey -p password1234
```

## 🔧 Advanced Usage

### Execute Commands in Containers

```bash
# Access backend shell
docker compose exec backend sh

# Access frontend shell
docker compose exec frontend sh

# Access MongoDB shell
docker compose exec mongodb mongosh -u melkey -p password1234

# Run Go tests in backend
docker compose exec backend go test ./...

# Run frontend tests
docker compose exec frontend npm test
```

### Build Individual Services

```bash
# Build only backend
docker compose build backend

# Build only frontend
docker compose build frontend
```

### Use Different Compose Files

```bash
# Development (default)
docker compose up

# Production
docker compose -f docker-compose.prod.yml up

# Custom
docker compose -f docker-compose.custom.yml up
```

## 📦 Production Deployment

The production Docker setup uses multi-stage builds for optimized images:

### Backend (Production)
- **Stage 1:** Build the Go binary with optimizations
- **Stage 2:** Minimal Alpine image with only the binary
- **Result:** ~20MB image size

### Frontend (Production)
- **Stage 1:** Build Vue.js app with Vite
- **Stage 2:** Serve static files with Nginx
- **Result:** ~25MB image size

### Deploy to Production

```bash
# Build production images
docker compose -f docker-compose.prod.yml build

# Start production services
docker compose -f docker-compose.prod.yml up -d

# Check status
docker compose -f docker-compose.prod.yml ps

# View logs
docker compose -f docker-compose.prod.yml logs -f
```

## 🔒 Security Best Practices

1. **Never commit `.env` files** - Use `.env.example` as template
2. **Use strong database passwords** in production
3. **Change default credentials** before deploying
4. **Use Docker secrets** for sensitive data in production
5. **Keep base images updated** regularly
6. **Scan images for vulnerabilities** using `docker scan`

## 📝 Notes

- The development setup prioritizes fast iteration and debugging
- The production setup prioritizes security and performance
- Go module cache is persisted in a named volume for faster builds
- Frontend node_modules are persisted in an anonymous volume
- MongoDB data is persisted across container restarts

## 🆘 Getting Help

If you encounter issues:

1. Check the logs: `docker compose logs -f [service]`
2. Verify environment variables in `.env`
3. Ensure all required ports are available
4. Try rebuilding: `docker compose up --build`
5. Try cleaning up: `docker compose down -v && docker compose up`

For more information, refer to the main [README.md](./README.md).



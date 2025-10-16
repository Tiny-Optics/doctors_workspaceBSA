# Quick Start Guide

Get the BLOODSA Doctor's Workspace up and running in under 5 minutes.

## 🚀 Fast Track (Docker - Recommended)

### 1. Prerequisites
- Docker Desktop installed ([Download here](https://www.docker.com/products/docker-desktop))
- Git

### 2. Clone & Setup
```bash
# Clone the repository
git clone <repository-url>
cd doctors_workspaceBSA

# Create environment file
cp .env.example .env
```

### 3. Start Everything
```bash
# Start all services (frontend, backend, MongoDB)
make dev

# OR if you don't have make installed
docker compose up
```

### 4. Access the Application
- **Frontend:** http://localhost:5173
- **Backend API:** http://localhost:8080
- **MongoDB:** localhost:27017

That's it! You're running! 🎉

### Hot-Reload is Active
- Edit any `.go` file in `backend/` → Backend auto-reloads
- Edit any `.vue` or `.ts` file in `frontend/` → Frontend updates instantly

### Useful Commands
```bash
make logs          # View all logs
make down          # Stop all services
make restart       # Restart all services
make clean         # Remove everything (fresh start)
```

---

## 🔧 Manual Setup (Without Docker)

If you prefer to run services individually:

### Prerequisites
- Go 1.24.4+
- Node.js 20.19.0+ or 22.12.0+
- MongoDB (running locally or remote)

### 1. Backend Setup
```bash
cd backend

# Install dependencies
go mod download

# Setup environment
cp .env.example .env
# Edit .env with your MongoDB credentials

# Start MongoDB (if using Docker)
docker-compose up -d

# Run backend (with hot-reload)
make watch

# OR run without hot-reload
make run
```

### 2. Frontend Setup
In a new terminal:
```bash
cd frontend

# Install dependencies
npm install

# Setup environment
cp .env.example .env

# Run development server
npm run dev
```

### 3. Access
- Frontend: http://localhost:5173
- Backend: http://localhost:8080

---

## 🐛 Troubleshooting

### Port Already in Use
```bash
# Check what's using the port
lsof -i :5173  # Frontend
lsof -i :8080  # Backend
lsof -i :27017 # MongoDB

# Stop Docker services
make down
```

### Docker Not Working
```bash
# Check Docker is running
docker ps

# Rebuild everything
docker compose down -v
docker compose up --build
```

### Module Not Found Errors

**Backend:**
```bash
cd backend
go mod tidy
go mod download
```

**Frontend:**
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Database Connection Issues
Check your `.env` file:
```env
# For Docker Compose
BLUEPRINT_DB_HOST=mongodb

# For local development (without Docker)
BLUEPRINT_DB_HOST=localhost
```

---

## 📚 Next Steps

- Read [README.md](./README.md) for full project documentation
- Read [DOCKER.md](./DOCKER.md) for detailed Docker usage
- Check out the [project structure](#) to understand the codebase

---

## 🆘 Need Help?

1. Check if all services are running: `docker compose ps`
2. View logs: `make logs` or `docker compose logs -f`
3. Try a clean restart: `make clean && make dev`
4. Check the main README for detailed troubleshooting

Happy coding! 🚀



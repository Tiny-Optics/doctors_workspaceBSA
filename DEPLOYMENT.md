# Deployment Guide

## 🚀 Automated Deployment with GitHub Actions

### Initial VPS Setup (One-time)

1. **Ensure the VPS is configured:**
   - Docker and Docker Compose installed ✅
   - Git repository cloned to `/opt/doctors_workspace` ✅
   - SSH deploy key added to GitHub ✅
   - Nginx reverse proxy configured ✅

2. **Create `.env` file on VPS:**
   ```bash
   ssh vmszjayo@169.255.58.102
   cd /opt/doctors_workspace
   cp .env.prod.example .env
   nano .env  # Fill in the actual values
   ```

   The `.env` file should contain:
   ```env
   MONGO_ROOT_USER=bloodsa_admin
   MONGO_ROOT_PASSWORD=<generate with: openssl rand -hex 32>
   ENCRYPTION_KEY=TgTskus4WT4YUN9y49++f5ivs9eZrD0jfWz86ILjHTk=
   JWT_SECRET=<generate with: openssl rand -hex 64>
   VITE_API_URL=http://169.255.58.102/api
   ```

3. **Add GitHub Secrets:**
   Go to your repository → Settings → Secrets and variables → Actions
   
   Add these secrets:
   - `SSH_HOST`: `169.255.58.102`
   - `SSH_USER`: `vmszjayo`
   - `SSH_PRIVATE_KEY`: (the private key from VPS at `~/.ssh/github_deploy`)

### How to Deploy

**Simply push to the `main` branch:**
```bash
git add .
git commit -m "Your changes"
git push origin main
```

GitHub Actions will automatically:
1. Connect to the VPS via SSH
2. Pull the latest code
3. Build Docker images using `docker-compose.prod.yml`
4. Restart containers
5. Report deployment status

### Manual Deployment (if needed)

SSH into the VPS and run:
```bash
cd /opt/doctors_workspace
git pull origin main
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml build --no-cache
docker compose -f docker-compose.prod.yml up -d
```

### Viewing Logs

```bash
# All services
docker compose -f docker-compose.prod.yml logs -f

# Specific service
docker compose -f docker-compose.prod.yml logs -f backend
docker compose -f docker-compose.prod.yml logs -f frontend
docker compose -f docker-compose.prod.yml logs -f mongodb
```

### Checking Status

```bash
docker compose -f docker-compose.prod.yml ps
```

### Accessing the Application

- **Frontend**: http://169.255.58.102
- **Backend API**: http://169.255.58.102/api

### Important Notes

1. **Volumes are persistent** - MongoDB data and uploads survive deployments
2. **Never commit `.env` files** - they contain secrets
3. **The `ENCRYPTION_KEY` must never change** - it will break encrypted data
4. **Database credentials should use hex format** (no special URL characters)

### Troubleshooting

1. **Check container logs:**
   ```bash
   docker compose -f docker-compose.prod.yml logs backend
   ```

2. **Restart a specific service:**
   ```bash
   docker compose -f docker-compose.prod.yml restart backend
   ```

3. **Rebuild from scratch:**
   ```bash
   docker compose -f docker-compose.prod.yml down
   docker compose -f docker-compose.prod.yml build --no-cache
   docker compose -f docker-compose.prod.yml up -d
   ```

4. **Check Nginx logs:**
   ```bash
   sudo tail -f /var/log/nginx/error.log
   ```


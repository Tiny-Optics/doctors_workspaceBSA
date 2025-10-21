# Deployment Guide

## 🚀 Automated Deployment with GitHub Actions

### How It Works

1. **GitHub Actions builds** Docker images on powerful GitHub servers
2. **Pushes images** to GitHub Container Registry (GHCR)
3. **VPS pulls** pre-built images and runs them (fast, low CPU)

**No more building on the VPS!** This prevents CPU overload and crashes.

---

## Initial VPS Setup (One-time)

### 1. **Ensure the VPS is configured:**
   - Docker and Docker Compose installed ✅
   - Git repository cloned to `/opt/doctors_workspace` ✅
   - SSH deploy key added to GitHub ✅
   - Nginx reverse proxy configured ✅

### 2. **Create `.env` file on VPS:**
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

### 3. **Make Docker images public:**
   Go to: `https://github.com/orgs/Tiny-Optics/packages`
   
   For each package (`doctors-backend`, `doctors-frontend`):
   - Click the package name
   - Go to **Package settings**
   - Scroll to **Danger Zone**
   - Click **Change visibility** → **Public**

### 4. **Add GitHub Secrets:**
   Go to: `https://github.com/Tiny-Optics/doctors_workspaceBSA/settings/secrets/actions`
   
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
1. **Build images** on GitHub servers (backend + frontend)
2. **Push to GHCR** (GitHub Container Registry)
3. **Connect to VPS** via SSH
4. **Pull latest code** and pre-built images
5. **Restart containers** with new images
6. **Report status**

**Total deployment time: ~2-3 minutes** (was crashing before!)

### Manual Deployment (if needed)

SSH into the VPS and run:
```bash
cd /opt/doctors_workspace
git pull origin main
docker compose -f docker-compose.ghcr.yml pull
docker compose -f docker-compose.ghcr.yml up -d --force-recreate
```

### Viewing Logs

```bash
# All services
docker compose -f docker-compose.ghcr.yml logs -f

# Specific service
docker compose -f docker-compose.ghcr.yml logs -f backend
docker compose -f docker-compose.ghcr.yml logs -f frontend
docker compose -f docker-compose.ghcr.yml logs -f mongodb
```

### Checking Status

```bash
docker compose -f docker-compose.ghcr.yml ps
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
   docker compose -f docker-compose.ghcr.yml logs backend
   ```

2. **Restart a specific service:**
   ```bash
   docker compose -f docker-compose.ghcr.yml restart backend
   ```

3. **Pull latest images and recreate:**
   ```bash
   docker compose -f docker-compose.ghcr.yml down
   docker compose -f docker-compose.ghcr.yml pull
   docker compose -f docker-compose.ghcr.yml up -d --force-recreate
   ```

4. **Check Nginx logs:**
   ```bash
   sudo tail -f /var/log/nginx/error.log
   ```

5. **VPS CPU overload?**
   - This should NOT happen anymore (no building on VPS)
   - If it does, restart from hosting panel
   - Check logs to see what's consuming CPU


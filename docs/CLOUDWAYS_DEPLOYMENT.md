# Cloudways Shared Hosting Deployment Guide

## ⚠️ CRITICAL: Shared Environment Considerations

This application will run on a **shared Cloudways server with 40+ other applications**. This setup is specifically designed to:

- ✅ Use custom ports to avoid conflicts
- ✅ Isolated Docker network (won't interfere with other apps)
- ✅ Minimal resource usage (conservative limits)
- ✅ Unique container names (prefixed with `bloodsa_doctors_`)
- ✅ Named volumes (prefixed to avoid collisions)
- ✅ No Nginx conflicts (Cloudways manages reverse proxy)

---

## 🔧 Architecture on Cloudways

```
Cloudways Nginx/Apache (Port 80/443)
         ↓
    [ Your App ]
         ↓
  Backend Container (Port 8880)
         ↓
  MongoDB Container (Internal only)
         ↓
  Frontend Container (Port 8881) - Optional
```

**Key Differences from Standard Setup**:
- No standalone Nginx container (Cloudways provides this)
- Custom ports (8880, 8881 instead of 80, 8080)
- Isolated Docker network with unique subnet
- Reduced resource limits

---

## 📋 Pre-Deployment Checklist

### 1. Check Available Ports on Cloudways

```bash
# SSH into your Cloudways server
ssh master@your-server-ip -p PORT

# Check what ports are in use
sudo netstat -tulpn | grep LISTEN | grep -E ":(88|89|90)[0-9]{2}"

# Or check Docker ports
docker ps --format "table {{.Names}}\t{{.Ports}}"
```

**Find 2 available ports** (e.g., 8880, 8881). If these are taken, update `.env.cloudways`:
```
BLOODSA_BACKEND_PORT=8882  # Change to available port
BLOODSA_FRONTEND_PORT=8883  # Change to available port
```

### 2. Check Docker Resource Usage

```bash
# See what's already running
docker stats --no-stream

# Check available resources
free -h
df -h
```

**Your app will use**:
- CPU: ~0.75 cores (max 1.0 in bursts)
- RAM: ~900 MB (max 1.15 GB in bursts)
- Disk: ~5-10 GB for images, ~10-20 GB for data

---

## 🚀 Deployment Steps

### Step 1: Clone Repository

```bash
cd ~/applications  # Or wherever Cloudways puts apps
git clone git@github.com:Tiny-Optics/doctors_workspaceBSA.git bloodsa-doctors
cd bloodsa-doctors
```

### Step 2: Configure Environment

```bash
# Copy Cloudways-specific template
cp .env.cloudways.example .env.cloudways

# Edit configuration
nano .env.cloudways
```

**MUST CONFIGURE**:
1. `BLOODSA_BACKEND_PORT` - Unique port (check availability!)
2. `BLOODSA_FRONTEND_PORT` - Unique port (check availability!)
3. `BLUEPRINT_DB_ROOT_PASSWORD` - Strong password
4. `JWT_SECRET` - Generate: `openssl rand -base64 32`
5. `DROPBOX_APP_API_ACCESS_TOKEN` - Your token
6. `VITE_API_URL` - Your Cloudways domain URL

### Step 3: Deploy

```bash
# Build and start containers
docker compose -f docker-compose.cloudways.yml --env-file .env.cloudways up -d --build

# Check status
docker compose -f docker-compose.cloudways.yml ps
```

### Step 4: Configure Cloudways Reverse Proxy

In **Cloudways Control Panel**:

1. **Add Application** → Custom Application
2. **Application URL**: doctors.bloodsa.org.za
3. **Configure Reverse Proxy**:

Add this to your Nginx/Apache config via Cloudways:

**For Nginx** (Cloudways → Server → Settings & Packages → Nginx Config):
```nginx
location /doctors-api/ {
    proxy_pass http://localhost:8880/api/;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

location /doctors-uploads/ {
    proxy_pass http://localhost:8880/uploads/;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
}

location /doctors/ {
    proxy_pass http://localhost:8881/;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_cache_bypass $http_upgrade;
}
```

**For Apache** (if Cloudways uses Apache):
```apache
ProxyPreserveHost On
ProxyPass /doctors-api/ http://localhost:8880/api/
ProxyPassReverse /doctors-api/ http://localhost:8880/api/

ProxyPass /doctors-uploads/ http://localhost:8880/uploads/
ProxyPassReverse /doctors-uploads/ http://localhost:8880/uploads/

ProxyPass /doctors/ http://localhost:8881/
ProxyPassReverse /doctors/ http://localhost:8881/
```

---

## 🔒 Security for Shared Environment

### 1. Network Isolation

The `bloodsa_isolated_network` uses a unique subnet (172.30.0.0/16) that won't conflict with other apps.

### 2. Container Naming

All containers prefixed with `bloodsa_doctors_` to avoid naming conflicts.

### 3. Volume Naming

All volumes prefixed with `bloodsa_doctors_` to prevent collisions.

### 4. Port Binding

**Internal ports are NOT exposed to host** except for the two custom ports you configure. This prevents conflicts.

### 5. Resource Limits

Conservative limits ensure your app doesn't starve other applications:
- MongoDB: 512 MB RAM max
- Backend: 512 MB RAM max
- Frontend: 128 MB RAM max
- Total: ~1.15 GB max (safe for shared environment)

---

## 🎯 Alternative: Managed Database Option

**Even Better for Shared Hosting**: Use Cloudways' managed MongoDB instead of running your own container.

### Benefits:
- No MongoDB container (saves ~512 MB RAM)
- Cloudways handles backups
- Better performance
- No resource competition

### Configuration:

If Cloudways provides managed MongoDB:

```yaml
# In .env.cloudways
BLUEPRINT_DB_HOST=localhost  # Or Cloudways MongoDB host
BLUEPRINT_DB_PORT=27017      # Cloudways MongoDB port
BLUEPRINT_DB_USERNAME=your_cloudways_db_user
BLUEPRINT_DB_ROOT_PASSWORD=your_cloudways_db_password
```

Then use this minimal `docker-compose.cloudways-minimal.yml`:

```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.prod
    container_name: bloodsa_doctors_backend
    restart: unless-stopped
    ports:
      - "${BLOODSA_BACKEND_PORT:-8880}:8080"
    environment:
      - PORT=8080
      - APP_ENV=production
      - GIN_MODE=release
      - BLUEPRINT_DB_HOST=${BLUEPRINT_DB_HOST}
      - BLUEPRINT_DB_PORT=${BLUEPRINT_DB_PORT}
      - BLUEPRINT_DB_USERNAME=${BLUEPRINT_DB_USERNAME}
      - BLUEPRINT_DB_ROOT_PASSWORD=${BLUEPRINT_DB_ROOT_PASSWORD}
      - BLUEPRINT_DB_DATABASE=${BLUEPRINT_DB_DATABASE}
      - JWT_SECRET=${JWT_SECRET}
      - DROPBOX_APP_API_ACCESS_TOKEN=${DROPBOX_APP_API_ACCESS_TOKEN}
      - DROPBOX_APP_KEY=${DROPBOX_APP_KEY}
      - DROPBOX_APP_SECRET=${DROPBOX_APP_SECRET}
      - DROPBOX_APP_PARENT_FOLDER=${DROPBOX_APP_PARENT_FOLDER}
    volumes:
      - bloodsa_sop_uploads:/app/uploads
    networks:
      - bloodsa_isolated_network
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
    logging:
      driver: "json-file"
      options:
        max-size: "5m"
        max-file: "2"

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.prod
      args:
        - VITE_API_URL=${VITE_API_URL}
    container_name: bloodsa_doctors_frontend
    restart: unless-stopped
    ports:
      - "${BLOODSA_FRONTEND_PORT:-8881}:80"
    networks:
      - bloodsa_isolated_network
    depends_on:
      - backend
    deploy:
      resources:
        limits:
          cpus: '0.25'
          memory: 128M
    logging:
      driver: "json-file"
      options:
        max-size: "5m"
        max-file: "2"

networks:
  bloodsa_isolated_network:
    name: bloodsa_doctors_workspace_network
    driver: bridge

volumes:
  bloodsa_sop_uploads:
    name: bloodsa_doctors_sop_uploads
    driver: local
```

---

## 📊 Resource Usage Summary

### Full Setup (with MongoDB)
- **RAM**: ~900 MB average, 1.15 GB peak
- **CPU**: ~0.75 cores average, 1.0 core peak
- **Disk**: ~15-25 GB total
- **Network**: Isolated (no conflicts)

### Minimal Setup (Cloudways MongoDB)
- **RAM**: ~400 MB average, 640 MB peak
- **CPU**: ~0.5 cores average, 0.75 core peak
- **Disk**: ~5-10 GB total
- **Network**: Isolated (no conflicts)

---

## 🔍 Verification After Deployment

### 1. Check No Port Conflicts

```bash
# Check your ports aren't used by others
sudo netstat -tulpn | grep :8880
sudo netstat -tulpn | grep :8881

# If taken, update .env.cloudways and redeploy
```

### 2. Test Backend

```bash
# Direct test (on server)
curl http://localhost:8880/health

# Should return: {"message":"It's healthy"}
```

### 3. Test Frontend

```bash
curl http://localhost:8881
# Should return HTML
```

### 4. Check Resource Usage

```bash
# Ensure not using too much
docker stats --no-stream bloodsa_doctors_backend bloodsa_doctors_mongodb
```

---

## 🚨 Troubleshooting Shared Environment Issues

### Issue: "Port already in use"

```bash
# Find what's using the port
sudo lsof -i :8880

# Change to different port in .env.cloudways
BLOODSA_BACKEND_PORT=8882  # Try different port
```

### Issue: "Out of memory"

```bash
# Check system memory
free -h

# If other apps using too much, reduce your limits
# Edit docker-compose.cloudways.yml
memory: 256M  # Instead of 512M
```

### Issue: "Network conflicts"

```bash
# Check existing Docker networks
docker network ls

# Our network should be isolated: 172.30.0.0/16
# If conflicts, change subnet in docker-compose.cloudways.yml
```

### Issue: "Container name conflicts"

```bash
# Check container names
docker ps -a | grep bloodsa

# All should be prefixed with bloodsa_doctors_
# If conflicts, they use unique names
```

---

## 💡 Best Practices for Shared Hosting

### 1. Monitor Resource Usage

```bash
# Daily check
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"

# Set up alerts if > 80% of limits
```

### 2. Keep Logs Small

Logs auto-rotate at 5MB (instead of 10MB) to save disk space.

```bash
# Check log sizes
docker inspect --format='{{.LogPath}}' bloodsa_doctors_backend | xargs ls -lh
```

### 3. Scheduled Restarts (Optional)

Some shared hosts prefer regular restarts:

```bash
# Add to crontab (weekly restart at 3 AM Sunday)
0 3 * * 0 cd /path/to/bloodsa-doctors && docker compose -f docker-compose.cloudways.yml restart
```

### 4. Backup Before Updates

```bash
# Always backup first on shared hosting
docker exec bloodsa_doctors_mongodb mongodump --out=/tmp/backup
docker cp bloodsa_doctors_mongodb:/tmp/backup ./backup_$(date +%Y%m%d)
```

---

## 🔄 Update Workflow for Shared Server

```bash
# 1. Check current resource usage
docker stats --no-stream

# 2. Backup
docker exec bloodsa_doctors_mongodb mongodump --out=/dump
docker cp bloodsa_doctors_mongodb:/dump ./backup

# 3. Pull updates
git pull origin main

# 4. Rebuild (one service at a time to minimize impact)
docker compose -f docker-compose.cloudways.yml build backend
docker compose -f docker-compose.cloudways.yml up -d --no-deps backend

docker compose -f docker-compose.cloudways.yml build frontend
docker compose -f docker-compose.cloudways.yml up -d --no-deps frontend

# 5. Verify
curl http://localhost:8880/health
```

---

## 📞 Cloudways-Specific Support

### Access Logs

```bash
# Application logs
docker logs bloodsa_doctors_backend --tail=100

# Cloudways system logs
# Use Cloudways control panel: Servers → Your Server → Logs
```

### Database Access

```bash
# Shell into your MongoDB
docker exec -it bloodsa_doctors_mongodb mongosh -u ${BLUEPRINT_DB_USERNAME} -p
```

### Application Shell

```bash
# Backend shell
docker exec -it bloodsa_doctors_backend sh

# Frontend shell
docker exec -it bloodsa_doctors_frontend sh
```

---

## ⚡ Quick Deploy Command

```bash
# Make the deployment script executable
chmod +x deploy-cloudways.sh

# Deploy
./deploy-cloudways.sh
```

---

## 🎯 Recommended Cloudways Server Specs

Since you're sharing with 40 apps:

### Server-Level (Total)
- **CPU**: 8+ vCPUs (your app uses ~0.75 cores)
- **RAM**: 16+ GB (your app uses ~900 MB)
- **Storage**: 200+ GB (your app uses ~15-25 GB)

### Your Application Allocation
- **RAM**: 1.15 GB max (with buffers)
- **CPU**: 1.0 core max (with bursts)
- **Disk**: 25 GB max (including data growth)

---

## 🔒 Security Considerations

### Network Isolation

Your app uses **172.30.0.0/16** subnet - completely isolated from other apps.

### Container Privileges

All containers run as **non-root users** - can't affect other apps.

### Resource Limits

Hard limits prevent your app from consuming resources needed by others.

### Firewall

Only your custom ports (8880, 8881) are exposed **locally** - Cloudways nginx handles external access.

---

## 📊 Monitoring on Shared Server

### Daily Checks

```bash
# Resource usage
docker stats --no-stream | grep bloodsa

# Health status
curl http://localhost:8880/health

# Disk usage
du -sh ~/applications/bloodsa-doctors
```

### Weekly Maintenance

```bash
# Cleanup unused Docker resources
docker system prune -f

# Check logs aren't too large
find . -name "*.log" -exec ls -lh {} \;

# Verify backups
ls -lh backups/
```

---

## ⚠️ Important Notes for Cloudways

1. **No Root Access**: Cloudways limits some operations - work within their constraints
2. **Use Cloudways SSL**: Let them manage SSL certificates via their panel
3. **No System Services**: Can't install systemd services - use Docker only
4. **Shared IP**: You'll share IP with other sites - use domain-based routing
5. **Resource Monitoring**: Cloudways may throttle if you exceed fair use

---

## 🎛️ Cloudways Control Panel Setup

### 1. Access Server

Cloudways → Servers → Your Server → Access Details

### 2. Add Domain

Cloudways → Applications → Add Domain → `doctors.bloodsa.org.za`

### 3. Configure SSL

Cloudways → Applications → Your App → SSL Certificate → Let's Encrypt (free)

### 4. Set Up Reverse Proxy

Cloudways → Applications → Application Settings → Advanced

Add the Nginx config provided in "Step 4" above.

---

## 🆘 Emergency Procedures

### App Using Too Much Memory

```bash
# Stop frontend (optional service)
docker stop bloodsa_doctors_frontend

# Reduce MongoDB cache
docker compose -f docker-compose.cloudways.yml down
# Edit docker-compose.cloudways.yml, reduce MongoDB memory to 256M
docker compose -f docker-compose.cloudways.yml up -d
```

### Port Conflict After Deployment

```bash
# Stop your app
docker compose -f docker-compose.cloudways.yml down

# Change ports in .env.cloudways
nano .env.cloudways

# Restart
docker compose -f docker-compose.cloudways.yml up -d
```

### Need to Remove Completely

```bash
# Stop and remove everything
docker compose -f docker-compose.cloudways.yml down -v

# Remove images
docker rmi bloodsa/doctors-backend bloodsa/doctors-frontend

# Clean up
docker system prune -af
```

---

## ✅ Deployment Checklist

### Before Deploying
- [ ] Checked available ports (8880, 8881 free)
- [ ] Verified server has 1 GB+ free RAM
- [ ] Configured .env.cloudways with all secrets
- [ ] Generated strong JWT_SECRET
- [ ] Confirmed Dropbox credentials work

### During Deployment
- [ ] Containers built successfully
- [ ] All services started
- [ ] Health checks passing
- [ ] No port conflict errors
- [ ] Resource usage acceptable

### After Deployment
- [ ] Backend health endpoint works
- [ ] Frontend loads
- [ ] Configured Cloudways reverse proxy
- [ ] SSL working (via Cloudways)
- [ ] Login works
- [ ] Can create users
- [ ] SOP categories display
- [ ] File downloads work

### Ongoing
- [ ] Monitor resource usage daily
- [ ] Check logs weekly
- [ ] Backup weekly (minimum)
- [ ] Update monthly
- [ ] Review Cloudways limits

---

## 💰 Cost Considerations

With shared hosting, you're paying for:
- **Server resources** (shared with 39 other apps)
- Your app's portion: ~2-5% of server resources
- Very cost-effective!

Monitor to ensure you don't exceed your fair share.

---

## 🎓 Pro Tips for Cloudways

1. **Use Cloudways CLI**: Easier management
   ```bash
   pip install cloudways-api
   ```

2. **Monitor via Cloudways Panel**: Built-in monitoring is excellent

3. **Leverage Cloudways Backups**: Enable automated server backups

4. **Use Cloudways CDN**: For static assets (optional)

5. **CloudwaysBot**: Slack integration for alerts

---

## 📚 Additional Resources

- **Cloudways Documentation**: https://support.cloudways.com
- **Docker on Cloudways**: Check their guides for Docker app deployment
- **General Deployment**: See `docs/PRODUCTION_DEPLOYMENT.md`

---

## ✨ Summary

Your Cloudways-optimized setup:
- ✅ **Isolated** - Won't interfere with other apps
- ✅ **Efficient** - Uses <1 GB RAM total
- ✅ **Secure** - Unique network, non-root containers
- ✅ **Portable** - Easy to move if needed
- ✅ **Safe** - Resource limits protect neighbors

**Deploy command**: 
```bash
docker compose -f docker-compose.cloudways.yml --env-file .env.cloudways up -d --build
```


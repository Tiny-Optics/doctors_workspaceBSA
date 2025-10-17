# Production Deployment Guide

## Prerequisites

- Ubuntu 22.04+ server
- Docker & Docker Compose installed
- Domain name configured (e.g., doctors.bloodsa.org.za)
- SSH access to server
- Minimum 2 vCPU / 4 GB RAM / 80 GB storage

## 🚀 Quick Deployment

### 1. Server Setup

```bash
# SSH into your server
ssh user@your-server-ip

# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo apt install docker-compose-plugin -y

# Logout and login again for docker group to take effect
exit
ssh user@your-server-ip
```

### 2. Clone Repository

```bash
# Clone your repository
git clone git@github.com:Tiny-Optics/doctors_workspaceBSA.git
cd doctors_workspaceBSA

# Or if using HTTPS
git clone https://github.com/Tiny-Optics/doctors_workspaceBSA.git
cd doctors_workspaceBSA
```

### 3. Configure Environment

```bash
# Copy production environment template
cp .env.production.example .env.production

# Edit with your production values
nano .env.production
```

**IMPORTANT**: Update these values:
- `BLUEPRINT_DB_USERNAME` - Database username
- `BLUEPRINT_DB_ROOT_PASSWORD` - Strong database password
- `JWT_SECRET` - Generate with: `openssl rand -base64 32`
- `DROPBOX_APP_API_ACCESS_TOKEN` - Your Dropbox token
- `VITE_API_URL` - Your production domain URL
- `DOMAIN` - Your domain name

### 4. Create Required Directories

```bash
# Create nginx directories
mkdir -p nginx/conf.d nginx/ssl

# Set permissions for uploads
sudo mkdir -p backend/uploads/sops
sudo chown -R $USER:$USER backend/uploads
```

### 5. Build and Start Containers

```bash
# Build images
docker compose -f docker-compose.prod.yml build

# Start services
docker compose -f docker-compose.prod.yml up -d

# Check status
docker compose -f docker-compose.prod.yml ps

# View logs
docker compose -f docker-compose.prod.yml logs -f
```

### 6. Verify Deployment

```bash
# Check health
curl http://localhost/health

# Should return: {"message":"It's healthy"}
```

### 7. Initialize Database

```bash
# Create first super admin user
# You'll need to create a seed script or use the API to create the first admin
# See backend/cmd/seed/ for examples
```

## 🔒 SSL/HTTPS Setup (Let's Encrypt)

### Install Certbot

```bash
# Add certbot
sudo snap install --classic certbot
sudo ln -s /snap/bin/certbot /usr/bin/certbot
```

### Get SSL Certificate

```bash
# Stop nginx temporarily
docker compose -f docker-compose.prod.yml stop nginx

# Get certificate
sudo certbot certonly --standalone -d your-domain.bloodsa.org.za

# Certificates will be in: /etc/letsencrypt/live/your-domain.bloodsa.org.za/
```

### Configure Nginx for SSL

```bash
# Copy certificates to nginx directory
sudo cp /etc/letsencrypt/live/your-domain.bloodsa.org.za/fullchain.pem nginx/ssl/
sudo cp /etc/letsencrypt/live/your-domain.bloodsa.org.za/privkey.pem nginx/ssl/
sudo chown -R $USER:$USER nginx/ssl

# Edit nginx.conf and uncomment the HTTPS server block
nano nginx/nginx.conf

# Restart nginx
docker compose -f docker-compose.prod.yml restart nginx
```

### Auto-Renewal

```bash
# Add to crontab
sudo crontab -e

# Add this line (runs twice daily)
0 0,12 * * * certbot renew --quiet && docker compose -f /path/to/doctors_workspaceBSA/docker-compose.prod.yml restart nginx
```

## 🔧 Management Commands

### View Logs

```bash
# All containers
docker compose -f docker-compose.prod.yml logs -f

# Specific container
docker compose -f docker-compose.prod.yml logs -f backend
docker compose -f docker-compose.prod.yml logs -f frontend
docker compose -f docker-compose.prod.yml logs -f mongodb
docker compose -f docker-compose.prod.yml logs -f nginx
```

### Restart Services

```bash
# Restart all
docker compose -f docker-compose.prod.yml restart

# Restart specific service
docker compose -f docker-compose.prod.yml restart backend
```

### Update Application

```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up -d

# Or zero-downtime update
docker compose -f docker-compose.prod.yml up -d --no-deps --build backend
```

### Stop Services

```bash
# Stop all containers
docker compose -f docker-compose.prod.yml down

# Stop and remove volumes (WARNING: deletes data!)
docker compose -f docker-compose.prod.yml down -v
```

## 💾 Backup & Restore

### Automated Backup Script

Create `/root/backup.sh`:

```bash
#!/bin/bash
BACKUP_DIR="/root/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Backup MongoDB
docker exec bloodsa_mongodb_prod mongodump \
  --username=$BLUEPRINT_DB_USERNAME \
  --password=$BLUEPRINT_DB_ROOT_PASSWORD \
  --out=/dump

docker cp bloodsa_mongodb_prod:/dump $BACKUP_DIR/mongo_$DATE

# Backup uploaded images
docker run --rm -v sop_uploads:/data -v $BACKUP_DIR:/backup \
  alpine tar czf /backup/uploads_$DATE.tar.gz /data

# Keep only last 7 days
find $BACKUP_DIR -type f -mtime +7 -delete

echo "Backup completed: $DATE"
```

### Schedule Backups

```bash
# Make executable
chmod +x /root/backup.sh

# Add to crontab (daily at 2 AM)
sudo crontab -e
0 2 * * * /root/backup.sh >> /var/log/backup.log 2>&1
```

### Restore from Backup

```bash
# Restore MongoDB
docker exec -i bloodsa_mongodb_prod mongorestore \
  --username=$BLUEPRINT_DB_USERNAME \
  --password=$BLUEPRINT_DB_ROOT_PASSWORD \
  --drop /dump

# Restore uploads
docker run --rm -v sop_uploads:/data -v /root/backups:/backup \
  alpine tar xzf /backup/uploads_YYYYMMDD_HHMMSS.tar.gz -C /
```

## 📊 Monitoring

### Resource Usage

```bash
# Real-time container stats
docker stats

# Disk usage
docker system df

# Volume sizes
docker volume ls -q | xargs docker volume inspect | grep -A 5 Mountpoint
```

### Health Checks

```bash
# Check all container health
docker compose -f docker-compose.prod.yml ps

# Manual health check
curl http://localhost/health
curl http://localhost/api/auth/me  # Should return 401 (expected)
```

### Log Rotation

Docker automatically rotates logs (configured to max 10MB x 3 files per container).

View current log sizes:
```bash
docker inspect --format='{{.LogPath}}' bloodsa_backend_prod | xargs ls -lh
```

## 🔥 Firewall Configuration

```bash
# Enable UFW firewall
sudo ufw enable

# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Check status
sudo ufw status
```

## 🚨 Troubleshooting

### Containers Won't Start

```bash
# Check logs
docker compose -f docker-compose.prod.yml logs

# Check specific container
docker logs bloodsa_backend_prod

# Verify environment variables
docker compose -f docker-compose.prod.yml config
```

### MongoDB Connection Issues

```bash
# Check MongoDB is running
docker exec bloodsa_mongodb_prod mongosh --eval "db.adminCommand('ping')"

# Check connection from backend
docker exec bloodsa_backend_prod wget -O- http://mongodb:27017
```

### Out of Disk Space

```bash
# Clean up Docker
docker system prune -a --volumes

# Check disk usage
df -h
docker system df
```

### High Memory Usage

```bash
# Check container stats
docker stats --no-stream

# Restart heavy containers
docker compose -f docker-compose.prod.yml restart mongodb
```

## 📈 Performance Tuning

### MongoDB Optimization

Add to `docker-compose.prod.yml` mongodb service:

```yaml
command: mongod --wiredTigerCacheSizeGB 1
```

### Backend Optimization

Set Go runtime variables:

```yaml
environment:
  - GOMAXPROCS=2
  - GOGC=100
```

## 🔐 Security Checklist

- [ ] Change all default passwords
- [ ] Generate strong JWT_SECRET
- [ ] Enable HTTPS with SSL certificate
- [ ] Configure firewall (UFW)
- [ ] Set up automated backups
- [ ] Restrict MongoDB port (don't expose externally)
- [ ] Regular security updates: `sudo apt update && sudo apt upgrade`
- [ ] Monitor audit logs
- [ ] Set up fail2ban for SSH protection

## 📱 Post-Deployment

### 1. Test Authentication

- Login as admin
- Create test user
- Test all features

### 2. Seed Initial Data

- Create SOP categories
- Upload category images
- Add test files to Dropbox

### 3. Monitor for 24-48 Hours

- Check logs regularly
- Monitor resource usage
- Test from different locations

### 4. Set Up Monitoring

Consider:
- UptimeRobot (free)
- Prometheus + Grafana
- CloudWatch (if on AWS)

## 📞 Support

**Logs Location**:
- Container logs: `docker compose logs`
- Nginx logs: Inside nginx container at `/var/log/nginx/`
- Application logs: Stdout/stderr captured by Docker

**Quick Diagnostics**:
```bash
# Check everything
docker compose -f docker-compose.prod.yml ps
curl http://localhost/health
docker stats --no-stream
df -h
```

## ⚡ Quick Commands Reference

```bash
# Start
docker compose -f docker-compose.prod.yml up -d

# Stop
docker compose -f docker-compose.prod.yml down

# Restart
docker compose -f docker-compose.prod.yml restart

# Update
git pull && docker compose -f docker-compose.prod.yml up -d --build

# Backup
/root/backup.sh

# Logs
docker compose -f docker-compose.prod.yml logs -f --tail=100

# Clean
docker system prune -f
```


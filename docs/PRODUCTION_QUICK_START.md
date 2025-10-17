# Production Quick Start Guide

## 🚀 One-Command Deployment

Once your server is set up with Docker, deploy with:

```bash
make prod
```

That's it! The deployment script handles everything.

## 📋 Pre-Deployment Checklist

### 1. Server Preparation (One-Time Setup)

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install docker-compose-plugin -y

# Add your user to docker group
sudo usermod -aG docker $USER
# Logout and login again

# Enable firewall
sudo ufw enable
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

### 2. Clone Repository

```bash
git clone git@github.com:Tiny-Optics/doctors_workspaceBSA.git
cd doctors_workspaceBSA
```

### 3. Configure Environment

```bash
# Copy template
cp .env.production.example .env.production

# Edit with your values
nano .env.production
```

**Required Changes**:
- Database password
- JWT secret (generate: `openssl rand -base64 32`)
- Dropbox credentials
- Domain name

### 4. Deploy

```bash
make prod
```

### 5. Set Up SSL (After Domain Points to Server)

```bash
# Install certbot
sudo snap install --classic certbot

# Get certificate
sudo certbot certonly --standalone -d your-domain.bloodsa.org.za

# Copy to nginx
sudo cp /etc/letsencrypt/live/your-domain.bloodsa.org.za/fullchain.pem nginx/ssl/
sudo cp /etc/letsencrypt/live/your-domain.bloodsa.org.za/privkey.pem nginx/ssl/
sudo chown -R $USER:$USER nginx/ssl

# Enable HTTPS in nginx.conf (uncomment the HTTPS server block)
nano nginx/nginx.conf

# Restart nginx
docker compose -f docker-compose.prod.yml restart nginx
```

## 🔄 Common Operations

### View Logs
```bash
make prod-logs
```

### Restart Services
```bash
make prod-restart
```

### Update Application
```bash
make update
```

### Backup Data
```bash
make backup
```

### Check Status
```bash
make prod-ps
```

## ✅ Verification

After deployment:

1. **Health Check**:
   ```bash
   curl http://localhost/health
   # Should return: {"message":"It's healthy"}
   ```

2. **Access Frontend**:
   ```
   http://your-server-ip
   ```

3. **Check Logs**:
   ```bash
   docker compose -f docker-compose.prod.yml logs -f
   ```

4. **Verify Containers**:
   ```bash
   docker compose -f docker-compose.prod.yml ps
   # All should show "healthy" status
   ```

## 🆘 Troubleshooting

### Services Won't Start
```bash
# Check logs
make prod-logs

# Check specific service
docker logs bloodsa_backend_prod
```

### Can't Access Application
```bash
# Check firewall
sudo ufw status

# Check nginx
docker logs bloodsa_nginx_prod

# Test backend directly
curl http://localhost:8080/health
```

### Database Issues
```bash
# Check MongoDB
docker exec bloodsa_mongodb_prod mongosh --eval "db.adminCommand('ping')"
```

## 📊 Resource Monitoring

```bash
# Container resource usage
docker stats

# Disk usage
df -h
docker system df
```

## 🔒 Security

After deployment:
- [ ] Change default passwords
- [ ] Set up SSL (see step 5 above)
- [ ] Configure firewall
- [ ] Set up automated backups
- [ ] Test authentication
- [ ] Monitor audit logs

## 📞 Support

For detailed information, see:
- Full deployment guide: `docs/PRODUCTION_DEPLOYMENT.md`
- Hosting recommendations: Ask for the hosting guide
- Troubleshooting: `docs/SOP_TROUBLESHOOTING.md`

## ⚡ Pro Tips

1. **Always backup before updates**: `make backup`
2. **Monitor logs after deployment**: `make prod-logs`
3. **Test in staging first** if possible
4. **Keep .env.production secure** - never commit it!
5. **Use strong passwords** for database and JWT secret

---

**Need help?** Check the logs first:
```bash
docker compose -f docker-compose.prod.yml logs --tail=100
```


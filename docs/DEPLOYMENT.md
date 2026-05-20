# Deployment Guide

## Automated deployment with GitHub Actions

### How it works

1. **GitHub Actions builds** Docker images on GitHub runners
2. **Pushes images** to GitHub Container Registry (GHCR)
3. **VPS pulls** pre-built images and runs them (no builds on the server)

Production uses **MongoDB Atlas** via `MONGO_URI`. There is no local MongoDB container on the VPS.

---

## Initial VPS setup (one-time)

### Phase A: App directory and deploy SSH key

```bash
sudo mkdir -p /opt/doctors_workspace
sudo chown $USER:$USER /opt/doctors_workspace
cd /opt/doctors_workspace
git clone git@github.com:Tiny-Optics/doctors_workspaceBSA.git .
```

Create a deploy key for GitHub Actions (if missing):

```bash
ssh-keygen -t ed25519 -f ~/.ssh/github_deploy -N ""
cat ~/.ssh/github_deploy.pub   # add to repo Deploy keys or ~/.ssh/authorized_keys
cat ~/.ssh/github_deploy         # paste into GitHub secret SSH_PRIVATE_KEY
```

Ensure `git pull origin main` works non-interactively from `/opt/doctors_workspace`.

### Phase B: Production `.env` on the VPS

```bash
cd /opt/doctors_workspace
cp .env.prod.example .env
nano .env
```

Minimum contents (use new secrets if rebuilding from scratch):

```env
MONGO_URI=mongodb+srv://<user>:<password>@<cluster>.mongodb.net/?appName=DoctorsWorkspace
BLUEPRINT_DB_DATABASE=doctors_workspace
JWT_SECRET=<openssl rand -hex 64>
ENCRYPTION_KEY=<openssl rand -base64 32>
```

**MongoDB Atlas:** In Atlas → **Network Access**, allow the VPS public IP (or `0.0.0.0/0` temporarily for testing).

**Secrets:** Do not put `MONGO_URI`, `JWT_SECRET`, or `ENCRYPTION_KEY` in GitHub Actions secrets. They live only in the VPS `.env` file. If you restore encrypted data from backup, reuse the original `ENCRYPTION_KEY` and `JWT_SECRET`.

The frontend API URL for GHCR builds is set in `.github/workflows/deploy.yml` (`VITE_API_URL=https://workspace.bloodsa.org.za/api`). The `.env` `VITE_API_URL` is only used for local `docker-compose.prod.yml` builds.

### Phase C: Host Nginx and SSL

```bash
sudo apt-get install -y nginx certbot python3-certbot-nginx
cd /opt/doctors_workspace
sudo cp nginx.conf.production /etc/nginx/sites-available/doctors_workspace
sudo ln -sf /etc/nginx/sites-available/doctors_workspace /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
```

DNS for `workspace.bloodsa.org.za` must point to this VPS. Obtain certificates before enabling the SSL server block (or use certbot with a temporary HTTP-only config):

```bash
sudo certbot certonly --nginx -d workspace.bloodsa.org.za
sudo nginx -t
sudo systemctl enable nginx
sudo systemctl reload nginx
```

Test renewal: `sudo certbot renew --dry-run`.

### Phase D: First container start

```bash
cd /opt/doctors_workspace
# If GHCR packages are private:
# echo "$TOKEN" | docker login ghcr.io -u <github-username> --password-stdin
docker compose -f docker-compose.ghcr.yml pull
docker compose -f docker-compose.ghcr.yml up -d
docker compose -f docker-compose.ghcr.yml ps
docker compose -f docker-compose.ghcr.yml logs -f backend
```

Verify on the VPS:

```bash
curl -s http://127.0.0.1:8080/health
curl -sI http://127.0.0.1:3000/
```

Then test: `https://workspace.bloodsa.org.za` and `https://workspace.bloodsa.org.za/api/...`.

### Phase E: Reseed empty Atlas database

If the database was wiped, run seed commands from the VPS (requires Go 1.24+):

```bash
cd /opt/doctors_workspace/backend
set -a && source /opt/doctors_workspace/.env && set +a

go run cmd/seed/main.go
go run cmd/init-registry/main.go
go run cmd/seed-institutions/main.go
```

Optional: `seed-users`, `migrate-roles` (see `backend/Makefile`).

The `sop_uploads` volume starts empty; re-upload SOP files as needed.

### Phase F: GHCR packages and GitHub secrets

**GHCR visibility:** Make these packages **public** (or authenticate on the VPS to pull private images):

- `ghcr.io/tiny-optics/doctors-backend`
- `ghcr.io/tiny-optics/doctors-frontend`

Org **Packages** → package → **Package settings** → **Change visibility** → Public.

**GitHub Actions secrets** (`Settings → Secrets and variables → Actions`):

| Secret | Example |
|--------|---------|
| `SSH_HOST` | `169.255.58.102` |
| `SSH_USER` | `vmszjayo` |
| `SSH_PRIVATE_KEY` | Contents of `~/.ssh/github_deploy` on the VPS |

`GITHUB_TOKEN` is provided automatically by GitHub Actions.

---

## Deploying

Push to `main`:

```bash
git push origin main
```

The workflow builds images, pushes to GHCR, SSHs to the VPS, runs `git pull`, and restarts containers with `docker-compose.ghcr.yml`.

### Manual deployment

```bash
cd /opt/doctors_workspace
git pull origin main
docker compose -f docker-compose.ghcr.yml pull
docker compose -f docker-compose.ghcr.yml up -d --force-recreate
```

### Logs and status

```bash
docker compose -f docker-compose.ghcr.yml logs -f
docker compose -f docker-compose.ghcr.yml logs -f backend
docker compose -f docker-compose.ghcr.yml logs -f frontend
docker compose -f docker-compose.ghcr.yml ps
```

### URLs

- **Site:** https://workspace.bloodsa.org.za
- **API:** https://workspace.bloodsa.org.za/api

---

## Important notes

1. **Uploads persist** in the `sop_uploads` Docker volume across deploys
2. **Never commit `.env`** — it contains production secrets
3. **`ENCRYPTION_KEY` must not change** if you need to read existing encrypted data
4. **URL-encode** special characters in Atlas passwords inside `MONGO_URI`

---

## Troubleshooting

| Symptom | Check |
|---------|--------|
| Backend crash loop | `docker compose logs backend`; Atlas IP whitelist; `MONGO_URI` |
| 502 from Nginx | Containers running; ports 8080/3000; `sudo tail -f /var/log/nginx/error.log` |
| Frontend wrong API | Rebuild frontend image after workflow deploy; hard-refresh browser |
| GHCR pull denied | Public packages or `docker login ghcr.io` on VPS |
| SSH deploy fails | `SSH_*` secrets; deploy key in `authorized_keys`; `/opt/doctors_workspace` exists |

```bash
docker compose -f docker-compose.ghcr.yml restart backend
docker compose -f docker-compose.ghcr.yml down && docker compose -f docker-compose.ghcr.yml pull && docker compose -f docker-compose.ghcr.yml up -d --force-recreate
```

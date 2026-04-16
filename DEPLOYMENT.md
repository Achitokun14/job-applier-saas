# DEPLOYMENT.md - VPS Deployment Guide

Complete guide for deploying Job Applier SaaS to a VPS with Caddy or Nginx.

## Prerequisites

- Ubuntu 20.04+ or Debian 11+ VPS
- Domain name pointed to your VPS IP
- Minimum 2GB RAM, 20GB disk
- Root or sudo access

## Quick Start (Caddy - Recommended)

### 1. Clone Repository

```bash
git clone <repository-url> /opt/job-applier-saas
cd /opt/job-applier-saas
```

### 2. Configure Environment

```bash
cp .env.production .env
nano .env
```

**Required changes:**
```env
DOMAIN=your-domain.com
EMAIL=admin@your-domain.com
DB_PASSWORD=<generate-strong-password>
JWT_SECRET=<generate-64-char-random>
LLM_API_KEY=your-api-key
```

Generate secure values:
```bash
# Database password
openssl rand -base64 32

# JWT secret
openssl rand -base64 64
```

### 3. Deploy

```bash
chmod +x scripts/deploy.sh
sudo ./scripts/deploy.sh
```

### 4. Verify

```bash
# Check services
docker compose -f docker-compose.production.yml ps

# Check logs
docker compose -f docker-compose.production.yml logs -f

# Test endpoint
curl https://your-domain.com/health
```

## Manual Deployment Steps

### Step 1: Server Setup

```bash
# Update system
apt update && apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
systemctl enable docker
systemctl start docker

# Install Docker Compose
apt install -y docker-compose-plugin

# Verify
docker --version
docker compose version
```

### Step 2: Configure Domain DNS

Point your domain to the VPS:
- A Record: `your-domain.com` → `YOUR_VPS_IP`
- A Record: `www.your-domain.com` → `YOUR_VPS_IP`

### Step 3: Configure Firewall

```bash
# Allow SSH, HTTP, HTTPS
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

### Step 4: Build and Start

```bash
cd /opt/job-applier-saas

# Build images
docker compose -f docker-compose.production.yml build

# Start services
docker compose -f docker-compose.production.yml up -d

# Check status
docker compose -f docker-compose.production.yml ps
```

### Step 5: SSL Certificates (Caddy)

Caddy automatically handles SSL:
1. Caddy detects your domain
2. Requests Let's Encrypt certificate
3. Configures HTTPS automatically

**No manual SSL setup needed!**

## Alternative: Nginx Deployment

### Step 1: Get SSL Certificates

```bash
chmod +x scripts/setup-ssl.sh
sudo ./scripts/setup-ssl.sh
```

### Step 2: Deploy with Nginx

```bash
docker compose -f docker-compose.nginx.yml up -d
```

### Step 3: Configure Nginx

Edit `nginx/conf.d/default.conf`:
```nginx
server_name your-domain.com www.your-domain.com;
```

Restart:
```bash
docker compose -f docker-compose.nginx.yml restart nginx
```

## Systemd Service (Auto-start on boot)

The deploy script creates a systemd service. Manual setup:

```bash
cat > /etc/systemd/system/jobapplier.service << EOF
[Unit]
Description=Job Applier SaaS
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/job-applier-saas
ExecStart=/usr/bin/docker compose -f docker-compose.production.yml up -d
ExecStop=/usr/bin/docker compose -f docker-compose.production.yml down

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable jobapplier
systemctl start jobapplier
```

## Maintenance Commands

### View Logs
```bash
# All services
docker compose -f docker-compose.production.yml logs -f

# Specific service
docker compose -f docker-compose.production.yml logs -f backend
```

### Restart Services
```bash
# All
systemctl restart jobapplier

# Specific
docker compose -f docker-compose.production.yml restart backend
```

### Update Application
```bash
cd /opt/job-applier-saas
chmod +x scripts/update.sh
sudo ./scripts/update.sh
```

### Backup Database
```bash
docker exec jobapplier-postgres pg_dump -U jobapplier jobapplier > backup.sql
```

### Restore Database
```bash
cat backup.sql | docker exec -i jobapplier-postgres psql -U jobapplier jobapplier
```

## Troubleshooting

### Services not starting
```bash
# Check logs
docker compose -f docker-compose.production.yml logs

# Check individual service
docker logs jobapplier-backend
```

### SSL not working
```bash
# Caddy logs
docker logs jobapplier-caddy

# Check Caddy config
docker exec jobapplier-caddy cat /etc/caddy/Caddyfile
```

### Database connection issues
```bash
# Check PostgreSQL
docker exec jobapplier-postgres pg_isready -U jobapplier

# Check network
docker network inspect jobapplier-network
```

### Port conflicts
```bash
# Check what's using ports
ss -tlnp | grep -E ':(80|443|8080)'

# Stop conflicting services
systemctl stop nginx  # if nginx is running on host
```

## Security Checklist

- [ ] Changed default passwords
- [ ] Set strong JWT_SECRET
- [ ] Enabled firewall (ufw)
- [ ] SSL/HTTPS working
- [ ] Security headers configured
- [ ] Regular backups scheduled
- [ ] Log rotation configured
- [ ] Docker running as non-root

## Performance Tuning

### PostgreSQL
```bash
# Edit postgresql.conf in container
docker exec -it jobapplier-postgres bash
# Edit /var/lib/postgresql/data/postgresql.conf
```

### Nginx Workers
Edit `nginx/nginx.conf`:
```nginx
worker_processes auto;  # Use all CPU cores
```

### Docker Resource Limits
Add to docker-compose:
```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
```

## URLs After Deployment

| Service | URL |
|---------|-----|
| Frontend | https://your-domain.com |
| API | https://your-domain.com/api |
| Health | https://your-domain.com/health |

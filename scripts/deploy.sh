#!/bin/bash
set -e

echo "============================================"
echo "  Job Applier SaaS - VPS Deployment Script"
echo "============================================"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Please run as root (use sudo)${NC}"
    exit 1
fi

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"

command -v docker >/dev/null 2>&1 || {
    echo -e "${YELLOW}Docker not found. Installing...${NC}"
    curl -fsSL https://get.docker.com | sh
    systemctl enable docker
    systemctl start docker
}

command -v docker compose >/dev/null 2>&1 || {
    echo -e "${YELLOW}Docker Compose not found. Installing...${NC}"
    apt-get update
    apt-get install -y docker-compose-plugin
}

echo -e "${GREEN}Prerequisites OK${NC}"

# Check for .env file
if [ ! -f .env ]; then
    echo -e "${YELLOW}No .env file found.${NC}"
    if [ -f .env.production ]; then
        echo -e "${YELLOW}Copying .env.production to .env${NC}"
        cp .env.production .env
        echo -e "${RED}IMPORTANT: Edit .env with your actual values!${NC}"
        echo -e "${RED}Then run this script again.${NC}"
        exit 1
    else
        echo -e "${RED}No .env.production found. Creating template...${NC}"
        cp .env.example .env
        echo -e "${RED}IMPORTANT: Edit .env with your actual values!${NC}"
        exit 1
    fi
fi

# Load environment variables
set -a
source .env
set +a

# Validate required variables
if [ -z "$DOMAIN" ]; then
    echo -e "${RED}DOMAIN not set in .env${NC}"
    exit 1
fi

if [ -z "$DB_PASSWORD" ] || [ "$DB_PASSWORD" = "change-me-to-strong-password" ]; then
    echo -e "${RED}DB_PASSWORD not set or using default value${NC}"
    exit 1
fi

if [ -z "$JWT_SECRET" ] || [ "$JWT_SECRET" = "change-this-to-a-very-long-random-string" ]; then
    echo -e "${RED}JWT_SECRET not set or using default value${NC}"
    exit 1
fi

echo -e "${GREEN}Configuration validated${NC}"

# Create necessary directories
echo -e "${YELLOW}Creating directories...${NC}"
mkdir -p caddy-logs
mkdir -p nginx/ssl
mkdir -p nginx/conf.d
mkdir -p backend/data

# Pull latest images
echo -e "${YELLOW}Pulling latest images...${NC}"
docker compose -f docker-compose.production.yml pull

# Build services
echo -e "${YELLOW}Building services...${NC}"
docker compose -f docker-compose.production.yml build --no-cache

# Stop existing containers
echo -e "${YELLOW}Stopping existing containers...${NC}"
docker compose -f docker-compose.production.yml down 2>/dev/null || true

# Start services
echo -e "${YELLOW}Starting services...${NC}"
docker compose -f docker-compose.production.yml up -d

# Wait for services
echo -e "${YELLOW}Waiting for services to start...${NC}"
sleep 10

# Check health
echo -e "${YELLOW}Checking service health...${NC}"
docker compose -f docker-compose.production.yml ps

# Setup systemd service (optional)
if [ -f /etc/systemd/system ]; then
    echo -e "${YELLOW}Creating systemd service...${NC}"
    cat > /etc/systemd/system/jobapplier.service << EOF
[Unit]
Description=Job Applier SaaS
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=$PROJECT_DIR
ExecStart=/usr/bin/docker compose -f docker-compose.production.yml up -d
ExecStop=/usr/bin/docker compose -f docker-compose.production.yml down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable jobapplier.service
    echo -e "${GREEN}Systemd service created and enabled${NC}"
fi

echo ""
echo -e "${GREEN}============================================${NC}"
echo -e "${GREEN}  Deployment Complete!${NC}"
echo -e "${GREEN}============================================${NC}"
echo ""
echo -e "Services running:"
echo -e "  - Frontend:  https://$DOMAIN"
echo -e "  - API:       https://$DOMAIN/api"
echo -e "  - Health:    https://$DOMAIN/health"
echo ""
echo -e "Useful commands:"
echo -e "  - View logs:     docker compose -f docker-compose.production.yml logs -f"
echo -e "  - Stop:          docker compose -f docker-compose.production.yml down"
echo -e "  - Restart:       systemctl restart jobapplier"
echo -e "  - Status:        systemctl status jobapplier"
echo ""

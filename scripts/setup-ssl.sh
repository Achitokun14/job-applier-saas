#!/bin/bash
set -e

echo "============================================"
echo "  SSL Certificate Setup (Let's Encrypt)"
echo "============================================"

# Check if certbot is installed
if ! command -v certbot &> /dev/null; then
    echo "Installing certbot..."
    apt-get update
    apt-get install -y certbot
fi

# Load domain from .env
if [ -f .env ]; then
    source .env
fi

if [ -z "$DOMAIN" ]; then
    echo "Error: DOMAIN not set in .env"
    exit 1
fi

if [ -z "$EMAIL" ]; then
    echo "Error: EMAIL not set in .env"
    exit 1
fi

echo "Setting up SSL for: $DOMAIN"

# Create webroot directory
mkdir -p /var/www/certbot

# Stop nginx/caddy temporarily
docker compose -f docker-compose.production.yml stop caddy 2>/dev/null || true

# Get certificate
certbot certonly --standalone \
    --preferred-challenges http \
    -d "$DOMAIN" \
    -d "www.$DOMAIN" \
    --email "$EMAIL" \
    --agree-tos \
    --non-interactive

# Copy certificates
cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem nginx/ssl/
cp /etc/letsencrypt/live/$DOMAIN/privkey.pem nginx/ssl/

# Setup auto-renewal
cat > /etc/cron.d/certbot-renew << EOF
0 3 * * * root certbot renew --quiet --deploy-hook "cd $(pwd) && docker compose -f docker-compose.production.yml restart caddy"
EOF

echo ""
echo "SSL certificates installed!"
echo "Restart services: docker compose -f docker-compose.production.yml up -d"

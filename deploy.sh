#!/bin/bash

set -e

GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m"

log()    { echo -e "${GREEN}[✔] $1${NC}"; }
warn()   { echo -e "${YELLOW}[!] $1${NC}"; }
error()  { echo -e "${RED}[✘] $1${NC}"; exit 1; }

BINARY_NAME="golangsmtp"
SERVICE_PORT="8000"
NGINX_CONF_NAME="golangsmtp"

echo ""
echo "========================================"
echo "   golangsmtp — Deploy Script"
echo "========================================"
echo ""

if [ "$EUID" -ne 0 ]; then
  error "Please run as root or with sudo"
fi

read -rp "Enter your domain or server IP (e.g. api.example.com): " DOMAIN
[ -z "$DOMAIN" ] && error "Domain/IP cannot be empty"

read -rp "Enable HTTPS with Let's Encrypt? (requires a valid domain) [y/N]: " ENABLE_SSL
ENABLE_SSL="${ENABLE_SSL,,}"

read -rp "Enter your email for SSL certificate notifications (leave blank to skip): " SSL_EMAIL

echo ""

log "Checking Go installation..."
export PATH="/usr/local/go/bin:$PATH"
if ! command -v go &>/dev/null; then
  warn "Go not found — installing latest stable version..."
  GO_VERSION=$(curl -fsSL "https://go.dev/VERSION?m=text" | head -1)
  GO_TARBALL="${GO_VERSION}.linux-amd64.tar.gz"
  curl -fsSL "https://go.dev/dl/${GO_TARBALL}" -o "/tmp/${GO_TARBALL}"
  rm -rf /usr/local/go
  tar -C /usr/local -xzf "/tmp/${GO_TARBALL}"
  rm "/tmp/${GO_TARBALL}"
  echo 'export PATH=/usr/local/go/bin:$PATH' > /etc/profile.d/go.sh
  chmod +x /etc/profile.d/go.sh
  log "Go installed: $(go version)"
else
  log "Go already installed: $(go version)"
fi

log "Building Go binary..."
go build -o "$BINARY_NAME" ./cmd/server/main.go
log "Binary built: ./$BINARY_NAME"

if [ ! -f ".env" ]; then
  warn ".env file not found. Copying from .env.example..."
  [ -f ".env.example" ] && cp .env.example .env || error ".env.example not found either. Create a .env file first."
  warn "Please edit .env with your real credentials, then re-run this script."
  exit 1
fi

DEPLOY_DIR="/opt/$BINARY_NAME"
log "Creating deploy directory: $DEPLOY_DIR"
mkdir -p "$DEPLOY_DIR"

if command -v pm2 &>/dev/null && pm2 describe "$BINARY_NAME" &>/dev/null; then
  warn "Stopping running service before file copy..."
  pm2 stop "$BINARY_NAME" || true
fi

cp "$BINARY_NAME" "$DEPLOY_DIR/"
cp .env "$DEPLOY_DIR/"
[ -f ecosystem.config.js ] && cp ecosystem.config.js "$DEPLOY_DIR/"
chmod +x "$DEPLOY_DIR/$BINARY_NAME"
log "Files deployed to $DEPLOY_DIR"

if command -v nginx &>/dev/null; then
  warn "Nginx already installed — skipping installation."
else
  log "Installing Nginx..."
  if command -v apt-get &>/dev/null; then
    apt-get update -qq
    apt-get install -y nginx
  elif command -v dnf &>/dev/null; then
    dnf install -y nginx
  elif command -v yum &>/dev/null; then
    yum install -y nginx
  else
    error "Unsupported package manager. Install Nginx manually."
  fi
  systemctl enable nginx
  systemctl start nginx
  log "Nginx installed and started."
fi

if command -v apt-get &>/dev/null; then
  OS_FAMILY="debian"
else
  OS_FAMILY="rhel"
fi

if [ "$OS_FAMILY" = "debian" ]; then
  NGINX_CONF="/etc/nginx/sites-available/$NGINX_CONF_NAME"
  NGINX_LINK="/etc/nginx/sites-enabled/$NGINX_CONF_NAME"
else
  NGINX_CONF="/etc/nginx/conf.d/$NGINX_CONF_NAME.conf"
  NGINX_LINK=""
fi

log "Writing Nginx config to $NGINX_CONF..."
cat > "$NGINX_CONF" <<EOF
server {
    listen 80;
    server_name $DOMAIN;

    location / {
        proxy_pass         http://127.0.0.1:$SERVICE_PORT;
        proxy_http_version 1.1;
        proxy_set_header   Host              \$host;
        proxy_set_header   X-Real-IP         \$remote_addr;
        proxy_set_header   X-Forwarded-For   \$proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto \$scheme;
        proxy_read_timeout 30s;
        proxy_connect_timeout 10s;

        add_header Access-Control-Allow-Origin  "*" always;
        add_header Access-Control-Allow-Methods "GET, POST, OPTIONS" always;
        add_header Access-Control-Allow-Headers "Content-Type, Authorization" always;

        if (\$request_method = OPTIONS) {
            return 204;
        }
    }
}
EOF

if [ "$OS_FAMILY" = "debian" ]; then
  [ -L "$NGINX_LINK" ] && rm "$NGINX_LINK"
  ln -s "$NGINX_CONF" "$NGINX_LINK"
  [ -f /etc/nginx/sites-enabled/default ] && rm /etc/nginx/sites-enabled/default
fi

nginx -t || error "Nginx config test failed. Check $NGINX_CONF"
if systemctl is-active --quiet nginx; then
  systemctl reload nginx
  log "Nginx reloaded."
else
  systemctl enable nginx
  systemctl start nginx
  log "Nginx started."
fi

if [ "$OS_FAMILY" = "rhel" ] && command -v setsebool &>/dev/null; then
  setsebool -P httpd_can_network_connect 1
  log "SELinux: httpd_can_network_connect enabled."
fi

if command -v pm2 &>/dev/null; then
  warn "PM2 already installed — skipping."
else
  log "Installing PM2..."
  if ! command -v node &>/dev/null; then
    warn "Node.js not found. Installing via NodeSource..."
    if [ "$OS_FAMILY" = "debian" ]; then
      curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
      apt-get install -y nodejs
    else
      curl -fsSL https://rpm.nodesource.com/setup_lts.x | bash -
      dnf install -y nodejs || yum install -y nodejs
    fi
  fi
  npm install -g pm2
  log "PM2 installed."
fi

log "Starting service with PM2..."
cd "$DEPLOY_DIR"

pm2 describe "$BINARY_NAME" &>/dev/null && pm2 delete "$BINARY_NAME"

pm2 start "./$BINARY_NAME" \
  --name "$BINARY_NAME" \
  --interpreter none \
  --restart-delay 3000 \
  --max-restarts 10

pm2 save

if ! pm2 startup | grep -q "already"; then
  pm2 startup | tail -1 | bash || warn "Could not auto-configure pm2 startup. Run 'pm2 startup' manually."
  pm2 save
fi

log "Service started with PM2."

if [ "$ENABLE_SSL" = "y" ]; then
  log "Setting up SSL with Let's Encrypt..."
  if ! command -v certbot &>/dev/null; then
    log "Installing Certbot..."
    if [ "$OS_FAMILY" = "debian" ]; then
      apt-get install -y certbot python3-certbot-nginx
    else
      dnf install -y epel-release 2>/dev/null || yum install -y epel-release 2>/dev/null || true
      dnf install -y certbot python3-certbot-nginx 2>/dev/null || \
        yum install -y certbot python3-certbot-nginx 2>/dev/null || \
        error "Cannot install Certbot. Install it manually and run: certbot --nginx -d $DOMAIN"
    fi
  fi

  CERTBOT_CMD="certbot --nginx -d $DOMAIN --non-interactive --agree-tos --redirect"
  [ -n "$SSL_EMAIL" ] && CERTBOT_CMD="$CERTBOT_CMD --email $SSL_EMAIL" || CERTBOT_CMD="$CERTBOT_CMD --register-unsafely-without-email"

  $CERTBOT_CMD || warn "Certbot failed. Check that DNS for $DOMAIN points to this server."
  log "SSL certificate issued and Nginx updated."
fi

echo ""
echo "========================================"
log "Deployment complete!"
echo ""
echo "  Service dir : $DEPLOY_DIR"
echo "  PM2 name    : $BINARY_NAME"
echo "  Nginx conf  : $NGINX_CONF"
if [ "$ENABLE_SSL" = "y" ]; then
  echo "  URL         : https://$DOMAIN/api/send-email"
else
  echo "  URL         : http://$DOMAIN/api/send-email"
fi
echo ""
echo "  pm2 logs $BINARY_NAME   → view logs"
echo "  pm2 status               → process status"
echo "========================================"
echo ""

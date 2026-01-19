#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Xferant VPN ASCII Art
echo -e "${PURPLE}"
cat << "EOF"
__  __ ______ _______ _____   _   _ _______ 
\ \/ /|  ____|__   __|  __ \ | \ | |__   __|
 \  / | |__     | |  | |__) ||  \| |  | |   
 /  \ |  __|    | |  |  _  / | . ` |  | |   
/ /\ \| |____   | |  | | \ \ | |\  |  | |   
/_/  \_\______|  |_|  |_|  \_\|_| \_|  |_|   
EOF
echo -e "${NC}"

echo -e "${BLUE}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                 Xferant VPN Installer                       ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}❌ Please run as root: sudo ./install.sh${NC}"
    exit 1
fi

# Default values
DOMAIN=""
EMAIL=""
INSTALL_DIR="/opt/xferant-vpn"
REPO_URL="https://github.com/R3G1ST/vpn-system.git"

log_info() { echo -e "${GREEN}✅ [INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}⚠️ [WARN]${NC} $1"; }
log_error() { echo -e "${RED}❌ [ERROR]${NC} $1"; }

# Function to generate safe password (without + and /)
generate_safe_password() {
    openssl rand -base64 $1 | tr -d '+/' | head -c $2
}

# Interactive input function
get_user_input() {
    echo -e "${CYAN}📝 Configuration Setup${NC}"
    echo ""
    
    # Domain input
    local input_domain=""
    while [ -z "$input_domain" ]; do
        echo -n "🌐 Enter your domain (e.g., vpn.yourdomain.com): "
        read -r input_domain < /dev/tty
        if [ -z "$input_domain" ]; then
            log_error "Domain name cannot be empty"
        fi
    done
    DOMAIN="$input_domain"
    
    # Email input
    local input_email=""
    while [ -z "$input_email" ]; do
        echo -n "📧 Enter your email (for SSL certificates): "
        read -r input_email < /dev/tty
        if [ -z "$input_email" ]; then
            log_error "Email cannot be empty"
        elif [[ ! "$input_email" =~ ^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$ ]]; then
            log_error "Please enter a valid email address"
            input_email=""
        fi
    done
    EMAIL="$input_email"
    
    # Confirmation
    echo ""
    echo -e "${YELLOW}📋 Installation Summary${NC}"
    echo "   Domain: https://$DOMAIN"
    echo "   Email: $EMAIL"
    echo "   Directory: $INSTALL_DIR"
    echo ""
    echo -n "Proceed with installation? (y/N): "
    read -r confirmation < /dev/tty
    if [ "$confirmation" != "y" ] && [ "$confirmation" != "Y" ]; then
        log_info "Installation cancelled by user"
        exit 0
    fi
}

main() {
    log_info "Starting Xferant VPN automated installation..."
    
    get_user_input
    
    # Installation steps
    check_system
    install_dependencies
    clone_repository
    setup_environment
    setup_ssl
    start_services
    finalize_installation
    
    log_info "🎉 Xferant VPN installed successfully!"
    show_success_message
}

check_system() {
    log_info "Checking system compatibility..."
    
    # Check OS
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS_NAME=$NAME
        OS_VERSION=$VERSION_ID
    else
        log_error "Cannot determine operating system"
        exit 1
    fi
    
    # Supported OS check
    case $ID in
        ubuntu|debian|centos|rocky|almalinux)
            log_info "Supported OS detected: $OS_NAME $OS_VERSION"
            ;;
        *)
            log_error "Unsupported OS: $OS_NAME"
            exit 1
            ;;
    esac
    
    # Check architecture
    ARCH=$(uname -m)
    if [ "$ARCH" != "x86_64" ]; then
        log_error "Unsupported architecture: $ARCH"
        exit 1
    fi
}

install_dependencies() {
    log_info "Installing system dependencies..."
    
    if ! command -v docker &> /dev/null; then
        log_info "Installing Docker..."
        curl -fsSL https://get.docker.com | sh
        systemctl enable docker
        systemctl start docker
    else
        log_info "Docker is already installed"
    fi
    
    if ! docker compose version &> /dev/null; then
        log_info "Installing Docker Compose Plugin..."
        apt-get update
        apt-get install -y docker-compose-plugin
    else
        log_info "Docker Compose is already installed"
    fi
}

clone_repository() {
    log_info "Cloning Xferant VPN repository..."
    
    mkdir -p $INSTALL_DIR
    
    CURRENT_DIR=$(pwd)
    cd $(dirname $INSTALL_DIR)
    
    rm -rf $INSTALL_DIR
    
    git clone $REPO_URL $INSTALL_DIR
    
    cd $CURRENT_DIR
    cd $INSTALL_DIR
}

setup_environment() {
    log_info "Setting up environment..."
    
    cd $INSTALL_DIR
    
    POSTGRES_PASSWORD=$(generate_safe_password 32 32)
    JWT_SECRET=$(generate_safe_password 64 64)
    API_SECRET_KEY=$(generate_safe_password 48 48)
    
    cat > .env << EOF
# Xferant VPN Configuration
DOMAIN=$DOMAIN
EMAIL=$EMAIL

# Database
POSTGRES_DB=xferant_vpn
POSTGRES_USER=xferant_user
POSTGRES_PASSWORD=$POSTGRES_PASSWORD

# Security
JWT_SECRET=$JWT_SECRET
API_SECRET_KEY=$API_SECRET_KEY

# Telegram Bot
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_WEBHOOK_URL=https://$DOMAIN/api/telegram/webhook
EOF

    log_info "Environment configuration created"
}

setup_ssl() {
    log_info "Setting up SSL certificates..."
    
    cd $INSTALL_DIR
    
    mkdir -p data/ssl/{certs,private}
    
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout data/ssl/private/key.pem \
        -out data/ssl/certs/cert.pem \
        -subj "/C=US/ST=State/L=City/O=Xferant/CN=$DOMAIN"
    
    log_info "SSL certificates generated"
}

start_services() {
    log_info "Starting Xferant VPN services..."
    
    cd $INSTALL_DIR
    
    docker compose up -d
    
    log_info "Waiting for services to start (30 seconds)..."
    sleep 30
    
    if docker compose ps 2>/dev/null | grep -q "Up"; then
        log_info "✅ All services are running"
    else
        log_warn "⚠️ Some services may need attention"
        docker compose logs --tail=50 2>/dev/null || true
    fi
}

finalize_installation() {
    log_info "Finalizing installation..."
    
    cd $INSTALL_DIR
    
    chmod 600 .env 2>/dev/null || true
    chmod 600 data/ssl/private/key.pem 2>/dev/null || true
    
    mkdir -p backups
    
    log_info "Installation finalized"
}

show_success_message() {
    echo ""
    echo -e "${GREEN}🎉 Xferant VPN successfully installed!${NC}"
    echo ""
    echo -e "${CYAN}🔗 Access URLs:${NC}"
    echo -e "   Control Panel: https://$DOMAIN"
    echo -e "   API Server: https://$DOMAIN/api"
    echo -e "   VPN Server Port: 4443"
    echo ""
    echo -e "${YELLOW}🔧 Management Commands:${NC}"
    echo -e "   Start: cd $INSTALL_DIR && docker compose start"
    echo -e "   Stop: cd $INSTALL_DIR && docker compose stop"
    echo -e "   Logs: cd $INSTALL_DIR && docker compose logs -f"
    echo -e "   Status: cd $INSTALL_DIR && docker compose ps"
    echo ""
    echo -e "${BLUE}📚 Next Steps:${NC}"
    echo -e "   1. Configure DNS for $DOMAIN"
    echo -e "   2. Wait 2-3 minutes for full startup"
    echo -e "   3. Access https://$DOMAIN"
    echo -e "   4. Open ports 80, 443, 4443 in firewall"
    echo ""
}

trap 'log_error "Installation failed at line $LINENO"; exit 1' ERR

main "$@"
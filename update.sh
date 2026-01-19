#!/bin/bash
set -e
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'
echo -e "${BLUE}🔄 Xferant VPN Update${NC}"
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}❌ Run: sudo ./update.sh${NC}"
    exit 1
fi
cd /opt/xferant-vpn || exit
echo -e "${GREEN}✅ Starting update...${NC}"
# Backup
mkdir -p backups
BACKUP_DIR="backups/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"
cp .env docker-compose.yml "$BACKUP_DIR/" 2>/dev/null || true
# Git
git pull origin main 2>/dev/null || {
    git fetch origin
    git reset --hard origin/main
}
# Rebuild
docker compose down
docker compose build --no-cache
docker compose up -d
sleep 15
echo -e "${GREEN}✅ Update complete!${NC}"
docker compose ps
echo ""
echo "🌐 Backend:  http://localhost:8080"
echo "🌐 Frontend: http://localhost"
echo "🔐 VPN Port: 4443"
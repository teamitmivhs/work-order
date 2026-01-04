#!/bin/bash

# ============================================
# Work Order System - Docker Setup Helper
# ============================================

set -e

# Colors untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
echo "╔════════════════════════════════════════════╗"
echo "║  Work Order System - Docker Setup Helper  ║"
echo "╚════════════════════════════════════════════╝"
echo -e "${NC}"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker is not installed${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}✗ Docker Compose is not installed${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker & Docker Compose found${NC}"

# ============================================
# Menu Selection
# ============================================

echo -e "\n${YELLOW}Choose deployment option:${NC}"
echo "1) MySQL External (Recommended for Production)"
echo "2) MySQL in Docker with Persistent Storage"
echo "3) Exit"
read -p "Enter choice [1-3]: " choice

case $choice in
    1)
        echo -e "\n${BLUE}Setting up with External MySQL...${NC}"
        COMPOSE_FILE="docker-compose.external-db.yml"
        ;;
    2)
        echo -e "\n${BLUE}Setting up with MySQL in Docker...${NC}"
        COMPOSE_FILE="docker-compose.persistent.yml"
        
        # Create persistent storage directory
        echo -e "${YELLOW}Creating persistent storage directory...${NC}"
        mkdir -p /docker/work-order/mysql-data
        sudo chown 999:999 /docker/work-order/mysql-data 2>/dev/null || echo -e "${YELLOW}Note: May need sudo for permissions${NC}"
        chmod 700 /docker/work-order/mysql-data 2>/dev/null || true
        echo -e "${GREEN}✓ Directory created${NC}"
        ;;
    3)
        echo "Exiting..."
        exit 0
        ;;
    *)
        echo -e "${RED}Invalid choice${NC}"
        exit 1
        ;;
esac

# ============================================
# Environment Setup
# ============================================

cd src

echo -e "\n${YELLOW}Setting up environment variables...${NC}"

if [ ! -f .env ]; then
    echo -e "${YELLOW}.env not found, creating from .env.example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ .env created${NC}"
    echo -e "${YELLOW}Please edit .env to update sensitive values (JWT_SECRET, DB credentials)${NC}"
else
    echo -e "${GREEN}✓ .env already exists${NC}"
fi

# ============================================
# Docker Compose Start
# ============================================

echo -e "\n${YELLOW}Starting services...${NC}"
docker-compose -f "$COMPOSE_FILE" up -d

echo -e "${GREEN}✓ Services starting...${NC}"

# Wait for backend to be ready
echo -e "\n${YELLOW}Waiting for services to be ready...${NC}"
sleep 10

# ============================================
# Health Check
# ============================================

echo -e "\n${BLUE}Health Check:${NC}"

if docker-compose ps | grep -q "work-order-backend.*Up"; then
    echo -e "${GREEN}✓ Backend is running${NC}"
else
    echo -e "${RED}✗ Backend is not running${NC}"
    docker-compose logs backend
    exit 1
fi

if docker-compose ps | grep -q "work-order-nginx.*Up"; then
    echo -e "${GREEN}✓ Nginx is running${NC}"
else
    echo -e "${RED}✗ Nginx is not running${NC}"
    exit 1
fi

if [ "$COMPOSE_FILE" == "docker-compose.persistent.yml" ]; then
    if docker-compose ps | grep -q "work-order-db.*Up"; then
        echo -e "${GREEN}✓ MySQL is running${NC}"
    else
        echo -e "${RED}✗ MySQL is not running${NC}"
        exit 1
    fi
fi

# ============================================
# Summary
# ============================================

echo -e "\n${GREEN}"
echo "╔════════════════════════════════════════════╗"
echo "║  ✓ Deployment Complete!                   ║"
echo "╚════════════════════════════════════════════╝"
echo -e "${NC}"

echo -e "\n${BLUE}Access Points:${NC}"
echo "  Frontend:  http://localhost"
echo "  Backend:   http://localhost/api"
echo "  API Docs:  See API_REFERENCE.md"

echo -e "\n${BLUE}Useful Commands:${NC}"
echo "  View logs:        docker-compose logs -f"
echo "  Stop services:    docker-compose down"
echo "  Restart services: docker-compose restart"

if [ "$COMPOSE_FILE" == "docker-compose.persistent.yml" ]; then
    echo "  Database access:  docker exec -it work-order-db mysql -u adminit2025 -pdatabaseit2045 dbwoit"
fi

echo -e "\n${YELLOW}Next Steps:${NC}"
echo "  1. Update .env with secure JWT_SECRET and database credentials"
echo "  2. Test login at http://localhost/login.html"
echo "  3. Review DEPLOYMENT.md for detailed configuration"

echo -e "\n${YELLOW}Important:${NC}"
echo "  - Change default database password in production"
echo "  - Generate strong JWT_SECRET using: openssl rand -base64 32"
echo "  - Setup HTTPS/SSL for production"
echo "  - Configure backup strategy"

echo ""

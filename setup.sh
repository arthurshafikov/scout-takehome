#!/bin/bash

# Scout Application Setup Script
# Automates environment setup for new clones

set -e  # Exit on any error

echo "🚀 Scout Setup - Automated Configuration"
echo "========================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: docker-compose.yml not found"
    echo "Please run this script from the project root directory"
    exit 1
fi

echo -e "${BLUE}Step 1: Setting up backend environment${NC}"

# Backend setup
if [ ! -f "backend/main.env" ]; then
    if [ -f "backend/main.env.example" ]; then
        cp backend/main.env.example backend/main.env
        echo -e "${GREEN}✓${NC} Created backend/main.env from template"
    else
        echo "❌ Error: backend/main.env.example not found"
        exit 1
    fi
else
    echo -e "${GREEN}✓${NC} backend/main.env already exists"
fi

echo ""
echo -e "${BLUE}Step 2: Setting up frontend environment${NC}"

# Frontend setup
if [ ! -f "frontend/.env" ]; then
    if [ -f "frontend/.env.example" ]; then
        cp frontend/.env.example frontend/.env
        echo -e "${GREEN}✓${NC} Created frontend/.env from template"
    else
        echo "❌ Error: frontend/.env.example not found"
        exit 1
    fi
else
    echo -e "${GREEN}✓${NC} frontend/.env already exists"
fi

echo ""
echo -e "${BLUE}Step 3: Verifying configuration${NC}"

# Verify key files exist
REQUIRED_FILES=(
    "backend/main.env"
    "frontend/.env"
    "docker-compose.yml"
    "dataset/predictions.db"
)

for file in "${REQUIRED_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✓${NC} $file"
    else
        echo "❌ Missing: $file"
        exit 1
    fi
done

echo ""
echo -e "${GREEN}✅ Setup Complete!${NC}"
echo ""
echo "Next steps:"
echo "==========="
echo ""
echo "1. Start the application:"
echo -e "   ${YELLOW}docker compose up${NC}"
echo ""
echo "2. Wait for containers to be healthy (~30 seconds)"
echo ""
echo "3. Open in browser:"
echo -e "   Frontend: ${YELLOW}http://localhost:5173${NC}"
echo -e "   Backend API: ${YELLOW}http://localhost:8080/api${NC}"
echo -e "   MinIO Console: ${YELLOW}http://localhost:9001${NC}"
echo ""
echo "4. Verify backend is running:"
echo -e "   ${YELLOW}curl http://localhost:8080/api/healthz${NC}"
echo ""
echo "5. (Optional) Re-seed dataset:"
echo -e "   ${YELLOW}cd backend && make seed${NC}"
echo ""
echo "For more information, see README.md"

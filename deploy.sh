#!/bin/bash

# Deployment Script for store.adrienbird.net
# Usage: ./deploy.sh [platform]
# Platforms: docker, railway, digitalocean, local

set -e

PLATFORM=${1:-local}
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Store.AdrienBird.net Deployment${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found!${NC}"
    echo "Please copy .env.example to .env and add your credentials."
    exit 1
fi

# Check for required environment variables
echo -e "${YELLOW}Checking environment variables...${NC}"
source .env

MISSING_VARS=()

if [ "$SQUARESPACE_SITE_ID" = "your-site-id-here" ]; then
    MISSING_VARS+=("SQUARESPACE_SITE_ID")
fi

if [ "$SQUARESPACE_API_KEY" = "your-api-key-here" ]; then
    MISSING_VARS+=("SQUARESPACE_API_KEY")
fi

if [ "$SQUARESPACE_ACCESS_TOKEN" = "your-access-token-here" ]; then
    MISSING_VARS+=("SQUARESPACE_ACCESS_TOKEN")
fi

if [ ${#MISSING_VARS[@]} -ne 0 ]; then
    echo -e "${RED}Error: Missing required environment variables:${NC}"
    for var in "${MISSING_VARS[@]}"; do
        echo "  - $var"
    done
    echo ""
    echo "Please edit .env file and add your Squarespace credentials."
    exit 1
fi

echo -e "${GREEN}✓ Environment variables configured${NC}"
echo ""

# Build based on platform
case $PLATFORM in
    docker)
        echo -e "${YELLOW}Deploying with Docker...${NC}"
        docker-compose down
        docker-compose build
        docker-compose up -d
        echo -e "${GREEN}✓ Deployed with Docker${NC}"
        echo "API available at: http://localhost:8080"
        echo "Health check: curl http://localhost:8080/health"
        ;;

    local)
        echo -e "${YELLOW}Building and running locally...${NC}"
        go build -o bin/api cmd/api/main.go
        echo -e "${GREEN}✓ Build complete${NC}"
        echo "Starting server on port 8080..."
        echo "Press Ctrl+C to stop"
        ./bin/api
        ;;

    railway)
        echo -e "${YELLOW}Deploying to Railway...${NC}"
        if ! command -v railway &> /dev/null; then
            echo -e "${RED}Error: railway CLI not installed${NC}"
            echo "Install with: npm install -g @railway/cli"
            exit 1
        fi
        railway login
        railway deploy
        echo -e "${GREEN}✓ Deployed to Railway${NC}"
        ;;

    digitalocean)
        echo -e "${YELLOW}Deploying to DigitalOcean...${NC}"
        if ! command -v doctl &> /dev/null; then
            echo -e "${RED}Error: doctl not installed${NC}"
            echo "Install with: brew install doctl"
            exit 1
        fi
        doctl apps create --spec .do/app.yaml
        echo -e "${GREEN}✓ Deployed to DigitalOcean${NC}"
        ;;

    *)
        echo -e "${RED}Error: Unknown platform '$PLATFORM'${NC}"
        echo "Available platforms: docker, local, railway, digitalocean"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Deployment Complete!${NC}"
echo -e "${GREEN}========================================${NC}"

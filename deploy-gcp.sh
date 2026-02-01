#!/bin/bash

# Zero-Cost Deployment Script for store.adrienbird.net
# Deploys to existing GCP e2-micro instance (taskflow-server)
# Uses Docker + Nginx reverse proxy

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Zero-Cost GCP Deployment${NC}"
echo -e "${BLUE}  store.adrienbird.net${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Configuration
GCP_PROJECT="taskflow-free-tier"
GCP_ZONE="us-central1-a"
INSTANCE_NAME="taskflow-server"
CONTAINER_NAME="store-api"
CONTAINER_PORT=8080
IMAGE_NAME="store-api"

echo -e "${YELLOW}Step 1: Building Docker image...${NC}"
docker build -t ${IMAGE_NAME} .

echo -e "${GREEN}✓ Build complete${NC}"
echo ""

echo -e "${YELLOW}Step 2: Saving Docker image...${NC}"
docker save ${IMAGE_NAME} | gzip > /tmp/${IMAGE_NAME}.tar.gz

echo -e "${GREEN}✓ Image saved${NC}"
echo ""

echo -e "${YELLOW}Step 3: Copying image to GCP instance...${NC}"
gcloud compute scp /tmp/${IMAGE_NAME}.tar.gz ${INSTANCE_NAME}:/tmp/ --project=${GCP_PROJECT} --zone=${GCP_ZONE}

echo -e "${GREEN}✓ Image copied${NC}"
echo ""

echo -e "${YELLOW}Step 4: Installing Docker on GCP instance (if needed)...${NC}"
gcloud compute ssh ${INSTANCE_NAME} --project=${GCP_PROJECT} --zone=${GCP_ZONE} --command='
  if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    echo "Docker installed successfully"
  else
    echo "Docker already installed"
  fi
'

echo -e "${GREEN}✓ Docker ready${NC}"
echo ""

echo -e "${YELLOW}Step 5: Loading Docker image on GCP...${NC}"
gcloud compute ssh ${INSTANCE_NAME} --project=${GCP_PROJECT} --zone=${GCP_ZONE} --command="
  docker load < /tmp/${IMAGE_NAME}.tar.gz
  rm /tmp/${IMAGE_NAME}.tar.gz
"

echo -e "${GREEN}✓ Image loaded${NC}"
echo ""

echo -e "${YELLOW}Step 6: Checking for existing container...${NC}"
EXISTING=$(gcloud compute ssh ${INSTANCE_NAME} --project=${GCP_PROJECT} --zone=${GCP_ZONE} --command="docker ps -q -f name=${CONTAINER_NAME}" 2>/dev/null || echo "")

if [ -n "$EXISTING" ]; then
  echo "Stopping existing container..."
  gcloud compute ssh ${INSTANCE_NAME} --project=${GCP_PROJECT} --zone=${GCP_ZONE} --command="docker stop ${CONTAINER_NAME} && docker rm ${CONTAINER_NAME}"
fi

echo -e "${GREEN}✓ Ready for deployment${NC}"
echo ""

echo -e "${YELLOW}Step 7: Starting container...${NC}"
gcloud compute ssh ${INSTANCE_NAME} --project=${GCP_PROJECT} --zone=${GCP_ZONE} --command="
  docker run -d \
    --name ${CONTAINER_NAME} \
    --restart unless-stopped \
    -p ${CONTAINER_PORT}:8080 \
    -e PORT=8080 \
    -e GIN_MODE=release \
    -e SQUARESPACE_SITE_ID=\${SQUARESPACE_SITE_ID} \
    -e SQUARESPACE_API_KEY=\${SQUARESPACE_API_KEY} \
    -e SQUARESPACE_ACCESS_TOKEN=\${SQUARESPACE_ACCESS_TOKEN} \
    -e SQUARESPACE_BASE_URL=\${SQUARESPACE_BASE_URL:-https://api.squarespace.com} \
    -e NODE_ENV=production \
    ${IMAGE_NAME}
"

echo -e "${GREEN}✓ Container started${NC}"
echo ""

echo -e "${YELLOW}Step 8: Installing Nginx (if needed)...${NC}"
gcloud compute ssh ${INSTANCE_NAME} --project=${GCP_PROJECT} --zone=${GCP_ZONE} --command='
  if ! command -v nginx &> /dev/null; then
    echo "Installing Nginx..."
    sudo apt-get update
    sudo apt-get install -y nginx
    echo "Nginx installed"
  else
    echo "Nginx already installed"
  fi
'

echo -e "${GREEN}✓ Nginx ready${NC}"
echo ""

echo -e "${YELLOW}Step 9: Configuring Nginx reverse proxy...${NC}"
cat > /tmp/store-api-nginx.conf <<'EOF'
server {
    listen 80;
    server_name store.adrienbird.net;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # CORS headers
        add_header Access-Control-Allow-Origin https://adrienbird.net always;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
        add_header Access-Control-Allow-Headers "Content-Type, Authorization" always;

        if ($request_method = OPTIONS) {
            return 204;
        }
    }

    # Health check endpoint
    location /health {
        proxy_pass http://localhost:8080/health;
        access_log off;
    }
}
EOF

gcloud compute scp /tmp/store-api-nginx.conf ${INSTANCE_NAME}:/tmp/store-api.conf --project=${GCP_PROJECT} --zone=${GCP_ZONE}

gcloud compute ssh ${INSTANCE_NAME} --project=${GCP_PROJECT} --zone=${GCP_ZONE} --command="
  sudo mv /tmp/store-api.conf /etc/nginx/sites-available/store-api.conf
  sudo ln -sf /etc/nginx/sites-available/store-api.conf /etc/nginx/sites-enabled/store-api.conf
  sudo nginx -t && sudo systemctl reload nginx
"

echo -e "${GREEN}✓ Nginx configured${NC}"
echo ""

echo -e "${YELLOW}Step 10: Testing deployment...${NC}"
sleep 5

gcloud compute ssh ${INSTANCE_NAME} --project=${GCP_PROJECT} --zone=${GCP_ZONE} --command="docker ps --filter name=${CONTAINER_NAME}"
gcloud compute ssh ${INSTANCE_NAME} --project=${GCP_PROJECT} --zone=${GCP_ZONE} --command="curl -s http://localhost:${CONTAINER_PORT}/health || echo 'Health check failed'"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Deployment Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}Next Steps:${NC}"
echo "1. Update DNS: Point store.adrienbird.net to ${EXTERNAL_IP}"
echo "2. Test API: curl http://store.adrienbird.net/health"
echo "3. Check logs: docker logs -f ${CONTAINER_NAME} (on GCP instance)"
echo ""
echo -e "${YELLOW}Note:${NC} Container will auto-start on instance reboot"
echo ""

# Cleanup
rm -f /tmp/${IMAGE_NAME}.tar.gz
rm -f /tmp/store-api-nginx.conf

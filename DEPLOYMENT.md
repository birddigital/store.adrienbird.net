# Deployment Guide - Store.AdrienBird.net

## Overview

This is a Go-based API backend for the `store.adrienbird.net` subdomain that integrates with Squarespace's Commerce APIs to provide a headless e-commerce experience.

## Architecture

```
store.adrienbird.net (Custom Frontend)
       ↓ API calls
Go Backend API (Port 8080)
       ↓ API calls
Squarespace Commerce APIs (Backend)
```

## Environment Setup

### 1. Environment Variables

Create a `.env` file with the following configuration:

```bash
# Server Configuration
PORT=8080
GIN_MODE=release                    # Use release for production
ENABLE_SWAGGER=false               # Disable swagger in production
ENABLE_HEALTH=true
ENABLE_METRICS=false
ENABLE_TRACING=false

# Squarespace API Configuration
SQUARESPACE_BASE_URL=https://api.squarespace.com
SQUARESPACE_SITE_ID=your-site-id-here
SQUARESPACE_API_KEY=your-api-key-here
SQUARESPACE_ACCESS_TOKEN=your-access-token-here

# Environment
NODE_ENV=production
```

### 2. Squarespace API Setup

**Step 1: Get API Credentials**
1. Go to your Squarespace admin panel
2. Navigate to Settings → Advanced → External API Keys
3. Create a new API key with permissions for:
   - Products (read)
   - Orders (read/write)
   - Inventory (read/write)
   - Profiles (read)

**Step 2: Get Your Site ID**
1. In Squarespace admin, go to Settings → Basic Information
2. Copy the "Site ID" (starts with your-site-name-xxxxx)

## Local Development

### Prerequisites
- Go 1.21+
- Git

### Setup
```bash
# Clone the repository
git clone https://github.com/birddigital/store.adrienbird.net.git
cd store.adrienbird.net

# Copy environment template
cp .env.example .env
# Edit .env with your actual credentials

# Install dependencies
go mod tidy

# Build the application
go build -o bin/api cmd/api/main.go

# Run the server
./bin/api
```

The API will be available at `http://localhost:8080`

### Development Commands
```bash
# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o bin/api-linux cmd/api/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/api-macos cmd/api/main.go
GOOS=windows GOARCH=amd64 go build -o bin/api-windows.exe cmd/api/main.go

# Run tests
go test ./...

# Run with live reload (for development)
go install github.com/air-verse/air@latest
air
```

## API Endpoints

### Health & Status
- `GET /health` - Health check with Squarespace connectivity
- `GET /` - Basic status information

### Products
- `GET /api/v1/products` - List products (with pagination, filtering)
- `GET /api/v1/products/{id}` - Get single product
- `GET /api/v1/products/{id}/variants` - Get product variants

**Query Parameters for Products:**
- `limit` - Number of products to return (default: 20)
- `offset` - Pagination offset (default: 0)
- `category` - Filter by category
- `tag` - Filter by tag

### Orders
- `GET /api/v1/orders` - List orders
- `GET /api/v1/orders/{id}` - Get single order
- `POST /api/v1/orders` - Create new order

**Query Parameters for Orders:**
- `limit` - Number of orders to return (default: 20)
- `offset` - Pagination offset (default: 0)
- `status` - Filter by order status
- `customerId` - Filter by customer ID

## Production Deployment

### Option 1: DigitalOcean App Platform

1. **Create App**
   ```bash
   # Install doctl
   brew install doctl

   # Authenticate
   doctl auth init

   # Create app
   doctl apps create --spec .do/app.yaml
   ```

2. **App Spec (`.do/app.yaml`)**
   ```yaml
   name: store-adrienbird-api
   services:
   - name: api
     source_dir: .
     github:
       repo: birddigital/store.adrienbird.net
       branch: main
     run_command: ./bin/api
     http_port: 8080
     instance_count: 1
     instance_size_slug: basic-xxs
     envs:
     - key: PORT
       value: "8080"
     - key: GIN_MODE
       value: "release"
     - key: SQUARESPACE_SITE_ID
       value: ${SQUARESPACE_SITE_ID}
     - key: SQUARESPACE_API_KEY
       value: ${SQUARESPACE_API_KEY}
   ```

### Option 2: Docker Deployment

1. **Dockerfile**
   ```dockerfile
   FROM golang:1.21-alpine AS builder
   WORKDIR /app
   COPY go.mod go.sum ./
   RUN go mod download
   COPY . .
   RUN go build -o main cmd/api/main.go

   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   WORKDIR /root/
   COPY --from=builder /app/main .
   COPY --from=builder /app/.env .env
   EXPOSE 8080
   CMD ["./main"]
   ```

2. **Deploy**
   ```bash
   # Build image
   docker build -t store-api .

   # Run container
   docker run -p 8080:8080 --env-file .env store-api
   ```

### Option 3: Traditional VPS

1. **Server Setup (Ubuntu/Debian)**
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y

   # Install Go
   wget https://golang.org/dl/go1.21.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc

   # Clone and build
   cd /var/www
   git clone https://github.com/birddigital/store.adrienbird.net.git
   cd store.adrienbird.net
   go mod tidy
   go build -o api cmd/api/main.go
   ```

2. **Systemd Service (`/etc/systemd/system/store-api.service`)**
   ```ini
   [Unit]
   Description=Store.AdrienBird.net API
   After=network.target

   [Service]
   Type=simple
   User=www-data
   WorkingDirectory=/var/www/store.adrienbird.net
   ExecStart=/var/www/store.adrienbird.net/api
   Restart=always
   RestartSec=5
   Environment=GIN_MODE=release
   EnvironmentFile=/var/www/store.adrienbird.net/.env

   [Install]
   WantedBy=multi-user.target
   ```

   ```bash
   # Enable and start service
   sudo systemctl enable store-api
   sudo systemctl start store-api
   sudo systemctl status store-api
   ```

## DNS Configuration

Set up your DNS records to point `store.adrienbird.net` to your deployment:

```bash
# A record for root domain
store.adrienbird.net.    A    <YOUR_SERVER_IP>

# Optional: CNAME for www
www.store.adrienbird.net.    CNAME    store.adrienbird.net.
```

## SSL/TLS Setup

### Using Let's Encrypt (Certbot)

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Get SSL certificate
sudo certbot --nginx -d store.adrienbird.net -d www.store.adrienbird.net

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

## Monitoring & Logging

### Application Monitoring
- Health endpoint: `https://store.adrienbird.net/health`
- Monitor Squarespace API connectivity
- Track response times and error rates

### Log Management
```bash
# View application logs
sudo journalctl -u store-api -f

# Log rotation
sudo nano /etc/logrotate.d/store-api
```

## Security Considerations

1. **Environment Variables**: Never commit `.env` files to version control
2. **API Keys**: Rotate Squarespace API keys regularly
3. **Rate Limiting**: Consider implementing rate limiting for production
4. **CORS**: Configure CORS to allow only your frontend domains
5. **HTTPS**: Always use HTTPS in production

## Performance Optimization

1. **Caching**: Implement Redis/Memcached for frequently accessed products
2. **Database**: Consider PostgreSQL for order history and customer data
3. **CDN**: Use CloudFlare for static assets and DDoS protection
4. **Load Balancing**: Use multiple instances for high availability

## Troubleshooting

### Common Issues

1. **401 Unauthorized**
   - Check Squarespace API credentials
   - Verify site ID is correct
   - Ensure API key has proper permissions

2. **Connection Timeout**
   - Check network connectivity to api.squarespace.com
   - Verify firewall settings
   - Increase timeout in HTTP client

3. **Build Failures**
   - Run `go mod tidy` to update dependencies
   - Check Go version compatibility
   - Verify all imports are correct

### Health Check Responses

- **healthy**: All systems operational
- **degraded**: Configuration issues, API partially working
- **unhealthy**: Critical failures, API not responding

## Support & Maintenance

- **Documentation**: Keep API documentation updated
- **Monitoring**: Set up alerts for API failures
- **Backups**: Regular backups of any local data
- **Updates**: Regular dependency updates and security patches

## Frontend Integration

Your frontend (React/Vue/Next.js) should:

1. **Make API calls to** `https://store.adrienbird.net/api/v1/`
2. **Handle authentication** via headers (if required)
3. **Implement proper error handling** for API responses
4. **Use pagination** for product listings
5. **Cache product data** appropriately
6. **Implement retry logic** for failed requests

Example API call from frontend:
```javascript
const response = await fetch('https://store.adrienbird.net/api/v1/products?limit=20');
const data = await response.json();
```

---

This API backend provides a solid foundation for your headless Squarespace commerce integration. The modular design makes it easy to extend with additional features like payment processing, inventory management, or customer accounts.
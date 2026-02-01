# Zero-Cost Deployment Guide - store.adrienbird.net

**Deployment Method**: Docker on existing GCP e2-micro instance
**Cost**: $0/month (within GCP Free Tier)
**Time**: ~10 minutes

---

## ✅ Prerequisites

### Required Credentials (BLOCKER)

The API **cannot function** without Squarespace credentials. You need:

```bash
SQUARESPACE_SITE_ID=your-site-id-here
SQUARESPACE_API_KEY=your-api-key-here
SQUARESPACE_ACCESS_TOKEN=your-access-token-here
```

**How to Get:**
1. Login to [Squarespace](https://www.squarespace.com)
2. Settings → Advanced → Developer API Keys
3. Create API key with permissions:
   - ✅ Products (read)
   - ✅ Orders (read/write)
   - ✅ Inventory (read/write)
   - ✅ Profiles (read)

---

## 🚀 Deployment Steps

### Option A: Automated Deployment (Recommended)

**After you have credentials:**

```bash
cd ~/sources/store.adrienbird.net

# 1. Add credentials to .env file
nano .env
# Add your Squarespace credentials

# 2. Run deployment script
./deploy-gcp.sh
```

**The script will:**
1. Build Docker image locally
2. Copy image to GCP instance (taskflow-server)
3. Install Docker on GCP (if needed)
4. Load and start container
5. Configure Nginx reverse proxy
6. Set up subdomain routing

### Option B: Manual Deployment

If you prefer manual control:

```bash
# 1. Build image
docker build -t store-api .

# 2. Save image
docker save store-api | gzip > store-api.tar.gz

# 3. Copy to GCP
gcloud compute scp store-api.tar.gz taskflow-server:/tmp/ \
  --project=taskflow-free-tier --zone=us-central1-a

# 4. SSH into instance
gcloud compute ssh taskflow-server \
  --project=taskflow-free-tier --zone=us-central1-a

# 5. Load and run (on instance)
docker load < store-api.tar.gz
docker run -d --name store-api --restart unless-stopped \
  -p 8080:8080 \
  -e SQUARESPACE_SITE_ID=... \
  -e SQUARESPACE_API_KEY=... \
  -e SQUARESPACE_ACCESS_TOKEN=... \
  store-api
```

---

## 🌐 DNS Configuration

After deployment, update DNS for `store.adrienbird.net`:

**Current DNS**: Points to Cloudflare (172.67.137.147, 104.21.64.230)

**Options:**

### Option 1: Use Cloudflare Tunnel (Recommended)

**Pros:**
- ✅ Free SSL certificate
- ✅ DDoS protection
- ✅ Hides GCP IP
- ✅ Easy setup

**Steps:**
1. Install `cloudflared` on GCP instance
2. Create tunnel in Cloudflare dashboard
3. Route `store.adrienbird.net` → `http://localhost:8080`

### Option 2: Direct A Record

**GCP Instance IP**: `136.119.217.179`

**Steps:**
1. Go to Cloudflare DNS settings
2. Add A record: `store.adrienbird.net` → `136.119.217.179`
3. Enable proxy (orange cloud icon)
4. Cloudflare will handle SSL automatically

### Option 3: CNAME to Existing Domain

If you want to keep the current setup, you can configure Cloudflare Workers to route `store.adrienbird.net` requests to the GCP instance.

---

## 📊 Deployment Architecture

```
store.adrienbird.net (DNS)
        ↓
    Cloudflare (Optional - SSL/DDoS)
        ↓
    GCP Firewall (allow HTTP/HTTPS)
        ↓
    taskflow-server (e2-micro, FREE TIER)
        ↓
    Nginx Reverse Proxy
        ↓
    Docker Container (store-api:8080)
        ↓
    Squarespace Commerce API
```

---

## ✅ Verification

After deployment, test:

```bash
# Health check (from your local machine)
curl http://store.adrienbird.net/health

# Expected response:
# {"status":"healthy","timestamp":"..."}

# Products endpoint
curl http://store.adrienbird.net/api/v1/products

# Check container status (SSH into GCP)
docker ps

# View logs
docker logs -f store-api

# Check Nginx status
sudo systemctl status nginx
```

---

## 🔧 Management Commands

### View Logs
```bash
# SSH into GCP instance
gcloud compute ssh taskflow-server \
  --project=taskflow-free-tier --zone=us-central1-a

# View container logs
docker logs -f store-api

# View Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Restart Container
```bash
# On GCP instance
docker restart store-api

# Or rebuild and redeploy
./deploy-gcp.sh
```

### Update Container
```bash
# 1. Make code changes
# 2. Rebuild and redeploy
./deploy-gcp.sh
```

### Stop Container
```bash
# On GCP instance
docker stop store-api
docker rm store-api
```

---

## 💰 Cost Analysis

**Monthly Cost: $0**

| Resource | Free Tier Allowance | Current Usage | Cost |
|----------|---------------------|---------------|------|
| e2-micro VM | 1 instance/month | 1 instance | $0 |
| Ephemeral IP | 1 IP/instance | 1 IP | $0 |
| Network Egress | 100 GB/month | <1 GB (estimated) | $0 |
| Docker | Free | 1 container | $0 |
| Nginx | Free | 1 reverse proxy | $0 |

**Total**: $0/month ✅

---

## 🚨 Troubleshooting

### Container Won't Start
```bash
# Check logs
docker logs store-api

# Common issues:
# - Missing Squarespace credentials
# - Port 8080 already in use
# - Invalid environment variables
```

### Nginx 502 Bad Gateway
```bash
# Check if container is running
docker ps | grep store-api

# Check if container is listening on port 8080
docker exec store-api netstat -tulpn

# Restart Nginx
sudo systemctl restart nginx
```

### DNS Not Resolving
```bash
# Check DNS propagation
dig store.adrienbird.net

# Should point to Cloudflare or GCP IP
```

### Health Check Fails
```bash
# Check Squarespace credentials
docker exec store-api env | grep SQUARESPACE

# Test Squarespace connectivity
docker exec store-api curl -I https://api.squarespace.com
```

---

## 📝 Next Steps

1. **Get Squarespace credentials** ← BLOCKER
2. **Add credentials to `.env` file**
3. **Run `./deploy-gcp.sh`**
4. **Update DNS (Cloudflare)**
5. **Test API endpoints**
6. **Monitor logs and performance**

---

## 🎯 Success Criteria

- ✅ Container running on GCP instance
- ✅ Nginx reverse proxy configured
- ✅ Health check returns `{"status":"healthy"}`
- ✅ Products endpoint returns Squarespace data
- ✅ DNS resolves correctly
- ✅ Zero monthly cost

---

**Ready to deploy?** Just add your Squarespace credentials to `.env` and run:

```bash
./deploy-gcp.sh
```

**Questions?** See `DEPLOYMENT_READINESS.md` for comprehensive details.

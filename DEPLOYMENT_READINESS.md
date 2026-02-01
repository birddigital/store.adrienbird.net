# Deployment Readiness Checklist - store.adrienbird.net

**Last Updated**: 2025-02-01
**Status**: 🟡 Ready for Deployment (with blockers)

---

## ✅ Ready Items

### Code & Build
- [x] Go code compiles successfully
- [x] All dependencies resolved
- [x] Binary builds for target platform (arm64)
- [x] Health check endpoint implemented
- [x] API routes properly configured
- [x] CORS configured for adrienbird.net
- [x] Error handling in place
- [x] Logging configured

### Infrastructure
- [x] Dockerfile created (multi-stage build)
- [x] docker-compose.yml created
- [x] Health check implemented
- [x] Non-root user configured
- [x] Resource limits defined
- [x] Restart policy configured

### API Endpoints
- [x] `GET /health` - Health check
- [x] `GET /` - API status
- [x] `GET /api/v1/products` - List products
- [x] `GET /api/v1/products/:id` - Get product
- [x] `GET /api/v1/products/:id/variants` - Get variants
- [x] `GET /api/v1/orders` - List orders
- [x] `GET /api/v1/orders/:id` - Get order
- [x] `POST /api/v1/orders` - Create order

### Documentation
- [x] README.md present
- [x] DEPLOYMENT.md comprehensive
- [x] Environment variables documented
- [x] API endpoints documented

---

## 🚫 Blockers

### Critical - Must Fix Before Deployment

#### 1. Missing Squarespace Credentials
**Status**: 🔴 Critical Blocker
**Impact**: API cannot communicate with Squarespace backend

**Required**:
- [ ] `SQUARESPACE_SITE_ID` - Get from Squarespace Settings → Basic Information
- [ ] `SQUARESPACE_API_KEY` - Generate in Squarespace Settings → Advanced → External API Keys
- [ ] `SQUARESPACE_ACCESS_TOKEN` - Generate OAuth token or use API key

**Steps to Fix**:
1. Login to Squarespace admin panel
2. Navigate to Settings → Advanced → Developer API Keys
3. Create new API key with permissions:
   - Products (read)
   - Orders (read/write)
   - Inventory (read/write)
   - Profiles (read)
4. Copy credentials to `.env` file or deployment environment variables

#### 2. Production Environment Not Configured
**Status**: 🟡 Medium Blocker
**Impact**: No deployment target configured

**Required Decisions**:
- [ ] Choose deployment platform:
  - [ ] DigitalOcean App Platform
  - [ ] AWS ECS
  - [ ] Google Cloud Run
  - [ ] VPS with Docker
  - [ ] Railway/Render
- [ ] Set up DNS for `store.adrienbird.net`
- [ ] Configure SSL/TLS certificate
- [ ] Set up monitoring/alerting

---

## ⚠️ Warnings - Recommended Before Deployment

### Security
- [ ] Review CORS settings (currently allows only adrienbird.net)
- [ ] Add rate limiting middleware
- [ ] Implement authentication for admin endpoints
- [ ] Add API versioning strategy
- [ ] Set up API key rotation schedule
- [ ] Configure webhook signatures for order updates

### Performance
- [ ] Implement Redis caching for products
- [ ] Add request timeout configurations
- [ ] Set up CDN for static assets
- [ ] Configure database connection pooling
- [ ] Add pagination to all list endpoints

### Reliability
- [ ] Add circuit breaker for Squarespace API
- [ ] Implement retry logic with exponential backoff
- [ ] Add request ID tracking for debugging
- [ ] Set up log aggregation (e.g., LogDNA, Papertrail)
- [ ] Configure error tracking (e.g., Sentry)

### Monitoring
- [ ] Set up uptime monitoring (e.g., Pingdom, UptimeRobot)
- [ ] Configure performance monitoring (e.g., Datadog, New Relic)
- [ ] Add alerting for API failures
- [ ] Create dashboard for metrics
- [ ] Set up cost tracking for Squarespace API calls

---

## 📋 Deployment Options

### Option 1: DigitalOcean App Platform (Recommended)

**Pros**:
- Easy setup
- Auto SSL
- Built-in monitoring
- Scales automatically

**Steps**:
1. Install doctl: `brew install doctl`
2. Authenticate: `doctl auth init`
3. Create app: `doctl apps create --spec .do/app.yaml`
4. Add environment variables in dashboard
5. Deploy!

**Estimated Cost**: $5-20/month

**Time**: ~30 minutes

### Option 2: Docker on VPS

**Pros**:
- Full control
- Cheaper for small scale
- Can run multiple services

**Steps**:
1. Provision VPS (DigitalOcean, Linode, etc.)
2. Install Docker & docker-compose
3. Copy files to server
4. Configure environment variables
5. Run: `docker-compose up -d`
6. Set up Nginx reverse proxy
7. Configure SSL with Let's Encrypt

**Estimated Cost**: $5-10/month

**Time**: ~2 hours

### Option 3: Railway/Render

**Pros**:
- Very simple setup
- GitHub integration
- Auto-deploys on push

**Steps**:
1. Connect GitHub repository
2. Configure environment variables
3. Deploy!

**Estimated Cost**: $5-20/month

**Time**: ~15 minutes

---

## 🚀 Quick Start Deployment

If you have credentials ready, here's the fastest path:

```bash
# 1. Update environment variables
cd ~/sources/store.adrienbird.net
nano .env  # Add your Squarespace credentials

# 2. Test locally with Docker
docker-compose up --build
curl http://localhost:8080/health

# 3. If health check passes, deploy to your chosen platform
# For Railway:
# - Connect repo at railway.app
# - Add environment variables
# - Deploy!

# For DigitalOcean:
doctl apps create --spec .do/app.yaml
```

---

## 📊 Post-Deployment Checklist

Once deployed, verify:

- [ ] Health endpoint returns `{"status": "healthy"}`
- [ ] Products endpoint returns data from Squarespace
- [ ] CORS headers allow requests from adrienbird.net
- [ ] SSL certificate is valid
- [ ] DNS resolves correctly
- [ ] Monitoring is collecting data
- [ ] Logs are accessible
- [ ] Error tracking is working

---

## 🔧 Troubleshooting

### Common Issues

**1. Health check fails**
- Check Squarespace credentials
- Verify network connectivity
- Check application logs

**2. CORS errors**
- Update allowed origins in middleware
- Verify frontend is making requests from correct domain

**3. 502/504 errors**
- Check Squarespace API status
- Verify timeout configurations
- Check rate limits

**4. High memory usage**
- Adjust resource limits in docker-compose
- Check for memory leaks in code
- Profile the application

---

## 📞 Support Resources

- **Squarespace API Docs**: https://developers.squarespace.com/
- **Gin Framework**: https://gin-gonic.com/docs/
- **Go Deployment**: https://golang.org/doc/deploy

---

## Next Steps

1. **Immediate**: Get Squarespace credentials
2. **Short-term**: Choose deployment platform
3. **Medium-term**: Set up monitoring and alerting
4. **Long-term**: Implement caching and performance optimization

**Current Recommendation**: Start with Railway or Render for fastest deployment, then migrate to DigitalOcean App Platform or AWS ECS for production scale.

# Deployment Guide

This document provides comprehensive information about deploying the APP system log management dashboard in production environments.

## 🚀 Overview

The APP system log management dashboard can be deployed in various environments including development, staging, and production. This guide covers containerized deployment using Docker and Docker Compose, as well as manual deployment options.

## 🛠️ Prerequisites

### Environment Requirements

- **Docker 20+** - Containerization platform
- **Docker Compose 2.0+** - Multi-container orchestration
- **SSL Certificates** (for production)
- **Domain Name** (for production)
- **Reverse Proxy** (nginx/Apache) (for production)

### Production Environment

```env
# Production environment variables
PORT=8080
QUICKWIT_URL=http://quickwit:7280
JWT_SECRET=your-super-secret-jwt-key-should-be-very-long
JWT_EXPIRES_IN=24h
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

## 📦 Deployment Options

### 1. Docker Compose (Recommended)

#### Production Setup

```yaml
# docker-compose.yml
version: '3.8'

services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - QUICKWIT_URL=http://quickwit:7280
      - JWT_SECRET=${JWT_SECRET}
      - JWT_EXPIRES_IN=24h
      - ALLOWED_ORIGINS=${ALLOWED_ORIGINS}
    depends_on:
      - quickwit
    restart: unless-stopped
    networks:
      - app-network

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - VITE_API_URL=http://localhost:8080/api/v1
    restart: unless-stopped
    networks:
      - app-network

  quickwit:
    image: quickwit/quickwit:latest
    ports:
      - "7280:7280"
    volumes:
      - ./quickwit-data:/data
    command: run --config /config/quickwit.yaml --data-dir /data
    restart: unless-stopped
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - backend
      - frontend
    restart: unless-stopped
    networks:
      - app-network

networks:
  app-network:
    driver: bridge

volumes:
  quickwit-data:
```

### 2. Manual Deployment

#### Backend Deployment

```bash
# Build backend
cd backend
go build -o app-backend main.go

# Run backend
./app-backend
```

#### Frontend Deployment

```bash
# Build frontend
cd frontend
npm run build

# Serve frontend (using nginx, apache, or similar)
```

## 🛡️ Security Configuration

### SSL/TLS Setup

```nginx
# nginx.conf
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/nginx/ssl/certificate.crt;
    ssl_certificate_key /etc/nginx/ssl/private.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Environment Variables

Create `.env` file for production:

```env
# Server Configuration
PORT=8080

# QuickWit Configuration
QUICKWIT_URL=http://quickwit:7280

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-should-be-very-long-and-random
JWT_EXPIRES_IN=24h

# CORS Configuration
ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com,https://app.yourdomain.com

# Database Configuration (if using database)
DB_HOST=localhost
DB_PORT=5432
DB_NAME=app_logs
DB_USER=app_user
DB_PASSWORD=secure_password
```

## 🚀 Deployment Steps

### 1. Prepare Environment

```bash
# Clone repository
git clone <repository-url>
cd app

# Create environment file
cp .env.example .env
# Edit .env with production values
```

### 2. Build Images

```bash
# Build Docker images
docker-compose build

# Or build specific services
docker-compose build backend
docker-compose build frontend
```

### 3. Deploy Services

```bash
# Start all services
docker-compose up -d

# Start specific services
docker-compose up -d backend
docker-compose up -d frontend

# View logs
docker-compose logs -f
```

### 4. Verify Deployment

```bash
# Check service status
docker-compose ps

# Test API endpoints
curl http://localhost:8080/api/v1/health

# Test frontend
curl http://localhost:3000
```

## 📊 Monitoring and Logging

### Health Checks

```yaml
# Health check endpoints
health:
  - GET /health
  - GET /api/v1/health
  - GET /api/v1/health/database
```

### Monitoring Setup

```yaml
# Add monitoring services
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3001:3000"
    depends_on:
      - prometheus
```

### Logging Configuration

```yaml
# Docker logging configuration
services:
  backend:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 🔄 Scaling and Load Balancing

### Horizontal Scaling

```yaml
# Scale backend services
docker-compose up -d --scale backend=3

# Scale frontend services
docker-compose up -d --scale frontend=2
```

### Load Balancer Setup

```nginx
upstream app_backend {
    server backend:8080;
    server backend:8080;
    server backend:8080;
}

server {
    listen 80;
    location / {
        proxy_pass http://app_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 🛡️ Security Hardening

### Production Security Measures

1. **Network Security**
   - Use Docker networks
   - Restrict container ports
   - Implement firewall rules

2. **Access Control**
   - Use strong JWT secrets
   - Implement rate limiting
   - Configure proper CORS

3. **Data Protection**
   - Enable HTTPS
   - Secure environment variables
   - Regular security updates

### Security Configuration

```yaml
# Security hardening in docker-compose
services:
  backend:
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
    user: 1000:1000
    cap_drop:
      - ALL
    cap_add:
      - CHOWN
      - SETGID
      - SETUID
```

## 📋 Backup and Recovery

### Data Backup Strategy

```bash
# Backup QuickWit data
docker run --rm \
  -v quickwit-data:/data \
  -v $(pwd)/backups:/backups \
  alpine tar czf /backups/quickwit-backup-$(date +%Y%m%d).tar.gz -C /data .

# Backup configuration
docker run --rm \
  -v app_env:/env \
  -v $(pwd)/backups:/backups \
  alpine tar czf /backups/config-backup-$(date +%Y%m%d).tar.gz -C /env .
```

### Recovery Process

```bash
# Restore QuickWit data
docker run --rm \
  -v quickwit-data:/data \
  -v $(pwd)/backups:/backups \
  alpine sh -c "cd /data && tar xzf /backups/quickwit-backup-*.tar.gz"
```

## 📈 Performance Optimization

### Backend Optimization

```go
// Connection pooling
db, err := sql.Open("postgres", connStr)
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

### Frontend Optimization

```javascript
// Build optimization
export default {
  build: {
    analyze: true,
    transpile: ['@astrojs/tailwind'],
    minify: true,
  }
}
```

## 📋 Deployment Best Practices

### Environment Management

1. **Use separate environments**
   - Development
   - Staging
   - Production

2. **Version control configuration**
   - Store `.env` files in version control (with placeholders)
   - Never store real secrets in version control

3. **Automated deployments**
   - Use CI/CD pipelines
   - Implement rollback strategies
   - Monitor deployment status

### Rollback Strategy

```bash
# Create backup before deployment
docker-compose stop
docker-compose rm -f
docker-compose up -d

# If deployment fails, rollback
docker-compose down
docker-compose up -d --force-recreate
```

## 📚 Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [QuickWit Deployment Guide](https://quickwit.io/docs/deployment/)
- [Go Production Deployment](https://golang.org/doc/effective_go.html#production)
# 🚀 Sauron-Pro Simple Deployment

**ONE clear way to deploy Sauron that actually works.**

## 📋 Prerequisites

- Ubuntu 20.04+ VPS with root access
- Domain name (e.g., `microsoftlogin365.com`)
- Cloudflare account (free)

## 🎯 Method 1: Source Build (Recommended)

**This is the original working method - builds Go directly on VPS**

### Step 1: Clone and Configure

```bash
# Clone repository
git clone https://github.com/Skillz147/Sauron-Pro.git
cd Sauron-Pro

# Configure environment
cp .env.example .env
nano .env  # Fill in your SAURON_DOMAIN, CLOUDFLARE_API_TOKEN, etc.
```

### Step 2: Install Everything

```bash
# Run the install script (this installs Go, builds binary, sets up systemd)
sudo bash install/setup.sh
```

**That's it!** Sauron is now running.

### Management Commands

```bash
# Check status
sudo systemctl status sauron

# View logs
sudo journalctl -u sauron -f

# Restart
sudo systemctl restart sauron

# Stop
sudo systemctl stop sauron
```

---

## 🐳 Method 2: Docker (Alternative)

**Only use this if you prefer Docker over the source build**

### Step 1: Install Docker

```bash
curl -fsSL https://get.docker.com | sh
sudo systemctl enable docker
sudo systemctl start docker
```

### Step 2: Deploy

```bash
# Clone and configure
git clone https://github.com/Skillz147/Sauron-Pro.git
cd Sauron-Pro
cp .env.example .env
nano .env  # Configure your settings

# Build and start
./scripts/build-docker-release.sh v2.0.0-pro
docker-compose -f docker-compose.pro.yml up -d
```

### Docker Management

```bash
# Check status
docker-compose -f docker-compose.pro.yml ps

# View logs
docker-compose -f docker-compose.pro.yml logs -f sauron-pro

# Restart
docker-compose -f docker-compose.pro.yml restart

# Stop
docker-compose -f docker-compose.pro.yml down
```

---

## 🔧 What Each Method Does

### Source Build (`install/setup.sh`)
1. Installs Go, Redis, and dependencies
2. Builds binary directly on VPS
3. Sets up systemd service
4. Configures Let's Encrypt
5. Starts service

### Docker Build
1. Builds Docker image with Sauron
2. Runs containerized with Redis
3. Handles dependencies in containers
4. Uses docker-compose for orchestration

---

## 🚨 Choose ONE Method

**Don't mix methods!** Pick either:
- **Source Build** → Use `install/setup.sh`
- **Docker** → Use docker-compose

**Most people should use Source Build** - it's simpler and your original working method.

---

## 🔍 Troubleshooting

### Source Build Issues
```bash
# Check if service is running
sudo systemctl status sauron

# Check binary exists
ls -la /usr/local/bin/sauron

# Rebuild if needed
cd /path/to/Sauron-Pro
go build -o /usr/local/bin/sauron
sudo systemctl restart sauron
```

### Docker Issues
```bash
# Check containers
docker ps

# Check Docker image
docker images | grep sauron

# Rebuild if needed
docker-compose -f docker-compose.pro.yml build --no-cache
```

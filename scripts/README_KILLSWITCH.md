# 🔴 KILL SWITCH SYSTEM - SCRIPTS OVERVIEW

This directory contains the complete kill switch system implementation for Sauron Pro.

## 🚨 CRITICAL LIMITATIONS DISCLAIMER

> **⚠️ IMPORTANT: This kill switch is NOT effective against cloud provider cooperation with authorities.**

**What Cloud Providers CAN Access:**
- Memory dumps, disk forensics, network logs, automated snapshots, hypervisor logs

**Kill Switch Protects Against:** Physical raids, customer detection, competitor analysis  
**Kill Switch Does NOT Protect Against:** Cloud provider cooperation, hypervisor forensics, logged traffic

**Use additional precautions and assume all cloud activity is permanently logged.**

## 📁 KILL SWITCH COMPONENTS

### Core Scripts

#### `kill-switch.sh` 🔴

**Purpose**: Emergency destruction tool for VPS instances  
**Usage**: `./kill-switch.sh [VPS_ID] [LEVEL] [DELAY] [REASON]`  
**Features**:

- Interactive and command-line modes
- Fleet-wide or single VPS destruction
- 5-level destruction intensity
- Triple confirmation safety
- Complete audit logging

#### `heartbeat-master.sh` 💓

**Purpose**: Master heartbeat service for dead man's switch  
**Usage**: `./heartbeat-master.sh {start|stop|restart|status}`  
**Features**:

- Automatic VPS fleet heartbeat monitoring
- Configurable intervals and timeouts
- Service-based operation with systemd
- Comprehensive logging and status reporting

#### `test-killswitch.sh` 🧪

**Purpose**: Integration testing for kill switch system  
**Usage**: `./test-killswitch.sh`  
**Features**:

- Comprehensive component validation
- Security and configuration testing
- Integration readiness assessment
- Detailed test reporting

### Service Files

#### `../install/sauron-heartbeat.service`

**Purpose**: Systemd service for heartbeat daemon  
**Location**: `/etc/systemd/system/sauron-heartbeat.service`  
**Features**:

- Automatic service management
- Security hardening
- Resource limits
- Restart policies

## 🎯 QUICK START

### 1. Emergency Kill Switch Activation

```bash
# Interactive mode (recommended)
./kill-switch.sh

# Direct fleet destruction (extreme emergency)
SAURON_ADMIN_KEY="your_key" ./kill-switch.sh ALL 5 0 "CODE RED"
```

### 2. Dead Man's Switch Setup

```bash
# Install heartbeat service
sudo cp ../install/sauron-heartbeat.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable sauron-heartbeat

# Configure environment
export SAURON_ADMIN_KEY="your_admin_key"
export SAURON_MASTER_URL="https://master.sauron.pro"

# Start heartbeat service
sudo systemctl start sauron-heartbeat
```

### 3. System Testing

```bash
# Run comprehensive tests
./test-killswitch.sh

# Check test results
cat /tmp/killswitch-test-*.log
```

## 🛠️ CONFIGURATION

### Environment Variables

```bash
# Required
export SAURON_ADMIN_KEY="your_admin_key"

# Optional
export SAURON_MASTER_URL="https://master.sauron.pro"
export HEARTBEAT_INTERVAL="30"  # seconds
export VPS_ID="unique_vps_id"
```

### Configuration Files

- **Dead Man's Switch**: Stored in secure_config database table
- **VPS Fleet**: Managed via master controller database
- **Heartbeat Service**: `/etc/systemd/system/sauron-heartbeat.service`

## 🔐 SECURITY FEATURES

### Authentication

- **Admin Key**: Primary authentication mechanism
- **Confirmation Code**: "Lord Sauron" required for activation
- **Triple Confirmation**: CLI requires three separate confirmations

### Destruction Levels

1. **Level 1**: Memory purge only (RAM cleanup)
2. **Level 2**: Data obliteration (logs, databases)
3. **Level 3**: System corruption (config files)
4. **Level 4**: Hardware destruction (disk wiping)
5. **Level 5**: Stealth exit (complete annihilation)

### Safety Mechanisms

- **Delayed Activation**: Configurable delays before execution
- **Audit Logging**: Complete operation tracking
- **Network Validation**: Connectivity verification
- **Permission Checks**: System access validation

## 🚨 EMERGENCY PROCEDURES

### Immediate Response (Code Red)

```bash
# Stop all operations immediately
./kill-switch.sh ALL 5 0 "EMERGENCY SHUTDOWN"
```

### Partial Response (Code Yellow)

```bash
# Destroy specific compromised VPS
./kill-switch.sh suspicious-vps-id 3 60 "Security breach"
```

### Dead Man's Switch Emergency

```bash
# Configure automatic destruction if master offline
curl -X POST https://vps.sauron.pro/admin/deadmans \
  -H "Authorization: Bearer $ADMIN_KEY" \
  -d '{"enabled":true,"master_timeout":"15m","auto_destruct":true}'
```

## 📊 MONITORING

### Service Status

```bash
# Check heartbeat service
sudo systemctl status sauron-heartbeat

# View service logs
sudo journalctl -u sauron-heartbeat -f

# Check heartbeat script status
./heartbeat-master.sh status
```

### Log Files

- **Kill Switch**: `/tmp/killswitch-*.log`
- **Heartbeat**: `/var/log/sauron/heartbeat.log`
- **System**: `/var/log/sauron/system.log`
- **Test Results**: `/tmp/killswitch-test-*.log`

### Fleet Health

```bash
# Master controller VPS list
curl -H "Authorization: Bearer $ADMIN_KEY" \
     https://master.sauron.pro/fleet/instances

# Individual VPS status
curl -H "Authorization: Bearer $ADMIN_KEY" \
     https://vps.sauron.pro/vps/status
```

## ⚠️ WARNINGS & LIMITATIONS

### Critical Warnings

- **IRREVERSIBLE**: All destruction operations are permanent
- **NO RECOVERY**: Zero data recovery possible after Level 3+
- **LEGAL COMPLIANCE**: Use only within legal boundaries
- **AUTHORIZATION**: Require proper operational authorization

### Technical Limitations

- **Hardware Access**: Level 4-5 require physical disk access
- **Cloud Providers**: Provider snapshots may survive destruction
- **Network Dependencies**: Remote activation requires connectivity
- **Permission Requirements**: Some operations need root access

### Known Issues

- **macOS Compatibility**: systemctl not available (use Docker/Linux)
- **Virtual Machines**: Limited hardware access in some environments
- **Cloud Snapshots**: Provider-level backups may persist
- **Hardware Logs**: Some hardware maintains internal audit logs

## 🔄 MAINTENANCE

### Regular Tasks

- **Weekly**: Test heartbeat connectivity
- **Monthly**: Rotate admin keys
- **Monthly**: Review destruction logs
- **Quarterly**: Update destruction methods

### Updates

- **Security Patches**: Keep destruction methods current
- **Key Rotation**: Rotate authentication keys regularly
- **Documentation**: Update procedures and protocols
- **Testing**: Regular integration testing

## 📞 SUPPORT

### Emergency Contacts

- **Security Team**: For immediate kill switch activation
- **Operations**: For service and maintenance issues
- **Development**: For system integration problems

### Troubleshooting

1. **Authentication Failures**: Verify admin key configuration
2. **Network Issues**: Check VPS connectivity and certificates
3. **Service Problems**: Review systemd service logs
4. **Compilation Errors**: Check Go dependencies and imports

---

> **⚠️ REMEMBER: This is a secured system designed for emergency use only.**
> **Improper use will result in complete system destruction.**
> **Always follow proper authorization and safety procedures.**

---

**Last Updated**: August 17, 2025  
**Version**: 2.1.0-omega  
**Classification**: RESTRICTED

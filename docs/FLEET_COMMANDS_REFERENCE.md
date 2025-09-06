# Sauron Fleet Management Commands Reference

Complete reference of all 29 available fleet management commands for your frontend development.

## API Endpoint

**POST** `/api/fleet/command`

## Request Format

```json
{
  "command": "command_name",
  "target_vps": "vps-id",
  "data": {
    "param1": "value1",
    "param2": "value2"
  }
}
```

## Authentication

All commands require **Firestore witness authentication** (no admin key bypass).

---

## 📊 Status & Information Commands

### 1. `status`

Get VPS status and basic information.

```json
{
  "command": "status",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Status, uptime, version, role

### 2. `health_check`

Comprehensive system health check.

```json
{
  "command": "health_check",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** System health status, configuration checks, network connectivity

### 3. `get_metrics`

Get system performance metrics.

```json
{
  "command": "get_metrics",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** CPU, memory, disk, network metrics

### 4. `get_ip_info`

Get detailed IP and location information.

```json
{
  "command": "get_ip_info",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** External IP, local IP, location, hostname

### 5. `get_version`

Get system version information.

```json
{
  "command": "get_version",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Binary version, Go version, system info, OS info, uptime

### 6. `get_stats`

Get comprehensive system statistics.

```json
{
  "command": "get_stats",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Slug stats, log stats, database stats, disk usage

---

## 🎛️ Fleet Management Commands

### 7. `promote_to_master`

Promote VPS instance to fleet master.

```json
{
  "command": "promote_to_master",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Success/error status

### 8. `demote_to_agent`

Demote VPS instance to agent.

```json
{
  "command": "demote_to_agent",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "new_master_id": "vps-master-id"
  }
}
```

**Returns:** Success/error status

### 9. `restart`

Restart the VPS instance.

```json
{
  "command": "restart",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Acknowledgment message

---

## ⚙️ Configuration Commands

### 10. `update_config`

Update VPS configuration dynamically.

```json
{
  "command": "update_config",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "payload": {
      "SAURON_DOMAIN": "newdomain.com",
      "ADMIN_KEY": "new-admin-key",
      "DEV_MODE": true,
      "heartbeat_interval": "300"
    }
  }
}
```

**Returns:** Update results, restart requirement status

### 11. `configure_env`

Configure environment variables with deployment type detection.

```json
{
  "command": "configure_env",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "SAURON_DOMAIN": "yourdomain.com",
    "ADMIN_KEY": "your-admin-key",
    "TELEGRAM_BOT_TOKEN": "your-bot-token",
    "action": "setup"
  }
}
```

**Returns:** Configuration status, deployment type, applied options

### 12. `export_config`

Export current configuration.

```json
{
  "command": "export_config",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Complete configuration export

### 13. `get_config`

Get full system configuration.

```json
{
  "command": "get_config",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Environment config, system config, deployment type, VPS config

### 14. `check_configuration`

Validate current configuration.

```json
{
  "command": "check_configuration",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Configuration validation results

---

## 🔧 Service Management Commands

### 15. `start_service`

Start the Sauron service.

```json
{
  "command": "start_service",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Service start status (Docker or systemd)

### 16. `stop_service`

Stop the Sauron service.

```json
{
  "command": "stop_service",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Service stop status

### 17. `restart_service`

Restart the Sauron service.

```json
{
  "command": "restart_service",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Service restart status

### 18. `service_status`

Get service status information.

```json
{
  "command": "service_status",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Service status, container status (Docker), or systemd status

### 19. `daemon_reload`

Reload systemd daemon.

```json
{
  "command": "daemon_reload",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Reload status

---

## 📊 Monitoring & Logs Commands

### 20. `get_logs`

Get system logs with filtering options.

```json
{
  "command": "get_logs",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "lines": 100,
    "type": "all"
  }
}
```

**Returns:** System logs, bot logs, emit logs, deployment-specific logs

### 21. `live_logs`

Get live/real-time logs.

```json
{
  "command": "live_logs",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "lines": 50,
    "follow": true
  }
}
```

**Returns:** Live log stream

### 22. `service_logs`

Get service-specific logs.

```json
{
  "command": "service_logs",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "service": "sauron",
    "lines": 100
  }
}
```

**Returns:** Service logs (Docker or systemd)

---

## 🔄 Update & Maintenance Commands

### 23. `check_updates`

Check for available updates.

```json
{
  "command": "check_updates",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Available updates, current version, update script availability

### 24. `update_system`

Update the system to latest version.

```json
{
  "command": "update_system",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "version": "latest"
  }
}
```

**Returns:** Update progress, deployment type detection, service restart status

### 25. `force_update`

Force system update regardless of version.

```json
{
  "command": "force_update",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "backup": true,
    "restart": true
  }
}
```

**Returns:** Force update status, backup creation, restart status

### 26. `backup_system`

Create system backup.

```json
{
  "command": "backup_system",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Backup creation status, backup file location

### 27. `clean_system`

Clean system (logs, Docker images, build artifacts).

```json
{
  "command": "clean_system",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Cleanup results, freed space

---

## 🌐 Network & Security Commands

### 28. `test_connectivity`

Test network connectivity.

```json
{
  "command": "test_connectivity",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Network connectivity test results

### 29. `test_domain`

Test domain connectivity and DNS resolution.

```json
{
  "command": "test_domain",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "domain": "yourdomain.com"
  }
}
```

**Returns:** DNS resolution, HTTP/HTTPS connectivity, endpoint tests

### 30. `check_ssl`

Check SSL certificate status.

```json
{
  "command": "check_ssl",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "domain": "yourdomain.com"
  }
}
```

**Returns:** Certificate validity, expiration, issuer information

### 31. `test_firebase`

Test Firebase connectivity.

```json
{
  "command": "test_firebase",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Firebase admin key status, endpoint connectivity

### 32. `verify_installation`

Verify system installation integrity.

```json
{
  "command": "verify_installation",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Installation verification results, missing components

---

## 🚨 System Management Commands

### 33. `kill_processes`

Kill all Sauron processes.

```json
{
  "command": "kill_processes",
  "target_vps": "vps-WebDev-Mac.local-1756812966"
}
```

**Returns:** Process termination results

### 34. `uninstall_system`

Uninstall Sauron system completely.

```json
{
  "command": "uninstall_system",
  "target_vps": "vps-WebDev-Mac.local-1756812966",
  "data": {
    "remove_data": false,
    "remove_certificates": false
  }
}
```

**Returns:** Uninstallation progress and results

---

## 🔧 Deployment Type Detection

The system automatically detects deployment type for each command:

### Docker Deployment

- **Detection**: Checks for `/.dockerenv` or running `docker-compose.pro.yml`
- **Service Commands**: Uses `docker-compose` commands
- **Config Updates**: Updates `.env` and restarts containers

### Systemd Deployment (Production)

- **Detection**: Checks for `/usr/local/bin/sauron` and `/etc/systemd/system/sauron.service`
- **Service Commands**: Uses `systemctl` commands
- **Config Updates**: Updates `.env` and restarts systemd service

### Development Mode

- **Detection**: Local `./sauron` binary or Go compiler available
- **Service Commands**: Manual restart required
- **Config Updates**: Updates `.env` only

---

## 📝 Response Format

All commands return a standardized response:

```json
{
  "status": "success|error",
  "action": "command_name",
  "deployment_type": "docker|systemd|development",
  "timestamp": "2025-09-02T12:34:56Z",
  "data": {
    // Command-specific response data
  },
  "error": "Error message if status is error"
}
```

---

## 🚀 Usage Examples for Frontend

### Basic Command Execution

```javascript
const executeCommand = async (command, targetVps, data = {}) => {
  const response = await fetch('/api/fleet/command', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer YOUR_FIRESTORE_TOKEN'
    },
    body: JSON.stringify({
      command,
      target_vps: targetVps,
      data
    })
  });
  return response.json();
};

// Get VPS status
const status = await executeCommand('status', 'vps-id');

// Update configuration
const configUpdate = await executeCommand('update_config', 'vps-id', {
  payload: {
    SAURON_DOMAIN: 'newdomain.com',
    DEV_MODE: false
  }
});

// Check for updates
const updates = await executeCommand('check_updates', 'vps-id');
```

### Batch Operations

```javascript
const executeMultipleCommands = async (vpsId, commands) => {
  const results = await Promise.all(
    commands.map(cmd => executeCommand(cmd.command, vpsId, cmd.data))
  );
  return results;
};

// Health check + metrics + logs
const healthReport = await executeMultipleCommands('vps-id', [
  { command: 'health_check' },
  { command: 'get_metrics' },
  { command: 'get_logs', data: { lines: 50 } }
]);
```

---

## 📊 Command Categories for UI Organization

### 📋 **Information & Status** (6 commands)

`status`, `health_check`, `get_metrics`, `get_ip_info`, `get_version`, `get_stats`

### 🎛️ **Fleet Management** (3 commands)  

`promote_to_master`, `demote_to_agent`, `restart`

### ⚙️ **Configuration** (5 commands)

`update_config`, `configure_env`, `export_config`, `get_config`, `check_configuration`

### 🔧 **Service Control** (5 commands)

`start_service`, `stop_service`, `restart_service`, `service_status`, `daemon_reload`

### 📊 **Monitoring & Logs** (3 commands)

`get_logs`, `live_logs`, `service_logs`

### 🔄 **Updates & Maintenance** (5 commands)

`check_updates`, `update_system`, `force_update`, `backup_system`, `clean_system`

### 🌐 **Network & Security** (5 commands)

`test_connectivity`, `test_domain`, `check_ssl`, `test_firebase`, `verify_installation`

### 🚨 **System Management** (2 commands)

`kill_processes`, `uninstall_system`

---

This reference contains all **29 fleet management commands** available in your Sauron system, organized for easy frontend development. Each command includes complete JSON examples and expected response formats.

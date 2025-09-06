# Fleet Remote Management Commands

## Overview

The fleet management system now provides comp#### health_check

Perform comprehensive health checks on the VPS.

```bash
curl -s -k -X POST https://localhost:443/fleet/command 
  -H "Host: admin.yourdomain.com" 
  -H "Content-Type: application/json" 
  -H "X-Admin-Key: your_admin_key" 
  -d '{
    "vps_id": "vps-001",
    "command": "health_check"
  }'
```

#### get_metrics

Get system performance metrics.

```bash
curl -s -k -X POST https://localhost:443/fleet/command 
  -H "Host: admin.yourdomain.com" 
  -H "Content-Type: application/json" 
  -H "X-Admin-Key: your_admin_key" 
  -d '{
    "vps_id": "vps-001",
    "command": "get_metrics"
  }'
```

#### test_connectivity

Test network connectivity to various endpoints.

```bash
curl -s -k -X POST https://localhost:443/fleet/command 
  -H "Host: admin.yourdomain.com" 
  -H "Content-Type: application/json" 
  -H "X-Admin-Key: your_admin_key" 
  -d '{
    "vps_id": "vps-001",
    "command": "test_connectivity"
  }'
```te VPS management capabilities, eliminating the need to SSH into individual VPS instances for basic operations. All commands are executed through the fleet command API with **Firestore witness authentication** properly implemented for sensitive operations.

## Prerequisites

- Fleet system must be properly configured
- **Firestore witness authentication is properly implemented and required for sensitive commands**
- Admin API key must be configured for basic operations
- Target VPS must be registered and sending heartbeats

## Authentication Levels

### Basic Commands (Admin Key Only)

- Information gathering commands (get_ip_info, health_check, get_metrics, etc.)
- Configuration viewing (get_config, export_config)
- Status monitoring (service_status, get_stats)

### Sensitive Commands (Firestore Witness Required)

- Service management (promote_to_master, demote_to_agent, restart)
- System modifications requiring elevated security
- Role changes and critical operations

## Fleet Command Structure

All fleet commands use these endpoints:

- **GET /fleet/register** - List all registered VPS instances
- **POST /fleet/command** - Send commands to specific VPS instances

### List All VPS Instances

```bash
curl -s -k -X GET https://localhost:443/fleet/register \
  -H "Host: admin.yourdomain.com" \
  -H "X-Admin-Key: your_admin_key" \
  -H "Content-Type: application/json"
```

### Send Commands to VPS

All fleet commands follow this structure:

```json
{
  "admin_key": "your_admin_api_key",
  "vps_id": "target_vps_id",
  "command": "command_name",
  "payload": {
    // Command-specific data
  },
  "timeout": 30
}
```

## Command Categories

### 1. Information & Monitoring Commands

#### get_ip_info

Get comprehensive IP information for the VPS.

```bash
curl -s -k -X POST https://localhost:443/fleet/command \
  -H "Host: admin.yourdomain.com" \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "get_ip_info"
  }'
```

**Response:**

```json
{
  "primary_ip": "216.131.72.129",
  "all_interfaces": {
    "eth0": "10.0.0.5",
    "lo": "127.0.0.1"
  },
  "external_services": {
    "ipify": "216.131.72.129",
    "ipinfo": "216.131.72.129"
  }
}
```

#### health_check

Perform comprehensive health checks on the VPS.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "health_check"
  }'
```

#### get_metrics

Get system performance metrics.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "get_metrics"
  }'
```

#### test_connectivity

Test network connectivity to various endpoints.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "test_connectivity"
  }'
```

### 2. Configuration Management Commands

#### update_config

Update VPS configuration in real-time.

```bash
# Update domain
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "update_config",
    "payload": {
      "domain": "newdomain.com"
    }
  }'

# Update multiple environment variables
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "update_config",
    "payload": {
      "domain": "newdomain.com",
      "heartbeat_interval": "600",
      "dev_mode": "false",
      "env_vars": {
        "CUSTOM_VAR": "custom_value",
        "ANOTHER_VAR": "another_value"
      }
    }
  }'
```

#### export_config

Export current configuration from VPS.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "export_config"
  }'
```

#### get_config

Get full system configuration.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "get_config"
  }'
```

### 3. Service Management Commands

#### start_service

Start the Sauron service.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "start_service"
  }'
```

#### stop_service

Stop the Sauron service.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "stop_service"
  }'
```

#### restart_service

Restart the Sauron service.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "restart_service"
  }'
```

#### service_status

Get detailed service status and resource usage.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "service_status"
  }'
```

**Response (Docker deployment):**

```json
{
  "status": "success",
  "deployment_type": "docker",
  "docker_status": "container status output",
  "resource_usage": "CPU and memory usage"
}
```

**Response (Binary deployment):**

```json
{
  "status": "success", 
  "deployment_type": "binary",
  "service_status": "systemctl status output",
  "process_info": "process information",
  "resource_usage": "CPU and memory usage"
}
```

### 4. System Operations Commands

#### get_logs

Retrieve system logs with customizable line count.

```bash
# Get last 100 lines (default)
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "get_logs"
  }'

# Get last 50 lines
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "get_logs",
    "payload": {
      "lines": 50
    }
  }'
```

**Response:**

```json
{
  "status": "success",
  "lines_requested": 50,
  "system_log": "system log content",
  "bot_log": "bot log content", 
  "emit_log": "emit log content",
  "docker_logs": "docker compose logs" // for Docker deployments
}
```

#### backup_system

Create a complete system backup.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "backup_system"
  }'
```

**Response:**

```json
{
  "status": "success",
  "backup_started": "20241222-143012",
  "backup_directory": "backup-20241222-143012",
  "backup_archive": "backup-20241222-143012.tar.gz",
  "success": true
}
```

#### update_system

Update the system to a specific version.

```bash
# Update to default version (v2.0.0-pro)
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "update_system"
  }'

# Update to specific version
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "update_system",
    "payload": {
      "version": "v2.1.0-pro"
    }
  }'
```

#### get_stats

Get comprehensive system statistics.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "get_stats"
  }'
```

**Response:**

```json
{
  "status": "success",
  "slug_stats": {
    "slugs": {
      "login": {"total_requests": 1423},
      "admin": {"total_requests": 89}
    }
  },
  "log_stats": {
    "system_lines": "2456",
    "bot_lines": "892", 
    "emit_lines": "1234"
  },
  "database_size": 2048576,
  "database_modified": "2024-12-22 14:30:12",
  "disk_usage": "1.2G"
}
```

#### clean_system

Clean logs, temporary files, and build artifacts.

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -d '{
    "vps_id": "vps-001",
    "command": "clean_system"
  }'
```

**Response:**

```json
{
  "status": "success",
  "cleaned_items": [
    "old logs",
    "docker system", 
    "release artifacts"
  ],
  "cleanup_completed": "2024-12-22 14:30:12",
  "success": true
}
```

### 5. Fleet Role Management Commands

These commands **require Firestore witness authentication** for security.

#### promote_to_master

Promote a VPS agent to master role. **Requires Firestore witness headers.**

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -H "X-Request-ID: firebase_request_id" \
  -H "X-Valid-Until: timestamp" \
  -d '{
    "vps_id": "vps-002",
    "command": "promote_to_master"
  }'
```

#### demote_to_agent

Demote a master VPS to agent role. **Requires Firestore witness headers.**

```bash
curl -X POST https://master.yourdomain.com/fleet/command \
  -H "Content-Type: application/json" \
  -H "X-Admin-Key: your_admin_key" \
  -H "X-Request-ID: firebase_request_id" \
  -H "X-Valid-Until: timestamp" \
  -d '{
    "vps_id": "vps-001",
    "command": "demote_to_agent",
    "payload": {
      "new_master_id": "vps-002"
    }
  }'
```

## Response Format

All commands return responses in this format:

```json
{
  "success": true,
  "vps_id": "vps-001",
  "command": "command_name",
  "result": {
    // Command-specific response data
  },
  "duration": "2.34s",
  "timestamp": "2024-12-22T14:30:12Z"
}
```

## Error Handling

If a command fails, the response will include error details:

```json
{
  "success": false,
  "vps_id": "vps-001", 
  "command": "command_name",
  "error": "Error description",
  "duration": "1.23s",
  "timestamp": "2024-12-22T14:30:12Z"
}
```

## Security Considerations

1. **Admin Authentication**: All commands require valid admin API key
2. **Firestore Witness**: Sensitive commands (promote_to_master, demote_to_agent, restart) require Firestore witness authentication with X-Request-ID and X-Valid-Until headers
3. **Network Security**: Commands are sent over HTTPS
4. **Timeout Protection**: All commands have configurable timeouts
5. **Audit Logging**: All fleet commands are logged for security auditing
6. **Command Classification**: Commands are automatically classified as basic (admin key only) or sensitive (Firestore witness required)

## Deployment Type Detection

The system automatically detects deployment type (Docker vs Binary) and adjusts commands accordingly:

- **Docker Deployment**: Uses `docker-compose` commands
- **Binary Deployment**: Uses `systemctl` commands

## Command Timeout

Default timeout is 30 seconds. For longer operations, specify custom timeout:

```json
{
  "vps_id": "vps-001",
  "command": "backup_system",
  "timeout": 120
}
```

## Best Practices

1. **Test Connectivity**: Use `test_connectivity` before other operations
2. **Health Checks**: Regular `health_check` commands to monitor VPS status
3. **Backup Before Updates**: Always run `backup_system` before `update_system`
4. **Monitor Resources**: Use `get_metrics` and `service_status` for monitoring
5. **Clean Regularly**: Use `clean_system` to maintain disk space
6. **Configuration Management**: Use `export_config` to backup configurations
7. **Authentication Security**: Ensure Firestore witness authentication is properly configured for sensitive operations
8. **Command Classification**: Understand which commands require elevated Firestore witness authentication

## Examples for Common Operations

### Deploy Configuration Change Across Fleet

```bash
# 1. Export current config from master
curl -X POST https://master.yourdomain.com/fleet/command \
  -d '{"vps_id": "master-vps", "command": "export_config"}'

# 2. Update all agents with new domain
for vps in vps-001 vps-002 vps-003; do
  curl -X POST https://master.yourdomain.com/fleet/command \
    -d "{\"vps_id\": \"$vps\", \"command\": \"update_config\", \"payload\": {\"domain\": \"newdomain.com\"}}"
done

# 3. Restart services on all VPS
for vps in vps-001 vps-002 vps-003; do
  curl -X POST https://master.yourdomain.com/fleet/command \
    -d "{\"vps_id\": \"$vps\", \"command\": \"restart_service\"}"
done
```

### Health Check All VPS

```bash
# Check health of all VPS instances
for vps in vps-001 vps-002 vps-003; do
  echo "Checking $vps..."
  curl -X POST https://master.yourdomain.com/fleet/command \
    -d "{\"vps_id\": \"$vps\", \"command\": \"health_check\"}"
done
```

### Emergency System Update

```bash
# 1. Create backups
for vps in vps-001 vps-002 vps-003; do
  curl -X POST https://master.yourdomain.com/fleet/command \
    -d "{\"vps_id\": \"$vps\", \"command\": \"backup_system\"}"
done

# 2. Update systems
for vps in vps-001 vps-002 vps-003; do
  curl -X POST https://master.yourdomain.com/fleet/command \
    -d "{\"vps_id\": \"$vps\", \"command\": \"update_system\", \"payload\": {\"version\": \"v2.1.0-security-patch\"}}"
done
```

This comprehensive remote management system eliminates the need to SSH into individual VPS instances for routine operations, providing centralized control through the fleet management interface.

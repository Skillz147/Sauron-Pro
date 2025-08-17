# Sauron-Pro Fleet Management System

## Overview

The Fleet Management System allows you to deploy and control multiple Sauron-Pro instances across different VPS servers from a centralized master controller. This creates a distributed MITM network that can be managed from a single admin interface.

## Architecture

```
Master Controller (master.example.com:8443)
├── Fleet Management API
├── VPS Registration & Heartbeat
├── Command Dispatch System
└── Centralized Monitoring

Connected VPS Instances
├── VPS-001 (vps1.example.com)
├── VPS-002 (vps2.example.com)
├── VPS-003 (vps3.example.com)
└── ... (up to 100+ VPS instances)
```

## Components

### 1. Master Controller

- **Location**: Single server (e.g., `master.example.com`)
- **Purpose**: Central command and control for all VPS instances
- **Features**:
  - VPS registration and heartbeat monitoring
  - Command dispatch to specific VPS instances
  - Fleet-wide statistics and monitoring
  - Script execution across multiple VPS
  - Real-time status dashboard

### 2. VPS Agents

- **Location**: Each individual VPS server
- **Purpose**: Local MITM operations + communication with master
- **Features**:
  - Automatic registration with master controller
  - Periodic heartbeat (every 5 minutes)
  - Command receiver for remote operations
  - Local credential capture and processing
  - Script execution via master commands

## Installation

### Deploy Master Controller

1. **Prerequisites**:
   - Linux server with root access
   - Go 1.19+ installed
   - Domain pointing to server (e.g., `master.example.com`)
   - SSL certificate configured

2. **Deployment**:

   ```bash
   # Set your domain
   export DOMAIN=master.example.com
   
   # Deploy master controller
   sudo ./scripts/deploy-fleet-master.sh
   ```

3. **Verification**:

   ```bash
   # Check service status
   systemctl status sauron-fleet-master
   
   # View logs
   journalctl -u sauron-fleet-master -f
   
   # Test API
   curl https://master.example.com:8443/fleet/instances
   ```

### Deploy VPS Agents

1. **Prerequisites**:
   - Linux VPS with root access
   - Go 1.19+ installed
   - Network connectivity to master controller

2. **Deployment** (on each VPS):

   ```bash
   # Set master controller URL
   export MASTER_URL=https://master.example.com:8443
   
   # Optional: Set custom VPS ID and domain
   export VPS_ID=vps-001
   export VPS_DOMAIN=vps1.example.com
   
   # Deploy VPS agent
   sudo ./scripts/deploy-vps-agent.sh
   ```

3. **Verification**:

   ```bash
   # Check VPS agent status
   /opt/sauron-pro/bin/vps-status
   
   # Test heartbeat to master
   /opt/sauron-pro/bin/vps-heartbeat
   
   # Check service
   systemctl status sauron-vps-agent
   ```

## Usage

### Fleet Management Commands

**View all VPS instances**:

```bash
/opt/sauron-pro/bin/fleet-status
```

**Send command to specific VPS**:

```bash
/opt/sauron-pro/bin/fleet-command <vps-id> <command> [payload]
```

**Available commands**:

- `status` - Get VPS status and statistics
- `restart` - Restart VPS service
- `script` - Execute script on VPS
- `config` - Update VPS configuration
- `update` - Update VPS software

### API Endpoints

#### Master Controller APIs

**VPS Registration**:

```http
POST /fleet/register
Content-Type: application/json
X-VPS-ID: vps-001

{
  "ip": "192.168.1.100",
  "domain": "vps-001.example.com",
  "admin_domain": "admin.vps-001.example.com",
  "version": "v2.0.1",
  "location": "US-East"
}
```

**List VPS Instances**:

```http
GET /fleet/instances

Response:
{
  "success": true,
  "instances": [
    {
      "id": "vps-001",
      "ip": "192.168.1.100",
      "domain": "vps-001.example.com",
      "status": "active",
      "last_seen": "2025-08-17T10:30:00Z",
      "location": "US-East",
      "version": "v2.0.1"
    }
  ],
  "fleet_stats": {
    "total_vps": 5,
    "active_vps": 4
  }
}
```

**Send Command to VPS**:

```http
POST /fleet/command
Content-Type: application/json

{
  "vps_id": "vps-001",
  "command": "status",
  "payload": {}
}
```

#### VPS Agent APIs

**Command Receiver**:

```http
POST /vps/command
Content-Type: application/json

{
  "vps_id": "vps-001",
  "command": "script",
  "payload": {
    "script": "update-sauron-template.sh"
  }
}
```

## Configuration

### Master Controller Configuration

**Location**: `/etc/sauron-pro/fleet-config.json`

```json
{
  "fleet_master": {
    "domain": "master.example.com",
    "port": 8443,
    "max_vps_instances": 100,
    "heartbeat_timeout": 600,
    "command_timeout": 30,
    "database": {
      "path": "/opt/sauron-pro/data/fleet.db",
      "backup_interval": 3600
    },
    "security": {
      "require_vps_auth": true,
      "max_command_rate": 10,
      "allowed_commands": ["status", "restart", "script", "config", "update"]
    }
  }
}
```

### VPS Agent Configuration

**Location**: `/etc/sauron-pro/vps-config.json`

```json
{
  "vps_agent": {
    "id": "vps-001",
    "domain": "vps-001.example.com",
    "location": "US-East",
    "master_url": "https://master.example.com:8443",
    "heartbeat_interval": 300,
    "command_port": 8444,
    "security": {
      "enable_auth": true,
      "max_command_rate": 5
    }
  }
}
```

## Monitoring and Maintenance

### Log Locations

**Master Controller**:

- Service logs: `/var/log/sauron-pro/fleet-master.log`
- Error logs: `/var/log/sauron-pro/fleet-master-error.log`
- System logs: `journalctl -u sauron-fleet-master`

**VPS Agent**:

- Service logs: `/var/log/sauron-pro/vps-agent.log`
- Error logs: `/var/log/sauron-pro/vps-agent-error.log`
- System logs: `journalctl -u sauron-vps-agent`

### Health Monitoring

**Check Fleet Status**:

```bash
# View all VPS instances and their status
curl -s https://master.example.com:8443/fleet/instances | jq '.'

# Check specific VPS
curl -s https://master.example.com:8443/fleet/instances | jq '.instances[] | select(.id=="vps-001")'
```

**VPS Heartbeat Monitoring**:

- VPS instances send heartbeat every 5 minutes
- Master marks VPS as inactive after 10 minutes without heartbeat
- Automatic alerts when VPS goes offline

### Troubleshooting

**VPS Not Registering**:

1. Check network connectivity to master controller
2. Verify `MASTER_URL` configuration
3. Check firewall rules (port 8443)
4. Review VPS agent logs

**Commands Not Executing**:

1. Verify VPS is active and responsive
2. Check command syntax and payload
3. Review VPS agent command receiver logs
4. Ensure command is in allowed commands list

**Performance Issues**:

1. Monitor heartbeat intervals
2. Check database performance on master
3. Review resource usage on VPS instances
4. Consider scaling master controller

## Security Considerations

### Network Security

- All communication over HTTPS
- VPS authentication via X-VPS-ID headers
- Rate limiting on command endpoints
- Firewall rules restricting access

### Access Control

- Master controller requires admin authentication
- VPS agents only accept commands from verified master
- Command rate limiting prevents abuse
- Audit logging of all fleet operations

### Data Protection

- No sensitive data stored in fleet database
- VPS identification only (no credential data)
- Encrypted communication channels
- Regular security audits

## Scaling

### Adding New VPS Instances

1. Deploy VPS agent on new server
2. Configure unique VPS_ID
3. VPS automatically registers with master
4. Begin receiving commands immediately

### Master Controller Scaling

- Single master handles 100+ VPS instances
- Database optimization for large fleets
- Load balancing for high-traffic scenarios
- Backup and disaster recovery procedures

## Examples

### Complete Deployment Example

**Step 1: Deploy Master Controller**

```bash
# On master server (master.example.com)
export DOMAIN=master.example.com
sudo ./scripts/deploy-fleet-master.sh
```

**Step 2: Deploy VPS Agents**

```bash
# On VPS-001 (vps1.example.com)
export MASTER_URL=https://master.example.com:8443
export VPS_ID=vps-001
export VPS_DOMAIN=vps1.example.com
sudo ./scripts/deploy-vps-agent.sh

# On VPS-002 (vps2.example.com)
export MASTER_URL=https://master.example.com:8443
export VPS_ID=vps-002
export VPS_DOMAIN=vps2.example.com
sudo ./scripts/deploy-vps-agent.sh
```

**Step 3: Verify Fleet**

```bash
# Check fleet status
/opt/sauron-pro/bin/fleet-status

# Send test command
/opt/sauron-pro/bin/fleet-command vps-001 status
```

### Script Execution Example

**Execute script on specific VPS**:

```bash
/opt/sauron-pro/bin/fleet-command vps-001 script '{"script": "update-sauron-template.sh"}'
```

**Execute script on all active VPS**:

```bash
# Get all active VPS IDs
VPS_LIST=$(curl -s https://master.example.com:8443/fleet/instances | jq -r '.instances[] | select(.status=="active") | .id')

# Execute script on each
for vps in $VPS_LIST; do
  /opt/sauron-pro/bin/fleet-command $vps script '{"script": "cleanup-logs.sh"}'
done
```

## Best Practices

1. **VPS Naming Convention**: Use descriptive IDs (`us-east-001`, `eu-west-002`)
2. **Regular Monitoring**: Check fleet status daily
3. **Staged Deployments**: Test commands on single VPS before fleet-wide
4. **Backup Strategy**: Regular database backups on master controller
5. **Security Updates**: Keep all instances updated with latest Sauron-Pro version
6. **Geographic Distribution**: Spread VPS across different regions for better coverage
7. **Capacity Planning**: Monitor resource usage and scale appropriately

## Support

For issues with the Fleet Management System:

1. Check the troubleshooting section above
2. Review log files on both master and VPS instances
3. Verify network connectivity and firewall rules
4. Ensure proper configuration of environment variables
5. Test with individual VPS before investigating fleet-wide issues

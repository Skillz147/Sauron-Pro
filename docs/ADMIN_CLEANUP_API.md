# Admin Cleanup API Documentation

## Overview

The Admin Cleanup API provides controlled cleanup operations for the Sauron MITM framework, allowing administrators to manage data retention, remove forensic traces, and maintain operational security.

## Authentication

The Admin Cleanup API supports two authentication methods:

### 🔐 Firestore-Witnessed Authentication (Recommended)

Generate a cryptographic proof and use the returned headers:

```bash
# Step 1: Generate proof via admin panel
curl -s "http://localhost:3000/api/auth/proof" | jq .

# Step 2: Use the returned headers
curl -k -X POST https://localhost:443/admin/cleanup \
  -H "X-Request-ID: [request_id_from_step_1]" \
  -H "X-Valid-Until: [timestamp_from_step_1]" \
  -H "Content-Type: application/json" \
  -d '{"operations": ["logs"], "retention_days": 30}'
```

### 🔑 Legacy Admin Key Authentication (Fallback)

Use the static admin key (deprecated):

```bash
curl -k -X POST https://localhost:443/admin/cleanup \
  -H "X-Admin-Key: your_admin_key_here" \
  -H "Content-Type: application/json" \
  -d '{"operations": ["logs"], "retention_days": 30}'
```

> **Note**: Firestore-witnessed authentication is strongly recommended for enhanced security. Admin keys are never transmitted over the network with this method.

## Endpoints

### POST /admin/cleanup

**Description:** Execute cleanup operations with fine-grained control.

**Request Body:**

```json
{
  "operations": ["logs", "database", "credentials", "firestore", "all"],
  "retention_days": 30,
  "dry_run": false
}
```

**Parameters:**

- `operations` (array, required): List of cleanup operations to perform
  - `"logs"`: Remove old log files (preserves active logs)
  - `"database"`: Clean old user_links records and vacuum database
  - `"credentials"`: Clear captured credentials from secure memory storage
  - `"firestore"`: Remove old documents from Firestore
  - `"all"`: Execute all cleanup operations
- `retention_days` (integer): Keep data newer than N days (0 = delete all)
- `dry_run` (boolean): Preview mode - shows what would be deleted without making changes

> **Authentication**: Use Firestore-witnessed authentication headers (X-Request-ID, X-Valid-Until) or legacy X-Admin-Key header

**Response:**

```json
{
  "success": true,
  "operations": {
    "logs": {
      "success": true,
      "items_removed": 15,
      "size_freed": 2048576,
      "details": "Processed 20 log files"
    },
    "database": {
      "success": true,
      "items_removed": 150,
      "size_freed": 15000,
      "details": "Found 150 old records for cleanup"
    }
  },
  "total_size_freed": 2063576,
  "message": "Cleanup completed. Freed 2.0 MB across 2 operations",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### GET /admin/cleanup/status

**Description:** Get current system status and cleanup statistics.

**Headers:**

- `X-Admin-Key`: Admin authentication key

**Response:**

```json
{
  "database": {
    "user_links": 1250,
    "config": 1,
    "banned_ips": 45
  },
  "logs": {
    "total_files": 12,
    "total_size": "15.2 MB"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Operation Details

### Logs Cleanup

- **Target:** Files in `logs/` directory
- **Behavior:**
  - Removes rotated log files older than retention period
  - Preserves active log files (`system.log`, `bot.log`)
  - Works with lumberjack rotation system
- **Safety:** High - only removes inactive log files

### Database Cleanup

- **Target:** SQLite `config.db` database
- **Behavior:**
  - Removes `user_links` records older than retention period
  - Executes `VACUUM` to reclaim space
  - Preserves `config`, `banned_ips`, and `licenses` tables
- **Safety:** Medium - affects operational data

### Credentials Cleanup

- **Target:** In-memory secure storage for captured credentials
- **Behavior:**
  - Clears captured credentials based on age
  - Works with capture package's secure storage
- **Safety:** High - memory-only operation

### Firestore Cleanup

- **Target:** Cloud-stored captured data
- **Behavior:**
  - Removes old documents from Firestore collections
  - Respects retention policies
- **Safety:** Low - affects persistent cloud data

## Security Considerations

### Operational Security

- All operations are logged with admin identification
- Dry run mode allows safe preview of cleanup operations
- Retention policies prevent accidental complete data deletion
- Admin key validation prevents unauthorized access

### Anti-Forensics

- Log cleanup removes traces of past operations
- Database vacuum eliminates deleted record fragments
- Credential cleanup clears memory evidence
- Firestore cleanup removes cloud evidence trails

## Usage Examples

### CLI Script Usage

```bash
# Preview log cleanup (30 days retention)
./scripts/admin_cleanup.sh logs 30 --dry-run

# Clean database records older than 7 days
./scripts/admin_cleanup.sh database 7

# Emergency cleanup - delete all cleanable data
./scripts/admin_cleanup.sh all 0

# Get current system status
./scripts/admin_cleanup.sh status
```

### Direct API Usage

```bash
# Dry run database cleanup
curl -X POST https://your-domain/admin/cleanup \
  -H "Content-Type: application/json" \
  -d '{
    "admin_key": "your_admin_key",
    "operations": ["database"],
    "retention_days": 7,
    "dry_run": true
  }'

# Get system status
curl -X GET https://your-domain/admin/cleanup/status \
  -H "X-Admin-Key: your_admin_key"
```

## Integration Notes

### Existing System Integration

- Uses existing admin authentication from WebSocket system
- Integrates with lumberjack log rotation
- Works with SQLite database structure
- Compatible with current capture and firestore packages

### Recommended Usage Patterns

1. **Regular Maintenance:** Weekly log cleanup with 30-day retention
2. **Security Cleanup:** Database cleanup after each campaign
3. **Emergency Cleanup:** Complete cleanup when operational security is compromised
4. **Status Monitoring:** Regular status checks to monitor data growth

## Error Handling

The API provides detailed error reporting for each operation:

- Authentication failures return 401 Unauthorized
- Invalid operations return specific error messages
- Partial failures include success/failure status per operation
- File system errors are captured and reported

## Recommendations

### Best Practices

- Run cleanup operations regularly to minimize forensic traces
- Use dry run mode before destructive operations
- Monitor system status to prevent data accumulation
- Consider automated cleanup policies for routine maintenance

### Emergency Procedures

- `all` operation with 0 retention for complete cleanup
- Verify cleanup completion with status endpoint
- Consider additional manual cleanup for complete trace removal

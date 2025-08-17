# 🔒 Secure SQLite Implementation Guide

## Overview

Sauron Pro now implements comprehensive SQLite security with multiple layers of protection:

1. **Database-level encryption** using SQLCipher
2. **Field-level encryption** for sensitive data
3. **Automatic data expiration** and cleanup
4. **Secure memory handling** with proper wiping
5. **Key rotation** and security auditing

## 🛡️ Security Features

### Database Encryption

- **SQLCipher**: Full database encryption at rest
- **AES-256**: Strong encryption for database files
- **Page-level encryption**: Every database page encrypted
- **Key derivation**: Environment-based or auto-generated keys

### Field-Level Protection

```go
// User IDs are hashed for indexing
userIDHash := sha256(userID + "sauron_salt")

// Sensitive fields are encrypted individually
encryptedSlug := AES-GCM(slug)
encryptedIP := AES-GCM(ipAddress)
```

### Data Minimization

- **Automatic expiry**: Records auto-delete after configurable time
- **Minimal retention**: Only essential data stored
- **Secure deletion**: Multiple-pass overwrite on cleanup

## 🗄️ Database Schema

### Secure User Links

```sql
CREATE TABLE secure_user_links (
    user_id_hash TEXT PRIMARY KEY,    -- SHA256 hash of user ID
    encrypted_slug TEXT,              -- AES-GCM encrypted slug
    created_at DATETIME,
    last_accessed DATETIME,
    expires_at DATETIME DEFAULT (datetime('now', '+30 days'))
);
```

### Secure Banned IPs

```sql
CREATE TABLE secure_banned_ips (
    ip_hash TEXT PRIMARY KEY,         -- SHA256 hash of IP
    encrypted_country TEXT,           -- AES-GCM encrypted country
    ban_reason TEXT,
    created_at DATETIME,
    expires_at DATETIME DEFAULT (datetime('now', '+7 days'))
);
```

### Security Audit Log

```sql
CREATE TABLE security_audit (
    id INTEGER PRIMARY KEY,
    event_type TEXT,                  -- 'user_link_stored', 'ip_banned', etc.
    event_data TEXT,                  -- Minimal audit information
    timestamp DATETIME,
    expires_at DATETIME DEFAULT (datetime('now', '+90 days'))
);
```

## 🔧 Implementation

### Environment Setup

```bash
# Set database encryption key (save this securely!)
export SAURON_DB_KEY="your-256-bit-base64-encoded-key"

# Or let system generate one (will be logged)
# System will auto-generate if not provided
```

### Initialization

```go
// Initialize secure databases
if err := configdb.InitDatabases(); err != nil {
    log.Fatal().Err(err).Msg("Failed to initialize secure databases")
}

// Migrate legacy data automatically
// Old config.db will be migrated to secure format
```

### Usage Examples

```go
// Store user link securely
secureDB := configdb.GetSecureDB()
err := secureDB.StoreSecureUserLink(userID, slug)

// Retrieve user link
slug, err := secureDB.GetSecureUserLink(userID)

// Ban IP with encryption
err := secureDB.BanIPSecure(ipAddress, country, reason)

// Check if IP is banned
isBanned := secureDB.IsIPBannedSecure(ipAddress)
```

## 🔄 Background Processes

### Data Cleanup

- **Frequency**: Every 1 hour
- **Action**: Delete expired records from all tables
- **Logging**: Reports number of records cleaned

### Key Rotation

- **Frequency**: Every 24 hours (configurable)
- **Trigger**: Automatic after 7 days
- **Process**: Logs warning to restart with new keys

### Security Audit

- **Event tracking**: All database operations logged
- **Retention**: 90 days (configurable)
- **Data**: Minimal information for forensics

## 🚨 Security Considerations

### What This Protects Against

✅ **Disk access**: Raw database files are encrypted  
✅ **Backup theft**: Stolen backups are unreadable  
✅ **Physical raids**: Database content is protected  
✅ **Forensic analysis**: Encrypted fields unreadable  

### What This Doesn't Protect Against

❌ **Runtime memory**: Decrypted data in application memory  
❌ **Process dumps**: Memory snapshots contain plaintext  
❌ **Cloud provider access**: Hypervisor can snapshot memory  
❌ **Application exploits**: Compromised app has access to decrypted data  

### Additional Recommendations

1. **Use VPNs**: Additional traffic obfuscation
2. **Rotate keys**: Regular key rotation schedule
3. **Monitor access**: Watch for unusual database activity
4. **Secure backups**: Encrypt database backups separately
5. **Network security**: Secure the VPS network layer

## 🔧 Migration Process

### Automatic Migration

```go
// System automatically detects legacy config.db
// Migrates data to encrypted storage
// Preserves existing functionality

// Migration includes:
// - user_links → secure_user_links
// - banned_ips → secure_banned_ips  
// - license_keys → secure_licenses
```

### Manual Migration

```bash
# If automatic migration fails
go run migrate_legacy.go

# Or use CLI tool
./sauron-cli migrate --from=config.db --to=config_secure.db
```

## 📊 Performance Impact

### Encryption Overhead

- **Database operations**: ~10-15% slower
- **Memory usage**: ~20% increase for encryption
- **Startup time**: Additional 1-2 seconds for initialization
- **Storage**: ~15% larger database files

### Optimization

- **Prepared statements**: Reuse encrypted queries
- **Connection pooling**: Minimize connection overhead
- **Indexing**: Hash-based indexes for performance
- **Caching**: In-memory cache for frequent lookups

## 🛠️ Troubleshooting

### Common Issues

**Database won't open**

```bash
# Check if key is set
echo $SAURON_DB_KEY

# Verify key format (should be base64)
echo $SAURON_DB_KEY | base64 -d | wc -c  # Should output 32
```

**Migration failed**

```bash
# Check legacy database
sqlite3 config.db ".tables"

# Manual migration
./sauron-cli migrate --force
```

**Performance issues**

```bash
# Check database size
ls -lh config_secure.db

# Analyze queries
PRAGMA query_planner = on;
```

### Emergency Recovery

```bash
# If secure database is corrupted
cp config_secure.db config_secure.db.backup

# Use legacy database temporarily
mv config.db config_legacy.db
./sauron-cli migrate --from=config_legacy.db
```

## 📋 Configuration Options

### Environment Variables

```bash
# Database encryption key (required)
SAURON_DB_KEY="base64-encoded-256-bit-key"

# Data retention (optional)
SAURON_DB_RETENTION_DAYS=30

# Cleanup frequency (optional)
SAURON_DB_CLEANUP_HOURS=1

# Key rotation warning (optional)
SAURON_DB_KEY_ROTATION_DAYS=7
```

### Runtime Configuration

```go
// Adjust retention periods
secureDB.SetRetentionDays("user_links", 30)
secureDB.SetRetentionDays("banned_ips", 7) 
secureDB.SetRetentionDays("audit_log", 90)

// Configure cleanup frequency
secureDB.SetCleanupInterval(1 * time.Hour)
```

## 🎯 Best Practices

1. **Key Management**
   - Store keys in secure environment variables
   - Never hardcode keys in source code
   - Use different keys for different environments
   - Rotate keys regularly

2. **Data Handling**
   - Minimize data retention periods
   - Use hashing where possible instead of encryption
   - Implement proper logging for security events
   - Monitor for unusual access patterns

3. **Operational Security**
   - Regular security audits
   - Monitor database size and performance
   - Implement proper backup encryption
   - Test disaster recovery procedures

4. **Development**
   - Use separate keys for development/production
   - Implement proper error handling
   - Test migration procedures thoroughly
   - Document key management procedures

---

**Remember**: This is defense in depth. The secure SQLite implementation significantly improves security but should be combined with other security measures including network security, access controls, and operational security practices.

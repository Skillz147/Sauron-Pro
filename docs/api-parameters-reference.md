# API Parameters Reference

## Authentication Methods

### Firestore-Witnessed Authentication (Recommended)

**Required Headers:**

- `X-Request-ID`: Unique request identifier from authentication proof
- `X-Valid-Until`: Unix timestamp (seconds) when proof expires
- `Content-Type`: application/json

**Example:**

```bash
curl -k -X GET https://localhost:443/admin/metrics \
  -H "X-Request-ID: a35905f9c363f8be3e297a5b59d8cf29" \
  -H "X-Valid-Until: 1755734576" \
  -H "Content-Type: application/json"
```

### Legacy Admin Key Authentication (Deprecated)

**Required Headers:**

- `X-Admin-Key`: Static admin key from configuration
- `Content-Type`: application/json

**Example:**

```bash
curl -k -X GET https://localhost:443/admin/metrics \
  -H "X-Admin-Key: your_admin_key_here" \
  -H "Content-Type: application/json"
```

## Common Response Format

All API endpoints return responses in this format:

```json
{
  "success": boolean,
  "data": object,          // Present on success
  "error": string,         // Present on error
  "authMethod": string,    // "firestore-witnessed" or "admin-key"
  "timestamp": string      // ISO 8601 timestamp
}
```

## Admin Endpoints

### GET /admin/metrics

Returns security metrics and alerts.

**Authentication:** Firestore-witnessed or Admin Key
**Parameters:** None

**Response:**

```json
{
  "success": true,
  "alerts": [],
  "summary": {
    "total_victims": 0,
    "risk_distribution": {
      "safe": 0,
      "suspicious": 0,
      "dangerous": 0,
      "law_enforcement": 0
    }
  },
  "timeframe_hours": 24,
  "authMethod": "firestore-witnessed",
  "timestamp": "2025-08-20T00:19:51.429Z"
}
```

### POST /admin/cleanup

Executes cleanup operations.

**Authentication:** Firestore-witnessed or Admin Key

**Request Body:**

```json
{
  "operations": ["logs", "database", "credentials", "firestore", "all"],
  "retention_days": 30,
  "dry_run": false
}
```

**Parameters:**

- `operations` (array, required): Operations to perform
- `retention_days` (integer, optional): Data retention period
- `dry_run` (boolean, optional): Preview mode if true

### GET /admin/status

Returns server status and health information.

**Authentication:** Firestore-witnessed or Admin Key
**Parameters:** None

### POST /admin/killswitch

Activates emergency shutdown procedures.

**Authentication:** Firestore-witnessed or Admin Key

**Request Body:**

```json
{
  "mode": "immediate",
  "reason": "security_breach",
  "preserve_logs": false
}
```

**Parameters:**

- `mode` (string, required): "immediate" or "graceful"
- `reason` (string, optional): Reason for activation
- `preserve_logs` (boolean, optional): Whether to keep logs

## Authentication Endpoints

### GET /api/auth/proof

Generates a new Firestore authentication proof.

**Authentication:** None (public endpoint for admin panel)

**Response:**

```json
{
  "success": true,
  "requestId": "a35905f9c363f8be3e297a5b59d8cf29",
  "validUntil": "2025-08-21T00:02:56.000Z",
  "validUntilTimestamp": 1755734576000,
  "signature": "3e69289a81eb9ef8774df5a5bfdef287ff8549c57614a5a3afd1a85dd4e17487",
  "message": "Authentication proof generated and stored in Firestore",
  "usage": {
    "headers": {
      "X-Request-ID": "a35905f9c363f8be3e297a5b59d8cf29",
      "X-Valid-Until": "1755734576000"
    },
    "note": "Use these headers instead of X-Admin-Key for API calls"
  }
}
```

### POST /api/auth/proof

Returns current cached authentication information.

**Authentication:** None (public endpoint for admin panel)

**Response:**

```json
{
  "success": true,
  "currentAuth": {
    "headers": {
      "X-Request-ID": "a35905f9c363f8be3e297a5b59d8cf29",
      "X-Valid-Until": "1755734576000",
      "Content-Type": "application/json"
    },
    "message": "Current authentication headers (cached or newly generated)"
  }
}
```

## Error Responses

### Authentication Errors

```json
{
  "success": false,
  "error": "Unauthorized - Invalid Firestore Authentication",
  "code": 401
}
```

```json
{
  "success": false,
  "error": "Auth proof expired",
  "code": 401
}
```

### Common Error Codes

- `400`: Bad Request - Invalid parameters
- `401`: Unauthorized - Authentication failed
- `403`: Forbidden - Access denied
- `404`: Not Found - Endpoint doesn't exist
- `500`: Internal Server Error - Server error

## Rate Limiting

- Authentication proof generation: 10 requests per minute
- Admin API endpoints: 100 requests per minute
- Global rate limit: 1000 requests per hour

## Best Practices

1. **Use Firestore Authentication**: Prefer Firestore-witnessed authentication over admin keys
2. **Cache Proofs**: Reuse authentication proofs for up to 24 hours
3. **Handle Expiration**: Implement automatic proof regeneration
4. **Error Handling**: Always check the `success` field in responses
5. **HTTPS Only**: Use HTTPS in production environments
6. **Monitor Usage**: Track API usage in server logs

## Migration from Admin Keys

Replace this pattern:

```bash
curl -H "X-Admin-Key: key" https://localhost:443/admin/endpoint
```

With this pattern:

```bash
# Generate proof
PROOF=$(curl -s "http://localhost:3000/api/auth/proof")
REQUEST_ID=$(echo $PROOF | jq -r '.usage.headers."X-Request-ID"')
VALID_UNTIL=$(echo $PROOF | jq -r '.usage.headers."X-Valid-Until"')

# Use proof
curl -H "X-Request-ID: $REQUEST_ID" -H "X-Valid-Until: $VALID_UNTIL" \
     https://localhost:443/admin/endpoint
```

## Support

For additional information, see:

- [Firestore Authentication Guide](FIRESTORE_AUTHENTICATION.md)
- [Admin Cleanup API](ADMIN_CLEANUP_API.md)
- [Security Documentation](security-features.html)

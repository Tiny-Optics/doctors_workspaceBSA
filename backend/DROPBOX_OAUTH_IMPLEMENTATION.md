# Dropbox OAuth Implementation - Complete Guide

## Overview

This implementation provides automatic Dropbox token refresh with admin UI management. The system automatically refreshes access tokens every ~4 hours using refresh tokens stored in the database.

## Architecture

### Components Created

1. **Database Layer**
   - `models/dropbox_config.go` - Configuration model with encrypted tokens
   - `repository/dropbox_config_repository.go` - Database operations

2. **Service Layer**
   - `service/encryption_service.go` - AES-256-GCM encryption for tokens
   - `service/dropbox_service.go` - Updated with auto-refresh logic
   - `service/dropbox_oauth_service.go` - OAuth flow management

3. **API Layer**
   - `handlers/dropbox_admin_handler.go` - Admin API endpoints
   - `server/routes.go` - Updated with new routes

## How It Works

### Token Lifecycle

```
1. Admin authorizes Dropbox (one-time)
   ↓
2. Backend receives refresh token (never expires)
   ↓
3. Backend encrypts and stores in MongoDB
   ↓
4. On every API call:
   - Check if access token expired
   - If expired, auto-refresh using refresh token
   - Continue with API call
   ↓
5. If refresh fails 3+ times → Alert admin (TODO: email)
```

### Automatic Refresh

The system checks token expiry on every Dropbox API call:
- If token expires in < 5 minutes → auto-refresh
- New access token valid for 4 hours
- Refresh token used to get new access tokens
- All happens transparently, no user interaction needed

## Setup Instructions

### Step 1: Environment Configuration

Add to your `.env` file:

```bash
# Generate a 32-byte encryption key
# Run: openssl rand -base64 32
ENCRYPTION_KEY=your-generated-key-here
```

### Step 2: Get Dropbox App Credentials

1. Go to https://www.dropbox.com/developers/apps
2. Create a new app or use existing
3. Settings → OAuth 2
4. Note your:
   - **App Key** (client_id)
   - **App Secret** (client_secret)
5. Add redirect URI (if using): `http://localhost:8080/admin/dropbox-success` (or your domain)

### Step 3: Initial Authorization (Via API)

#### Option A: Using curl

**Step 3a: Initiate OAuth**
```bash
curl -X POST http://localhost:8080/api/admin/dropbox/authorize \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "appKey": "YOUR_DROPBOX_APP_KEY",
    "appSecret": "YOUR_DROPBOX_APP_SECRET",
    "parentFolder": "/SOPS",
    "redirectUri": ""
  }'
```

Response:
```json
{
  "authUrl": "https://www.dropbox.com/oauth2/authorize?...",
  "message": "Visit this URL to authorize Dropbox access",
  "instructions": "After authorizing, you will receive a code..."
}
```

**Step 3b: Visit the Authorization URL**
- Copy the `authUrl` from the response
- Open in browser
- Sign in to Dropbox
- Click "Allow" to authorize
- Copy the authorization code displayed

**Step 3c: Complete Authorization**
```bash
curl -X POST http://localhost:8080/api/admin/dropbox/callback \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "AUTHORIZATION_CODE_FROM_DROPBOX",
    "appKey": "YOUR_DROPBOX_APP_KEY",
    "appSecret": "YOUR_DROPBOX_APP_SECRET",
    "parentFolder": "/SOPS",
    "redirectUri": ""
  }'
```

Response:
```json
{
  "message": "Dropbox successfully authorized and configured",
  "status": {
    "isConnected": true,
    "tokenExpiry": "2025-10-19T14:30:00Z",
    "consecutiveFailures": 0,
    ...
  }
}
```

**Done!** Dropbox is now configured and will auto-refresh tokens.

## API Endpoints

All endpoints require super admin permissions (`PermManageSystem`).

### GET /api/admin/dropbox/status
Get current Dropbox connection status.

**Response:**
```json
{
  "isConnected": true,
  "tokenExpiry": "2025-10-19T14:30:00Z",
  "lastRefreshSuccess": "2025-10-19T10:30:00Z",
  "consecutiveFailures": 0,
  "lastError": "",
  "needsReconnection": false,
  "parentFolder": "/SOPS"
}
```

### POST /api/admin/dropbox/authorize
Initiate OAuth flow.

**Request:**
```json
{
  "appKey": "your_app_key",
  "appSecret": "your_app_secret",
  "parentFolder": "/SOPS",
  "redirectUri": ""
}
```

**Response:**
```json
{
  "authUrl": "https://www.dropbox.com/oauth2/authorize?...",
  "message": "Visit this URL to authorize Dropbox access"
}
```

### POST /api/admin/dropbox/callback
Complete OAuth flow with authorization code.

**Request:**
```json
{
  "code": "authorization_code_from_dropbox",
  "appKey": "your_app_key",
  "appSecret": "your_app_secret",
  "parentFolder": "/SOPS",
  "redirectUri": ""
}
```

### POST /api/admin/dropbox/refresh
Manually trigger token refresh.

**Response:**
```json
{
  "message": "Token refresh successful"
}
```

### POST /api/admin/dropbox/test
Test Dropbox connection.

**Response:**
```json
{
  "success": true,
  "message": "Dropbox connection is working correctly"
}
```

### DELETE /api/admin/dropbox/configuration
Delete Dropbox configuration (requires re-authorization).

**Response:**
```json
{
  "message": "Dropbox configuration deleted successfully"
}
```

## Monitoring & Alerts

### Health Monitoring

The system tracks:
- Last successful refresh timestamp
- Consecutive failure count
- Last error message
- Token expiry time

### Alert Threshold

After 3 consecutive refresh failures:
- System logs critical alert
- **TODO:** Email notification to admin

### Manual Intervention Required When:

1. **Refresh token revoked** by user in Dropbox settings
2. **App credentials changed** in Dropbox console
3. **Security issues** detected by Dropbox

In these cases:
- System marks as disconnected
- Admin receives alert (TODO: email)
- Admin uses UI to re-authorize

## Security Features

### Encryption
- All tokens encrypted with AES-256-GCM
- Encryption key from environment variable
- Tokens never exposed in API responses

### Access Control
- Only super admins can configure Dropbox
- All actions logged in audit trail
- IP address tracking

### Database Security
- Sensitive fields encrypted at rest
- Singleton configuration (one record only)
- Automatic cleanup on failure

## Troubleshooting

### Problem: "Dropbox not configured"
**Solution:** Run initial authorization flow via API

### Problem: "Failed to refresh token"
**Possible causes:**
1. Refresh token revoked → Re-authorize via admin UI
2. App credentials invalid → Check Dropbox console
3. Network issues → Check connectivity

**Solution:** Check `/api/admin/dropbox/status` for details

### Problem: "Encryption key error"
**Solution:** Set `ENCRYPTION_KEY` in `.env` file
```bash
openssl rand -base64 32
```

### Problem: Access token expired error
**Expected behavior:** System should auto-refresh
**If persisting:** Check consecutive failures in status endpoint

## Database Structure

### Collection: `dropbox_config`

```javascript
{
  "_id": ObjectId("..."),
  "app_key": "your_app_key",
  "app_secret": "encrypted_secret",
  "refresh_token": "encrypted_refresh_token",
  "access_token": "encrypted_access_token",
  "token_expiry": ISODate("2025-10-19T14:30:00Z"),
  "parent_folder": "/SOPS",
  "is_connected": true,
  "last_refresh_success": ISODate("2025-10-19T10:30:00Z"),
  "last_refresh_attempt": ISODate("2025-10-19T10:30:00Z"),
  "consecutive_failures": 0,
  "last_error": "",
  "created_at": ISODate("2025-10-19T08:00:00Z"),
  "updated_at": ISODate("2025-10-19T10:30:00Z"),
  "created_by": ObjectId("...")
}
```

## Future Enhancements (TODO)

1. **Email Alerts**
   - Configure email service
   - Send alerts on 3+ failures
   - Include re-authorization link

2. **Admin UI Dashboard**
   - Visual token status
   - One-click reconnect button
   - Health metrics graph

3. **Backup Tokens**
   - Multiple Dropbox accounts
   - Automatic failover
   - Redundancy

4. **Metrics**
   - Track refresh frequency
   - Monitor API usage
   - Performance metrics

## Testing

### Test Connection
```bash
curl -X POST http://localhost:8080/api/admin/dropbox/test \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### Test Manual Refresh
```bash
curl -X POST http://localhost:8080/api/admin/dropbox/refresh \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### Check Status
```bash
curl -X GET http://localhost:8080/api/admin/dropbox/status \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

## Support

For issues or questions:
1. Check `/api/admin/dropbox/status` for health info
2. Review backend logs for detailed errors
3. Check audit logs for configuration changes
4. Verify Dropbox app settings in console

---

**Last Updated:** October 19, 2025
**Version:** 1.0.0


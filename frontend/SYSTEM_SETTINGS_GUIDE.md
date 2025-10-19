# System Settings - Admin Panel Guide

## Overview

The System Settings page provides a centralized location for managing all system-wide configurations through an intuitive tabbed interface.

## Location

**Path:** `/admin/settings`  
**Access:** Super Admin only (requires `PermManageSystem`)

## Features

### Horizontal Tab Navigation

The settings are organized into four main categories:

1. **Dropbox** - OAuth configuration and file storage management
2. **Email & Notifications** - Email settings (Coming Soon)
3. **Security** - Security policies and password requirements (Coming Soon)
4. **General** - Application preferences (Coming Soon)

## Dropbox Configuration Tab

### Features

#### 1. Connection Status Card
- **Visual Status Indicator**: Green (connected) or Red (disconnected)
- **Last Refresh Time**: When the token was last refreshed
- **Health Metrics Display**:
  - Token expiry time
  - Consecutive failures count
  - Parent folder path
- **Actions Button**: Quick access to management functions

#### 2. OAuth Configuration Flow

**Step 1: Enter Credentials**
- Dropbox App Key
- Dropbox App Secret
- Parent Folder (default: `/SOPS`)
- Click "Get Authorization URL"

**Step 2: Authorize in Dropbox**
- Click the "Open Dropbox Authorization" button
- Sign in to your Dropbox account
- Click "Allow" to grant permissions
- Copy the authorization code shown

**Step 3: Complete Authorization**
- Paste the authorization code
- Click "Complete Authorization"
- Configuration saved and token refresh begins automatically

#### 3. Management Actions

When connected, you can:
- **Test Connection**: Verify Dropbox is accessible
- **Force Refresh**: Manually refresh the access token
- **Disconnect**: Remove Dropbox configuration

#### 4. Health Monitoring

The system displays:
- ✅ Connection status (Connected/Disconnected)
- ⏰ Token expiry countdown
- 🔄 Last successful refresh time
- ⚠️ Failure count (alerts if 3+)
- 📁 Configured parent folder

#### 5. Warnings & Alerts

- **Yellow Warning**: Shows when 3+ consecutive failures occur
- **Error Messages**: Displays last error if refresh fails
- **Success Messages**: Confirms successful operations

## How to Use

### Initial Setup

1. Navigate to Admin Panel → Settings
2. Click on the "Dropbox" tab
3. Enter your Dropbox App credentials:
   - Get these from https://www.dropbox.com/developers/apps
   - Create an app or use existing one
4. Follow the OAuth flow (Step 1 → Step 2)
5. Done! Token will auto-refresh every ~4 hours

### Reconnecting

If the connection fails:
1. A yellow warning will appear
2. Click "Actions" → "Disconnect"
3. Re-enter credentials
4. Follow OAuth flow again

### Testing

To verify everything works:
1. Click "Actions" button
2. Click "Test Connection"
3. Success message confirms Dropbox is accessible

## File Structure

```
frontend/src/
├── types/
│   └── dropbox.ts                          # TypeScript types
├── services/
│   └── dropboxAdminService.ts              # API client
└── views/
    └── admin/
        ├── SystemSettings.vue              # Main page with tabs
        └── settings/
            ├── DropboxSettings.vue         # Dropbox config
            ├── EmailSettings.vue           # Placeholder
            ├── SecuritySettings.vue        # Placeholder
            └── GeneralSettings.vue         # Placeholder
```

## API Endpoints Used

All endpoints require super admin authentication:

- `GET /api/admin/dropbox/status` - Get connection status
- `POST /api/admin/dropbox/authorize` - Initiate OAuth
- `POST /api/admin/dropbox/callback` - Complete OAuth
- `POST /api/admin/dropbox/refresh` - Force refresh
- `POST /api/admin/dropbox/test` - Test connection
- `DELETE /api/admin/dropbox/configuration` - Disconnect

## Security

- ✅ All tokens encrypted in database (AES-256-GCM)
- ✅ Super admin only access
- ✅ Never exposes refresh tokens in UI
- ✅ Audit logging for all actions
- ✅ IP address tracking

## Future Settings

Additional tabs will be added for:
- **Email**: SMTP configuration, notification templates
- **Security**: Password policies, 2FA settings
- **General**: App name, logo, default preferences

## Troubleshooting

### "Failed to load Dropbox status"
**Cause:** Backend not running or not authenticated  
**Fix:** Ensure backend is running and you're logged in as super admin

### "Invalid authorization code"
**Cause:** Code expired or already used  
**Fix:** Get a new authorization URL and code

### "Connection test failed"
**Cause:** Token expired or network issue  
**Fix:** Try "Force Refresh" or reconnect

### OAuth flow not working
**Cause:** Incorrect app credentials  
**Fix:** Verify App Key and Secret from Dropbox console

## Tips

1. **Keep credentials safe**: Store Dropbox App Secret securely
2. **Test after setup**: Always click "Test Connection" after configuring
3. **Monitor failures**: Check the failure count regularly
4. **Auto-refresh**: System refreshes tokens automatically, no manual intervention needed
5. **Reconnect if needed**: Only required if tokens are manually revoked

---

**Need Help?** Check the backend documentation at `backend/DROPBOX_QUICKSTART.md`


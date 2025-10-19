# Dropbox OAuth - Quick Start Guide

## 🚀 5-Minute Setup

### 1. Generate Encryption Key
```bash
openssl rand -base64 32
```
Copy the output.

### 2. Update .env
```bash
ENCRYPTION_KEY=paste-your-key-here
```

### 3. Restart Backend
```bash
# Stop and restart your Go server
go run cmd/api/main.go
```

### 4. Get Dropbox App Credentials
1. Visit: https://www.dropbox.com/developers/apps
2. Click "Create app"
3. Choose: "Scoped access" → "Full Dropbox" → Name it
4. Go to Settings tab
5. Copy:
   - **App key**
   - **App secret**

### 5. Authorize Dropbox (One Time Only)

**Step A: Get your admin token**
```bash
# Login to get your token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@bloodsa.org.za","password":"your-password"}'
  
# Copy the "token" from response
```

**Step B: Start OAuth**
```bash
curl -X POST http://localhost:8080/api/admin/dropbox/authorize \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "appKey": "YOUR_DROPBOX_APP_KEY",
    "appSecret": "YOUR_DROPBOX_APP_SECRET",
    "parentFolder": "/SOPS"
  }'
```

**Step C: Authorize in Browser**
1. Copy the `authUrl` from the response
2. Open it in your browser
3. Sign in to Dropbox
4. Click "Allow"
5. **Copy the authorization code** shown on the page

**Step D: Complete Authorization**
```bash
curl -X POST http://localhost:8080/api/admin/dropbox/callback \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "PASTE_AUTHORIZATION_CODE_HERE",
    "appKey": "YOUR_DROPBOX_APP_KEY",
    "appSecret": "YOUR_DROPBOX_APP_SECRET",
    "parentFolder": "/SOPS"
  }'
```

**Expected Response:**
```json
{
  "message": "Dropbox successfully authorized and configured",
  "status": {
    "isConnected": true,
    "tokenExpiry": "2025-10-19T14:30:00Z",
    ...
  }
}
```

### 6. Test It Works
```bash
# Test connection
curl -X POST http://localhost:8080/api/admin/dropbox/test \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# Check status
curl -X GET http://localhost:8080/api/admin/dropbox/status \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"

# Try listing SOPs (should work now!)
curl -X GET "http://localhost:8080/api/sops/categories" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

## ✅ You're Done!

Your Dropbox integration now:
- ✅ Automatically refreshes tokens every ~4 hours
- ✅ Stores everything securely in MongoDB
- ✅ Never expires (unless you revoke it manually)
- ✅ Alerts you if something goes wrong

## 🔧 Common Issues

### "Failed to initialize encryption service"
**Fix:** Make sure `ENCRYPTION_KEY` is set in `.env`

### "Dropbox not configured"
**Fix:** Complete the authorization steps above

### "Invalid authorization code"
**Fix:** The code expires quickly. Get a new one by repeating Step B-D

### "Insufficient permissions"
**Fix:** Make sure you're using a super admin account

## 📖 More Info

- **Full Documentation:** `DROPBOX_OAUTH_IMPLEMENTATION.md`
- **Implementation Details:** `IMPLEMENTATION_SUMMARY.md`

## 🆘 Still Having Issues?

1. Check backend logs for detailed errors
2. Verify your Dropbox app settings
3. Make sure MongoDB is running
4. Ensure you're using super admin credentials

---

**That's it! The expired token error is now permanently solved! 🎉**


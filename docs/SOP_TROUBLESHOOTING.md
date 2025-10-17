# SOP Management Troubleshooting Guide

## Issue: Page Shows Loading Icon Forever

### Quick Fix Steps

1. **Open Browser Developer Console** (F12 or Right-click → Inspect → Console tab)
   - Look for any red error messages
   - Check what the actual error is

2. **Check if you're logged in as Super Admin**
   - Only super admins can view the SOP Management page
   - Regular users can only view public SOP pages

3. **Hard Refresh the Page**
   - Press `Ctrl+Shift+R` (Windows/Linux) or `Cmd+Shift+R` (Mac)
   - This clears the cache and reloads everything

4. **Check Backend Logs**
   ```bash
   docker-compose logs backend | tail -50
   ```
   Look for any errors related to `/api/sops/categories`

5. **Restart Containers**
   ```bash
   docker-compose restart frontend backend
   ```

---

## Understanding "Seed Categories"

**What it does**: The "Seed Initial Categories" button automatically creates 4 pre-configured SOP categories so you don't have to create them manually:

1. **Anemia** - For anemia-related procedures
2. **Lymphoma** - For lymphoma management
3. **Myeloma** - For myeloma treatment
4. **General Business** - For business procedures

**When to use it**:
- First time setting up the system
- After database is reset
- When you have 0 categories

**Note**: It's idempotent - clicking it multiple times won't create duplicates. If categories already exist, it will just tell you how many exist.

---

## Common Issues and Solutions

### Issue: Buttons Don't Work

**Cause**: Frontend JavaScript error or page still loading

**Solution**:
1. Check browser console for errors
2. Ensure you have a valid auth token (try logging out and back in)
3. Clear browser cache and reload

### Issue: "Failed to load categories" Error

**Cause**: Backend not responding or authentication issue

**Solution**:
1. Check if backend is running:
   ```bash
   curl http://localhost:8080/health
   ```
   Should return: `{"message":"It's healthy"}`

2. Check if you're authenticated - look in browser console for:
   ```
   Authorization: Bearer xxx...
   ```

3. Verify CORS settings in `backend/internal/server/routes.go`:
   ```go
   AllowOrigins: []string{"http://localhost:5173"}
   ```

### Issue: Categories Load But Show as Empty

**Cause**: No categories have been created yet (this is normal!)

**Solution**:
- Click "Seed Initial Categories" button
- OR click "Create Custom Category" to make your own

### Issue: Can't Upload Images

**Cause**: Upload directory doesn't exist or permissions issue

**Solution**:
1. Create upload directory:
   ```bash
   mkdir -p backend/uploads/sops
   chmod 755 backend/uploads
   ```

2. Restart backend:
   ```bash
   docker-compose restart backend
   ```

---

## Testing the API Directly

### 1. Get Your Auth Token

Open browser console on the admin page and run:
```javascript
localStorage.getItem('token')
```

Copy the token (everything after "Bearer " if present)

### 2. Test List Categories Endpoint

```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  http://localhost:8080/api/sops/categories?limit=100
```

**Expected Response**:
```json
{
  "categories": [],
  "total": 0,
  "page": 1,
  "limit": 100
}
```

### 3. Test Seed Endpoint

```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  http://localhost:8080/api/sops/seed
```

**Expected Response**:
```json
{
  "message": "categories seeded successfully",
  "count": 4,
  "categories": [...]
}
```

---

## Debugging Frontend

### Check Network Requests

1. Open browser Developer Tools (F12)
2. Go to "Network" tab
3. Reload the page
4. Look for request to `/api/sops/categories`
5. Click on it to see:
   - Request headers (check Authorization header)
   - Response (check what data was returned)
   - Status code (should be 200)

### Common Network Issues

| Status | Meaning | Solution |
|--------|---------|----------|
| 401 | Unauthorized | Token expired, log out and back in |
| 403 | Forbidden | Not a super admin |
| 404 | Not Found | Backend route not registered |
| 500 | Server Error | Check backend logs |
| CORS Error | Cross-origin blocked | Check CORS settings |

---

## Fresh Start (Nuclear Option)

If nothing works, reset everything:

```bash
# Stop containers
docker-compose down

# Remove volumes (WARNING: deletes all data!)
docker volume rm $(docker volume ls -q | grep doctors_workspace)

# Rebuild and start
docker-compose up -d --build

# Watch logs
docker-compose logs -f
```

Then:
1. Login as super admin
2. Go to `/admin/sops`
3. Click "Seed Initial Categories"

---

## Contact/Support

If issues persist:

1. **Check browser console** - screenshot any errors
2. **Check backend logs**:
   ```bash
   docker-compose logs backend > backend_logs.txt
   ```
3. **Check frontend logs**:
   ```bash
   docker-compose logs frontend > frontend_logs.txt
   ```

4. Provide:
   - Browser console errors
   - Backend logs (relevant parts)
   - Steps to reproduce
   - What you expected vs what happened


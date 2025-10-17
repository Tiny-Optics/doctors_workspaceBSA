# SOP Backend Testing Guide

## Prerequisites

1. Backend server running on `http://localhost:8080`
2. MongoDB running and connected
3. Valid JWT tokens for:
   - Regular user (haematologist, physician, or data_capturer)
   - Super admin user
4. Dropbox configured with valid credentials

## Getting Authentication Tokens

### 1. Login as Super Admin
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "your_password"
  }'
```

Save the `token` from the response:
```bash
export ADMIN_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 2. Login as Regular User
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "doctor@example.com",
    "password": "your_password"
  }'
```

Save the token:
```bash
export USER_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Test Scenarios

### Test 1: Create Categories (Admin Only)

#### 1.1 Create "Anemia" Category
```bash
curl -X POST http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Anemia",
    "description": "Standard operating procedures for anemia diagnosis and treatment",
    "imagePath": "/uploads/sops/anemia.jpg",
    "displayOrder": 1
  }'
```

**Expected Result:**
- Status: `201 Created`
- Returns category object with generated `id` and `slug`
- Dropbox folder created at `/home/Doctors_Workspace/SOPS/Anemia`

Save the category ID:
```bash
export ANEMIA_ID="507f1f77bcf86cd799439011"
```

#### 1.2 Create "Lymphoma" Category
```bash
curl -X POST http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Lymphoma",
    "description": "Lymphoma treatment protocols",
    "imagePath": "/uploads/sops/lymphoma.jpg",
    "displayOrder": 2
  }'
```

#### 1.3 Create "Myeloma" Category
```bash
curl -X POST http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Myeloma",
    "description": "Multiple myeloma management procedures",
    "imagePath": "/uploads/sops/myeloma.jpg",
    "displayOrder": 3
  }'
```

#### 1.4 Create "General Business" Category
```bash
curl -X POST http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "General Business",
    "description": "General business procedures and policies",
    "imagePath": "/uploads/sops/general-business.jpg",
    "displayOrder": 4
  }'
```

#### 1.5 Test Create as Regular User (Should Fail)
```bash
curl -X POST http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Should Fail",
    "displayOrder": 5
  }'
```

**Expected Result:**
- Status: `403 Forbidden`
- Error message about insufficient permissions

### Test 2: List Categories

#### 2.1 List as Admin (See All)
```bash
curl http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**Expected Result:**
- Status: `200 OK`
- Returns all categories (active and inactive)
- Sorted by `displayOrder`

#### 2.2 List as Regular User (See Active Only)
```bash
curl http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Expected Result:**
- Status: `200 OK`
- Returns only active categories

#### 2.3 Search Categories
```bash
curl "http://localhost:8080/api/sops/categories?search=anemia" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Expected Result:**
- Returns categories matching "anemia" in name or description

#### 2.4 Paginated List
```bash
curl "http://localhost:8080/api/sops/categories?page=1&limit=2" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Expected Result:**
- Returns first 2 categories
- Includes `total`, `page`, and `limit` in response

### Test 3: Get Single Category

#### 3.1 Get by ID
```bash
curl http://localhost:8080/api/sops/categories/$ANEMIA_ID \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Expected Result:**
- Status: `200 OK`
- Returns full category details

#### 3.2 Get with Invalid ID
```bash
curl http://localhost:8080/api/sops/categories/invalid-id \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Expected Result:**
- Status: `400 Bad Request`
- Error: "invalid category ID"

#### 3.3 Get Non-existent Category
```bash
curl http://localhost:8080/api/sops/categories/507f1f77bcf86cd799439999 \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Expected Result:**
- Status: `404 Not Found`

### Test 4: Update Category (Admin Only)

#### 4.1 Update Description
```bash
curl -X PUT http://localhost:8080/api/sops/categories/$ANEMIA_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description for anemia procedures"
  }'
```

**Expected Result:**
- Status: `200 OK`
- Returns updated category
- `updatedAt` timestamp changed

#### 4.2 Update Name (Test Dropbox Rename)
```bash
curl -X PUT http://localhost:8080/api/sops/categories/$ANEMIA_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Anaemia"
  }'
```

**Expected Result:**
- Status: `200 OK`
- `slug` changed to "anaemia"
- `dropboxPath` updated to "SOPS/Anaemia"
- Dropbox folder renamed from "Anemia" to "Anaemia"

#### 4.3 Deactivate Category
```bash
curl -X PUT http://localhost:8080/api/sops/categories/$ANEMIA_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "isActive": false
  }'
```

**Expected Result:**
- Category no longer visible to regular users
- Still visible to admins

#### 4.4 Reactivate Category
```bash
curl -X PUT http://localhost:8080/api/sops/categories/$ANEMIA_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "isActive": true
  }'
```

#### 4.5 Test Update as Regular User (Should Fail)
```bash
curl -X PUT http://localhost:8080/api/sops/categories/$ANEMIA_ID \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Should fail"
  }'
```

**Expected Result:**
- Status: `403 Forbidden`

### Test 5: File Operations

**Note:** First manually upload some test PDFs to Dropbox folder `/home/Doctors_Workspace/SOPS/Anemia/`

Example files to upload:
- `Iron_Deficiency_Anemia.pdf`
- `B12_Deficiency.pdf`
- `Aplastic_Anemia.pdf`

#### 5.1 List Files in Category
```bash
curl http://localhost:8080/api/sops/categories/$ANEMIA_ID/files \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Expected Result:**
- Status: `200 OK`
- Returns array of files with metadata (name, size, modifiedTime)

#### 5.2 Get Download Link
```bash
curl "http://localhost:8080/api/sops/categories/$ANEMIA_ID/files/download?path=Iron_Deficiency_Anemia.pdf" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Expected Result:**
- Status: `200 OK`
- Returns `downloadLink` (Dropbox temporary URL)
- Link valid for 4 hours

#### 5.3 Test Download Link
Copy the download link from the previous response and open in browser:
```bash
# Copy the downloadLink from response
curl "https://dl.dropboxusercontent.com/apitl/..."
```

**Expected Result:**
- PDF file downloads successfully

#### 5.4 Get Download Link for Non-existent File
```bash
curl "http://localhost:8080/api/sops/categories/$ANEMIA_ID/files/download?path=NonExistent.pdf" \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Expected Result:**
- Status: `404 Not Found`

### Test 6: Delete Category (Admin Only)

#### 6.1 Create Test Category
```bash
curl -X POST http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Category",
    "displayOrder": 99
  }'
```

Save the ID:
```bash
export TEST_ID="..."
```

#### 6.2 Delete Category
```bash
curl -X DELETE http://localhost:8080/api/sops/categories/$TEST_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**Expected Result:**
- Status: `200 OK`
- Message: "category deleted successfully"
- Category removed from database
- Dropbox folder remains intact

#### 6.3 Verify Deletion
```bash
curl http://localhost:8080/api/sops/categories/$TEST_ID \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Expected Result:**
- Status: `404 Not Found`

#### 6.4 Check Dropbox Folder Still Exists
Manually check Dropbox or use Dropbox API to verify folder `/home/Doctors_Workspace/SOPS/Test Category` still exists.

### Test 7: Validation

#### 7.1 Create with Empty Name
```bash
curl -X POST http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "",
    "displayOrder": 1
  }'
```

**Expected Result:**
- Status: `400 Bad Request`
- Validation error for name

#### 7.2 Create with Long Description
```bash
curl -X POST http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Test\",
    \"description\": \"$(printf 'a%.0s' {1..1001})\",
    \"displayOrder\": 1
  }"
```

**Expected Result:**
- Status: `400 Bad Request`
- Error about description length

#### 7.3 Create with Invalid Image Format
```bash
curl -X POST http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test",
    "imagePath": "/uploads/test.pdf",
    "displayOrder": 1
  }'
```

**Expected Result:**
- Status: `400 Bad Request`
- Error about invalid image format

#### 7.4 Create Duplicate Name
```bash
# Try to create another "Anemia" category
curl -X POST http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Anemia",
    "displayOrder": 10
  }'
```

**Expected Result:**
- Status: `409 Conflict`
- Error about duplicate category

## Verification Checklist

After running all tests, verify:

### Database
- [ ] Categories created in `sop_categories` collection
- [ ] Indexes created (slug, is_active, display_order)
- [ ] Audit logs created for admin actions

### Dropbox
- [ ] Folders created in `/home/Doctors_Workspace/SOPS/`
- [ ] Folder renamed when category name changed
- [ ] Folder remains after category deletion

### Permissions
- [ ] Regular users can view active categories
- [ ] Regular users can list and download files
- [ ] Regular users cannot create/update/delete categories
- [ ] Admins can see all categories (active and inactive)
- [ ] Admins can perform all operations

### API Responses
- [ ] Proper status codes (200, 201, 400, 403, 404, 409)
- [ ] Consistent error message format
- [ ] Pagination works correctly
- [ ] Search filters work

### Edge Cases
- [ ] Empty file list returns empty array (not error)
- [ ] Invalid IDs return 400
- [ ] Non-existent resources return 404
- [ ] Unauthorized actions return 403

## Troubleshooting

### Issue: "Dropbox is not configured"
**Solution:** Check environment variables:
```bash
echo $DROPBOX_APP_API_ACCESS_TOKEN
echo $DROPBOX_APP_PARENT_FOLDER
```

### Issue: "Category not found"
**Solution:** Verify category ID is correct:
```bash
curl http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.categories[].id'
```

### Issue: "Unauthorized"
**Solution:** Check token is valid:
```bash
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $USER_TOKEN"
```

### Issue: Files not listing
**Solution:** 
1. Verify files exist in Dropbox folder
2. Check folder path matches category's `dropboxPath`
3. Verify Dropbox API access token has file read permissions

## Performance Testing

### Load Test - List Categories
```bash
# Using Apache Bench
ab -n 100 -c 10 -H "Authorization: Bearer $USER_TOKEN" \
  http://localhost:8080/api/sops/categories
```

### Load Test - Download Links
```bash
# Generate 100 download links
for i in {1..100}; do
  curl "http://localhost:8080/api/sops/categories/$ANEMIA_ID/files/download?path=test.pdf" \
    -H "Authorization: Bearer $USER_TOKEN"
done
```

## Cleanup

### Remove Test Categories
```bash
# List all categories
curl http://localhost:8080/api/sops/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.categories[] | {id, name}'

# Delete test categories
curl -X DELETE http://localhost:8080/api/sops/categories/$TEST_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## Next Steps

After successful testing:
1. Deploy to staging environment
2. Test with real Dropbox folders and files
3. Implement frontend integration
4. Conduct user acceptance testing
5. Deploy to production


# SOP Backend API Documentation

## Overview

The SOP (Standard Operating Procedures) system allows authenticated users to browse and download SOP files organized by category. Super admins can manage categories through CRUD operations.

## Architecture

- **Model**: `SOPCategory` - Represents a category with metadata and Dropbox integration
- **Repository**: Database operations for categories
- **Service**: Business logic, Dropbox integration, and permission checks
- **Handler**: HTTP request handlers
- **Routes**: RESTful API endpoints

## Database Schema

### Collection: `sop_categories`

```json
{
  "_id": "ObjectId",
  "name": "string",              // Category name (e.g., "Anemia")
  "slug": "string",              // URL-friendly slug (e.g., "anemia")
  "description": "string",       // Optional description
  "image_path": "string",        // Local path to category image
  "dropbox_path": "string",      // Path in Dropbox (e.g., "SOPS/Anemia")
  "display_order": "number",     // For sorting categories
  "is_active": "boolean",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "created_by": "ObjectId"       // Reference to admin who created it
}
```

**Indexes:**
- `slug` (unique)
- `is_active`
- `display_order`
- `created_at`

## API Endpoints

Base URL: `/api/sops`

All endpoints require authentication (Bearer token).

### 1. List Categories

**GET** `/api/sops/categories`

Lists all SOP categories. Non-admin users only see active categories.

**Query Parameters:**
- `search` (optional) - Search term to filter by name or description
- `page` (optional) - Page number (default: 1)
- `limit` (optional) - Items per page (default: 20, max: 100)

**Response:**
```json
{
  "categories": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Anemia",
      "slug": "anemia",
      "description": "Standard operating procedures for anemia treatment",
      "imagePath": "/uploads/sops/anemia.jpg",
      "dropboxPath": "SOPS/Anemia",
      "displayOrder": 1,
      "isActive": true,
      "createdAt": "2025-10-17T10:00:00Z",
      "updatedAt": "2025-10-17T10:00:00Z",
      "createdBy": "507f1f77bcf86cd799439012"
    }
  ],
  "total": 4,
  "page": 1,
  "limit": 20
}
```

### 2. Get Category by ID

**GET** `/api/sops/categories/:id`

Retrieves a specific category by ID.

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Anemia",
  "slug": "anemia",
  "description": "Standard operating procedures for anemia treatment",
  "imagePath": "/uploads/sops/anemia.jpg",
  "dropboxPath": "SOPS/Anemia",
  "displayOrder": 1,
  "isActive": true,
  "createdAt": "2025-10-17T10:00:00Z",
  "updatedAt": "2025-10-17T10:00:00Z",
  "createdBy": "507f1f77bcf86cd799439012"
}
```

**Errors:**
- `400` - Invalid category ID
- `404` - Category not found

### 3. Get Category Files

**GET** `/api/sops/categories/:id/files`

Lists all files in a category's Dropbox folder.

**Response:**
```json
{
  "files": [
    {
      "name": "Iron_Deficiency_Anemia.pdf",
      "path": "/Doctors_Workspace/SOPS/Anemia/Iron_Deficiency_Anemia.pdf",
      "size": 1048576,
      "modifiedTime": "2025-10-15T14:30:00Z",
      "isFolder": false
    },
    {
      "name": "B12_Deficiency.pdf",
      "path": "/Doctors_Workspace/SOPS/Anemia/B12_Deficiency.pdf",
      "size": 2097152,
      "modifiedTime": "2025-10-16T09:15:00Z",
      "isFolder": false
    }
  ]
}
```

**Errors:**
- `400` - Invalid category ID
- `404` - Category not found

### 4. Get File Download Link

**GET** `/api/sops/categories/:id/files/download?path=filename.pdf`

Generates a temporary download link for a specific file (valid for 4 hours).

**Query Parameters:**
- `path` (required) - File name or relative path within the category folder

**Response:**
```json
{
  "downloadLink": "https://dl.dropboxusercontent.com/apitl/..."
}
```

**Errors:**
- `400` - Invalid category ID or missing file path
- `404` - Category or file not found

### 5. Create Category (Admin Only)

**POST** `/api/sops/categories`

Creates a new SOP category. Requires super admin permission.

**Request Body:**
```json
{
  "name": "General Business",
  "description": "General business procedures and policies",
  "imagePath": "/uploads/sops/general-business.jpg",
  "displayOrder": 4
}
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "General Business",
  "slug": "general-business",
  "description": "General business procedures and policies",
  "imagePath": "/uploads/sops/general-business.jpg",
  "dropboxPath": "SOPS/General Business",
  "displayOrder": 4,
  "isActive": true,
  "createdAt": "2025-10-17T10:00:00Z",
  "updatedAt": "2025-10-17T10:00:00Z",
  "createdBy": "507f1f77bcf86cd799439012"
}
```

**Notes:**
- Slug is automatically generated from the name
- A Dropbox folder is automatically created at `/Doctors_Workspace/SOPS/{name}`
- Action is logged in audit trail

**Errors:**
- `400` - Validation error (invalid name, description too long, invalid image format)
- `403` - Insufficient permissions (not a super admin)
- `409` - Category with this name already exists

### 6. Update Category (Admin Only)

**PUT** `/api/sops/categories/:id`

Updates an existing category. Requires super admin permission.

**Request Body:**
All fields are optional:
```json
{
  "name": "Anaemia",
  "description": "Updated description",
  "imagePath": "/uploads/sops/anaemia-updated.jpg",
  "displayOrder": 2,
  "isActive": false
}
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "Anaemia",
  "slug": "anaemia",
  "description": "Updated description",
  "imagePath": "/uploads/sops/anaemia-updated.jpg",
  "dropboxPath": "SOPS/Anaemia",
  "displayOrder": 2,
  "isActive": false,
  "createdAt": "2025-10-17T10:00:00Z",
  "updatedAt": "2025-10-17T11:00:00Z",
  "createdBy": "507f1f77bcf86cd799439012"
}
```

**Notes:**
- If name is changed, the slug and Dropbox path are automatically updated
- Dropbox folder is renamed (files remain intact)
- Action is logged in audit trail

**Errors:**
- `400` - Invalid category ID or validation error
- `403` - Insufficient permissions
- `404` - Category not found
- `409` - New name conflicts with existing category

### 7. Delete Category (Admin Only)

**DELETE** `/api/sops/categories/:id`

Deletes a category from the database. **Note: Dropbox folder and files are NOT deleted.**

Requires super admin permission.

**Response:**
```json
{
  "message": "category deleted successfully"
}
```

**Notes:**
- Only removes database entry
- Dropbox folder and all files remain intact
- Action is logged in audit trail

**Errors:**
- `400` - Invalid category ID
- `403` - Insufficient permissions
- `404` - Category not found

## Permissions

- **All Authenticated Users**: Can view active categories, list files, and download files
- **Super Admins** (users with `PermDeleteUsers`): Can create, update, and delete categories

## Dropbox Integration

### Configuration

Required environment variables:
```
DROPBOX_APP_API_ACCESS_TOKEN=your_token_here
DROPBOX_APP_PARENT_FOLDER=/home/Doctors_Workspace
```

### Folder Structure

Categories are stored in Dropbox with this structure:
```
/home/Doctors_Workspace/
└── SOPS/
    ├── Anemia/
    │   ├── Iron_Deficiency_Anemia.pdf
    │   └── B12_Deficiency.pdf
    ├── Lymphoma/
    │   └── Hodgkin_Lymphoma.pdf
    ├── Myeloma/
    │   └── Multiple_Myeloma.pdf
    └── General Business/
        └── Business_Procedures.pdf
```

### Operations

1. **Create Category**: Automatically creates folder in Dropbox
2. **Rename Category**: Automatically renames Dropbox folder
3. **Delete Category**: Database entry removed, Dropbox folder remains
4. **List Files**: Reads directly from Dropbox folder
5. **Download File**: Generates temporary signed URL (4-hour expiry)

### Error Handling

- If Dropbox is not configured, endpoints return appropriate error
- If folder doesn't exist when listing files, returns empty array
- Failed folder operations are logged but don't fail the database operation

## Image Storage

Category images are stored locally on the server. The `imagePath` field contains the relative path to the image file.

**Valid image formats:**
- `.jpg`, `.jpeg`
- `.png`
- `.webp`

**Example:**
```
/uploads/sops/anemia.jpg
```

## Validation Rules

### Category Name
- Required
- 1-100 characters
- Unique (slug-based)

### Description
- Optional
- Maximum 1000 characters

### Image Path
- Optional
- Must end with valid image extension

### Display Order
- Must be >= 0
- Used for sorting categories

## Example Usage

### cURL Examples

**List all categories:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/sops/categories
```

**Get category files:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/sops/categories/507f1f77bcf86cd799439011/files
```

**Download a file:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/api/sops/categories/507f1f77bcf86cd799439011/files/download?path=Iron_Deficiency_Anemia.pdf"
```

**Create category (admin):**
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Lymphoma",
    "description": "Lymphoma treatment procedures",
    "displayOrder": 2
  }' \
  http://localhost:8080/api/sops/categories
```

## Testing Checklist

- [ ] List categories as regular user (only see active)
- [ ] List categories as admin (see all)
- [ ] Get category by ID
- [ ] Get files from category
- [ ] Generate download link for file
- [ ] Create category as admin
- [ ] Create category as non-admin (should fail)
- [ ] Update category name (check Dropbox rename)
- [ ] Deactivate category (check visibility)
- [ ] Delete category (check Dropbox folder remains)
- [ ] Search categories by name
- [ ] Test pagination
- [ ] Test with invalid category ID
- [ ] Test with non-existent file

## Audit Logging

All admin actions are logged with:
- Action type: `sop_category.create`, `sop_category.update`, `sop_category.delete`
- User who performed action
- Category details
- Timestamp
- IP address

## Future Enhancements

- [ ] File upload endpoint
- [ ] Bulk category import
- [ ] Category statistics (view count, download count)
- [ ] File search across categories
- [ ] Version control for files
- [ ] Category access restrictions per user role


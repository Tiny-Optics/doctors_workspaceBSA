# SMTP-Only Configuration Endpoints

## Overview
Dedicated API endpoints for managing SMTP email settings independently from the full registry configuration. These endpoints are designed for the EmailSettings.vue admin page.

## Endpoints

### 1. Get SMTP Configuration
**GET** `/api/admin/registry/smtp-config`

**Description:** Retrieve current SMTP configuration (admin only)

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "host": "smtp.gmail.com",
  "port": 587,
  "username": "noreply@bloodsa.org.za",
  "fromEmail": "noreply@bloodsa.org.za",
  "fromName": "BLOODSA Registry",
  "isComplete": true
}
```

**Status Codes:**
- `200` - Success
- `401` - Unauthorized
- `500` - Internal Server Error

### 2. Update SMTP Configuration
**PUT** `/api/admin/registry/smtp-config`

**Description:** Update SMTP configuration (admin only)

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "host": "smtp.gmail.com",
  "port": 587,
  "username": "noreply@bloodsa.org.za",
  "password": "your-app-password",
  "fromEmail": "noreply@bloodsa.org.za",
  "fromName": "BLOODSA Registry"
}
```

**Response:**
```json
{
  "host": "smtp.gmail.com",
  "port": 587,
  "username": "noreply@bloodsa.org.za",
  "fromEmail": "noreply@bloodsa.org.za",
  "fromName": "BLOODSA Registry",
  "isComplete": true
}
```

**Status Codes:**
- `200` - Success
- `400` - Bad Request (validation error)
- `401` - Unauthorized
- `403` - Forbidden (insufficient permissions)
- `500` - Internal Server Error

### 3. Test Email (Existing)
**POST** `/api/admin/registry/test-email`

**Description:** Send test email using current SMTP configuration

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "test@example.com"
}
```

**Response:**
```json
{
  "message": "test email sent successfully"
}
```

## Features

### 🔐 **Security**
- **Password Encryption:** SMTP passwords are encrypted before storage
- **Permission Control:** Requires `PermManageSystem` permission
- **Audit Logging:** All SMTP changes are logged with user details

### 🎯 **Independence**
- **Separate from Registry Config:** SMTP settings can be managed independently
- **Partial Updates:** Only send the fields you want to update
- **Backward Compatibility:** Existing registry configuration remains unchanged

### 📊 **Data Management**
- **Singleton Pattern:** Only one SMTP configuration exists
- **Upsert Operations:** Creates new config if none exists
- **Validation:** Server-side validation of all SMTP fields

## Usage Examples

### Get Current SMTP Settings
```bash
curl -X GET "http://localhost:8080/api/admin/registry/smtp-config" \
  -H "Authorization: Bearer <your-token>"
```

### Update SMTP Settings
```bash
curl -X PUT "http://localhost:8080/api/admin/registry/smtp-config" \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "host": "smtp.gmail.com",
    "port": 587,
    "username": "noreply@bloodsa.org.za",
    "password": "your-app-password",
    "fromEmail": "noreply@bloodsa.org.za",
    "fromName": "BLOODSA Registry"
  }'
```

### Test Email Configuration
```bash
curl -X POST "http://localhost:8080/api/admin/registry/test-email" \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@bloodsa.org.za"
  }'
```

## Implementation Details

### Backend Components
- **Models:** `UpdateSMTPConfigRequest`, `SMTPConfigResponse`
- **Service:** `GetSMTPConfig()`, `UpdateSMTPConfig()`
- **Handler:** `GetSMTPConfig()`, `UpdateSMTPConfig()`
- **Routes:** `/api/admin/registry/smtp-config` (GET, PUT)

### Database Schema
The SMTP configuration is stored within the existing `registry_config` collection as part of the `SMTPConfig` embedded document, maintaining data consistency while providing independent access.

### Audit Trail
All SMTP configuration changes are logged with:
- User who made the change
- Timestamp
- IP address
- Action type: `smtp_config_updated`

# REDCap Referral System - Backend API Endpoints

## Overview
The REDCap Referral system allows admins to configure a referral link and users to access it with audit logging.

## Endpoints

### User Endpoints (Authenticated Users)

#### 1. Get Referral Configuration
**GET** `/api/referrals/config`

**Authentication:** Required (Bearer token)

**Response:**
```json
{
  "isConfigured": true,
  "isEnabled": true,
  "redcapUrl": "https://redcap.example.com/survey"
}
```

#### 2. Log Referral Access
**POST** `/api/referrals/access`

**Authentication:** Required (Bearer token)

**Description:** Logs the user's access to the referral link and returns the URL to redirect to.

**Response:**
```json
{
  "redirectUrl": "https://redcap.example.com/survey"
}
```

### Admin Endpoints (Super Admin Only)

#### 3. Get Full Referral Configuration
**GET** `/api/admin/referrals/config`

**Authentication:** Required (Bearer token + PermManageSystem)

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "redcapUrl": "https://redcap.example.com/survey",
  "isEnabled": true,
  "createdAt": "2025-10-22T10:00:00Z",
  "updatedAt": "2025-10-22T14:30:00Z",
  "updatedBy": "507f1f77bcf86cd799439012"
}
```

#### 4. Update Referral Configuration
**PUT** `/api/admin/referrals/config`

**Authentication:** Required (Bearer token + PermManageSystem)

**Request Body:**
```json
{
  "redcapUrl": "https://redcap.example.com/survey",
  "isEnabled": true
}
```

**Response:** Same as Get Full Configuration

## Database Model

### Collection: `referral_config`
- Singleton pattern (only one document)
- Fields:
  - `_id`: ObjectID
  - `redcap_url`: String (validated URL)
  - `is_enabled`: Boolean
  - `created_at`: DateTime
  - `updated_at`: DateTime
  - `updated_by`: ObjectID (reference to user)

## Audit Logging

### New Audit Actions:
1. `referral_config_updated` - When admin updates the configuration
2. `referral_accessed` - When user accesses the referral link

## Testing Commands

### 1. Login as Admin
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@bloodsa.org.za","password":"password123"}'
```

### 2. Get Referral Config (User)
```bash
curl -X GET http://localhost:8080/api/referrals/config \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. Update Referral Config (Admin)
```bash
curl -X PUT http://localhost:8080/api/admin/referrals/config \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "redcapUrl": "https://redcap.example.com/surveys/12345",
    "isEnabled": true
  }'
```

### 4. Get Full Config (Admin)
```bash
curl -X GET http://localhost:8080/api/admin/referrals/config \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 5. Log Access and Get URL (User)
```bash
curl -X POST http://localhost:8080/api/referrals/access \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Next Steps

1. Test all endpoints with the backend running
2. Implement frontend components
3. Test end-to-end flow


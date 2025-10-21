# African HOPeR Registry Backend Implementation

## Overview

Complete backend implementation for the African HOPeR Registry section with custom form builder, document management, and email notifications.

## Implemented Components

### 1. Database Models (`backend/internal/models/`)

#### `registry_config.go`
- **RegistryConfig** - Singleton configuration model
  - Video URL for training videos
  - Documents path for example documents in Dropbox
  - Notification emails list
  - SMTP configuration (with encrypted password)
- **SMTPConfig** - Email server configuration
- **UpdateRegistryConfigRequest** - Request model for configuration updates
- Full validation methods

#### `registry_form.go`
- **FormFieldType** enum - Supports: text, textarea, select, radio, date, number, email, file
- **FormField** - Custom field definition with validation rules
- **RegistryFormSchema** - Complete form schema with multiple fields
- **CreateFormSchemaRequest** / **UpdateFormSchemaRequest** - Request models
- Support for multiple file uploads via `allow_multiple` flag
- Validation ensures at least one file field exists

#### `registry_submission.go`
- **SubmissionStatus** enum - submitted, reviewed, approved, rejected
- **RegistrySubmission** - User submission with form data and uploaded documents
- **CreateSubmissionRequest** / **UpdateSubmissionStatusRequest** - Request models
- Tracks documents path and uploaded file names

### 2. Repositories (`backend/internal/repository/`)

#### `registry_config_repository.go`
- Singleton pattern for configuration management
- `GetConfig()` - Retrieve configuration
- `CreateOrUpdate()` - Upsert configuration
- `AddNotificationEmail()` / `RemoveNotificationEmail()` - Manage email list

#### `registry_form_repository.go`
- Complete CRUD operations for form schemas
- `FindActive()` - Get currently active form
- `SetActive()` - Atomic activation (deactivates all others)
- `List()` - Paginated listing with filters

#### `registry_submission_repository.go`
- CRUD operations for submissions
- `FindByUser()` - Get user's submissions with pagination
- `List()` - Admin view with filters
- `UpdateStatus()` - Change submission status with review notes
- `CountByStatus()` / `CountByUser()` - Statistics

### 3. Services (`backend/internal/service/`)

#### `email_service.go`
- Email notification system using gomail
- `SendSubmissionNotification()` - Notify admins of new submissions
  - Professional HTML email template
  - Includes user info, form details, and Dropbox link
- `SendTestEmail()` - Verify SMTP configuration
- Automatic password decryption using existing EncryptionService

#### `registry_service.go`
Complete business logic layer with:

**Configuration Management:**
- `GetConfiguration()` - Full config (admin only)
- `GetPublicConfiguration()` - Public data only (video, documents)
- `UpdateConfiguration()` - Update with automatic password encryption

**Form Schema Management:**
- `CreateFormSchema()` - Create new form with validation
- `UpdateFormSchema()` - Update existing form
- `GetFormSchema()` / `GetActiveFormSchema()` - Retrieve forms
- `ListFormSchemas()` - Paginated listing
- `DeleteFormSchema()` - Delete with audit trail

**Submission Management:**
- `ValidatePreSubmission()` - **Critical validation**:
  - Active form exists
  - Notification emails configured
  - SMTP config complete
- `SubmitForm()` - Complete submission flow:
  1. Validate form data against schema
  2. Upload files to Dropbox
  3. Create submission record
  4. Send email notifications
- `GetUserSubmissions()` - User's own submissions
- `GetAllSubmissions()` - Admin view with filters
- `UpdateSubmissionStatus()` - Review submissions

#### `dropbox_service.go` (Updated)
- Added `UploadFile()` method for document uploads
- Supports streaming file uploads to Dropbox

### 4. Handlers (`backend/internal/handlers/`)

#### `registry_handler.go`
Complete REST API with 17 endpoints:

**Public/User Endpoints:**
- `GET /api/registry/config` - Get public configuration
- `GET /api/registry/form-schema` - Get active form schema
- `POST /api/registry/submit` - Submit form with documents (multipart)
- `GET /api/registry/submissions` - Get user's submissions
- `GET /api/registry/submissions/:id` - Get specific submission

**Admin Endpoints (PermManageUsers - Admins & User Managers):**
- `POST /api/admin/registry/form-schema` - Create form schema
- `GET /api/admin/registry/form-schemas` - List all schemas
- `GET /api/admin/registry/form-schema/:id` - Get specific schema
- `PUT /api/admin/registry/form-schema/:id` - Update schema
- `DELETE /api/admin/registry/form-schema/:id` - Delete schema

**Super Admin Endpoints (PermManageSystem):**
- `GET /api/admin/registry/config` - Get full configuration
- `PUT /api/admin/registry/config` - Update configuration
- `GET /api/admin/registry/submissions` - List all submissions
- `PATCH /api/admin/registry/submissions/:id/status` - Update submission status

### 5. Routes (`backend/internal/server/routes.go`)

Integrated all registry routes with proper middleware:
- Authentication required for all endpoints
- Role-based permissions enforced
- Three access levels:
  1. All authenticated users (view, submit)
  2. Admins & User Managers (manage forms)
  3. Super Admins (manage configuration)

## Key Features

### 1. Custom Form Builder
- Admins create completely custom forms
- Supported field types: text, textarea, select, radio, date, number, email, file
- Field-level validation rules
- Required/optional fields
- Multiple file uploads per form

### 2. Pre-Submission Validation
System validates before accepting submissions:
- ✅ Active form schema exists
- ✅ Notification emails configured
- ✅ SMTP configuration complete
- Returns clear error if prerequisites not met

### 3. Document Upload Flow
1. User submits form with multiple documents
2. Backend uploads to Dropbox: `BLDS_approvals/{username}/{submission_id}/`
3. Stores submission with file list in MongoDB
4. Sends email notifications with Dropbox link

### 4. Email Notifications
- Professional HTML email template
- Sent to all configured notification emails
- Includes:
  - User information
  - Form name
  - Submission ID and timestamp
  - Direct link to Dropbox folder
- Graceful error handling (logs but doesn't fail submission)

### 5. Security
- SMTP passwords encrypted in database
- Automatic decryption when sending emails
- Role-based access control
- Audit logging for all actions
- Owner/admin validation for viewing submissions

## Database Collections

### `registry_config` (Singleton)
```json
{
  "_id": ObjectId,
  "video_url": "string",
  "documents_path": "string",
  "notification_emails": ["string"],
  "smtp_config": {
    "host": "string",
    "port": 587,
    "username": "string",
    "password": "encrypted_string",
    "from_email": "string",
    "from_name": "string"
  },
  "created_at": ISODate,
  "updated_at": ISODate,
  "updated_by": ObjectId
}
```

### `registry_form_schemas`
```json
{
  "_id": ObjectId,
  "form_name": "string",
  "description": "string",
  "fields": [
    {
      "id": "string",
      "label": "string",
      "type": "text|textarea|select|radio|date|number|email|file",
      "required": true,
      "placeholder": "string",
      "help_text": "string",
      "options": ["string"],
      "allow_multiple": true,
      "validation_rules": {
        "min_length": 10,
        "max_length": 500,
        "pattern": "regex"
      },
      "display_order": 1
    }
  ],
  "is_active": true,
  "created_at": ISODate,
  "updated_at": ISODate,
  "created_by": ObjectId,
  "updated_by": ObjectId
}
```

### `registry_submissions`
```json
{
  "_id": ObjectId,
  "user_id": ObjectId,
  "form_schema_id": ObjectId,
  "form_data": {
    "field_id": "value"
  },
  "documents_path": "BLDS_approvals/username/submission_id",
  "uploaded_documents": ["file1.pdf", "file2.pdf"],
  "status": "submitted|reviewed|approved|rejected",
  "created_at": ISODate,
  "updated_at": ISODate,
  "reviewed_by": ObjectId,
  "reviewed_at": ISODate,
  "review_notes": "string"
}
```

## API Endpoints Summary

### Public Access (Authenticated Users)
- `GET /api/registry/config` - Public configuration
- `GET /api/registry/form-schema` - Active form
- `POST /api/registry/submit` - Submit with files
- `GET /api/registry/submissions` - User's submissions
- `GET /api/registry/submissions/:id` - View submission

### Admin Access (Admins & User Managers)
- `POST /api/admin/registry/form-schema` - Create form
- `GET /api/admin/registry/form-schemas` - List forms
- `GET /api/admin/registry/form-schema/:id` - Get form
- `PUT /api/admin/registry/form-schema/:id` - Update form
- `DELETE /api/admin/registry/form-schema/:id` - Delete form

### Super Admin Access
- `GET /api/admin/registry/config` - Full config
- `PUT /api/admin/registry/config` - Update config
- `GET /api/admin/registry/submissions` - All submissions
- `PATCH /api/admin/registry/submissions/:id/status` - Update status

## Dependencies Added

- `gopkg.in/gomail.v2` - Email sending library

## Testing Checklist

### Configuration
- [ ] Create/update registry configuration
- [ ] Add/remove notification emails
- [ ] Configure SMTP settings
- [ ] Test email with SendTestEmail

### Form Management
- [ ] Create form schema with various field types
- [ ] Activate/deactivate forms
- [ ] Update form schema
- [ ] Delete form schema
- [ ] Verify only one form can be active

### Submissions
- [ ] Submit form with single file
- [ ] Submit form with multiple files
- [ ] Verify pre-submission validation
- [ ] Check Dropbox upload
- [ ] Verify email notifications sent
- [ ] View own submissions as user
- [ ] View all submissions as admin
- [ ] Update submission status

### Security
- [ ] Verify SMTP password encryption
- [ ] Test role-based permissions
- [ ] Verify users can only see own submissions
- [ ] Verify admins can see all submissions

## Notes

1. **SMTP Configuration**: Stored in database (not environment variables) for admin manageability
2. **File Uploads**: Handled via multipart/form-data with JSON form data
3. **Dropbox Path**: Automatically organized by username and submission ID
4. **Email Failures**: Logged but don't prevent submission from succeeding
5. **Form Activation**: Atomic operation - only one form can be active at a time
6. **Permissions**: 
   - Super Admins: Full access to everything
   - Admins & User Managers: Can manage forms and view submissions
   - All Users: Can view config, submit forms, view own submissions

## Future Enhancements

- File type validation (PDF, DOC, etc.)
- File size limits
- Form field dependencies (conditional fields)
- Submission edit/update capability
- Email templates customization
- Submission statistics dashboard
- Export submissions to CSV
- Dropbox shared link generation


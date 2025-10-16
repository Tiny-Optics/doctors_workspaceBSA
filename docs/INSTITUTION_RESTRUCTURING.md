# Institution Restructuring

## Overview

Institutions have been restructured from a simple string field to a separate, normalized entity with proper relationships. This provides better data integrity, easier management, and scalability.

## Changes Made

### 1. New Institution Model

**File**: `backend/internal/models/institution.go`

```go
type Institution struct {
    ID          primitive.ObjectID
    Name        string              // Full institution name
    ShortName   string              // Abbreviated name (e.g., "UCT")
    Type        InstitutionType     // Type of institution
    Country     string              // Country
    Province    string              // Province/State
    City        string              // City
    Address     string              // Physical address
    PostalCode  string              // Postal/ZIP code
    Phone       string              // Contact phone
    Email       string              // Contact email
    Website     string              // Institution website
    IsActive    bool                // Active status
    CreatedAt   time.Time
    UpdatedAt   time.Time
    CreatedBy   *primitive.ObjectID
}
```

**Institution Types**:
- `university`
- `hospital`
- `laboratory`
- `research_center`
- `government`
- `private_practice`
- `ngo`
- `other`

### 2. Updated User Model

**Before**:
```go
type UserProfile struct {
    Institution string  // "University of Cape Town"
    Location    string  // "Cape Town, South Africa"
}
```

**After**:
```go
type UserProfile struct {
    InstitutionID *primitive.ObjectID  // Reference to Institution
}
```

The location is now derived from the institution's city, province, and country fields.

### 3. New Repository Layer

**File**: `backend/internal/repository/institution_repository.go`

**Features**:
- Full CRUD operations
- `FindByID()` - Get institution by ID
- `FindByName()` - Get institution by name
- `List()` - Get paginated list with filters
- `Count()` - Count institutions
- `Activate()` / `Deactivate()` - Status management
- `NameExists()` - Check for duplicates

### 4. New Service Layer

**File**: `backend/internal/service/institution_service.go`

**Features**:
- Business logic and validation
- Permission checks
- Audit logging for all operations
- Search functionality (name, city, province, country)
- `ValidateInstitutionID()` - Verify institution exists and is active

## Migration Strategy

### Step 1: Create Institutions (COMPLETED)
- ✅ Created Institution model
- ✅ Created Institution repository
- ✅ Created Institution service
- ✅ Updated User model

### Step 2: Backend Integration (TODO)
- [ ] Create institution API handlers
- [ ] Add institution routes to server
- [ ] Update user service to validate institution IDs
- [ ] Create seed script for institutions

### Step 3: Frontend Integration (TODO)
- [ ] Create institution types
- [ ] Create institution store
- [ ] Update user types to use institutionId
- [ ] Update UserManagement component
- [ ] Add institution selector dropdown

### Step 4: Data Migration (TODO)
- [ ] Create migration script to:
  - Extract unique institutions from existing users
  - Create Institution documents
  - Update users with institution IDs
  - Remove old institution/location strings

## API Endpoints (TODO)

### Institution Management

```http
GET    /api/institutions          # List institutions
GET    /api/institutions/:id      # Get institution by ID
POST   /api/institutions          # Create institution (requires PermManageUsers)
PUT    /api/institutions/:id      # Update institution (requires PermManageUsers)
DELETE /api/institutions/:id      # Delete institution (requires PermDeleteUsers)
POST   /api/institutions/:id/activate    # Activate
POST   /api/institutions/:id/deactivate  # Deactivate
```

### User Management (Updated)

```http
POST /api/users
Body: {
    "institutionId": "507f1f77bcf86cd799439011",  // Required ObjectID
    // other fields...
}
```

## Benefits

### 1. Data Normalization
- Single source of truth for institution details
- Update institution info in one place, affects all users
- No duplicate or inconsistent data

### 2. Better Data Integrity
- Referential integrity with ObjectID relationships
- Validate institution exists before assigning to user
- Prevent orphaned data

### 3. Enhanced Functionality
- Search and filter institutions independently
- Manage institution lifecycle (activate/deactivate)
- Track institution usage and statistics
- Institution-specific reporting

### 4. Improved User Experience
- Dropdown selection instead of free-text entry
- Auto-complete search for institutions
- Consistent institution names across system
- Display full institution details on demand

### 5. Scalability
- Easy to add new institution fields
- Support for institution hierarchies (future)
- Institution-level permissions (future)
- Multi-institution users (future)

## Example Usage

### Creating an Institution

```go
institution := &CreateInstitutionRequest{
    Name:       "University of Cape Town",
    ShortName:  "UCT",
    Type:       InstitutionTypeUniversity,
    Country:    "South Africa",
    Province:   "Western Cape",
    City:       "Cape Town",
    Address:    "Private Bag X3, Rondebosch",
    PostalCode: "7701",
    Phone:      "+27 21 650 9111",
    Email:      "info@uct.ac.za",
    Website:    "https://www.uct.ac.za",
}
```

### Creating a User with Institution

```go
user := &CreateUserRequest{
    Username:      "john.doe",
    Email:         "john.doe@uct.ac.za",
    InstitutionID: "507f1f77bcf86cd799439011", // UCT's ID
    // other fields...
}
```

### Displaying User with Institution

```typescript
// Frontend - user object now includes institution details
{
  firstName: "John",
  lastName: "Doe",
  institutionId: "507f1f77bcf86cd799439011",
  institution: {  // Populated from join/lookup
    name: "University of Cape Town",
    shortName: "UCT",
    city: "Cape Town",
    country: "South Africa"
  }
}
```

## Testing Checklist

### Backend
- [ ] Create institution
- [ ] Get institution by ID
- [ ] Update institution
- [ ] Delete institution
- [ ] List institutions with pagination
- [ ] Search institutions
- [ ] Activate/deactivate institution
- [ ] Validate institution ID in user creation
- [ ] Reject invalid institution IDs
- [ ] Permission checks work correctly

### Frontend
- [ ] Load institutions list
- [ ] Display institution dropdown
- [ ] Search institutions
- [ ] Select institution when creating user
- [ ] Display institution name in user list
- [ ] Display full institution details in user profile
- [ ] Handle institution not found errors

## Future Enhancements

1. **Institution Hierarchy**
   - Parent-child relationships (e.g., hospital → department)
   - Multi-level organization structures

2. **Institution-Level Permissions**
   - Users can only see data from their institution
   - Cross-institution collaboration features

3. **Institution Statistics**
   - User count per institution
   - Activity metrics
   - Usage reports

4. **Multi-Institution Users**
   - Users affiliated with multiple institutions
   - Primary vs secondary affiliations

5. **Institution Settings**
   - Custom branding per institution
   - Institution-specific configurations
   - Institutional admin roles

## Migration Script (TODO)

```go
// Pseudo-code for data migration
func MigrateUsersToInstitutions() {
    // 1. Get all unique institution names from users
    institutions := extractUniqueInstitutions()
    
    // 2. Create Institution documents
    institutionMap := createInstitutions(institutions)
    
    // 3. Update each user with institution ID
    for user := range users {
        institutionID := institutionMap[user.Profile.Institution]
        updateUserInstitutionID(user.ID, institutionID)
    }
    
    // 4. Remove old fields (optional, for clean migration)
    removeOldInstitutionFields()
}
```

## Notes

- **Breaking Change**: This is a breaking change requiring data migration
- **Backwards Compatibility**: Old API endpoints using institution strings will need updating
- **Data Loss Prevention**: Keep old institution/location fields during migration period
- **Rollback Plan**: Maintain backup before migration in case rollback is needed

## Support

For questions or issues with the institution restructuring, please contact the development team or refer to:
- Backend models: `backend/internal/models/institution.go`
- Repository: `backend/internal/repository/institution_repository.go`
- Service: `backend/internal/service/institution_service.go`


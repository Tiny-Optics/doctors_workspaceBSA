# Database Seed Script

This document explains how to use the database seed script to create the initial Super Admin account for the Doctor's Workspace platform.

## Overview

The seed script (`cmd/seed/main.go`) creates:
- Database indexes for users, sessions, and audit logs
- An initial Super Admin account with full system permissions
- An audit log entry for the admin creation

## Prerequisites

1. MongoDB must be running and accessible
2. Environment variables must be configured (see `.env.example`)

## Running the Seed Script

### Option 1: Using Go directly

```bash
cd backend
go run cmd/seed/main.go
```

### Option 2: Build and run

```bash
cd backend
go build -o seed cmd/seed/main.go
./seed
```

### Option 3: Add to Makefile

Add this target to your `backend/Makefile`:

```makefile
.PHONY: seed
seed:
	@echo "Running database seed..."
	@go run cmd/seed/main.go
```

Then run:

```bash
make seed
```

## Configuration

The seed script uses environment variables to configure the Super Admin account. If not set, it will use default values.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SUPER_ADMIN_EMAIL` | `admin@bloodsa.org.za` | Email address for super admin |
| `SUPER_ADMIN_USERNAME` | `superadmin` | Username for super admin |
| `SUPER_ADMIN_PASSWORD` | `BloodSA2025!` | Initial password (CHANGE AFTER FIRST LOGIN) |
| `SUPER_ADMIN_FIRST_NAME` | `Super` | First name |
| `SUPER_ADMIN_LAST_NAME` | `Admin` | Last name |

### Custom Configuration Example

```bash
export SUPER_ADMIN_EMAIL=admin@bloodsa.org.za
export SUPER_ADMIN_USERNAME=bloodsa_admin
export SUPER_ADMIN_PASSWORD='MySecure@Pass123!'
export SUPER_ADMIN_FIRST_NAME=John
export SUPER_ADMIN_LAST_NAME=Smith

go run cmd/seed/main.go
```

## Output

The script will output the created admin credentials:

```
✅ Super admin created successfully!
================================================
Username: superadmin
Email: admin@bloodsa.org.za
Password: BloodSA2025!
Role: admin
Admin Level: super_admin
================================================
⚠️  IMPORTANT: Please change the password after first login!
Seed completed successfully
```

## Important Notes

1. **Idempotent**: The script checks if a super admin with the specified email already exists. If found, it will skip creation.

2. **Password Security**: The default password is intended for initial setup only. **You MUST change it after first login!**

3. **Production**: In production environments:
   - Always set a strong custom password via environment variables
   - Never commit `.env` files with real credentials
   - Change the password immediately after first login
   - Consider creating backup Super Admin accounts

4. **Audit Trail**: The script logs the admin creation in the audit logs collection for security tracking.

## After Running the Seed

1. Log in to the platform using the created credentials
2. **Immediately change the password** via the user profile settings
3. Create additional admin accounts (User Managers or Super Admins) as needed
4. Create clinical user accounts for haematologists, physicians, and data capturers

## Troubleshooting

### Error: "Failed to connect to MongoDB"
- Ensure MongoDB is running: `docker-compose up -d` or check your MongoDB service
- Verify `BLUEPRINT_DB_HOST` and `BLUEPRINT_DB_PORT` environment variables

### Error: "Super admin already exists"
- This is not an error - the admin account was already created
- To reset, you can delete the user from MongoDB and run the seed again

### Error: "Failed to hash password"
- Check that the password meets minimum requirements (8+ chars, mixed case, numbers, special chars)

## Creating Additional Super Admins

After the initial setup, you can create additional Super Admin accounts:

1. Log in as the initial Super Admin
2. Navigate to User Management
3. Create a new user with:
   - Role: Admin
   - Admin Level: Super Admin

Only Super Admins can create other Super Admin accounts.

## Security Best Practices

1. **Minimum Super Admins**: Maintain at least 2 Super Admin accounts to prevent lockout
2. **Strong Passwords**: Use strong, unique passwords for all admin accounts
3. **Regular Audits**: Periodically review the audit logs for suspicious activity
4. **Access Control**: Only grant Super Admin access to trusted personnel
5. **Password Rotation**: Implement a password rotation policy for admin accounts


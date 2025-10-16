# Dashboard & Navigation Guide

## 🎉 New Features

### Dashboard Page
A personalized landing page for authenticated users with feature selection.

### Updated Header
Smart navigation that adapts based on authentication state and user role.

## 📍 Routes

### Public Routes
- **`/`** - Home/Landing page (public)
- **`/login`** - Login page (redirects to dashboard if already authenticated)

### Protected Routes (Require Authentication)
- **`/dashboard`** - Main dashboard with feature selection
- **`/sops`** - Standard Operating Procedures (coming soon)
- **`/referrals`** - Transplant Referrals (coming soon)
- **`/registry`** - African HOPeR Registry (coming soon)
- **`/documents`** - Document Library (coming soon)
- **`/training`** - Training Resources (coming soon)
- **`/profile`** - User Profile (coming soon)
- **`/settings`** - Settings (coming soon)

### Admin-Only Routes
- **`/admin`** - Admin Panel (coming soon, requires admin role)

## 🎨 Dashboard Features

### Welcome Section
- **User greeting** - "Welcome back, [First Name]"
- **Institution & role** display
- **Last login** timestamp with smart formatting:
  - "X minutes ago" (< 1 hour)
  - "X hours ago" (< 24 hours)
  - "X days ago" (24+ hours)

### Quick Stats (Admin Only)
Four stat cards showing:
- Total Users
- SOPs
- Referrals
- Registry Entries

*Note: Currently showing placeholder values (--), will be populated when backend endpoints are connected.*

### Feature Cards
Six interactive cards for main features:

1. **Standard Operating Procedures**
   - Access SOPs for Anemia, Lymphoma, Myeloma, General Business
   - Links to `/sops`

2. **Transplant Referrals**
   - Submit and track referrals via REDCap
   - Links to `/referrals`

3. **African HOPeR Registry**
   - Documentation, training, ethics approvals
   - Links to `/registry`

4. **Document Library**
   - Access all clinical documents
   - Links to `/documents`

5. **Training Resources**
   - Videos, guidelines, tutorials
   - Links to `/training`

6. **My Profile**
   - View and update profile
   - Links to `/profile`

### Card Features
- ✅ Hover effects with scale animation
- ✅ Color-coded icons (red, green, purple, blue, yellow, gray)
- ✅ Arrow animation on hover
- ✅ Border highlight on hover
- ✅ Smooth transitions

## 📱 Header Navigation

### When Not Authenticated
Shows:
- BLOODSA Logo (links to home)
- **Login** button (goes to `/login`)

### When Authenticated

#### For All Users
- **Dashboard** link
- **User menu dropdown** with:
  - User avatar (initials)
  - First name
  - Email address
  - Profile link
  - Settings link
  - **Sign out** button

#### For Admin Users ONLY
- **Admin Panel** button (red, with gear icon)
  - Only visible when user has admin role
  - Links to `/admin`
  - Protected by route guard

### User Menu Dropdown
- Click avatar/name to open
- Shows user info at top
- Links to Profile and Settings
- Sign out button (red text)
- Click outside or on backdrop to close
- Automatically closes when navigating

## 🔐 Route Protection

### Authentication Guard
```typescript
// Routes with meta: { requiresAuth: true }
// Redirect to /login if not authenticated
```

### Guest Guard
```typescript
// Routes with meta: { requiresGuest: true }
// Redirect to /dashboard if already authenticated
```

### Admin Guard
```typescript
// Routes with meta: { requiresAdmin: true }
// Redirect to /dashboard if not an admin
```

## 🎯 User Flow

### First-Time Login
1. User lands on `/login`
2. Enters credentials
3. Successfully authenticates
4. Redirected to `/dashboard`
5. Sees personalized welcome message
6. Can select feature to use

### Returning User
1. User accesses any URL
2. If authenticated (token in localStorage):
   - Loads user data
   - Shows authenticated header
   - Access allowed to protected routes
3. If not authenticated:
   - Redirected to `/login`

### Admin User
1. Logs in as admin
2. Sees `/dashboard` with:
   - Quick stats section
   - All feature cards
3. Header shows:
   - **Admin Panel** button
   - User menu
4. Can access `/admin` route

### Non-Admin User
1. Logs in as clinical user
2. Sees `/dashboard` with:
   - Feature cards (no stats)
3. Header shows:
   - User menu (no admin button)
4. Cannot access `/admin` (redirected to dashboard)

## 🚀 Testing

### Test Dashboard
```bash
# 1. Start backend
cd backend && make run

# 2. Start frontend
cd frontend && npm run dev

# 3. Login
http://localhost:5173/login
Email: admin@bloodsa.org.za
Password: BloodSA2025!

# 4. View dashboard
# Should redirect to: http://localhost:5173/dashboard
```

### Test Admin Features
1. Login with admin account
2. Check header shows "Admin Panel" button
3. Check dashboard shows "Quick Stats" section
4. Click "Admin Panel" button
5. Should navigate to "Coming Soon" page

### Test Non-Admin
1. Create a non-admin user via API or seed script
2. Login with that account
3. Verify NO "Admin Panel" button in header
4. Verify NO "Quick Stats" on dashboard
5. Verify cannot access `/admin` URL directly

### Test Logout
1. Click on user avatar/name
2. Menu drops down
3. Click "Sign out"
4. Should redirect to `/login`
5. Session cleared from localStorage

## 🎨 Styling & Theme

### Colors Used
- **Primary (BLOODSA Red)**: `#8B0000`
- **Green**: Stats & referrals
- **Purple**: Registry
- **Blue**: Documents
- **Yellow**: Training
- **Gray**: Profile/neutral

### Components
- **Cards**: White background, shadow on hover
- **Buttons**: Red primary, white text
- **Avatar**: Red circle with white initials
- **Icons**: Tailwind Heroicons

## 📝 Coming Soon Features

All feature routes currently show a "Coming Soon" page with:
- Clock icon
- "Under development" message
- "Back to Dashboard" button
- Development status card

## 🔧 Next Steps

1. **Build feature pages**
   - Replace `ComingSoon.vue` with actual feature pages
   - SOPs, Referrals, Registry, etc.

2. **Connect backend stats**
   - Create API endpoints for dashboard stats
   - Display real data instead of placeholders

3. **Build admin panel**
   - User management interface
   - System settings
   - Audit logs viewer

4. **Enhance profile page**
   - Edit profile information
   - Change password
   - View account details

5. **Add notifications**
   - Toast notifications
   - Real-time updates

## 🐛 Troubleshooting

### "Admin Panel" button not showing
- Verify user role is `"admin"` in database
- Check `authStore.isAdmin` returns `true`
- Clear localStorage and re-login

### Redirecting to login constantly
- Check token in localStorage
- Verify backend is running
- Check token hasn't expired (24h)

### Dashboard blank or error
- Check browser console for errors
- Verify all imports are correct
- Check auth store has user data

### Feature cards not clickable
- Check router has routes defined
- Verify `ComingSoon.vue` exists
- Check browser console for errors

## 📚 File Structure

```
frontend/src/
├── views/
│   ├── Dashboard.vue          # Main dashboard
│   ├── Login.vue              # Login page
│   ├── Home.vue               # Public landing
│   └── ComingSoon.vue         # Placeholder
├── components/
│   └── Header.vue             # Updated navigation
├── stores/
│   ├── auth.ts                # Authentication store
│   └── users.ts               # User management
├── router/
│   └── index.ts               # Route definitions
└── types/
    └── user.ts                # TypeScript types
```

---

**Status**: ✅ Dashboard and navigation complete

All routing, guards, and UI components are working. Ready for feature implementation!


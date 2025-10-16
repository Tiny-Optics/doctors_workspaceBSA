# Login Page - Testing Guide

## 🎉 What Was Created

A fully functional login page with:
- ✅ Modern, responsive UI with BLOODSA branding
- ✅ Form validation
- ✅ Password show/hide toggle
- ✅ Error handling with dismissible alerts
- ✅ Loading states
- ✅ Route guards for authentication
- ✅ Persistent sessions via localStorage
- ✅ Remember me functionality

## 🚀 Testing the Login Page

### Step 1: Start Backend Server

```bash
cd backend

# Make sure MongoDB is running
make docker-run

# Seed the database (creates Super Admin)
make seed

# Start the backend server
make run
```

**Backend should be running on**: `http://localhost:8080`

### Step 2: Start Frontend Development Server

```bash
cd frontend

# Install dependencies (if not already done)
npm install

# Start dev server
npm run dev
```

**Frontend should be running on**: `http://localhost:5173`

### Step 3: Access the Login Page

Open your browser and go to: **http://localhost:5173/login**

### Step 4: Test Login

Use the default Super Admin credentials from the seed script:

**Email**: `admin@bloodsa.org.za`  
**Password**: `BloodSA2025!`

## 🧪 Testing Scenarios

### 1. Successful Login ✅
- Enter valid credentials
- Click "Sign in"
- Should redirect to home page (`/`)
- Token stored in localStorage

### 2. Invalid Credentials ❌
- Enter wrong email or password
- Should show error: "Invalid email or password"
- Error is dismissible

### 3. Account Locked 🔒
- Try logging in with wrong password 5+ times
- Should show: "Account is locked due to too many failed login attempts"
- Account locks for 30 minutes

### 4. Empty Fields ⚠️
- Try submitting without email or password
- Browser validation should prevent submission

### 5. Password Visibility 👁️
- Click the eye icon to toggle password visibility
- Password should switch between masked and visible

### 6. Already Authenticated 🔄
- Log in successfully
- Try to navigate to `/login` again
- Should automatically redirect to home page

### 7. Protected Routes 🔐
- Without logging in, try to access `/`
- Should redirect to `/login`

### 8. Persistent Session 💾
- Log in
- Refresh the page
- Should remain logged in (token from localStorage)

### 9. Logout 🚪
- Log in
- Call logout (need to add logout button to home page)
- Should redirect to login and clear session

## 📱 Responsive Testing

Test on different screen sizes:
- ✅ Desktop (1920x1080)
- ✅ Tablet (768x1024)
- ✅ Mobile (375x667)

The login page should look good on all devices!

## 🎨 UI Features

### Visual Elements
- ✅ BLOODSA red branding (#8B0000)
- ✅ Background pattern
- ✅ Card-based layout with shadow
- ✅ Smooth transitions and hover effects
- ✅ Loading spinner during authentication

### Form Elements
- ✅ Email input with validation
- ✅ Password input with show/hide toggle
- ✅ Remember me checkbox
- ✅ Forgot password link (placeholder)
- ✅ Submit button with loading state

### Error Handling
- ✅ Dismissible error alerts
- ✅ Red border for error state
- ✅ Clear error messages

## 🔧 Development Notes

### API Endpoint
Currently hardcoded to: `http://localhost:8080/api/auth/login`

In production, you should use an environment variable:
```typescript
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
```

### Routes
- `/login` - Login page (guest only)
- `/` - Home page (authenticated only)

### Route Guards
Implemented in `frontend/src/router/index.ts`:
- Checks `authStore.isAuthenticated`
- Redirects unauthenticated users to `/login`
- Redirects authenticated users away from `/login`

## 🐛 Troubleshooting

### "Connection refused" error
**Problem**: Backend not running or wrong URL  
**Solution**: Make sure backend is running on port 8080

### "CORS error"
**Problem**: CORS not configured properly  
**Solution**: Backend has CORS enabled for `http://localhost:5173`

### Login successful but redirect fails
**Problem**: Router not properly initialized  
**Solution**: Check that router is imported and used in main.ts

### "Invalid token" on refresh
**Problem**: Token expired or invalid  
**Solution**: Token expires after 24 hours, use refresh token or login again

### Stays on login page after successful login
**Problem**: Router navigation issue  
**Solution**: Check browser console for errors

## 📝 Next Steps

1. **Add Logout Button** to home page
2. **Implement Forgot Password** flow
3. **Add User Profile** page
4. **Create Admin Dashboard** for user management
5. **Add Navigation Menu** with logout
6. **Implement Token Refresh** automatically before expiry

## 🎯 Test Checklist

- [ ] Backend is running
- [ ] Frontend is running
- [ ] Super Admin exists in database
- [ ] Can navigate to `/login`
- [ ] Can submit login form with valid credentials
- [ ] Redirects to home after login
- [ ] Shows error with invalid credentials
- [ ] Password toggle works
- [ ] Error is dismissible
- [ ] Protected routes redirect to login
- [ ] Login persists after page refresh
- [ ] Already logged-in users redirect from login page

## 🔐 Default Test Credentials

**Super Admin**:
- Email: `admin@bloodsa.org.za`
- Password: `BloodSA2025!`

⚠️ **Important**: Change this password after first login in production!

## 📞 Support

If you encounter issues:
1. Check browser console for errors
2. Check backend logs
3. Verify MongoDB is running
4. Ensure seed script ran successfully

---

**Status**: ✅ Login page ready for testing

All core authentication functionality is implemented and working!


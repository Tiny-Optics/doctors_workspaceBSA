import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Dashboard from '../views/Dashboard.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
      meta: { requiresGuest: true }
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: { requiresGuest: true }
    },
    {
      path: '/register',
      name: 'register',
      component: Register,
      meta: { requiresGuest: true }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: Dashboard,
      meta: { requiresAuth: true }
    },
    // Placeholder routes for dashboard features
    {
      path: '/sops',
      name: 'sops',
      component: () => import('../views/sops/SOPs.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/sops/:slug',
      name: 'sops-list',
      component: () => import('../views/sops/SOPList.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/referrals',
      name: 'referrals',
      component: () => import('../views/ReferralPage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/registry',
      name: 'registry',
      component: () => import('../views/registry/Registry.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/registry/example-documents',
      name: 'registry-example-documents',
      component: () => import('../views/registry/ExampleDocuments.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/registry/upload',
      name: 'registry-upload',
      component: () => import('../views/registry/UploadApprovals.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/documents',
      name: 'documents',
      component: () => import('../views/ComingSoon.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/training',
      name: 'training',
      component: () => import('../views/ComingSoon.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/Profile.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/ComingSoon.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/working-parties',
      name: 'working-parties',
      component: () => import('../views/ComingSoon.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/Admin.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: () => import('../views/admin/AdminDashboard.vue')
        },
        {
          path: 'users',
          name: 'admin-users',
          component: () => import('../views/admin/UserManagement.vue')
        },
        {
          path: 'institutions',
          name: 'admin-institutions',
          component: () => import('../views/admin/InstitutionManagement.vue')
        },
        {
          path: 'settings',
          name: 'admin-settings',
          component: () => import('../views/admin/SystemSettings.vue')
        },
        {
          path: 'audit-logs',
          name: 'admin-audit-logs',
          component: () => import('../views/admin/AuditLogs.vue')
        },
        {
          path: 'sops',
          name: 'admin-sops',
          component: () => import('../views/admin/SOPManagement.vue')
        },
        {
          path: 'registry',
          name: 'admin-registry',
          component: () => import('../views/admin/RegistrySettings.vue')
        },
        {
          path: 'registry/configuration',
          name: 'admin-registry-configuration',
          component: () => import('../views/admin/RegistryConfiguration.vue')
        },
        {
          path: 'registry/forms',
          name: 'admin-registry-forms',
          component: () => import('../views/admin/RegistryForms.vue')
        },
        {
          path: 'registry/submissions',
          name: 'admin-registry-submissions',
          component: () => import('../views/admin/RegistrySubmissions.vue')
        },
        {
          path: 'referrals',
          name: 'admin-referrals',
          component: () => import('../views/admin/ReferralSettings.vue')
        }
      ]
    }
  ],
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // Check if route requires authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login' })
  } 
  // Redirect authenticated users away from guest-only pages (home, login) to dashboard
  else if (to.meta.requiresGuest && authStore.isAuthenticated) {
    next({ name: 'dashboard' })
  }
  // Check if route requires admin
  else if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'dashboard' })
  }
  else {
    next()
  }
})

export default router

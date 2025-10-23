<template>
  <aside 
    :class="[
      'bg-white shadow-lg border-r border-gray-200 flex flex-col h-screen transition-all duration-300 ease-in-out overflow-hidden',
      isCollapsed ? 'w-20' : 'w-64'
    ]"
  >
    <!-- Logo Section -->
    <div class="p-6 border-b border-gray-200 relative">
      <div class="flex items-center" :class="isCollapsed ? 'justify-center' : 'space-x-3'">
        <div class="w-10 h-10 bg-bloodsa-red rounded-lg flex items-center justify-center flex-shrink-0">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </div>
        <div v-if="!isCollapsed" class="overflow-hidden">
          <h2 class="text-lg font-bold text-gray-900 whitespace-nowrap">Admin Panel</h2>
          <p class="text-sm text-gray-500 whitespace-nowrap">Doctor's Workspace</p>
        </div>
      </div>
      
      <!-- Toggle Button -->
      <button
        @click="toggleSidebar"
        class="absolute -right-3 top-8 bg-white border-2 border-gray-200 rounded-full p-1 hover:bg-gray-50 transition-colors shadow-md"
        :title="isCollapsed ? 'Expand sidebar' : 'Collapse sidebar'"
      >
        <svg 
          class="w-4 h-4 text-gray-600 transition-transform duration-300"
          :class="{ 'rotate-180': isCollapsed }"
          fill="none" 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
    </div>

    <!-- Navigation Menu -->
    <nav class="flex-1 px-4 py-6 space-y-2 overflow-y-auto scrollbar-thin scrollbar-thumb-gray-300 scrollbar-track-gray-100 hover:scrollbar-thumb-gray-400">
      <template v-for="item in navigationItems" :key="item.name">
        <router-link
          :to="item.to"
          class="flex items-center rounded-lg transition-all duration-200 group relative overflow-visible"
          :class="[
            { 
              'bg-bloodsa-red text-white shadow-md': isActiveRoute(item),
              'text-gray-700 hover:bg-bloodsa-red hover:text-white': !isActiveRoute(item)
            },
            isCollapsed ? 'justify-center px-3 py-3' : 'space-x-3 px-4 py-3'
          ]"
          :title="isCollapsed ? item.name : ''"
        >
          <!-- Active indicator bar (left side) -->
          <div 
            v-if="isCollapsed && isActiveRoute(item)"
            class="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-white rounded-r-full"
          ></div>
          
          <!-- Dashboard Icon -->
          <svg 
            v-if="item.name === 'Admin Dashboard'"
            class="flex-shrink-0 transition-colors duration-200"
            :class="[
              isCollapsed ? 'w-6 h-6' : 'w-5 h-5',
              { 
                'text-white': isActiveRoute(item),
                'text-gray-600 group-hover:text-white': !isActiveRoute(item)
              }
            ]"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          
          <!-- Users Icon -->
          <svg 
            v-else-if="item.name === 'User Management'"
            class="flex-shrink-0 transition-colors duration-200"
            :class="[
              isCollapsed ? 'w-6 h-6' : 'w-5 h-5',
              { 
                'text-white': isActiveRoute(item),
                'text-gray-600 group-hover:text-white': !isActiveRoute(item)
              }
            ]"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
          </svg>
          
          <!-- Institution Icon -->
          <svg 
            v-else-if="item.name === 'Institution Management'"
            class="flex-shrink-0 transition-colors duration-200"
            :class="[
              isCollapsed ? 'w-6 h-6' : 'w-5 h-5',
              { 
                'text-white': isActiveRoute(item),
                'text-gray-600 group-hover:text-white': !isActiveRoute(item)
              }
            ]"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
          </svg>
          
          <!-- Settings Icon -->
          <svg 
            v-else-if="item.name === 'System Settings'"
            class="flex-shrink-0 transition-colors duration-200"
            :class="[
              isCollapsed ? 'w-6 h-6' : 'w-5 h-5',
              { 
                'text-white': isActiveRoute(item),
                'text-gray-600 group-hover:text-white': !isActiveRoute(item)
              }
            ]"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          
          <!-- SOP Management Icon -->
          <svg 
            v-else-if="item.name === 'SOP Management'"
            class="flex-shrink-0 transition-colors duration-200"
            :class="[
              isCollapsed ? 'w-6 h-6' : 'w-5 h-5',
              { 
                'text-white': isActiveRoute(item),
                'text-gray-600 group-hover:text-white': !isActiveRoute(item)
              }
            ]"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          
          <!-- Registry Settings Icon -->
          <svg 
            v-else-if="item.name === 'Registry Settings'"
            class="flex-shrink-0 transition-colors duration-200"
            :class="[
              isCollapsed ? 'w-6 h-6' : 'w-5 h-5',
              { 
                'text-white': isActiveRoute(item),
                'text-gray-600 group-hover:text-white': !isActiveRoute(item)
              }
            ]"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
          </svg>
          
          <!-- Referral Settings Icon -->
          <svg 
            v-else-if="item.name === 'Referral Settings'"
            class="flex-shrink-0 transition-colors duration-200"
            :class="[
              isCollapsed ? 'w-6 h-6' : 'w-5 h-5',
              { 
                'text-white': isActiveRoute(item),
                'text-gray-600 group-hover:text-white': !isActiveRoute(item)
              }
            ]"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
          </svg>
          
          <!-- Audit Logs Icon -->
          <svg 
            v-else-if="item.name === 'Audit Logs'"
            class="flex-shrink-0 transition-colors duration-200"
            :class="[
              isCollapsed ? 'w-6 h-6' : 'w-5 h-5',
              { 
                'text-white': isActiveRoute(item),
                'text-gray-600 group-hover:text-white': !isActiveRoute(item)
              }
            ]"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          <span v-if="!isCollapsed" class="font-medium whitespace-nowrap">{{ item.name }}</span>
          
          <!-- Tooltip for collapsed state -->
          <div 
            v-if="isCollapsed"
            class="absolute left-full ml-2 px-3 py-2 bg-gray-900 text-white text-sm rounded-md opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 whitespace-nowrap z-50 shadow-lg pointer-events-none"
          >
            {{ item.name }}
            <div class="absolute right-full top-1/2 -translate-y-1/2 border-4 border-transparent border-r-gray-900"></div>
          </div>
        </router-link>
      </template>
    </nav>

    <!-- User Info Footer -->
    <div class="p-4 border-t border-gray-200">
      <div class="flex items-center" :class="isCollapsed ? 'justify-center' : 'space-x-3'">
        <div class="w-8 h-8 bg-bloodsa-red rounded-full flex items-center justify-center text-white text-sm font-semibold flex-shrink-0">
          {{ userInitials }}
        </div>
        <div v-if="!isCollapsed" class="flex-1 min-w-0 overflow-hidden">
          <p class="text-sm font-medium text-gray-900 truncate">{{ userFullName }}</p>
          <p class="text-xs text-gray-500 truncate">{{ userRole }}</p>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getUserRoleDisplayName } from '@/types/user'

const route = useRoute()

// Sidebar collapse state
const isCollapsed = ref(false)

// Helper function to check if a route is active (including sub-routes)
const isActiveRoute = (item: { name: string; to: string }) => {
  // For Registry Settings, match all /admin/registry/* paths
  if (item.to === '/admin/registry') {
    return route.path.startsWith('/admin/registry')
  }
  // For Referral Settings, match all /admin/referrals/* paths
  if (item.to === '/admin/referrals') {
    return route.path.startsWith('/admin/referrals')
  }
  // For other routes, exact match
  return route.path === item.to
}

// Load collapse state from localStorage
onMounted(() => {
  const savedState = localStorage.getItem('adminSidebarCollapsed')
  if (savedState !== null) {
    isCollapsed.value = savedState === 'true'
  }
})

// Toggle sidebar and save state
const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
  localStorage.setItem('adminSidebarCollapsed', isCollapsed.value.toString())
}

const authStore = useAuthStore()

// Navigation items - easily configurable
const navigationItems = [
  {
    name: 'Admin Dashboard',
    to: '/admin'
  },
  {
    name: 'User Management',
    to: '/admin/users'
  },
  {
    name: 'Institution Management',
    to: '/admin/institutions'
  },
  {
    name: 'SOP Management',
    to: '/admin/sops'
  },
  {
    name: 'Registry Settings',
    to: '/admin/registry'
  },
  {
    name: 'Referral Settings',
    to: '/admin/referrals'
  },
  {
    name: 'System Settings',
    to: '/admin/settings'
  },
  {
    name: 'Audit Logs',
    to: '/admin/audit-logs'
  }
]

const user = computed(() => authStore.user)

const userFullName = computed(() => {
  if (!user.value) return 'Admin User'
  return `${user.value.profile.firstName} ${user.value.profile.lastName}`
})

const userInitials = computed(() => {
  if (!user.value) return 'AU'
  const firstName = user.value.profile.firstName || ''
  const lastName = user.value.profile.lastName || ''
  return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase()
})

const userRole = computed(() => {
  if (!user.value) return 'Admin'
  return getUserRoleDisplayName(user.value.role as any)
})
</script>

<style scoped>
/* Custom scrollbar styling for the navigation menu */
nav::-webkit-scrollbar {
  width: 6px;
}

nav::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 3px;
}

nav::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
  transition: background-color 0.2s ease;
}

nav::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* Firefox scrollbar styling */
nav {
  scrollbar-width: thin;
  scrollbar-color: #cbd5e1 #f1f5f9;
}

nav:hover {
  scrollbar-color: #94a3b8 #f1f5f9;
}
</style>

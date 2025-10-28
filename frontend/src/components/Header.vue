<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo and Title -->
        <router-link to="/dashboard" class="flex items-center space-x-3 hover:opacity-80 transition-opacity">
          <div class="relative">
            <img 
              src="/BLOODSA-Logo.svg" 
              alt="BloodSA Logo" 
              class="w-12 h-12 object-contain"
            />
          </div>
          <h1 class="text-xl font-bold text-bloodsa-red">Doctor's Workspace</h1>
        </router-link>

        <!-- Navigation -->
        <nav class="flex items-center space-x-6">
          <!-- Show these only when authenticated -->
          <template v-if="authStore.isAuthenticated">
            <router-link 
              to="/dashboard" 
              class="text-gray-700 hover:text-bloodsa-red font-medium transition-colors duration-200"
              :class="{ 'text-bloodsa-red font-semibold': $route.path === '/dashboard' }"
            >
              Dashboard
            </router-link>

            <!-- Admin Panel Button (only for admins) -->
            <router-link 
              v-if="authStore.isAdmin"
              to="/admin"
              class="flex items-center space-x-2 px-4 py-2 bg-bloodsa-red text-white rounded-md hover:bg-opacity-90 transition-all duration-200 font-medium"
            >
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span>Admin Panel</span>
            </router-link>

            <!-- User Menu Dropdown -->
            <div class="relative">
              <button
                @click="showUserMenu = !showUserMenu"
                class="flex items-center space-x-3 px-3 py-2 rounded-md hover:bg-gray-100 transition-colors"
              >
                <div class="w-8 h-8 bg-bloodsa-red rounded-full flex items-center justify-center text-white font-semibold">
                  {{ userInitials }}
                </div>
                <span class="text-gray-700 font-medium">{{ userName }}</span>
                <svg class="w-4 h-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>

              <!-- Dropdown Menu -->
              <div
                v-if="showUserMenu"
                @click.stop
                class="absolute right-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-50"
              >
                <div class="py-1">
                  <div class="px-4 py-3 border-b border-gray-200">
                    <p class="text-sm text-gray-500">Signed in as</p>
                    <p class="text-sm font-medium text-gray-900 truncate">{{ userEmail }}</p>
                  </div>
                  <router-link
                    to="/profile"
                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    @click="showUserMenu = false"
                  >
                    Your Profile
                  </router-link>
                  <router-link
                    to="/settings"
                    class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    @click="showUserMenu = false"
                  >
                    Settings
                  </router-link>
                  <button
                    @click="handleLogout"
                    class="block w-full text-left px-4 py-2 text-sm text-red-700 hover:bg-gray-100"
                  >
                    Sign out
                  </button>
                </div>
              </div>
            </div>
          </template>

          <!-- Login Button (only when not authenticated) -->
          <template v-else>
            <router-link to="/login" class="btn-primary">
              Login
            </router-link>
          </template>
        </nav>
      </div>
    </div>

    <!-- Click outside to close dropdown -->
    <div
      v-if="showUserMenu"
      @click="showUserMenu = false"
      class="fixed inset-0 z-40"
    ></div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const showUserMenu = ref(false)

const userName = computed(() => {
  if (!authStore.user) return 'Guest'
  return authStore.user.profile.firstName
})

const userEmail = computed(() => {
  return authStore.user?.email || ''
})

const userInitials = computed(() => {
  if (!authStore.user) return 'G'
  const first = authStore.user.profile.firstName.charAt(0).toUpperCase()
  const last = authStore.user.profile.lastName.charAt(0).toUpperCase()
  return `${first}${last}`
})

const handleLogout = async () => {
  showUserMenu.value = false
  await authStore.logout()
  router.push('/login')
}
</script>

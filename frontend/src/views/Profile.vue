<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
    <!-- Header Section -->
    <section class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-6">
            <!-- Avatar -->
            <div class="w-20 h-20 bg-bloodsa-red rounded-full flex items-center justify-center text-white text-3xl font-bold">
              {{ userInitials }}
            </div>
            <!-- User Info -->
            <div>
              <h1 class="text-3xl font-bold text-gray-900">{{ userFullName }}</h1>
              <div class="flex items-center space-x-3 mt-2">
                <span class="px-3 py-1 bg-bloodsa-red bg-opacity-10 text-bloodsa-red rounded-full text-sm font-medium">
                  {{ getRoleDisplayName(user?.role) }}
                </span>
                <span v-if="user?.adminLevel" class="px-3 py-1 bg-purple-100 text-purple-700 rounded-full text-sm font-medium">
                  {{ getAdminLevelDisplayName(user?.adminLevel) }}
                </span>
              </div>
            </div>
          </div>
          
          <!-- Back to Dashboard -->
          <router-link
            to="/dashboard"
            class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors flex items-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Back to Dashboard
          </router-link>
        </div>
      </div>
    </section>

    <!-- Main Content -->
    <section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Tabs Navigation -->
      <div class="bg-white rounded-t-lg border-b border-gray-200">
        <nav class="flex space-x-8 px-6" aria-label="Tabs">
          <button
            @click="activeTab = 'personal'"
            :class="[
              activeTab === 'personal'
                ? 'border-bloodsa-red text-bloodsa-red'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'py-4 px-1 border-b-2 font-medium text-sm transition-colors flex items-center gap-2'
            ]"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            Personal Information
          </button>
          <button
            @click="activeTab = 'security'"
            :class="[
              activeTab === 'security'
                ? 'border-bloodsa-red text-bloodsa-red'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'py-4 px-1 border-b-2 font-medium text-sm transition-colors flex items-center gap-2'
            ]"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
            Security
          </button>
        </nav>
      </div>

      <!-- Tab Content -->
      <div class="bg-white rounded-b-lg shadow-lg p-6 min-h-[500px]">
        <PersonalInfoTab v-if="activeTab === 'personal'" />
        <SecurityTab v-else-if="activeTab === 'security'" />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getUserRoleDisplayName as getRoleDisplayName, getAdminLevelDisplayName } from '@/types/user'
import PersonalInfoTab from '@/components/profile/PersonalInfoTab.vue'
import SecurityTab from '@/components/profile/SecurityTab.vue'

const authStore = useAuthStore()
const activeTab = ref('personal')

const user = computed(() => authStore.user)

const userFullName = computed(() => {
  if (!user.value) return 'User'
  return `${user.value.profile.firstName} ${user.value.profile.lastName}`
})

const userInitials = computed(() => {
  if (!user.value) return 'U'
  const firstName = user.value.profile.firstName || ''
  const lastName = user.value.profile.lastName || ''
  return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase()
})
</script>


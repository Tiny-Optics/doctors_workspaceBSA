<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
    <!-- Hero Section -->
    <section class="relative overflow-hidden bg-white">
      <!-- Background Pattern -->
      <div class="absolute inset-0 opacity-5">
        <div class="absolute inset-0" style="background-image: repeating-linear-gradient(45deg, #8B0000 0, #8B0000 1px, transparent 0, transparent 50%); background-size: 10px 10px;"></div>
      </div>

      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-20">
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <!-- Welcome Message -->
            <h1 class="text-4xl md:text-5xl font-bold text-gray-900 mb-4">
              Welcome back,
              <span class="text-bloodsa-red">{{ userFullName }}</span>
            </h1>
            
            <!-- User Info -->
            <div class="flex flex-col sm:flex-row sm:items-center gap-4 mb-6">
              <div class="flex items-center space-x-2">
                <div class="w-2 h-2 bg-green-500 rounded-full"></div>
                <span class="text-lg text-gray-600">{{ getInstitutionName(user?.profile.institutionId) }}</span>
              </div>
              <div class="flex items-center space-x-2">
                <span class="px-3 py-1 bg-bloodsa-red bg-opacity-10 text-bloodsa-red rounded-full text-sm font-medium">
                  {{ getRoleDisplayName(user?.role) }}
                </span>
                <span v-if="authStore.isAdmin" class="px-3 py-1 bg-purple-100 text-purple-700 rounded-full text-sm font-medium">
                  {{ getAdminLevelDisplayName(user?.adminLevel) }}
                </span>
              </div>
            </div>

            <!-- Quick Stats -->
            <div v-if="authStore.isAdmin" class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
              <div class="bg-gradient-to-r from-bloodsa-red to-red-700 text-white p-4 rounded-lg">
                <div class="text-2xl font-bold">{{ systemStats.totalUsers }}</div>
                <div class="text-sm opacity-90">Total Users</div>
              </div>
              <div class="bg-gradient-to-r from-blue-500 to-blue-700 text-white p-4 rounded-lg">
                <div class="text-2xl font-bold">{{ systemStats.totalSOPs }}</div>
                <div class="text-sm opacity-90">SOPs</div>
              </div>
              <div class="bg-gradient-to-r from-green-500 to-green-700 text-white p-4 rounded-lg">
                <div class="text-2xl font-bold">{{ systemStats.totalReferrals }}</div>
                <div class="text-sm opacity-90">Referrals</div>
              </div>
              <div class="bg-gradient-to-r from-purple-500 to-purple-700 text-white p-4 rounded-lg">
                <div class="text-2xl font-bold">{{ systemStats.totalRegistry }}</div>
                <div class="text-sm opacity-90">Registry</div>
              </div>
            </div>
          </div>

          <!-- User Avatar & Quick Actions -->
          <div class="flex flex-col items-end space-y-4">
            <div class="flex items-center space-x-4">
              <div class="text-right">
                <p class="text-sm text-gray-500">Last login</p>
                <p class="text-sm font-medium text-gray-900">
                  {{ formatLastLogin(user?.lastLoginAt) }}
                </p>
              </div>
              <div class="w-16 h-16 bg-bloodsa-red rounded-full flex items-center justify-center text-white text-xl font-bold">
                {{ userInitials }}
              </div>
            </div>
            
            <!-- Quick Action Buttons -->
            <div class="flex space-x-2">
              <router-link
                to="/profile"
                class="px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 transition-colors text-sm font-medium"
              >
                Profile
              </router-link>
              <router-link
                v-if="authStore.isAdmin"
                to="/admin"
                class="px-4 py-2 bg-bloodsa-red text-white rounded-md hover:bg-opacity-90 transition-colors text-sm font-medium"
              >
                Admin Panel
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Main Content Section -->
    <section class="py-12 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            What would you like to do?
          </h2>
          <p class="text-lg text-gray-600 max-w-2xl mx-auto">
            Access your clinical tools and resources to manage workflows efficiently
          </p>
        </div>

        <!-- Feature Cards Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <!-- SOPs Card -->
          <router-link
            to="/sops"
            class="group bg-gradient-to-br from-white to-gray-50 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1"
          >
            <div class="p-8">
              <div class="w-16 h-16 bg-bloodsa-red bg-opacity-10 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-8 h-8 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-3">Standard Operating Procedures</h3>
              <p class="text-gray-600 mb-6 leading-relaxed">
                Access comprehensive SOPs for Anemia, Lymphoma, Myeloma, and General Business procedures.
              </p>
              <div class="flex items-center text-bloodsa-red font-medium group-hover:text-red-700">
                View SOPs
                <svg class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </div>
          </router-link>

          <!-- Transplant Referrals Card -->
          <router-link
            to="/referrals"
            class="group bg-gradient-to-br from-white to-gray-50 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-green-500 transform hover:-translate-y-1"
          >
            <div class="p-8">
              <div class="w-16 h-16 bg-green-100 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-3">Transplant Referrals</h3>
              <p class="text-gray-600 mb-6 leading-relaxed">
                Submit and track transplant referrals via REDCap integration with streamlined workflow.
              </p>
              <div class="flex items-center text-green-600 font-medium group-hover:text-green-700">
                Make Referral
                <svg class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </div>
          </router-link>

          <!-- African HOPeR Registry Card -->
          <router-link
            to="/registry"
            class="group bg-gradient-to-br from-white to-gray-50 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-purple-500 transform hover:-translate-y-1"
          >
            <div class="p-8">
              <div class="w-16 h-16 bg-purple-100 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-3">African HOPeR Registry</h3>
              <p class="text-gray-600 mb-6 leading-relaxed">
                Access registry documentation, training materials, and upload ethics approvals.
              </p>
              <div class="flex items-center text-purple-600 font-medium group-hover:text-purple-700">
                View Registry
                <svg class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </div>
          </router-link>

          <!-- Document Library Card -->
          <router-link
            to="/documents"
            class="group bg-gradient-to-br from-white to-gray-50 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-blue-500 transform hover:-translate-y-1"
          >
            <div class="p-8">
              <div class="w-16 h-16 bg-blue-100 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-3">Document Library</h3>
              <p class="text-gray-600 mb-6 leading-relaxed">
                Access and manage all clinical documents and resources with secure storage.
              </p>
              <div class="flex items-center text-blue-600 font-medium group-hover:text-blue-700">
                Browse Documents
                <svg class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </div>
          </router-link>

          <!-- Training Resources Card -->
          <router-link
            to="/training"
            class="group bg-gradient-to-br from-white to-gray-50 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-yellow-500 transform hover:-translate-y-1"
          >
            <div class="p-8">
              <div class="w-16 h-16 bg-yellow-100 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-8 h-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-3">Training Resources</h3>
              <p class="text-gray-600 mb-6 leading-relaxed">
                Access video tutorials, guidelines, and training materials to enhance your practice.
              </p>
              <div class="flex items-center text-yellow-600 font-medium group-hover:text-yellow-700">
                View Training
                <svg class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </div>
          </router-link>

          <!-- Working Parties Card -->
          <router-link
            to="/working-parties"
            class="group bg-gradient-to-br from-white to-gray-50 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-indigo-500 transform hover:-translate-y-1"
          >
            <div class="p-8">
              <div class="w-16 h-16 bg-indigo-100 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-8 h-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-3">Working Parties</h3>
              <p class="text-gray-600 mb-6 leading-relaxed">
                Collaborate with colleagues and participate in working groups and committees.
              </p>
              <div class="flex items-center text-indigo-600 font-medium group-hover:text-indigo-700">
                Join Groups
                <svg class="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </div>
          </router-link>
        </div>

        <!-- Recent Activity Section (for admins) -->
        <div v-if="authStore.isAdmin" class="mt-16">
          <h3 class="text-2xl font-bold text-gray-900 mb-8">Recent Activity</h3>
          <div class="bg-white rounded-xl shadow-lg p-6">
            <div class="space-y-4">
              <div class="flex items-center space-x-4 p-4 bg-gray-50 rounded-lg">
                <div class="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">
                  <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                  </svg>
                </div>
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">New user registered</p>
                  <p class="text-sm text-gray-500">Dr. Sarah Johnson joined as Haematologist</p>
                </div>
                <span class="text-sm text-gray-400">2 hours ago</span>
              </div>
              
              <div class="flex items-center space-x-4 p-4 bg-gray-50 rounded-lg">
                <div class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                  <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">SOP updated</p>
                  <p class="text-sm text-gray-500">Anemia Management Guidelines v2.1</p>
                </div>
                <span class="text-sm text-gray-400">1 day ago</span>
              </div>
              
              <div class="flex items-center space-x-4 p-4 bg-gray-50 rounded-lg">
                <div class="w-10 h-10 bg-purple-100 rounded-full flex items-center justify-center">
                  <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                  </svg>
                </div>
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">New referral submitted</p>
                  <p class="text-sm text-gray-500">Transplant referral from Dr. Michael Brown</p>
                </div>
                <span class="text-sm text-gray-400">3 days ago</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useInstitutionsStore } from '@/stores/institutions'
import { getUserRoleDisplayName } from '@/types/user'

const authStore = useAuthStore()
const institutionsStore = useInstitutionsStore()

const user = computed(() => authStore.user)

// System stats for admin dashboard
const systemStats = ref({
  totalUsers: 0,
  totalSOPs: 0,
  totalReferrals: 0,
  totalRegistry: 0
})

const userFullName = computed(() => {
  if (!user.value) return 'Guest'
  return `${user.value.profile.firstName} ${user.value.profile.lastName}`
})

const userInitials = computed(() => {
  if (!user.value) return 'G'
  const firstName = user.value.profile.firstName || ''
  const lastName = user.value.profile.lastName || ''
  return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase()
})

const getRoleDisplayName = (role?: string) => {
  if (!role) return ''
  return getUserRoleDisplayName(role as any)
}

const getAdminLevelDisplayName = (adminLevel?: string) => {
  if (!adminLevel) return ''
  switch (adminLevel) {
    case 'super_admin':
      return 'Super Admin'
    case 'user_manager':
      return 'User Manager'
    default:
      return 'Admin'
  }
}

const getInstitutionName = (institutionId?: string): string => {
  if (!institutionId) return 'No Institution'
  const institution = institutionsStore.institutions.find(i => i.id === institutionId)
  return institution ? institution.name : 'Unknown Institution'
}

const formatLastLogin = (lastLogin?: string) => {
  if (!lastLogin) return 'First login'
  const date = new Date(lastLogin)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  
  if (diffMins < 60) return `${diffMins} minutes ago`
  const diffHours = Math.floor(diffMins / 60)
  if (diffHours < 24) return `${diffHours} hours ago`
  const diffDays = Math.floor(diffHours / 24)
  return `${diffDays} days ago`
}

// Load system stats for admin users
const loadSystemStats = async () => {
  if (!authStore.isAdmin) return
  
  try {
    // TODO: Replace with actual API calls when endpoints are available
    // For now, using mock data
    systemStats.value = {
      totalUsers: 24,
      totalSOPs: 12,
      totalReferrals: 8,
      totalRegistry: 15
    }
  } catch (error) {
    console.error('Failed to load system stats:', error)
  }
}

onMounted(async () => {
  // Load institutions for institution name lookup
  await institutionsStore.fetchInstitutions({ isActive: true, limit: 1000 })
  loadSystemStats()
})
</script>


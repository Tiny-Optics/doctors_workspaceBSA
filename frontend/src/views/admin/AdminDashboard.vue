<template>
  <div class="p-6">
    <!-- Page Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Admin Dashboard</h1>
      <p class="text-gray-600 mt-2">Overview of system statistics and user activity</p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <!-- Total Users -->
      <div class="bg-gradient-to-r from-bloodsa-red to-red-700 rounded-xl p-6 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-red-100 text-sm font-medium">Total Users</p>
            <p class="text-3xl font-bold">{{ stats.totalUsers }}</p>
            <p class="text-red-100 text-sm mt-1">+{{ stats.newUsersThisMonth }} this month</p>
          </div>
          <div class="w-12 h-12 bg-white bg-opacity-20 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Active Users -->
      <div class="bg-gradient-to-r from-green-500 to-green-700 rounded-xl p-6 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-green-100 text-sm font-medium">Active Users</p>
            <p class="text-3xl font-bold">{{ stats.activeUsers }}</p>
            <p class="text-green-100 text-sm mt-1">{{ Math.round((stats.activeUsers / stats.totalUsers) * 100) }}% of total</p>
          </div>
          <div class="w-12 h-12 bg-white bg-opacity-20 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- New Registrations -->
      <div class="bg-gradient-to-r from-blue-500 to-blue-700 rounded-xl p-6 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-blue-100 text-sm font-medium">New This Week</p>
            <p class="text-3xl font-bold">{{ stats.newUsersThisWeek }}</p>
            <p class="text-blue-100 text-sm mt-1">+{{ stats.newUsersToday }} today</p>
          </div>
          <div class="w-12 h-12 bg-white bg-opacity-20 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Total Institutions -->
      <div class="bg-gradient-to-r from-purple-500 to-purple-700 rounded-xl p-6 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-purple-100 text-sm font-medium">Total Institutions</p>
            <p class="text-3xl font-bold">{{ stats.totalInstitutions }}</p>
            <p class="text-purple-100 text-sm mt-1">Registered</p>
          </div>
          <div class="w-12 h-12 bg-white bg-opacity-20 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts Section -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
      <!-- Role Distribution Chart -->
      <div class="bg-white rounded-xl shadow-lg p-6">
        <h3 class="text-xl font-semibold text-gray-900 mb-6">User Role Distribution</h3>
        <div class="space-y-4">
          <div v-for="role in roleDistribution" :key="role.name" class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <div 
                class="w-4 h-4 rounded-full" 
                :style="{ backgroundColor: role.color }"
              ></div>
              <span class="text-gray-700 font-medium">{{ role.name }}</span>
            </div>
            <div class="flex items-center space-x-3">
              <div class="w-32 bg-gray-200 rounded-full h-2">
                <div 
                  class="h-2 rounded-full" 
                  :style="{ 
                    backgroundColor: role.color, 
                    width: `${(role.count / stats.totalUsers) * 100}%` 
                  }"
                ></div>
              </div>
              <span class="text-gray-900 font-semibold w-8 text-right">{{ role.count }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="bg-white rounded-xl shadow-lg p-6">
        <h3 class="text-xl font-semibold text-gray-900 mb-6">Recent Activity</h3>
        <div v-if="recentActivity.length === 0" class="text-center py-8">
          <p class="text-gray-500">No recent activity</p>
        </div>
        <div v-else class="space-y-4">
          <div v-for="activity in recentActivity" :key="activity.id" class="flex items-start space-x-3">
            <div 
              class="w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0"
              :class="activity.iconBg"
            >
              <svg class="w-4 h-4" :class="activity.iconColor" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="activity.icon === 'user-plus'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                <path v-else-if="activity.icon === 'login'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
                <path v-else-if="activity.icon === 'settings'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path v-else-if="activity.icon === 'user-check'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path v-else-if="activity.icon === 'user-x'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                <path v-else-if="activity.icon === 'trash'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                <path v-else-if="activity.icon === 'key'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                <path v-else-if="activity.icon === 'shield'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900">{{ activity.title }}</p>
              <p class="text-sm text-gray-500">{{ activity.description }}</p>
              <p class="text-xs text-gray-400 mt-1">{{ activity.time }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="bg-white rounded-xl shadow-lg p-6">
      <h3 class="text-xl font-semibold text-gray-900 mb-6">Quick Actions</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <router-link
          to="/admin/users"
          class="flex items-center space-x-3 p-4 bg-gray-50 rounded-lg hover:bg-bloodsa-red hover:text-white transition-all duration-200 group"
        >
          <div class="w-10 h-10 bg-bloodsa-red bg-opacity-10 rounded-lg flex items-center justify-center group-hover:bg-white group-hover:bg-opacity-20">
            <svg class="w-5 h-5 text-bloodsa-red group-hover:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
            </svg>
          </div>
          <div>
            <p class="font-medium">Manage Users</p>
            <p class="text-sm opacity-75">View and manage all users</p>
          </div>
        </router-link>

        <button
          @click="refreshStats"
          class="flex items-center space-x-3 p-4 bg-gray-50 rounded-lg hover:bg-blue-500 hover:text-white transition-all duration-200 group"
        >
          <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center group-hover:bg-white group-hover:bg-opacity-20">
            <svg class="w-5 h-5 text-blue-600 group-hover:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </div>
          <div>
            <p class="font-medium">Refresh Data</p>
            <p class="text-sm opacity-75">Update statistics</p>
          </div>
        </button>

        <router-link
          to="/admin/audit-logs"
          class="flex items-center space-x-3 p-4 bg-gray-50 rounded-lg hover:bg-purple-500 hover:text-white transition-all duration-200 group"
        >
          <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center group-hover:bg-white group-hover:bg-opacity-20">
            <svg class="w-5 h-5 text-purple-600 group-hover:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <div>
            <p class="font-medium">View Logs</p>
            <p class="text-sm opacity-75">Check audit logs</p>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// Stats data
const stats = ref({
  totalUsers: 0,
  activeUsers: 0,
  newUsersThisMonth: 0,
  newUsersThisWeek: 0,
  newUsersToday: 0,
  totalInstitutions: 0
})

// Role distribution data
const roleDistribution = ref([
  { name: 'Haematologists', count: 0, color: '#DC2626' },
  { name: 'Physicians', count: 0, color: '#059669' },
  { name: 'Data Capturers', count: 0, color: '#7C3AED' },
  { name: 'Admins', count: 0, color: '#EA580C' }
])

// Recent activity data
const recentActivity = ref<any[]>([])

const refreshStats = async () => {
  try {
    console.log('Refreshing stats...')
    await Promise.all([
      loadSystemStats(),
      loadRecentActivity()
    ])
  } catch (error) {
    console.error('Failed to refresh stats:', error)
  }
}

const loadSystemStats = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/stats/admin', {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch admin stats')
    }
    
    const data = await response.json()
    
    // Update stats
    stats.value = {
      totalUsers: data.totalUsers || 0,
      activeUsers: data.activeUsers || 0,
      newUsersThisMonth: data.newUsersThisMonth || 0,
      newUsersThisWeek: data.newUsersThisWeek || 0,
      newUsersToday: data.newUsersToday || 0,
      totalInstitutions: data.totalInstitutions || 0
    }
    
    // Update role distribution
    if (data.roleDistribution && Array.isArray(data.roleDistribution)) {
      const roleMap: Record<string, { name: string; color: string }> = {
        'haematologist': { name: 'Haematologists', color: '#DC2626' },
        'physician': { name: 'Physicians', color: '#059669' },
        'data_capturer': { name: 'Data Capturers', color: '#7C3AED' },
        'admin': { name: 'Admins', color: '#EA580C' }
      }
      
      roleDistribution.value = data.roleDistribution.map((item: { role: string; count: number }) => ({
        name: roleMap[item.role]?.name || item.role,
        count: item.count,
        color: roleMap[item.role]?.color || '#6B7280'
      }))
    }
  } catch (error) {
    console.error('Failed to load system stats:', error)
  }
}

const loadRecentActivity = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/stats/recent-activity?limit=10', {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch recent activity')
    }
    
    const data = await response.json()
    recentActivity.value = data || []
  } catch (error) {
    console.error('Failed to load recent activity:', error)
    recentActivity.value = []
  }
}

onMounted(() => {
  loadSystemStats()
  loadRecentActivity()
})
</script>

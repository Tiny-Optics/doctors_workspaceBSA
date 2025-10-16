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

      <!-- System Health -->
      <div class="bg-gradient-to-r from-purple-500 to-purple-700 rounded-xl p-6 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-purple-100 text-sm font-medium">System Health</p>
            <p class="text-3xl font-bold">{{ stats.systemUptime }}%</p>
            <p class="text-purple-100 text-sm mt-1">Uptime</p>
          </div>
          <div class="w-12 h-12 bg-white bg-opacity-20 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
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
        <div class="space-y-4">
          <div v-for="activity in recentActivity" :key="activity.id" class="flex items-start space-x-3">
            <div 
              class="w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0"
              :class="activity.iconBg"
            >
              <component :is="activity.icon" class="w-4 h-4" :class="activity.iconColor" />
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
import { useUsersStore } from '@/stores/users'

// Icons for recent activity
const UserIcon = {
  template: `
    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
    </svg>
  `
}

const LoginIcon = {
  template: `
    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
    </svg>
  `
}

const SettingsIcon = {
  template: `
    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
    </svg>
  `
}

// Stats data
const stats = ref({
  totalUsers: 24,
  activeUsers: 18,
  newUsersThisMonth: 5,
  newUsersThisWeek: 2,
  newUsersToday: 1,
  systemUptime: 99.9
})

// Role distribution data
const roleDistribution = ref([
  { name: 'Haematologists', count: 8, color: '#DC2626' },
  { name: 'Physicians', count: 6, color: '#059669' },
  { name: 'Data Capturers', count: 7, color: '#7C3AED' },
  { name: 'Admins', count: 3, color: '#EA580C' }
])

// Recent activity data
const recentActivity = ref([
  {
    id: 1,
    title: 'New user registered',
    description: 'Dr. Sarah Johnson joined as Haematologist',
    time: '2 hours ago',
    icon: UserIcon,
    iconBg: 'bg-green-100',
    iconColor: 'text-green-600'
  },
  {
    id: 2,
    title: 'User login',
    description: 'Dr. Michael Brown logged in',
    time: '4 hours ago',
    icon: LoginIcon,
    iconBg: 'bg-blue-100',
    iconColor: 'text-blue-600'
  },
  {
    id: 3,
    title: 'Profile updated',
    description: 'Dr. Emily Davis updated their profile',
    time: '6 hours ago',
    icon: SettingsIcon,
    iconBg: 'bg-purple-100',
    iconColor: 'text-purple-600'
  },
  {
    id: 4,
    title: 'New user registered',
    description: 'Dr. James Wilson joined as Physician',
    time: '1 day ago',
    icon: UserIcon,
    iconBg: 'bg-green-100',
    iconColor: 'text-green-600'
  }
])

const usersStore = useUsersStore()

const refreshStats = async () => {
  try {
    console.log('Refreshing stats...')
    await loadSystemStats()
  } catch (error) {
    console.error('Failed to refresh stats:', error)
  }
}

const loadSystemStats = async () => {
  try {
    // Load real user data to get accurate stats
    const usersData = await usersStore.fetchUsers()
    const totalUsers = usersData.users?.length || 0
    const activeUsers = usersData.users?.filter((user: any) => user.isActive).length || 0
    
    // Update stats with real data
    stats.value = {
      ...stats.value,
      totalUsers,
      activeUsers,
      newUsersThisMonth: Math.floor(totalUsers * 0.2), // 20% of total as new this month
      newUsersThisWeek: Math.floor(totalUsers * 0.08), // 8% of total as new this week
      newUsersToday: Math.floor(Math.random() * 3), // Random 0-2 new today
    }
  } catch (error) {
    console.error('Failed to load system stats:', error)
  }
}

onMounted(() => {
  loadSystemStats()
})
</script>

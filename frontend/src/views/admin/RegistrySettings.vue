<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Registry Settings</h1>
      <p class="mt-2 text-gray-600">
        Manage all aspects of the African HOPeR Registry system.
      </p>
    </div>

    <!-- Settings Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Configuration Card -->
      <router-link
        to="/admin/registry/configuration"
        class="group bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-lg transition-all duration-300 hover:border-bloodsa-red transform hover:-translate-y-1"
      >
        <div class="flex items-center mb-4">
          <div class="w-12 h-12 bg-red-50 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform duration-300">
            <svg class="w-7 h-7 text-bloodsa-red" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-4">
            <h3 class="text-lg font-semibold text-gray-900 group-hover:text-bloodsa-red transition-colors">
              Configuration
            </h3>
            <p class="text-sm text-gray-500">System settings & email</p>
          </div>
        </div>
        <p class="text-gray-600 text-sm mb-4">
          Manage training video, example documents, notification emails, and SMTP settings.
        </p>
        <div class="flex items-center text-bloodsa-red font-medium group-hover:text-red-700">
          <span>Configure Settings</span>
          <svg class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </router-link>

      <!-- Forms Card -->
      <router-link
        to="/admin/registry/forms"
        class="group bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-lg transition-all duration-300 hover:border-blue-500 transform hover:-translate-y-1"
      >
        <div class="flex items-center mb-4">
          <div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform duration-300">
            <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <div class="ml-4">
            <h3 class="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">
              Form Schemas
            </h3>
            <p class="text-sm text-gray-500">Custom form builder</p>
          </div>
        </div>
        <p class="text-gray-600 text-sm mb-4">
          Create and manage custom forms for registry submissions with drag-and-drop field builder.
        </p>
        <div class="flex items-center text-blue-600 font-medium group-hover:text-blue-700">
          <span>Manage Forms</span>
          <svg class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </router-link>

      <!-- Submissions Card -->
      <router-link
        to="/admin/registry/submissions"
        class="group bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-lg transition-all duration-300 hover:border-green-500 transform hover:-translate-y-1"
      >
        <div class="flex items-center mb-4">
          <div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform duration-300">
            <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
            </svg>
          </div>
          <div class="ml-4">
            <h3 class="text-lg font-semibold text-gray-900 group-hover:text-green-600 transition-colors">
              Submissions
            </h3>
            <p class="text-sm text-gray-500">Review & approve</p>
          </div>
        </div>
        <p class="text-gray-600 text-sm mb-4">
          Review, approve, and manage all registry form submissions with advanced filtering.
        </p>
        <div class="flex items-center text-green-600 font-medium group-hover:text-green-700">
          <span>View Submissions</span>
          <svg class="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </router-link>

      <!-- Quick Stats Card -->
      <div class="bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Stats</h3>
        <div class="space-y-3">
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-600">Total Forms</span>
            <span class="text-sm font-medium text-gray-900">{{ stats.totalForms || 0 }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-600">Total Submissions</span>
            <span class="text-sm font-medium text-gray-900">{{ stats.totalSubmissions || 0 }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-600">Pending Reviews</span>
            <span class="text-sm font-medium text-red-600">{{ stats.pendingSubmissions || 0 }}</span>
          </div>
        </div>
      </div>

      <!-- Recent Activity Card -->
      <div class="bg-gradient-to-br from-blue-50 to-blue-100 rounded-xl border border-blue-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h3>
        <div class="space-y-3">
          <div v-if="recentActivity.length === 0" class="text-sm text-gray-500">
            No recent activity
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="activity in recentActivity.slice(0, 3)"
              :key="activity.id"
              class="text-sm text-gray-700"
            >
              {{ activity.description }}
              <span class="text-gray-500 text-xs block">{{ formatTime(activity.timestamp) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- System Status Card -->
      <div class="bg-gradient-to-br from-green-50 to-green-100 rounded-xl border border-green-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">System Status</h3>
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">Registry Status</span>
            <span class="flex items-center text-sm font-medium" :class="systemStatus.registryActive ? 'text-green-600' : 'text-yellow-600'">
              <div class="w-2 h-2 rounded-full mr-2" :class="systemStatus.registryActive ? 'bg-green-500' : 'bg-yellow-500'"></div>
              {{ systemStatus.registryActive ? 'Active' : 'No Active Forms' }}
            </span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">Email Service</span>
            <span class="flex items-center text-sm font-medium" :class="systemStatus.emailConfigured ? 'text-green-600' : 'text-red-600'">
              <div class="w-2 h-2 rounded-full mr-2" :class="systemStatus.emailConfigured ? 'bg-green-500' : 'bg-red-500'"></div>
              {{ systemStatus.emailConfigured ? 'Configured' : 'Not Configured' }}
            </span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">Form Builder</span>
            <span class="flex items-center text-sm font-medium text-green-600">
              <div class="w-2 h-2 bg-green-500 rounded-full mr-2"></div>
              Ready
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { registryService } from '@/services/registryService'
import type { FormSchema, Submission } from '@/services/registryService'

interface Activity {
  id: string
  description: string
  timestamp: Date
}

const stats = ref({
  totalForms: 0,
  totalSubmissions: 0,
  pendingSubmissions: 0
})

const recentActivity = ref<Activity[]>([])

const systemStatus = ref({
  registryActive: false,
  emailConfigured: false,
  formBuilderReady: true
})

async function loadStats() {
  try {
    // Load form schemas (handle case where no forms exist yet)
    let forms: FormSchema[] = []
    let totalFormsCount = 0
    try {
      const formsResponse = await registryService.listFormSchemas({ page: 1, limit: 100 })
      forms = Array.isArray(formsResponse.schemas) ? formsResponse.schemas : []
      totalFormsCount = formsResponse.total || 0
      console.log('Loaded forms:', forms.length)
      console.log('Total forms from API:', totalFormsCount)
      console.log('Forms:', forms.map(f => ({ name: f.formName, active: f.isActive })))
    } catch (error) {
      console.log('No forms found yet:', error)
      forms = []
      totalFormsCount = 0
    }
    const activeForms = forms.filter(f => f.isActive)
    console.log('Active forms:', activeForms.length, activeForms.map(f => f.formName))
    stats.value.totalForms = totalFormsCount
    
    // Load submissions (handle case where no submissions exist yet)
    let submissions: Submission[] = []
    let totalSubmissions = 0
    try {
      const submissionsResponse = await registryService.getAllSubmissions()
      submissions = Array.isArray(submissionsResponse.submissions) ? submissionsResponse.submissions : []
      totalSubmissions = submissionsResponse.total || 0
    } catch (error) {
      console.log('No submissions found yet:', error)
      submissions = []
      totalSubmissions = 0
    }
    
    stats.value.totalSubmissions = totalSubmissions
    stats.value.pendingSubmissions = submissions.filter(s => s.status === 'pending').length
    
    // Generate recent activity from forms and submissions
    await loadRecentActivity(forms, submissions)
    
    // Check system status (pass activeForms info)
    await loadSystemStatus(activeForms.length > 0)
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

async function loadRecentActivity(forms: FormSchema[], submissions: Submission[]) {
  const activities: Activity[] = []
  
  // Add recent form creation activities
  forms
    .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
    .slice(0, 2)
    .forEach(form => {
      activities.push({
        id: `form-${form.id}`,
        description: `Form "${form.formName}" created`,
        timestamp: new Date(form.createdAt)
      })
    })
  
  // Add recent submission activities
  submissions
    .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
    .slice(0, 3)
    .forEach(submission => {
      activities.push({
        id: `submission-${submission.id}`,
        description: `New submission received - Status: ${submission.status}`,
        timestamp: new Date(submission.createdAt)
      })
    })
  
  // Sort all activities by timestamp and keep top 5
  recentActivity.value = activities
    .sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime())
    .slice(0, 5)
}

async function loadSystemStatus(hasActiveForm: boolean) {
  try {
    // Check if email is configured
    const config = await registryService.getConfiguration()
    
    console.log('SMTP Config:', config.smtpConfig)
    
    // Check if all required SMTP fields are configured
    const isEmailConfigured = !!(
      config.smtpConfig &&
      config.smtpConfig.host && 
      config.smtpConfig.host.trim() !== '' &&
      config.smtpConfig.username && 
      config.smtpConfig.username.trim() !== '' &&
      config.smtpConfig.fromEmail && 
      config.smtpConfig.fromEmail.trim() !== '' &&
      config.smtpConfig.fromName && 
      config.smtpConfig.fromName.trim() !== ''
    )
    
    console.log('Email configured status:', isEmailConfigured)
    
    systemStatus.value.emailConfigured = isEmailConfigured
    
    // Registry is active if there's at least one active form (passed from loadStats)
    systemStatus.value.registryActive = hasActiveForm
  } catch (error) {
    console.error('Failed to load system status:', error)
    systemStatus.value.emailConfigured = false
  }
}

function formatTime(timestamp: Date) {
  const now = new Date()
  const diff = now.getTime() - timestamp.getTime()
  const hours = Math.floor(diff / (1000 * 60 * 60))
  
  if (hours < 1) return 'Just now'
  if (hours === 1) return '1 hour ago'
  if (hours < 24) return `${hours} hours ago`
  
  const days = Math.floor(hours / 24)
  if (days === 1) return '1 day ago'
  return `${days} days ago`
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.bg-bloodsa-red {
  background-color: #8B0000;
}
.text-bloodsa-red {
  color: #8B0000;
}
.border-bloodsa-red {
  border-color: #8B0000;
}
</style>

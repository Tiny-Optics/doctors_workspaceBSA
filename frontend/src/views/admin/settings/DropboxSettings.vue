<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h2 class="text-2xl font-bold text-gray-900">Dropbox Configuration</h2>
      <p class="text-gray-600 mt-1">Manage Dropbox OAuth connection for SOP document storage</p>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="bg-white rounded-lg border border-gray-200 p-8">
      <div class="flex items-center justify-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-bloodsa-red"></div>
        <span class="ml-3 text-gray-600">Loading configuration...</span>
      </div>
    </div>

    <!-- Content -->
    <template v-else>
      <!-- Status Card -->
      <div class="bg-white rounded-lg border border-gray-200 p-6">
        <div class="flex items-start justify-between">
          <div class="flex items-start space-x-4">
            <!-- Status Icon -->
            <div 
              :class="[
                'w-12 h-12 rounded-full flex items-center justify-center',
                status?.isConnected ? 'bg-green-100' : 'bg-red-100'
              ]"
            >
              <svg 
                v-if="status?.isConnected" 
                class="w-6 h-6 text-green-600" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <svg 
                v-else 
                class="w-6 h-6 text-red-600" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>

            <!-- Status Info -->
            <div>
              <h3 class="text-lg font-semibold text-gray-900">
                {{ status?.isConnected ? 'Connected' : 'Not Connected' }}
              </h3>
              <p class="text-sm text-gray-600 mt-1">
                {{ status?.isConnected 
                  ? `Last refresh: ${formatDate(status.lastRefreshSuccess)}` 
                  : status?.message || 'Dropbox needs to be configured'
                }}
              </p>
              
              <!-- Error Message -->
              <div v-if="status?.lastError" class="mt-2 text-sm text-red-600 bg-red-50 px-3 py-2 rounded">
                {{ status.lastError }}
              </div>
            </div>
          </div>

          <!-- Action Button -->
          <button
            v-if="status?.isConnected"
            @click="showActions = !showActions"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors"
          >
            Actions
          </button>
        </div>

        <!-- Health Metrics (when connected) -->
        <div v-if="status?.isConnected" class="mt-6 grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="bg-gray-50 rounded-lg p-4">
            <p class="text-sm text-gray-600">Token Expires</p>
            <p class="text-lg font-semibold text-gray-900 mt-1">{{ formatDate(status.tokenExpiry) }}</p>
          </div>
          <div class="bg-gray-50 rounded-lg p-4">
            <p class="text-sm text-gray-600">Consecutive Failures</p>
            <p 
              :class="[
                'text-lg font-semibold mt-1',
                status.consecutiveFailures === 0 ? 'text-green-600' : 'text-red-600'
              ]"
            >
              {{ status.consecutiveFailures }}
            </p>
          </div>
          <div class="bg-gray-50 rounded-lg p-4">
            <p class="text-sm text-gray-600">Parent Folder</p>
            <p class="text-lg font-semibold text-gray-900 mt-1">{{ status.parentFolder }}</p>
          </div>
        </div>

        <!-- Actions Dropdown -->
        <div v-if="showActions && status?.isConnected" class="mt-6 pt-6 border-t border-gray-200 flex flex-wrap gap-3">
          <button
            @click="handleTestConnection"
            :disabled="actionLoading"
            class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Test Connection
          </button>
          <button
            @click="handleForceRefresh"
            :disabled="actionLoading"
            class="px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Force Refresh
          </button>
          <button
            @click="confirmDelete"
            :disabled="actionLoading"
            class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Disconnect
          </button>
        </div>
      </div>

      <!-- Warning (when needs reconnection) -->
      <div v-if="status?.needsReconnection" class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
        <div class="flex items-start">
          <svg class="w-5 h-5 text-yellow-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-yellow-800">Action Required</h3>
            <p class="text-sm text-yellow-700 mt-1">
              Dropbox connection has failed multiple times. Please reconnect using the form below.
            </p>
          </div>
        </div>
      </div>

      <!-- OAuth Configuration Form -->
      <div v-if="!status?.isConnected || status?.needsReconnection" class="bg-white rounded-lg border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">
          {{ status?.isConnected ? 'Reconnect Dropbox' : 'Connect Dropbox' }}
        </h3>

        <!-- Step 1: Get Authorization URL -->
        <div v-if="oauthStep === 1" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Dropbox App Key <span class="text-red-500">*</span>
            </label>
            <input
              v-model="authConfig.appKey"
              type="text"
              placeholder="e.g., abc123def456"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Dropbox App Secret <span class="text-red-500">*</span>
            </label>
            <input
              v-model="authConfig.appSecret"
              type="password"
              placeholder="Enter your app secret"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Parent Folder <span class="text-red-500">*</span>
            </label>
            <input
              v-model="authConfig.parentFolder"
              type="text"
              placeholder="/SOPS"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
            <p class="text-sm text-gray-500 mt-1">The Dropbox folder path where SOPs are stored</p>
          </div>

          <button
            @click="handleInitiateAuth"
            :disabled="!canInitiateAuth || actionLoading"
            class="w-full px-4 py-2 text-sm font-medium text-white bg-bloodsa-red rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {{ actionLoading ? 'Generating Authorization URL...' : 'Get Authorization URL' }}
          </button>

          <!-- Help Text -->
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <div class="flex items-start">
              <svg class="w-5 h-5 text-blue-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div class="ml-3 text-sm text-blue-700">
                <p class="font-medium">How to get Dropbox credentials:</p>
                <ol class="list-decimal list-inside mt-2 space-y-1">
                  <li>Visit <a href="https://www.dropbox.com/developers/apps" target="_blank" class="underline">Dropbox App Console</a></li>
                  <li>Create a new app or select existing</li>
                  <li>Copy the App Key and App Secret from Settings</li>
                </ol>
              </div>
            </div>
          </div>
        </div>

        <!-- Step 2: Enter Authorization Code -->
        <div v-if="oauthStep === 2" class="space-y-4">
          <div class="bg-green-50 border border-green-200 rounded-lg p-4">
            <h4 class="text-sm font-medium text-green-800 mb-2">Step 1: Authorize in Dropbox</h4>
            <a 
              :href="authorizationUrl" 
              target="_blank"
              class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700 transition-colors"
            >
              Open Dropbox Authorization
              <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
            </a>
            <p class="text-sm text-green-700 mt-2">After authorizing, Dropbox will display an authorization code. Copy it and paste below.</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Authorization Code <span class="text-red-500">*</span>
            </label>
            <input
              v-model="authorizationCode"
              type="text"
              placeholder="Paste the code from Dropbox"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
          </div>

          <div class="flex gap-3">
            <button
              @click="oauthStep = 1"
              class="flex-1 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors"
            >
              Back
            </button>
            <button
              @click="handleCompleteAuth"
              :disabled="!authorizationCode || actionLoading"
              class="flex-1 px-4 py-2 text-sm font-medium text-white bg-bloodsa-red rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {{ actionLoading ? 'Connecting...' : 'Complete Authorization' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Success/Error Messages -->
      <div v-if="successMessage" class="bg-green-50 border border-green-200 rounded-lg p-4">
        <div class="flex items-start">
          <svg class="w-5 h-5 text-green-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          <p class="ml-3 text-sm text-green-700">{{ successMessage }}</p>
        </div>
      </div>

      <div v-if="errorMessage" class="bg-red-50 border border-red-200 rounded-lg p-4">
        <div class="flex items-start">
          <svg class="w-5 h-5 text-red-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
          <p class="ml-3 text-sm text-red-700">{{ errorMessage }}</p>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { dropboxAdminService } from '@/services/dropboxAdminService'
import type { DropboxStatus } from '@/types/dropbox'

// State
const loading = ref(false)
const actionLoading = ref(false)
const status = ref<DropboxStatus | null>(null)
const showActions = ref(false)

// OAuth flow
const oauthStep = ref(1)
const authConfig = ref({
  appKey: '',
  appSecret: '',
  parentFolder: '/SOPS',
  redirectUri: ''
})
const authorizationUrl = ref('')
const authorizationCode = ref('')

// Messages
const successMessage = ref('')
const errorMessage = ref('')

// Computed
const canInitiateAuth = computed(() => {
  return authConfig.value.appKey && 
         authConfig.value.appSecret && 
         authConfig.value.parentFolder
})

// Methods
async function loadStatus() {
  loading.value = true
  errorMessage.value = ''
  
  try {
    status.value = await dropboxAdminService.getStatus()
  } catch (error: any) {
    console.error('Failed to load Dropbox status:', error)
    errorMessage.value = error.message || 'Failed to load Dropbox status'
  } finally {
    loading.value = false
  }
}

async function handleInitiateAuth() {
  actionLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const response = await dropboxAdminService.initiateAuth(authConfig.value)
    authorizationUrl.value = response.authUrl
    oauthStep.value = 2
    successMessage.value = 'Authorization URL generated. Please click the button to authorize.'
  } catch (error: any) {
    console.error('Failed to initiate auth:', error)
    errorMessage.value = error.message || 'Failed to initiate authorization'
  } finally {
    actionLoading.value = false
  }
}

async function handleCompleteAuth() {
  actionLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const response = await dropboxAdminService.completeAuth({
      code: authorizationCode.value,
      ...authConfig.value
    })
    
    successMessage.value = response.message
    
    // Reset form and reload status
    oauthStep.value = 1
    authConfig.value = { appKey: '', appSecret: '', parentFolder: '/SOPS', redirectUri: '' }
    authorizationCode.value = ''
    authorizationUrl.value = ''
    
    await loadStatus()
  } catch (error: any) {
    console.error('Failed to complete auth:', error)
    errorMessage.value = error.message || 'Failed to complete authorization'
  } finally {
    actionLoading.value = false
  }
}

async function handleTestConnection() {
  actionLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const response = await dropboxAdminService.testConnection()
    if (response.success) {
      successMessage.value = response.message
    } else {
      errorMessage.value = response.message
    }
  } catch (error: any) {
    console.error('Connection test failed:', error)
    errorMessage.value = error.message || 'Connection test failed'
  } finally {
    actionLoading.value = false
  }
}

async function handleForceRefresh() {
  actionLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const response = await dropboxAdminService.forceRefresh()
    successMessage.value = response.message
    await loadStatus()
  } catch (error: any) {
    console.error('Force refresh failed:', error)
    errorMessage.value = error.message || 'Force refresh failed'
  } finally {
    actionLoading.value = false
  }
}

function confirmDelete() {
  if (confirm('Are you sure you want to disconnect Dropbox? You will need to re-authorize to reconnect.')) {
    handleDelete()
  }
}

async function handleDelete() {
  actionLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const response = await dropboxAdminService.deleteConfiguration()
    successMessage.value = response.message
    showActions.value = false
    await loadStatus()
  } catch (error: any) {
    console.error('Delete configuration failed:', error)
    errorMessage.value = error.message || 'Delete configuration failed'
  } finally {
    actionLoading.value = false
  }
}

function formatDate(dateString: string | undefined): string {
  if (!dateString) return 'Never'
  const date = new Date(dateString)
  if (isNaN(date.getTime())) return 'Never'
  return date.toLocaleString()
}

// Lifecycle
onMounted(() => {
  loadStatus()
})
</script>


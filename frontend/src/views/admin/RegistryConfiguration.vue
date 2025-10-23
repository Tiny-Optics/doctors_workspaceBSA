<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Back Button -->
    <button 
      @click="$router.push('/admin/registry')"
      class="mb-6 inline-flex items-center text-gray-600 hover:text-bloodsa-red transition-colors"
    >
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      Back to Registry Settings
    </button>

    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Registry Configuration</h1>
      <p class="mt-2 text-gray-600">
        Manage the African HOPeR Registry settings, training materials, and email notifications.
      </p>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-bloodsa-red"></div>
    </div>

    <!-- Configuration Form -->
    <div v-else class="space-y-8">
      <!-- Training Video Configuration -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Training Video</h2>
        <div class="space-y-4">
          <div>
            <label for="videoUrl" class="block text-sm font-medium text-gray-700 mb-2">
              YouTube Video URL
              <span class="text-gray-400 text-xs ml-1">(Optional)</span>
            </label>
            <input
              id="videoUrl"
              v-model="config.videoUrl"
              type="url"
              placeholder="https://www.youtube.com/watch?v=..."
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
            <p class="mt-1 text-sm text-gray-500">
              The YouTube video that will be displayed on the registry page
            </p>
          </div>
        </div>
      </div>

      <!-- Example Documents Configuration -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Example Documents</h2>
        <div class="space-y-4">
          <div>
            <label for="documentsPath" class="block text-sm font-medium text-gray-700 mb-2">
              Dropbox Path for Example Documents
              <span class="text-gray-400 text-xs ml-1">(Optional)</span>
            </label>
            <input
              id="documentsPath"
              v-model="config.documentsPath"
              type="text"
              placeholder="/examples/"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
            <p class="mt-1 text-sm text-gray-500">
              The Dropbox path where example documents are stored
            </p>
          </div>
        </div>
      </div>

      <!-- Notification Emails Configuration -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">
          Notification Emails
          <span class="text-red-600 text-lg ml-1">*</span>
        </h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Admin Notification Emails
              <span class="text-red-600">*</span>
            </label>
            <div class="space-y-2">
              <div
                v-for="(email, index) in config.notificationEmails"
                :key="index"
                class="flex items-center space-x-2"
              >
                <input
                  v-model="config.notificationEmails[index]"
                  type="email"
                  class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                />
                <button
                  @click="removeEmail(index)"
                  class="px-3 py-2 text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
            <button
              @click="addEmail"
              class="mt-2 px-4 py-2 text-bloodsa-red hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors flex items-center space-x-2"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              <span>Add Email</span>
            </button>
            <p class="mt-1 text-sm text-gray-500">
              Emails that will receive notifications when users submit forms
            </p>
          </div>
        </div>
      </div>

      <!-- SMTP Configuration -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">
          Email Settings (SMTP)
          <span class="text-red-600 text-lg ml-1">*</span>
        </h2>
        <p class="text-sm text-gray-600 mb-4">
          Required for sending notification emails when users submit forms.
        </p>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="smtpHost" class="block text-sm font-medium text-gray-700 mb-2">
              SMTP Host
              <span class="text-red-600">*</span>
            </label>
            <input
              id="smtpHost"
              v-model="config.smtpConfig.host"
              type="text"
              placeholder="smtp.gmail.com"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
          </div>
          <div>
            <label for="smtpPort" class="block text-sm font-medium text-gray-700 mb-2">
              SMTP Port
              <span class="text-red-600">*</span>
            </label>
            <input
              id="smtpPort"
              v-model.number="config.smtpConfig.port"
              type="number"
              placeholder="587"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
          </div>
          <div>
            <label for="smtpUsername" class="block text-sm font-medium text-gray-700 mb-2">
              SMTP Username
              <span class="text-red-600">*</span>
            </label>
            <input
              id="smtpUsername"
              v-model="config.smtpConfig.username"
              type="text"
              placeholder="noreply@bloodsa.org.za"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
          </div>
          <div>
            <label for="smtpPassword" class="block text-sm font-medium text-gray-700 mb-2">
              SMTP Password
              <span v-if="!config.id" class="text-red-600">*</span>
              <span v-else class="text-gray-400 text-xs ml-1">(Leave blank to keep existing)</span>
            </label>
            <div class="relative">
              <input
                id="smtpPassword"
                v-model="config.smtpConfig.password"
                :type="showPassword ? 'text' : 'password'"
                :placeholder="config.id ? 'Leave blank to keep existing password' : 'Enter password'"
                class="w-full px-4 py-2 pr-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
              />
              <button
                @click="showPassword = !showPassword"
                type="button"
                class="absolute inset-y-0 right-0 pr-3 flex items-center"
              >
                <svg v-if="!showPassword" class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                <svg v-else class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21" />
                </svg>
              </button>
            </div>
          </div>
          <div>
            <label for="fromEmail" class="block text-sm font-medium text-gray-700 mb-2">
              From Email
              <span class="text-red-600">*</span>
            </label>
            <input
              id="fromEmail"
              v-model="config.smtpConfig.fromEmail"
              type="email"
              placeholder="noreply@bloodsa.org.za"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
          </div>
          <div>
            <label for="fromName" class="block text-sm font-medium text-gray-700 mb-2">
              From Name
              <span class="text-red-600">*</span>
            </label>
            <input
              id="fromName"
              v-model="config.smtpConfig.fromName"
              type="text"
              placeholder="BLOODSA Registry"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
          </div>
        </div>
      </div>

      <!-- Test Email Section -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Test Email Configuration</h2>
        <div class="space-y-4">
          <div>
            <label for="testEmail" class="block text-sm font-medium text-gray-700 mb-2">
              Test Email Address
            </label>
            <div class="flex space-x-4">
              <input
                id="testEmail"
                v-model="testEmail"
                type="email"
                placeholder="test@example.com"
                class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
              />
              <button
                @click="sendTestEmail"
                :disabled="sendingTest"
                class="px-6 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
              >
                <svg v-if="sendingTest" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                <span>{{ sendingTest ? 'Sending...' : 'Send Test Email' }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex justify-end space-x-4">
        <button
          @click="resetForm"
          class="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
        >
          Reset
        </button>
        <button
          @click="saveConfiguration"
          :disabled="saving"
          class="px-6 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
        >
          <svg v-if="saving" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span>{{ saving ? 'Saving...' : 'Save Configuration' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { registryService } from '@/services/registryService'
import type { RegistryConfig } from '@/services/registryService'
import { useToast } from '@/composables/useToast'

const loading = ref(false)
const saving = ref(false)
const sendingTest = ref(false)
const showPassword = ref(false)
const toast = useToast()
const testEmail = ref('')

const config = ref<RegistryConfig>({
  id: '',
  videoUrl: '',
  documentsPath: '',
  notificationEmails: [],
  smtpConfig: {
    host: '',
    port: 587,
    username: '',
    password: '',
    fromEmail: '',
    fromName: ''
  },
  createdAt: '',
  updatedAt: ''
})

async function loadConfiguration() {
  loading.value = true
  try {
    const response = await registryService.getConfiguration()
    config.value = response
    
    // Ensure smtpConfig is always defined
    if (!config.value.smtpConfig) {
      config.value.smtpConfig = {
        host: '',
        port: 587,
        username: '',
        password: '',
        fromEmail: '',
        fromName: ''
      }
    }
    
    // Ensure notificationEmails is always an array
    if (!config.value.notificationEmails) {
      config.value.notificationEmails = []
    }
  } catch (error) {
    console.error('Failed to load configuration:', error)
    toast.error('Failed to load configuration. Please try again.')
  } finally {
    loading.value = false
  }
}

async function saveConfiguration() {
  // Remove empty email addresses before saving
  config.value.notificationEmails = config.value.notificationEmails.filter(email => email.trim() !== '')
  
  // Collect all missing required fields
  const missingFields: string[] = []
  
  // Validate notification emails
  if (config.value.notificationEmails.length === 0) {
    missingFields.push('At least one Notification Email')
  }
  
  // Validate email format
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  const invalidNotificationEmails = config.value.notificationEmails.filter(email => !emailRegex.test(email))
  if (invalidNotificationEmails.length > 0) {
    toast.error(`Invalid notification email format: ${invalidNotificationEmails.join(', ')}`)
    return
  }
  
  // Validate SMTP configuration
  if (!config.value.smtpConfig.host || config.value.smtpConfig.host.trim() === '') {
    missingFields.push('SMTP Host')
  }
  
  if (!config.value.smtpConfig.port || config.value.smtpConfig.port <= 0) {
    missingFields.push('SMTP Port')
  }
  
  if (!config.value.smtpConfig.username || config.value.smtpConfig.username.trim() === '') {
    missingFields.push('SMTP Username')
  }
  
  // Only require password if config doesn't have an ID (new config) or password is explicitly cleared
  // If updating existing config, password might be encrypted in the response
  if (!config.value.id && (!config.value.smtpConfig.password || config.value.smtpConfig.password.trim() === '')) {
    missingFields.push('SMTP Password')
  }
  
  if (!config.value.smtpConfig.fromEmail || config.value.smtpConfig.fromEmail.trim() === '') {
    missingFields.push('From Email')
  } else if (!emailRegex.test(config.value.smtpConfig.fromEmail)) {
    toast.error('From Email has an invalid format')
    return
  }
  
  if (!config.value.smtpConfig.fromName || config.value.smtpConfig.fromName.trim() === '') {
    missingFields.push('From Name')
  }
  
  // If there are missing fields, show detailed error message
  if (missingFields.length > 0) {
    const fieldList = missingFields.join(', ')
    toast.error(`Missing required fields: ${fieldList}`)
    return
  }
  
  saving.value = true
  try {
    await registryService.updateConfiguration(config.value)
    toast.success('Configuration saved successfully!')
    // Reload to get latest data
    await loadConfiguration()
  } catch (error) {
    console.error('Failed to save configuration:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    toast.error(`Failed to save configuration: ${errorMessage}`)
  } finally {
    saving.value = false
  }
}

async function sendTestEmail() {
  if (!testEmail.value) {
    toast.warning('Please enter an email address')
    return
  }
  
  // Validate email format
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(testEmail.value)) {
    toast.error('Please enter a valid email address')
    return
  }
  
  sendingTest.value = true
  try {
    await registryService.sendTestEmail(testEmail.value)
    toast.success(`Test email sent successfully to ${testEmail.value}`)
    testEmail.value = ''
  } catch (error) {
    console.error('Failed to send test email:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    toast.error(`Failed to send test email: ${errorMessage}`)
  } finally {
    sendingTest.value = false
  }
}

function addEmail() {
  config.value.notificationEmails.push('')
}

function removeEmail(index: number) {
  config.value.notificationEmails.splice(index, 1)
}

function resetForm() {
  loadConfiguration()
}

onMounted(() => {
  loadConfiguration()
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

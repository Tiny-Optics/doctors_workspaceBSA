<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h2 class="text-2xl font-bold text-gray-900">Email & SMTP Settings</h2>
      <p class="text-gray-600 mt-1">Configure SMTP settings for email notifications and system communications</p>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-bloodsa-red"></div>
    </div>

    <!-- Configuration Form -->
    <div v-else class="space-y-6">
      <!-- SMTP Configuration -->
      <div class="bg-white rounded-lg border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">SMTP Configuration</h3>
        <p class="text-sm text-gray-600 mb-6">
          Configure SMTP settings for sending email notifications. These settings are used for system emails and notifications.
        </p>

        <form @submit.prevent="saveSMTPConfig" class="space-y-6">
          <!-- SMTP Host and Port -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label for="smtpHost" class="block text-sm font-medium text-gray-700 mb-2">
                SMTP Host
                <span class="text-red-600">*</span>
              </label>
              <input
                id="smtpHost"
                v-model="smtpConfig.host"
                type="text"
                placeholder="smtp.gmail.com"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                required
              />
            </div>
            <div>
              <label for="smtpPort" class="block text-sm font-medium text-gray-700 mb-2">
                SMTP Port
                <span class="text-red-600">*</span>
              </label>
              <input
                id="smtpPort"
                v-model.number="smtpConfig.port"
                type="number"
                placeholder="587"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                required
              />
            </div>
          </div>

          <!-- SMTP Username and Password -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label for="smtpUsername" class="block text-sm font-medium text-gray-700 mb-2">
                SMTP Username
                <span class="text-red-600">*</span>
              </label>
              <input
                id="smtpUsername"
                v-model="smtpConfig.username"
                type="text"
                placeholder="noreply@bloodsa.org.za"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                required
              />
            </div>
            <div>
              <label for="smtpPassword" class="block text-sm font-medium text-gray-700 mb-2">
                SMTP Password
                <span v-if="!hasExistingConfig" class="text-red-600">*</span>
                <span v-else class="text-gray-400 text-xs ml-1">(Leave blank to keep existing)</span>
              </label>
              <div class="relative">
                <input
                  id="smtpPassword"
                  v-model="smtpConfig.password"
                  :type="showPassword ? 'text' : 'password'"
                  :placeholder="hasExistingConfig ? 'Leave blank to keep existing password' : 'Enter password'"
                  class="w-full px-4 py-2 pr-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  :required="!hasExistingConfig"
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
          </div>

          <!-- From Email and Name -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label for="fromEmail" class="block text-sm font-medium text-gray-700 mb-2">
                From Email
                <span class="text-red-600">*</span>
              </label>
              <input
                id="fromEmail"
                v-model="smtpConfig.fromEmail"
                type="email"
                placeholder="noreply@bloodsa.org.za"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                required
              />
            </div>
            <div>
              <label for="fromName" class="block text-sm font-medium text-gray-700 mb-2">
                From Name
                <span class="text-red-600">*</span>
              </label>
              <input
                id="fromName"
                v-model="smtpConfig.fromName"
                type="text"
                placeholder="BLOODSA Registry"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                required
              />
            </div>
          </div>

          <!-- Configuration Status -->
          <div class="bg-gray-50 rounded-lg p-4">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div v-if="smtpConfig.isComplete" class="w-3 h-3 bg-green-400 rounded-full"></div>
                <div v-else class="w-3 h-3 bg-yellow-400 rounded-full"></div>
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium text-gray-900">
                  {{ smtpConfig.isComplete ? 'Configuration Complete' : 'Configuration Incomplete' }}
                </p>
                <p class="text-sm text-gray-600">
                  {{ smtpConfig.isComplete 
                    ? 'All required SMTP settings are configured and ready to use.' 
                    : 'Please fill in all required fields to complete the configuration.' 
                  }}
                </p>
              </div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex justify-end space-x-4">
            <button
              @click="resetForm"
              type="button"
              class="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Reset
            </button>
            <button
              type="submit"
              :disabled="saving"
              class="px-6 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
            >
              <svg v-if="saving" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <span>{{ saving ? 'Saving...' : 'Save Configuration' }}</span>
            </button>
          </div>
        </form>
      </div>

      <!-- Test Email Section -->
      <div class="bg-white rounded-lg border border-gray-200 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Test Email Configuration</h3>
        <p class="text-sm text-gray-600 mb-4">
          Send a test email to verify that your SMTP configuration is working correctly.
        </p>

        <form @submit.prevent="sendTestEmail" class="space-y-4">
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
                required
              />
              <button
                type="submit"
                :disabled="sendingTest || !smtpConfig.isComplete"
                class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
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
            <p v-if="!smtpConfig.isComplete" class="mt-2 text-sm text-yellow-600">
              Please complete the SMTP configuration above before sending a test email.
            </p>
          </div>
        </form>
      </div>
    </div>

    <!-- Toast Notifications -->
    <ToastNotification />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useToast } from '@/composables/useToast'
import ToastNotification from '@/components/ToastNotification.vue'

const loading = ref(false)
const saving = ref(false)
const sendingTest = ref(false)
const showPassword = ref(false)
const toast = useToast()
const testEmail = ref('')

// SMTP Configuration
const smtpConfig = ref({
  host: '',
  port: 587,
  username: '',
  password: '',
  fromEmail: '',
  fromName: '',
  isComplete: false
})

// Check if we have existing configuration (similar to RegistryConfiguration.vue)
const hasExistingConfig = computed(() => {
  return smtpConfig.value.host !== '' || smtpConfig.value.username !== '' || smtpConfig.value.fromEmail !== ''
})

// Load SMTP configuration
async function loadSMTPConfig() {
  loading.value = true
  try {
    const response = await fetch('/api/admin/registry/smtp-config', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error(`HTTP error ${response.status}`)
    }

    const config = await response.json()
    smtpConfig.value = {
      host: config.host || '',
      port: config.port || 587,
      username: config.username || '',
      password: '', // Never load password from server
      fromEmail: config.fromEmail || '',
      fromName: config.fromName || '',
      isComplete: config.isComplete || false
    }
  } catch (error) {
    console.error('Failed to load SMTP configuration:', error)
    toast.error('Failed to load SMTP configuration. Please try again.')
  } finally {
    loading.value = false
  }
}

// Save SMTP configuration
async function saveSMTPConfig() {
  // Validate required fields
  if (!smtpConfig.value.host || !smtpConfig.value.username || !smtpConfig.value.fromEmail || !smtpConfig.value.fromName) {
    toast.error('Please fill in all required fields')
    return
  }

  // Only require password if no existing config or password is explicitly provided
  if (!hasExistingConfig.value && !smtpConfig.value.password) {
    toast.error('Please enter the SMTP password')
    return
  }

  saving.value = true
  try {
    const response = await fetch('/api/admin/registry/smtp-config', {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        host: smtpConfig.value.host,
        port: smtpConfig.value.port,
        username: smtpConfig.value.username,
        ...(smtpConfig.value.password && { password: smtpConfig.value.password }),
        fromEmail: smtpConfig.value.fromEmail,
        fromName: smtpConfig.value.fromName
      })
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || `HTTP error ${response.status}`)
    }

    const updatedConfig = await response.json()
    smtpConfig.value.isComplete = updatedConfig.isComplete
    
    toast.success('SMTP configuration saved successfully!')
    
    // Clear password field after successful save (only if password was provided)
    if (smtpConfig.value.password) {
      smtpConfig.value.password = ''
    }
  } catch (error) {
    console.error('Failed to save SMTP configuration:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    toast.error(`Failed to save SMTP configuration: ${errorMessage}`)
  } finally {
    saving.value = false
  }
}

// Send test email
async function sendTestEmail() {
  if (!testEmail.value) {
    toast.warning('Please enter an email address')
    return
  }

  if (!smtpConfig.value.isComplete) {
    toast.warning('Please complete the SMTP configuration first')
    return
  }

  sendingTest.value = true
  try {
    const response = await fetch('/api/admin/registry/test-email', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        email: testEmail.value
      })
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || `HTTP error ${response.status}`)
    }

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

// Reset form to original values
function resetForm() {
  loadSMTPConfig()
}

onMounted(() => {
  loadSMTPConfig()
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


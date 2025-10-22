<template>
  <div class="p-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Referral Settings</h1>
        <p class="text-gray-600 mt-2">Configure REDCap referral system settings</p>
      </div>
      <button
        @click="showTestModal = true"
        :disabled="loading || !config.isConfigured || !config.isEnabled"
        class="bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-opacity-90 transition-colors font-medium flex items-center space-x-2 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
        </svg>
        <span>Test Link</span>
      </button>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error</h3>
          <div class="mt-2 text-sm text-red-700">
            {{ error }}
          </div>
        </div>
        <div class="ml-auto pl-3">
          <button @click="error = null" class="text-red-400 hover:text-red-600">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Configuration Form -->
    <div class="bg-white rounded-xl shadow-lg p-6">
      <div class="space-y-6">
        <!-- Status Overview -->
        <div class="bg-gray-50 rounded-lg p-4">
          <h3 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
            <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            System Status
          </h3>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="text-center">
              <div class="text-2xl font-bold" :class="config.isConfigured ? 'text-green-600' : 'text-red-600'">
                {{ config.isConfigured ? 'Configured' : 'Not Configured' }}
              </div>
              <div class="text-sm text-gray-500">Configuration</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold" :class="config.isEnabled ? 'text-green-600' : 'text-yellow-600'">
                {{ config.isEnabled ? 'Enabled' : 'Disabled' }}
              </div>
              <div class="text-sm text-gray-500">Status</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold" :class="config.redcapUrl ? 'text-green-600' : 'text-red-600'">
                {{ config.redcapUrl ? 'Set' : 'Not Set' }}
              </div>
              <div class="text-sm text-gray-500">REDCap URL</div>
            </div>
          </div>
        </div>

        <!-- Configuration Form -->
        <form @submit.prevent="saveConfiguration">
          <div class="space-y-6">
            <!-- REDCap URL -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                REDCap URL <span class="text-red-500">*</span>
              </label>
              <input
                v-model="formData.redcapUrl"
                type="url"
                required
                placeholder="https://redcap.example.com/surveys/12345"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
              />
              <p class="text-xs text-gray-500 mt-1">
                Enter the full URL to your REDCap survey or form
              </p>
            </div>

            <!-- Enable/Disable Toggle -->
            <div>
              <label class="flex items-center space-x-3">
                <input
                  v-model="formData.isEnabled"
                  type="checkbox"
                  class="h-4 w-4 text-bloodsa-red focus:ring-bloodsa-red border-gray-300 rounded"
                />
                <div>
                  <span class="text-sm font-medium text-gray-700">Enable Referral System</span>
                  <p class="text-xs text-gray-500">
                    When enabled, users can access the referral link. When disabled, users will see a "not available" message.
                  </p>
                </div>
              </label>
            </div>

            <!-- Current Configuration Info -->
            <div v-if="config.isConfigured" class="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <h4 class="text-sm font-medium text-blue-800 mb-2">Current Configuration</h4>
              <div class="space-y-2 text-sm text-blue-700">
                <div class="flex justify-between">
                  <span>URL:</span>
                  <span class="font-mono break-all">{{ config.redcapUrl || 'Not set' }}</span>
                </div>
                <div class="flex justify-between">
                  <span>Status:</span>
                  <span class="font-medium">{{ config.isEnabled ? 'Enabled' : 'Disabled' }}</span>
                </div>
                <div v-if="config.updatedAt" class="flex justify-between">
                  <span>Last Updated:</span>
                  <span>{{ formatDate(config.updatedAt) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Form Actions -->
          <div class="flex justify-end space-x-3 mt-8 pt-6 border-t border-gray-200">
            <button
              type="button"
              @click="showResetModal = true"
              class="px-6 py-2 text-gray-700 bg-gray-200 rounded-lg hover:bg-gray-300 transition-colors"
            >
              Reset
            </button>
            <button
              type="submit"
              :disabled="loading || !formData.redcapUrl"
              class="px-6 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-opacity-90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
            >
              <svg v-if="loading" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>{{ loading ? 'Saving...' : 'Save Configuration' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Test Link Modal -->
    <div v-if="showTestModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="closeTestModal">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <!-- Test Icon -->
          <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100 mb-4">
            <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
          </div>
          
          <!-- Modal Content -->
          <div class="text-center">
            <h3 class="text-lg font-medium text-gray-900 mb-2">Test Referral Link</h3>
            <div class="mt-2 px-7 py-3">
              <p class="text-sm text-gray-500 mb-4">
                This will open the REDCap referral link in a new tab to test if it's working correctly.
              </p>
              
              <!-- Current URL Display -->
              <div class="bg-gray-50 rounded-lg p-4 mb-4">
                <p class="text-xs text-gray-600 mb-1">Current URL:</p>
                <p class="text-sm font-mono text-gray-900 break-all">{{ config.redcapUrl }}</p>
              </div>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex justify-center space-x-3 mt-6">
              <button
                @click="closeTestModal"
                :disabled="testing"
                class="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Cancel
              </button>
              <button
                @click="testReferralLink"
                :disabled="testing"
                class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
              >
                <svg v-if="testing" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ testing ? 'Testing...' : 'Open Link' }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Reset Confirmation Modal -->
    <div v-if="showResetModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="closeResetModal">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <!-- Warning Icon -->
          <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-yellow-100 mb-4">
            <svg class="h-6 w-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 18.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          
          <!-- Modal Content -->
          <div class="text-center">
            <h3 class="text-lg font-medium text-gray-900 mb-2">Reset Configuration</h3>
            <div class="mt-2 px-7 py-3">
              <p class="text-sm text-gray-500 mb-4">
                Are you sure you want to reset the form? This will discard any unsaved changes and restore the form to the last saved configuration.
              </p>
              
              <!-- Current vs Form Data Display -->
              <div class="bg-gray-50 rounded-lg p-4 mb-4 text-left">
                <p class="text-xs text-gray-600 mb-2 font-medium">Current saved configuration:</p>
                <div class="text-sm text-gray-700 space-y-1">
                  <div>URL: {{ config.redcapUrl || 'Not set' }}</div>
                  <div>Status: {{ config.isEnabled ? 'Enabled' : 'Disabled' }}</div>
                </div>
              </div>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex justify-center space-x-3 mt-6">
              <button
                @click="closeResetModal"
                class="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors"
              >
                Cancel
              </button>
              <button
                @click="confirmReset"
                class="px-4 py-2 bg-yellow-600 text-white rounded-md hover:bg-yellow-700 transition-colors"
              >
                Reset Form
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { referralService, type ReferralConfig, type UpdateReferralConfigRequest } from '@/services/referralService'
import { useToast } from '@/composables/useToast'

const toast = useToast()

// Reactive data
const loading = ref(false)
const testing = ref(false)
const error = ref<string | null>(null)
const showTestModal = ref(false)
const showResetModal = ref(false)

// Configuration data
const config = ref<ReferralConfig>({
  redcapUrl: '',
  isEnabled: false,
  isConfigured: false
})

// Form data
const formData = ref<UpdateReferralConfigRequest>({
  redcapUrl: '',
  isEnabled: false
})

// Computed properties
const hasChanges = computed(() => {
  return formData.value.redcapUrl !== config.value.redcapUrl || 
         formData.value.isEnabled !== config.value.isEnabled
})

// Methods
const loadConfiguration = async () => {
  try {
    loading.value = true
    error.value = null
    
    const response = await referralService.getAdminReferralConfig()
    config.value = response
    
    // Initialize form data
    formData.value = {
      redcapUrl: response.redcapUrl || '',
      isEnabled: response.isEnabled || false
    }
  } catch (err) {
    console.error('Failed to load referral configuration:', err)
    if (err instanceof Error && err.message.includes('not configured')) {
      // Configuration doesn't exist yet, that's okay
      config.value = {
        redcapUrl: '',
        isEnabled: false,
        isConfigured: false
      }
      formData.value = {
        redcapUrl: '',
        isEnabled: false
      }
    } else {
      error.value = err instanceof Error ? err.message : 'Failed to load configuration'
      toast.error('Failed to load referral configuration')
    }
  } finally {
    loading.value = false
  }
}

const saveConfiguration = async () => {
  try {
    loading.value = true
    error.value = null
    
    const updateData: UpdateReferralConfigRequest = {
      redcapUrl: formData.value.redcapUrl,
      isEnabled: formData.value.isEnabled
    }
    
    const response = await referralService.updateReferralConfig(updateData)
    config.value = response
    
    toast.success('Referral configuration saved successfully')
  } catch (err) {
    console.error('Failed to save referral configuration:', err)
    error.value = err instanceof Error ? err.message : 'Failed to save configuration'
    toast.error('Failed to save referral configuration')
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  formData.value = {
    redcapUrl: config.value.redcapUrl || '',
    isEnabled: config.value.isEnabled || false
  }
  error.value = null
}

const testReferralLink = async () => {
  try {
    testing.value = true
    
    // Log the access and get the redirect URL
    const response = await referralService.logReferralAccess()
    
    // Open the link in a new tab
    window.open(response.redirectUrl, '_blank')
    
    toast.success('Referral link opened in new tab')
    closeTestModal()
  } catch (err) {
    console.error('Failed to test referral link:', err)
    toast.error('Failed to test referral link')
  } finally {
    testing.value = false
  }
}

const closeTestModal = () => {
  showTestModal.value = false
}

const closeResetModal = () => {
  showResetModal.value = false
}

const confirmReset = () => {
  resetForm()
  closeResetModal()
  toast.success('Form reset to last saved configuration')
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

// Lifecycle
onMounted(() => {
  loadConfiguration()
})
</script>

<style scoped>
/* Custom styles if needed */
</style>

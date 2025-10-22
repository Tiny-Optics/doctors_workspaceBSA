<template>
  <div class="min-h-screen bg-bloodsa-red">
    <!-- Header Section -->
    <section class="relative overflow-hidden">
      <div class="absolute inset-0 opacity-5">
        <div class="absolute inset-0" style="background-image: repeating-linear-gradient(45deg, #8B0000 0, #8B0000 1px, transparent 0, transparent 50%); background-size: 10px 10px;"></div>
      </div>

      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-20">
        <!-- Back Button -->
        <button 
          @click="$router.push({ name: 'dashboard' })"
          class="mb-6 inline-flex items-center text-white hover:text-white/80 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        
        <div class="text-center">
          <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">Transplant Referrals</h1>
          <p class="text-white max-w-2xl mx-auto">
            Access the REDCap referral system to submit patient referrals for bone marrow biopsies and other transplant-related procedures.
          </p>
        </div>
      </div>
    </section>

    <!-- Loading State -->
    <section v-if="loading" class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-bloodsa-red mx-auto"></div>
          <p class="text-gray-600 mt-4">Loading referral system...</p>
        </div>
      </div>
    </section>

    <!-- Error State -->
    <section v-else-if="error" class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <svg class="w-16 h-16 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">{{ error }}</h3>
          <button
            @click="loadReferralConfig"
            class="mt-4 px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors"
          >
            Try Again
          </button>
        </div>
      </div>
    </section>

    <!-- Referral System Available -->
    <section v-else-if="config.isConfigured && config.isEnabled" class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-6 text-center">
          <!-- Instructions -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-gray-900 mb-8">How to Submit a Referral</h2>
            <div class="max-w-3xl mx-auto space-y-8">
              <div class="text-center">
                <div class="inline-flex items-center justify-center h-12 w-12 rounded-full bg-bloodsa-red text-white text-lg font-bold mb-4">
                  1
                </div>
                <h3 class="text-xl font-semibold text-gray-900 mb-2">Click "Access REDCap Form"</h3>
                <p class="text-gray-600">This will open the REDCap referral form in a new tab.</p>
              </div>
              <div class="text-center">
                <div class="inline-flex items-center justify-center h-12 w-12 rounded-full bg-bloodsa-red text-white text-lg font-bold mb-4">
                  2
                </div>
                <h3 class="text-xl font-semibold text-gray-900 mb-2">Complete the Form</h3>
                <p class="text-gray-600">Fill out all required patient information and clinical details.</p>
              </div>
              <div class="text-center">
                <div class="inline-flex items-center justify-center h-12 w-12 rounded-full bg-bloodsa-red text-white text-lg font-bold mb-4">
                  3
                </div>
                <h3 class="text-xl font-semibold text-gray-900 mb-2">Submit the Referral</h3>
                <p class="text-gray-600">Review your information and submit the form. You'll receive confirmation once submitted.</p>
              </div>
            </div>
          </div>

          <!-- Access Button -->
          <div class="text-center">
            <button
              @click="accessReferralForm"
              :disabled="accessing"
              class="inline-flex items-center px-8 py-4 border border-transparent text-lg font-medium rounded-lg text-white bg-bloodsa-red hover:bg-opacity-90 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-bloodsa-red disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <svg v-if="accessing" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <svg v-else class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
              {{ accessing ? 'Opening Form...' : 'Access REDCap Form' }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- Referral System Unavailable -->
    <section v-else class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <svg class="w-16 h-16 text-yellow-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 18.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">Referral System Unavailable</h3>
          <p class="text-gray-600 mb-4">
            The referral system is currently not available. This may be due to maintenance or configuration updates.
          </p>
          <p class="text-sm text-gray-500">
            Please try again later or contact your system administrator if the issue persists.
          </p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { referralService, type ReferralConfigResponse } from '@/services/referralService'
import { useToast } from '@/composables/useToast'

const toast = useToast()

// Reactive data
const loading = ref(true)
const accessing = ref(false)
const error = ref<string | null>(null)
const config = ref<ReferralConfigResponse>({
  isConfigured: false,
  isEnabled: false,
  redcapUrl: ''
})

// Methods
const loadReferralConfig = async () => {
  try {
    loading.value = true
    error.value = null
    
    const response = await referralService.getReferralConfig()
    config.value = response
  } catch (err) {
    console.error('Failed to load referral configuration:', err)
    error.value = err instanceof Error ? err.message : 'Failed to load referral system status'
  } finally {
    loading.value = false
  }
}

const accessReferralForm = async () => {
  try {
    accessing.value = true
    
    // Log the access and get the redirect URL
    const response = await referralService.logReferralAccess()
    
    // Open the link in a new tab
    window.open(response.redirectUrl, '_blank')
    
    toast.success('REDCap form opened in new tab')
  } catch (err) {
    console.error('Failed to access referral form:', err)
    toast.error('Failed to open referral form')
  } finally {
    accessing.value = false
  }
}

// Lifecycle
onMounted(() => {
  loadReferralConfig()
})
</script>

<style scoped>
/* Custom styles if needed */
</style>

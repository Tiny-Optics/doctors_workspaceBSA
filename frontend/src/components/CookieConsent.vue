<template>
  <Transition name="cookie-banner">
    <div
      v-if="showBanner"
      class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 shadow-lg z-50"
    >
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
          <!-- Content -->
          <div class="flex-1">
            <div class="flex items-start space-x-3">
              <div class="flex-shrink-0 mt-1">
                <svg class="w-6 h-6 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
                </svg>
              </div>
              <div class="flex-1">
                <h3 class="text-sm font-semibold text-gray-900 mb-1">
                  Cookie Consent
                </h3>
                <p class="text-sm text-gray-600 leading-relaxed">
                  We use cookies to enhance your experience, analyze site usage, and assist in our marketing efforts. 
                  By clicking "Accept All", you consent to our use of cookies. 
                  <router-link 
                    to="/privacy-policy" 
                    class="text-bloodsa-red hover:text-bloodsa-light-red underline transition-colors duration-200"
                  >
                    Learn more in our Privacy Policy
                  </router-link>.
                </p>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex flex-col sm:flex-row gap-3 w-full sm:w-auto">
            <button
              @click="acceptCookies"
              class="px-6 py-2 bg-bloodsa-red text-white rounded-md font-medium hover:bg-bloodsa-light-red transition-colors duration-200 text-sm whitespace-nowrap"
            >
              Accept All
            </button>
            <button
              @click="rejectCookies"
              class="px-6 py-2 bg-gray-100 text-gray-700 rounded-md font-medium hover:bg-gray-200 transition-colors duration-200 text-sm whitespace-nowrap"
            >
              Reject
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const COOKIE_CONSENT_KEY = 'cookie_consent'
const showBanner = ref(false)

onMounted(() => {
  // Check if user has already made a choice
  const consent = localStorage.getItem(COOKIE_CONSENT_KEY)
  if (!consent) {
    // Show banner after a short delay for better UX
    setTimeout(() => {
      showBanner.value = true
    }, 500)
  }
})

function acceptCookies() {
  localStorage.setItem(COOKIE_CONSENT_KEY, 'accepted')
  showBanner.value = false
}

function rejectCookies() {
  localStorage.setItem(COOKIE_CONSENT_KEY, 'rejected')
  showBanner.value = false
}
</script>

<style scoped>
/* Cookie banner animations */
.cookie-banner-enter-active,
.cookie-banner-leave-active {
  transition: all 0.3s ease-in-out;
}

.cookie-banner-enter-from {
  opacity: 0;
  transform: translateY(100%);
}

.cookie-banner-leave-to {
  opacity: 0;
  transform: translateY(100%);
}
</style>


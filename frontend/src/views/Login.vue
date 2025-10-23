<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center px-4 sm:px-6 lg:px-8">
    <!-- Background Pattern -->
    <div class="absolute inset-0 opacity-5">
      <div class="absolute inset-0" style="background-image: repeating-linear-gradient(45deg, #8B0000 0, #8B0000 1px, transparent 0, transparent 50%); background-size: 10px 10px;"></div>
    </div>

    <!-- Login Card -->
    <div class="relative max-w-md w-full">
      <!-- Logo/Header -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-gray-900 mb-2">
          <span class="text-bloodsa-red">Doctor's Workspace</span>
        </h1>
        <p class="text-gray-600">Sign in to access your workspace</p>
      </div>

      <!-- Login Form Card -->
      <div class="bg-white rounded-xl shadow-xl p-8 border-2 border-gray-100">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- Error Alert -->
          <div
            v-if="error"
            class="bg-red-50 border-l-4 border-red-500 p-4 rounded"
          >
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm text-red-700">{{ error }}</p>
              </div>
              <div class="ml-auto pl-3">
                <button
                  type="button"
                  @click="clearError"
                  class="inline-flex text-red-400 hover:text-red-500"
                >
                  <span class="sr-only">Dismiss</span>
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <!-- Email Input -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              Email Address
            </label>
            <input
              id="email"
              v-model="credentials.email"
              type="email"
              required
              autocomplete="email"
              class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors"
              placeholder="you@example.com"
              :disabled="isLoading"
            />
          </div>

          <!-- Password Input -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <div class="relative">
              <input
                id="password"
                v-model="credentials.password"
                :type="showPassword ? 'text' : 'password'"
                required
                autocomplete="current-password"
                class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors pr-12"
                placeholder="Enter your password"
                :disabled="isLoading"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
                :disabled="isLoading"
              >
                <svg v-if="!showPassword" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                <svg v-else class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Forgot Password -->
          <div v-if="smtpStore.isConfigured" class="flex items-center justify-end">
            <div class="text-sm">
              <router-link to="/forgot-password" class="font-medium text-bloodsa-red hover:text-bloodsa-light-red transition-colors">
                Forgot password?
              </router-link>
            </div>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-white bg-bloodsa-red hover:bg-bloodsa-light-red focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-bloodsa-red transition-all duration-200 font-semibold disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg
              v-if="isLoading"
              class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ isLoading ? 'Signing in...' : 'Sign in' }}
          </button>
        </form>

        <!-- Additional Info -->
        <div class="mt-6 text-center text-sm text-gray-600 space-y-2">
          <p>
            Don't have an account? 
            <router-link to="/register" class="font-medium text-bloodsa-red hover:text-red-700 transition-colors">
              Register here
            </router-link>
            or contact your administrator.
          </p>
        </div>
      </div>

      <!-- Footer Info -->
      <div class="mt-8 text-center text-xs text-gray-500">
        <p>By signing in, you agree to BLOODSA's terms of service and privacy policy.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSMTPStore } from '@/stores/smtp'

const router = useRouter()
const authStore = useAuthStore()
const smtpStore = useSMTPStore()

// Form state
const credentials = ref({
  email: '',
  password: ''
})

const showPassword = ref(false)
const error = ref<string | null>(null)
const isLoading = ref(false)

// Methods
const handleLogin = async () => {
  error.value = null
  isLoading.value = true

  try {
    await authStore.login(credentials.value)
    
    // Redirect to dashboard on success
    router.push({ name: 'dashboard' })
  } catch (err) {
    // Error is already set in the auth store
    error.value = authStore.error || 'Login failed. Please try again.'
  } finally {
    isLoading.value = false
  }
}

const clearError = () => {
  error.value = null
  authStore.clearError()
}

// Check SMTP configuration on mount
onMounted(async () => {
  try {
    await smtpStore.checkSMTPConfiguration()
  } catch (err) {
    // Silently fail - forgot password link will be hidden
    console.warn('Failed to check SMTP configuration:', err)
  }
})
</script>

<style scoped>
/* Custom styles if needed */
</style>


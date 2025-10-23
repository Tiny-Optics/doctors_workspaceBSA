<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center px-4 sm:px-6 lg:px-8 py-8">
    <!-- Background Pattern -->
    <div class="absolute inset-0 opacity-5">
      <div class="absolute inset-0" style="background-image: repeating-linear-gradient(45deg, #8B0000 0, #8B0000 1px, transparent 0, transparent 50%); background-size: 10px 10px;"></div>
    </div>

    <div class="max-w-md w-full space-y-8 relative z-10">
      <!-- Header -->
      <div class="text-center">
        <img src="/BLOODSA-SVG-Logo.svg" alt="BloodSA Logo" class="w-24 h-24 object-contain mx-auto mb-4" />
        <h2 class="text-3xl font-bold text-gray-900">Reset Your Password</h2>
        <p class="mt-2 text-sm text-gray-600">
          Follow the steps below to reset your password
        </p>
      </div>

      <!-- Progress Steps -->
      <div class="flex items-center justify-center space-x-4 mb-8">
        <div class="flex items-center">
          <div :class="[
            'w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium',
            currentStep === 'email' ? 'bg-bloodsa-red text-white' : 
            currentStep === 'code' || currentStep === 'password' || currentStep === 'success' ? 'bg-green-500 text-white' : 'bg-gray-300 text-gray-600'
          ]">
            1
          </div>
          <span class="ml-2 text-sm font-medium text-gray-700">Enter Email</span>
        </div>
        
        <div class="w-8 h-0.5 bg-gray-300"></div>
        
        <div class="flex items-center">
          <div :class="[
            'w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium',
            currentStep === 'code' ? 'bg-bloodsa-red text-white' : 
            currentStep === 'password' || currentStep === 'success' ? 'bg-green-500 text-white' : 'bg-gray-300 text-gray-600'
          ]">
            2
          </div>
          <span class="ml-2 text-sm font-medium text-gray-700">Verify Code</span>
        </div>
        
        <div class="w-8 h-0.5 bg-gray-300"></div>
        
        <div class="flex items-center">
          <div :class="[
            'w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium',
            currentStep === 'password' ? 'bg-bloodsa-red text-white' : 
            currentStep === 'success' ? 'bg-green-500 text-white' : 'bg-gray-300 text-gray-600'
          ]">
            3
          </div>
          <span class="ml-2 text-sm font-medium text-gray-700">New Password</span>
        </div>
      </div>

      <!-- Main Form Card -->
      <div class="bg-white rounded-lg shadow-lg p-8">
        <!-- Step 1: Email Input -->
        <div v-if="currentStep === 'email'" class="space-y-6">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 mb-2">Enter Your Email Address</h3>
            <p class="text-sm text-gray-600 mb-6">
              We'll send a verification code to your email address to reset your password.
            </p>
          </div>

          <form @submit.prevent="handleEmailSubmit" class="space-y-4">
            <div>
              <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
                Email Address
              </label>
              <input
                id="email"
                v-model="email"
                type="email"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent transition-colors"
                placeholder="Enter your email address"
                :disabled="isLoading"
              />
            </div>

            <button
              type="submit"
              :disabled="isLoading || !email"
              class="w-full bg-bloodsa-red text-white py-3 px-4 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center space-x-2"
            >
              <svg v-if="isLoading" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <span>{{ isLoading ? 'Sending Code...' : 'Send Reset Code' }}</span>
            </button>
          </form>
        </div>

        <!-- Step 2: Code Verification -->
        <div v-if="currentStep === 'code'" class="space-y-6">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 mb-2">Enter Verification Code</h3>
            <p class="text-sm text-gray-600 mb-6">
              We've sent a 6-digit verification code to <strong>{{ email }}</strong>. Please enter it below.
            </p>
          </div>

          <form @submit.prevent="handleCodeSubmit" class="space-y-4">
            <div>
              <label for="code" class="block text-sm font-medium text-gray-700 mb-2">
                Verification Code
              </label>
              <input
                id="code"
                v-model="code"
                type="text"
                maxlength="6"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent transition-colors text-center text-2xl font-mono tracking-widest"
                placeholder="000000"
                :disabled="isLoading"
                @input="formatCode"
              />
              <p class="mt-2 text-xs text-gray-500">
                Didn't receive the code? Check your spam folder or 
                <button type="button" @click="resendCode" class="text-bloodsa-red hover:underline" :disabled="isLoading">
                  resend code
                </button>
              </p>
            </div>

            <button
              type="submit"
              :disabled="isLoading || code.length !== 6"
              class="w-full bg-bloodsa-red text-white py-3 px-4 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center space-x-2"
            >
              <svg v-if="isLoading" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <span>{{ isLoading ? 'Verifying...' : 'Verify Code' }}</span>
            </button>
          </form>
        </div>

        <!-- Step 3: New Password -->
        <div v-if="currentStep === 'password'" class="space-y-6">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 mb-2">Set New Password</h3>
            <p class="text-sm text-gray-600 mb-6">
              Enter your new password below. Make sure it's secure and easy to remember.
            </p>
          </div>

          <form @submit.prevent="handlePasswordSubmit" class="space-y-4">
            <div>
              <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-2">
                New Password
              </label>
              <div class="relative">
                <input
                  id="newPassword"
                  v-model="newPassword"
                  :type="showPassword ? 'text' : 'password'"
                  required
                  minlength="8"
                  class="w-full px-4 py-3 pr-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent transition-colors"
                  placeholder="Enter your new password"
                  :disabled="isLoading"
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
              <p class="mt-1 text-xs text-gray-500">
                Password must be at least 8 characters long
              </p>
            </div>

            <div>
              <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
                Confirm New Password
              </label>
              <input
                id="confirmPassword"
                v-model="confirmPassword"
                :type="showPassword ? 'text' : 'password'"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent transition-colors"
                placeholder="Confirm your new password"
                :disabled="isLoading"
              />
            </div>

            <button
              type="submit"
              :disabled="isLoading || !isPasswordValid"
              class="w-full bg-bloodsa-red text-white py-3 px-4 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center space-x-2"
            >
              <svg v-if="isLoading" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <span>{{ isLoading ? 'Resetting Password...' : 'Reset Password' }}</span>
            </button>
          </form>
        </div>

        <!-- Step 4: Success -->
        <div v-if="currentStep === 'success'" class="text-center space-y-6">
          <div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto">
            <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          
          <div>
            <h3 class="text-lg font-semibold text-gray-900 mb-2">Password Reset Successful!</h3>
            <p class="text-sm text-gray-600 mb-6">
              Your password has been successfully reset. You can now log in with your new password.
            </p>
          </div>

          <router-link
            to="/login"
            class="w-full bg-bloodsa-red text-white py-3 px-4 rounded-lg hover:bg-red-700 transition-colors inline-block text-center"
          >
            Go to Login
          </router-link>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
          <div class="flex">
            <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div class="ml-3">
              <p class="text-sm text-red-800">{{ error }}</p>
            </div>
            <button @click="clearError" class="ml-auto text-red-400 hover:text-red-600">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Success Message -->
        <div v-if="success" class="mt-4 p-4 bg-green-50 border border-green-200 rounded-lg">
          <div class="flex">
            <svg class="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div class="ml-3">
              <p class="text-sm text-green-800">{{ success }}</p>
            </div>
            <button @click="clearSuccess" class="ml-auto text-green-400 hover:text-green-600">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Back to Login -->
        <div class="mt-6 text-center">
          <router-link to="/login" class="text-sm text-bloodsa-red hover:text-red-700 transition-colors">
            ← Back to Login
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useForgotPasswordStore } from '@/stores/forgotPassword'

const router = useRouter()
const forgotPasswordStore = useForgotPasswordStore()

// Form data
const email = ref('')
const code = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)

// Computed properties
const isLoading = computed(() => forgotPasswordStore.isLoading)
const error = computed(() => forgotPasswordStore.error)
const success = computed(() => forgotPasswordStore.success)
const currentStep = computed(() => forgotPasswordStore.currentStep)

const isPasswordValid = computed(() => {
  return newPassword.value.length >= 8 && 
         newPassword.value === confirmPassword.value
})

// Methods
const handleEmailSubmit = async () => {
  try {
    await forgotPasswordStore.requestPasswordReset(email.value)
  } catch (err) {
    // Error is handled by the store
  }
}

const handleCodeSubmit = async () => {
  try {
    await forgotPasswordStore.validateResetCode(code.value)
  } catch (err) {
    // Error is handled by the store
  }
}

const handlePasswordSubmit = async () => {
  try {
    await forgotPasswordStore.resetPassword(newPassword.value)
  } catch (err) {
    // Error is handled by the store
  }
}

const formatCode = (event: Event) => {
  const target = event.target as HTMLInputElement
  target.value = target.value.replace(/\D/g, '').slice(0, 6)
  code.value = target.value
}

const resendCode = async () => {
  try {
    await forgotPasswordStore.requestPasswordReset(email.value)
  } catch (err) {
    // Error is handled by the store
  }
}

const clearError = () => {
  forgotPasswordStore.clearError()
}

const clearSuccess = () => {
  forgotPasswordStore.clearSuccess()
}

// Reset store on mount
onMounted(() => {
  forgotPasswordStore.reset()
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

<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center px-4 sm:px-6 lg:px-8 py-8">
    <!-- Background Pattern -->
    <div class="absolute inset-0 opacity-5">
      <div class="absolute inset-0" style="background-image: repeating-linear-gradient(45deg, #8B0000 0, #8B0000 1px, transparent 0, transparent 50%); background-size: 10px 10px;"></div>
    </div>

    <!-- Registration Card -->
    <div class="relative max-w-2xl w-full">
      <!-- Logo/Header -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-gray-900 mb-2">
          <span class="text-bloodsa-red">Doctor's Workspace</span>
        </h1>
        <p class="text-gray-600">Create your account to access the workspace</p>
      </div>

      <!-- Registration Form Card -->
      <div class="bg-white rounded-xl shadow-xl p-8 border-2 border-gray-100">
        <form @submit.prevent="handleRegister" class="space-y-6">
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

          <!-- Success Alert -->
          <div
            v-if="success"
            class="bg-green-50 border-l-4 border-green-500 p-4 rounded"
          >
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm text-green-700">{{ success }}</p>
              </div>
            </div>
          </div>

          <!-- Form Grid -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- First Name -->
            <div>
              <label for="firstName" class="block text-sm font-medium text-gray-700 mb-2">
                First Name <span class="text-red-500">*</span>
              </label>
              <input
                id="firstName"
                v-model="formData.firstName"
                type="text"
                required
                class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors"
                placeholder="Enter your first name"
                :disabled="isLoading"
              />
            </div>

            <!-- Last Name -->
            <div>
              <label for="lastName" class="block text-sm font-medium text-gray-700 mb-2">
                Last Name <span class="text-red-500">*</span>
              </label>
              <input
                id="lastName"
                v-model="formData.lastName"
                type="text"
                required
                class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors"
                placeholder="Enter your last name"
                :disabled="isLoading"
              />
            </div>

            <!-- Username -->
            <div>
              <label for="username" class="block text-sm font-medium text-gray-700 mb-2">
                Username <span class="text-red-500">*</span>
              </label>
              <input
                id="username"
                v-model="formData.username"
                type="text"
                required
                class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors"
                placeholder="Choose a username"
                :disabled="isLoading"
              />
            </div>

            <!-- Email -->
            <div>
              <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
                Email Address <span class="text-red-500">*</span>
              </label>
              <input
                id="email"
                v-model="formData.email"
                type="email"
                required
                autocomplete="email"
                class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors"
                placeholder="you@example.com"
                :disabled="isLoading"
              />
            </div>

            <!-- Password -->
            <div>
              <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
                Password <span class="text-red-500">*</span>
              </label>
              <div class="relative">
                <input
                  id="password"
                  v-model="formData.password"
                  :type="showPassword ? 'text' : 'password'"
                  required
                  autocomplete="new-password"
                  class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors pr-12"
                  placeholder="Create a strong password"
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

            <!-- Confirm Password -->
            <div>
              <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
                Confirm Password <span class="text-red-500">*</span>
              </label>
              <div class="relative">
                <input
                  id="confirmPassword"
                  v-model="confirmPassword"
                  :type="showConfirmPassword ? 'text' : 'password'"
                  required
                  autocomplete="new-password"
                  :class="['appearance-none block w-full px-4 py-3 border rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors pr-12', passwordMatch === false ? 'border-red-300' : 'border-gray-300']"
                  placeholder="Confirm your password"
                  :disabled="isLoading"
                />
                <button
                  type="button"
                  @click="showConfirmPassword = !showConfirmPassword"
                  class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
                  :disabled="isLoading"
                >
                  <svg v-if="!showConfirmPassword" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  <svg v-else class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                  </svg>
                </button>
              </div>
              <p v-if="passwordMatch === false" class="mt-1 text-sm text-red-600">Passwords do not match</p>
            </div>
          </div>

          <!-- Institution (Full Width) -->
          <div class="mt-6">
            <label for="institutionId" class="block text-sm font-medium text-gray-700 mb-2">
              Institution <span class="text-red-500">*</span>
            </label>
            <select
              id="institutionId"
              v-model="formData.institutionId"
              required
              class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors"
              :disabled="isLoading || institutionsStore.isLoading"
            >
              <option value="">{{ institutionsStore.isLoading ? 'Loading institutions...' : 'Select your institution' }}</option>
              <option v-for="institution in institutionsStore.institutions" :key="institution.id" :value="institution.id">
                {{ institution.name }}
              </option>
            </select>
            <p class="mt-2 text-sm text-gray-600">
              <span class="font-medium">Can't find your institution?</span> Select "Other" from the dropdown. Once your account is approved, go to your profile section to add your institution and update it.
            </p>
          </div>

          <!-- Form Grid (continued) -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6">

            <!-- Specialty -->
            <div>
              <label for="specialty" class="block text-sm font-medium text-gray-700 mb-2">
                Specialty
              </label>
              <input
                id="specialty"
                v-model="formData.specialty"
                type="text"
                class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors"
                placeholder="e.g., Hematology, Oncology"
                :disabled="isLoading"
              />
            </div>

            <!-- Registration Number -->
            <div>
              <label for="registrationNumber" class="block text-sm font-medium text-gray-700 mb-2">
                Registration Number
              </label>
              <input
                id="registrationNumber"
                v-model="formData.registrationNumber"
                type="text"
                class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors"
                placeholder="Your professional registration number"
                :disabled="isLoading"
              />
            </div>

            <!-- Phone Number -->
            <div>
              <label for="phoneNumber" class="block text-sm font-medium text-gray-700 mb-2">
                Phone Number
              </label>
              <input
                id="phoneNumber"
                v-model="formData.phoneNumber"
                type="tel"
                class="appearance-none block w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors"
                placeholder="+27 XX XXX XXXX"
                :disabled="isLoading"
              />
            </div>
          </div>

          <!-- Submit Button -->
          <div class="pt-4">
            <button
              type="submit"
              :disabled="isLoading || !isFormValid"
              class="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-bloodsa-red hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-bloodsa-red disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <svg v-if="isLoading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ isLoading ? 'Creating Account...' : 'Create Account' }}
            </button>
          </div>

          <!-- Login Link -->
          <div class="text-center">
            <p class="text-sm text-gray-600">
              Already have an account? 
              <router-link to="/login" class="font-medium text-bloodsa-red hover:text-red-700 transition-colors">
                Sign in here
              </router-link>
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useInstitutionsStore } from '@/stores/institutions'
import type { RegisterUserRequest } from '@/types/user'

const router = useRouter()
const authStore = useAuthStore()
const institutionsStore = useInstitutionsStore()

// Form state
const formData = ref<RegisterUserRequest>({
  username: '',
  email: '',
  password: '',
  firstName: '',
  lastName: '',
  institutionId: '',
  specialty: '',
  registrationNumber: '',
  phoneNumber: ''
})

const showPassword = ref(false)
const showConfirmPassword = ref(false)
const confirmPassword = ref('')
const error = ref<string | null>(null)
const success = ref<string | null>(null)
const isLoading = ref(false)

// Computed
const passwordMatch = computed(() => {
  if (!confirmPassword.value) return null // Don't show error until user starts typing
  return formData.value.password === confirmPassword.value
})

const isFormValid = computed(() => {
  return formData.value.username &&
         formData.value.email &&
         formData.value.password &&
         confirmPassword.value &&
         passwordMatch.value === true &&
         formData.value.firstName &&
         formData.value.lastName &&
         formData.value.institutionId
})

// Methods
const handleRegister = async () => {
  error.value = null
  success.value = null
  
  // Validate passwords match
  if (formData.value.password !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }
  
  isLoading.value = true

  try {
    await authStore.register(formData.value)
    
    // Show success message
    success.value = 'Account created successfully! Your account is pending approval from an administrator. You will receive an email once your account is activated.'
    
    // Clear form
    formData.value = {
      username: '',
      email: '',
      password: '',
      firstName: '',
      lastName: '',
      institutionId: '',
      specialty: '',
      registrationNumber: '',
      phoneNumber: ''
    }
    confirmPassword.value = ''
  } catch (err) {
    // Error is already set in the auth store
    error.value = authStore.error || 'Registration failed. Please try again.'
  } finally {
    isLoading.value = false
  }
}

const clearError = () => {
  error.value = null
  authStore.clearError()
}

const loadInstitutions = async () => {
  try {
    await institutionsStore.fetchPublicInstitutions()
  } catch (err) {
    console.error('Failed to load institutions:', err)
    // Don't show error to user as this is not critical for registration
  }
}

onMounted(() => {
  loadInstitutions()
})
</script>

<style scoped>
/* Custom styles if needed */
</style>

<template>
  <div class="space-y-6">
    <!-- Change Password Section -->
    <div>
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Change Password</h2>
      <p class="text-gray-600 mb-6">
        Update your password to keep your account secure. Password must be at least 8 characters and include uppercase, lowercase, number, and special character.
      </p>

      <form @submit.prevent="changePassword" class="max-w-md space-y-4">
        <!-- Old Password -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Current Password <span class="text-red-500">*</span>
          </label>
          <div class="relative">
            <input
              v-model="passwordForm.oldPassword"
              :type="showOldPassword ? 'text' : 'password'"
              required
              class="w-full px-4 py-2 pr-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
              placeholder="Enter current password"
            />
            <button
              type="button"
              @click="showOldPassword = !showOldPassword"
              class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
            >
              <svg v-if="!showOldPassword" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
        </div>

        <!-- New Password -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            New Password <span class="text-red-500">*</span>
          </label>
          <div class="relative">
            <input
              v-model="passwordForm.newPassword"
              :type="showNewPassword ? 'text' : 'password'"
              required
              class="w-full px-4 py-2 pr-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
              placeholder="Enter new password"
              @input="checkPasswordStrength"
            />
            <button
              type="button"
              @click="showNewPassword = !showNewPassword"
              class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
            >
              <svg v-if="!showNewPassword" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
          
          <!-- Password Strength Indicator -->
          <div v-if="passwordForm.newPassword" class="mt-2">
            <div class="flex items-center gap-2">
              <div class="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
                <div 
                  :class="[
                    'h-full transition-all duration-300',
                    passwordStrength.color
                  ]"
                  :style="{ width: passwordStrength.width }"
                ></div>
              </div>
              <span class="text-xs font-medium" :class="passwordStrength.textColor">
                {{ passwordStrength.label }}
              </span>
            </div>
            <p class="text-xs text-gray-500 mt-1">{{ passwordStrength.message }}</p>
          </div>
        </div>

        <!-- Confirm New Password -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Confirm New Password <span class="text-red-500">*</span>
          </label>
          <div class="relative">
            <input
              v-model="passwordForm.confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              required
              class="w-full px-4 py-2 pr-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
              placeholder="Re-enter new password"
            />
            <button
              type="button"
              @click="showConfirmPassword = !showConfirmPassword"
              class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
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
          
          <!-- Password Match Indicator -->
          <p v-if="passwordForm.confirmPassword" class="text-xs mt-1" :class="passwordsMatch ? 'text-green-600' : 'text-red-600'">
            {{ passwordsMatch ? '✓ Passwords match' : '✗ Passwords do not match' }}
          </p>
        </div>

        <!-- Submit Button -->
        <div class="pt-4">
          <button
            type="submit"
            :disabled="!canSubmitPassword || changingPassword"
            class="w-full px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <svg v-if="changingPassword" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ changingPassword ? 'Updating Password...' : 'Update Password' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Security Tips -->
    <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
      <h3 class="text-sm font-medium text-blue-800 mb-2">🔒 Password Security Tips</h3>
      <ul class="text-sm text-blue-700 space-y-1">
        <li>• Use a unique password you don't use elsewhere</li>
        <li>• Include a mix of uppercase, lowercase, numbers, and symbols</li>
        <li>• Avoid common words or personal information</li>
        <li>• Consider using a password manager</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useToast } from '@/composables/useToast'

const toast = useToast()

const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const showOldPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)
const changingPassword = ref(false)

const passwordsMatch = computed(() => {
  return passwordForm.value.newPassword === passwordForm.value.confirmPassword
})

const canSubmitPassword = computed(() => {
  return (
    passwordForm.value.oldPassword &&
    passwordForm.value.newPassword &&
    passwordForm.value.confirmPassword &&
    passwordsMatch.value &&
    passwordForm.value.newPassword.length >= 8
  )
})

// Password strength checker
const passwordStrength = computed(() => {
  const password = passwordForm.value.newPassword
  if (!password) return { width: '0%', color: '', label: '', message: '', textColor: '' }
  
  let score = 0
  let message = ''
  
  // Length check
  if (password.length >= 8) score++
  if (password.length >= 12) score++
  
  // Character variety
  if (/[a-z]/.test(password)) score++ // lowercase
  if (/[A-Z]/.test(password)) score++ // uppercase
  if (/[0-9]/.test(password)) score++ // number
  if (/[^a-zA-Z0-9]/.test(password)) score++ // special char
  
  if (score <= 2) {
    return {
      width: '33%',
      color: 'bg-red-500',
      textColor: 'text-red-600',
      label: 'Weak',
      message: 'Add more character variety and length'
    }
  } else if (score <= 4) {
    return {
      width: '66%',
      color: 'bg-yellow-500',
      textColor: 'text-yellow-600',
      label: 'Medium',
      message: 'Good, but could be stronger'
    }
  } else {
    return {
      width: '100%',
      color: 'bg-green-500',
      textColor: 'text-green-600',
      label: 'Strong',
      message: 'Excellent password!'
    }
  }
})

async function changePassword() {
  // Validate passwords match
  if (!passwordsMatch.value) {
    toast.error('New passwords do not match')
    return
  }
  
  // Validate password strength
  if (passwordForm.value.newPassword.length < 8) {
    toast.error('Password must be at least 8 characters')
    return
  }
  
  changingPassword.value = true
  
  try {
    const response = await fetch('http://localhost:8080/api/auth/change-password', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        oldPassword: passwordForm.value.oldPassword,
        newPassword: passwordForm.value.newPassword
      })
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to change password')
    }
    
    // Reset form
    passwordForm.value = {
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
    
    toast.success('Password changed successfully!')
  } catch (error: any) {
    console.error('Failed to change password:', error)
    toast.error(error.message || 'Failed to change password')
  } finally {
    changingPassword.value = false
  }
}

function checkPasswordStrength() {
  // This reactive computed will update automatically
}
</script>


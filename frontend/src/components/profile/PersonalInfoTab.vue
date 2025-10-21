<template>
  <div class="space-y-6">
    <!-- Read-Only Account Information -->
    <div>
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Account Information</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-gray-50 rounded-lg p-4">
          <label class="block text-sm font-medium text-gray-600 mb-1">Username</label>
          <p class="text-lg text-gray-900">{{ user?.username }}</p>
        </div>
        <div class="bg-gray-50 rounded-lg p-4">
          <label class="block text-sm font-medium text-gray-600 mb-1">Email</label>
          <p class="text-lg text-gray-900">{{ user?.email }}</p>
        </div>
        <div class="bg-gray-50 rounded-lg p-4">
          <label class="block text-sm font-medium text-gray-600 mb-1">Role</label>
          <p class="text-lg text-gray-900">{{ getRoleDisplayName(user?.role) }}</p>
        </div>
        <div class="bg-gray-50 rounded-lg p-4">
          <label class="block text-sm font-medium text-gray-600 mb-1">Institution</label>
          <p class="text-lg text-gray-900">{{ institutionName }}</p>
        </div>
      </div>
      <p class="text-sm text-gray-500 mt-3">
        ℹ️ Contact an administrator to change your username, email, role, or institution.
      </p>
    </div>

    <!-- Editable Information -->
    <div>
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-semibold text-gray-900">Personal Information</h2>
        <button
          v-if="!isEditing"
          @click="startEditing"
          class="px-4 py-2 text-sm font-medium text-bloodsa-red border border-bloodsa-red rounded-lg hover:bg-bloodsa-red hover:text-white transition-colors"
        >
          Edit Profile
        </button>
      </div>

      <form @submit.prevent="saveChanges" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- First Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              First Name <span class="text-red-500">*</span>
            </label>
            <input
              v-model="formData.firstName"
              type="text"
              required
              :disabled="!isEditing"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent disabled:bg-gray-50 disabled:text-gray-600"
              placeholder="Enter your first name"
            />
          </div>

          <!-- Last Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Last Name <span class="text-red-500">*</span>
            </label>
            <input
              v-model="formData.lastName"
              type="text"
              required
              :disabled="!isEditing"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent disabled:bg-gray-50 disabled:text-gray-600"
              placeholder="Enter your last name"
            />
          </div>

          <!-- Phone Number -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
            <input
              v-model="formData.phoneNumber"
              type="tel"
              :disabled="!isEditing"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent disabled:bg-gray-50 disabled:text-gray-600"
              placeholder="+27 11 123 4567"
            />
          </div>

          <!-- Specialty -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Specialty</label>
            <input
              v-model="formData.specialty"
              type="text"
              :disabled="!isEditing"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent disabled:bg-gray-50 disabled:text-gray-600"
              placeholder="e.g., Haematology, Internal Medicine"
            />
          </div>

          <!-- Registration Number -->
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">Professional Registration Number</label>
            <input
              v-model="formData.registrationNumber"
              type="text"
              :disabled="!isEditing"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent disabled:bg-gray-50 disabled:text-gray-600"
              placeholder="e.g., HPCSA MP 12345, BHF 67890"
            />
            <p class="text-sm text-gray-500 mt-1">HPCSA or other professional registration number</p>
          </div>
        </div>

        <!-- Action Buttons (only show when editing) -->
        <div v-if="isEditing" class="flex justify-end gap-3 pt-4 border-t border-gray-200">
          <button
            type="button"
            @click="cancelEditing"
            :disabled="saving"
            class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="saving"
            class="px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 flex items-center gap-2"
          >
            <svg v-if="saving" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ saving ? 'Saving...' : 'Save Changes' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Account Metadata -->
    <div>
      <h2 class="text-xl font-semibold text-gray-900 mb-4">Account Details</h2>
      <div class="bg-gray-50 rounded-lg p-4 space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-gray-600">Account Created</span>
          <span class="text-sm text-gray-900">{{ formatDate(user?.createdAt) }}</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-gray-600">Last Updated</span>
          <span class="text-sm text-gray-900">{{ formatDate(user?.updatedAt) }}</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-gray-600">Last Login</span>
          <span class="text-sm text-gray-900">{{ formatLastLogin(user?.lastLoginAt) }}</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-gray-600">Account Status</span>
          <span 
            :class="[
              'px-2 py-1 text-xs font-medium rounded-full',
              user?.isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
            ]"
          >
            {{ user?.isActive ? 'Active' : 'Inactive' }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useInstitutionsStore } from '@/stores/institutions'
import { getUserRoleDisplayName as getRoleDisplayName } from '@/types/user'
import type { UpdateUserRequest } from '@/types/user'
import { useToast } from '@/composables/useToast'

const authStore = useAuthStore()
const institutionsStore = useInstitutionsStore()
const toast = useToast()

const user = computed(() => authStore.user)
const isEditing = ref(false)
const saving = ref(false)

// Load institutions on mount
onMounted(async () => {
  await institutionsStore.fetchInstitutions({ isActive: true, limit: 1000 })
})

// Form data
const formData = ref({
  firstName: '',
  lastName: '',
  phoneNumber: '',
  specialty: '',
  registrationNumber: ''
})

// Initialize form data from user
watch(user, (newUser) => {
  if (newUser) {
    formData.value = {
      firstName: newUser.profile.firstName,
      lastName: newUser.profile.lastName,
      phoneNumber: newUser.profile.phoneNumber || '',
      specialty: newUser.profile.specialty || '',
      registrationNumber: newUser.profile.registrationNumber || ''
    }
  }
}, { immediate: true })

const institutionName = computed(() => {
  if (!user.value?.profile.institutionId) return 'Not Set'
  const institution = institutionsStore.institutions.find(i => i.id === user.value.profile.institutionId)
  return institution ? institution.name : 'Unknown Institution'
})

function startEditing() {
  isEditing.value = true
}

function cancelEditing() {
  // Reset form data to original values
  if (user.value) {
    formData.value = {
      firstName: user.value.profile.firstName,
      lastName: user.value.profile.lastName,
      phoneNumber: user.value.profile.phoneNumber || '',
      specialty: user.value.profile.specialty || '',
      registrationNumber: user.value.profile.registrationNumber || ''
    }
  }
  isEditing.value = false
}

async function saveChanges() {
  if (!user.value) return
  
  saving.value = true
  
  try {
    // Build update request with only changed fields
    const updateData: UpdateUserRequest = {
      firstName: formData.value.firstName,
      lastName: formData.value.lastName,
      phoneNumber: formData.value.phoneNumber || undefined,
      specialty: formData.value.specialty || undefined,
      registrationNumber: formData.value.registrationNumber || undefined
    }
    
    // Call the update API
    const response = await fetch(`/api/users/${user.value.id}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(updateData)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to update profile')
    }
    
    const updatedUser = await response.json()
    
    // Update the auth store with new user data
    authStore.user = updatedUser
    
    isEditing.value = false
    toast.success('Profile updated successfully!')
  } catch (error: any) {
    console.error('Failed to update profile:', error)
    toast.error(error.message || 'Failed to update profile')
  } finally {
    saving.value = false
  }
}

function formatDate(dateString?: string): string {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
}

function formatLastLogin(lastLogin?: string): string {
  if (!lastLogin) return 'Never'
  const date = new Date(lastLogin)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  
  if (diffMins < 60) return `${diffMins} minutes ago`
  const diffHours = Math.floor(diffMins / 60)
  if (diffHours < 24) return `${diffHours} hours ago`
  const diffDays = Math.floor(diffHours / 24)
  if (diffDays === 1) return 'Yesterday'
  if (diffDays < 30) return `${diffDays} days ago`
  return formatDate(lastLogin)
}
</script>


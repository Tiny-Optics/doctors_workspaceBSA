<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <section class="bg-bloodsa-red relative overflow-hidden">
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-16">
        <!-- Back Button -->
        <button 
          @click="$router.push({ name: 'registry' })"
          class="mb-6 inline-flex items-center text-white hover:text-white/80 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        
        <div class="text-center">
          <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">Upload Final Approvals</h1>
          <p class="text-white max-w-2xl mx-auto">
            Submit your ethics approval documents for review
          </p>
        </div>
      </div>
    </section>

    <!-- Loading State -->
    <section v-if="loading" class="py-12">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-bloodsa-red mx-auto"></div>
          <p class="text-gray-600 mt-4">Loading form...</p>
        </div>
      </div>
    </section>

    <!-- Error State -->
    <section v-else-if="error" class="py-12">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <svg class="w-16 h-16 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">{{ error }}</h3>
          <button
            @click="loadActiveForm"
            class="mt-4 px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors"
          >
            Try Again
          </button>
        </div>
      </div>
    </section>

    <!-- No Active Form State -->
    <section v-else-if="!activeForm" class="py-12">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">No Active Form Available</h3>
          <p class="text-gray-600">
            There is currently no active submission form. Please contact the administrator.
          </p>
        </div>
      </div>
    </section>

    <!-- Success State -->
    <section v-else-if="submissionSuccess" class="py-12">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <div class="mb-6">
            <div class="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-green-100">
              <svg class="h-10 w-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
          </div>
          <h3 class="text-2xl font-bold text-gray-900 mb-4">Application Submitted Successfully!</h3>
          <p class="text-gray-600 mb-6">
            Your submission has been received and will be reviewed by the administrators.
          </p>
          <button
            @click="$router.push({ name: 'registry' })"
            class="px-6 py-3 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors font-medium"
          >
            Return to Registry
          </button>
        </div>
      </div>
    </section>

    <!-- Submission Error State -->
    <section v-else-if="submissionError" class="py-12">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <div class="mb-6">
            <div class="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-red-100">
              <svg class="h-10 w-10 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
          </div>
          <h3 class="text-2xl font-bold text-gray-900 mb-4">Submission Failed</h3>
          <p class="text-gray-600 mb-2">{{ submissionErrorMessage }}</p>
          <p class="text-sm text-gray-500 mb-6">Please try again or contact support if the problem persists.</p>
          <div class="flex gap-4 justify-center">
            <button
              @click="resetSubmission"
              class="px-6 py-3 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors font-medium"
            >
              Try Again
            </button>
            <button
              @click="$router.push({ name: 'registry' })"
              class="px-6 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium"
            >
              Return to Registry
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- Form Section -->
    <section v-else-if="activeForm && !submissionSuccess && !submissionError" class="py-12">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-8">
          <!-- Form Header -->
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-gray-900 mb-2">{{ activeForm.formName }}</h2>
            <p v-if="activeForm.description" class="text-gray-600">{{ activeForm.description }}</p>
            <p class="text-sm text-gray-500 mt-2">* Required fields</p>
          </div>

          <!-- Form Fields -->
          <form @submit.prevent="handleSubmit" class="space-y-6">
            <div v-for="field in sortedFields" :key="field.id" class="form-field">
              <!-- Text Input -->
              <div v-if="field.type === 'text'">
                <label :for="field.id" class="block text-sm font-medium text-gray-700 mb-1">
                  {{ field.label }}
                  <span v-if="field.required" class="text-red-500">*</span>
                </label>
                <input
                  :id="field.id"
                  v-model="formData[field.id]"
                  type="text"
                  :placeholder="field.placeholder"
                  :required="field.required"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  :class="{ 'border-red-500': validationErrors[field.id] }"
                />
                <p v-if="field.helpText" class="text-sm text-gray-500 mt-1">{{ field.helpText }}</p>
                <p v-if="validationErrors[field.id]" class="text-sm text-red-500 mt-1">{{ validationErrors[field.id] }}</p>
              </div>

              <!-- Textarea -->
              <div v-else-if="field.type === 'textarea'">
                <label :for="field.id" class="block text-sm font-medium text-gray-700 mb-1">
                  {{ field.label }}
                  <span v-if="field.required" class="text-red-500">*</span>
                </label>
                <textarea
                  :id="field.id"
                  v-model="formData[field.id]"
                  :placeholder="field.placeholder"
                  :required="field.required"
                  rows="4"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  :class="{ 'border-red-500': validationErrors[field.id] }"
                ></textarea>
                <p v-if="field.helpText" class="text-sm text-gray-500 mt-1">{{ field.helpText }}</p>
                <p v-if="validationErrors[field.id]" class="text-sm text-red-500 mt-1">{{ validationErrors[field.id] }}</p>
              </div>

              <!-- Select -->
              <div v-else-if="field.type === 'select'">
                <label :for="field.id" class="block text-sm font-medium text-gray-700 mb-1">
                  {{ field.label }}
                  <span v-if="field.required" class="text-red-500">*</span>
                </label>
                <select
                  :id="field.id"
                  v-model="formData[field.id]"
                  :required="field.required"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  :class="{ 'border-red-500': validationErrors[field.id] }"
                >
                  <option value="">Select an option</option>
                  <option v-for="option in field.options" :key="option" :value="option">
                    {{ option }}
                  </option>
                </select>
                <p v-if="field.helpText" class="text-sm text-gray-500 mt-1">{{ field.helpText }}</p>
                <p v-if="validationErrors[field.id]" class="text-sm text-red-500 mt-1">{{ validationErrors[field.id] }}</p>
              </div>

              <!-- Email -->
              <div v-else-if="field.type === 'email'">
                <label :for="field.id" class="block text-sm font-medium text-gray-700 mb-1">
                  {{ field.label }}
                  <span v-if="field.required" class="text-red-500">*</span>
                </label>
                <input
                  :id="field.id"
                  v-model="formData[field.id]"
                  type="email"
                  :placeholder="field.placeholder"
                  :required="field.required"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  :class="{ 'border-red-500': validationErrors[field.id] }"
                />
                <p v-if="field.helpText" class="text-sm text-gray-500 mt-1">{{ field.helpText }}</p>
                <p v-if="validationErrors[field.id]" class="text-sm text-red-500 mt-1">{{ validationErrors[field.id] }}</p>
              </div>

              <!-- Number -->
              <div v-else-if="field.type === 'number'">
                <label :for="field.id" class="block text-sm font-medium text-gray-700 mb-1">
                  {{ field.label }}
                  <span v-if="field.required" class="text-red-500">*</span>
                </label>
                <input
                  :id="field.id"
                  v-model.number="formData[field.id]"
                  type="number"
                  :placeholder="field.placeholder"
                  :required="field.required"
                  :min="field.validationRules?.minValue"
                  :max="field.validationRules?.maxValue"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  :class="{ 'border-red-500': validationErrors[field.id] }"
                />
                <p v-if="field.helpText" class="text-sm text-gray-500 mt-1">{{ field.helpText }}</p>
                <p v-if="validationErrors[field.id]" class="text-sm text-red-500 mt-1">{{ validationErrors[field.id] }}</p>
              </div>

              <!-- Date -->
              <div v-else-if="field.type === 'date'">
                <label :for="field.id" class="block text-sm font-medium text-gray-700 mb-1">
                  {{ field.label }}
                  <span v-if="field.required" class="text-red-500">*</span>
                </label>
                <input
                  :id="field.id"
                  v-model="formData[field.id]"
                  type="date"
                  :required="field.required"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  :class="{ 'border-red-500': validationErrors[field.id] }"
                />
                <p v-if="field.helpText" class="text-sm text-gray-500 mt-1">{{ field.helpText }}</p>
                <p v-if="validationErrors[field.id]" class="text-sm text-red-500 mt-1">{{ validationErrors[field.id] }}</p>
              </div>

              <!-- File Upload -->
              <div v-else-if="field.type === 'file'">
                <label :for="field.id" class="block text-sm font-medium text-gray-700 mb-1">
                  {{ field.label }}
                  <span v-if="field.required" class="text-red-500">*</span>
                </label>
                <div class="mt-1">
                  <input
                    :id="field.id"
                    type="file"
                    multiple
                    :required="field.required && (!uploadedFiles[field.id] || (uploadedFiles[field.id] && uploadedFiles[field.id]!.length === 0))"
                    @change="handleFileChange($event, field.id)"
                    class="block w-full text-sm text-gray-500
                      file:mr-4 file:py-2 file:px-4
                      file:rounded-lg file:border-0
                      file:text-sm file:font-semibold
                      file:bg-bloodsa-red file:text-white
                      hover:file:bg-red-700
                      file:cursor-pointer cursor-pointer"
                    :class="{ 'border-red-500': validationErrors[field.id] }"
                  />
                  <p v-if="field.helpText" class="text-sm text-gray-500 mt-1">{{ field.helpText }}</p>
                  <p class="text-sm text-gray-500 mt-1">You can select multiple files</p>
                  
                  <!-- Display selected files -->
                  <div v-if="uploadedFiles[field.id] && uploadedFiles[field.id]!.length > 0" class="mt-2 space-y-1">
                    <div v-for="(file, index) in uploadedFiles[field.id]!" :key="index" class="flex items-center justify-between bg-gray-50 px-3 py-2 rounded">
                      <span class="text-sm text-gray-700 truncate">{{ file.name }}</span>
                      <button
                        type="button"
                        @click="removeFile(field.id, index)"
                        class="text-red-500 hover:text-red-700 ml-2"
                      >
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </div>
                  </div>
                  
                  <p v-if="validationErrors[field.id]" class="text-sm text-red-500 mt-1">{{ validationErrors[field.id] }}</p>
                </div>
              </div>
            </div>

            <!-- Submit Button -->
            <div class="flex items-center justify-between pt-6 border-t">
              <button
                type="button"
                @click="$router.push({ name: 'registry' })"
                class="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="submitting"
                class="px-6 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
              >
                <svg v-if="submitting" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ submitting ? 'Submitting...' : 'Submit Application' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </section>

    <!-- Toast Notification -->
    <ToastNotification
      v-if="toast.show"
      :message="toast.message"
      :type="toast.type"
      @close="toast.show = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { registryService } from '@/services/registryService'
import type { FormSchema, FormField } from '@/services/registryService'
import ToastNotification from '@/components/ToastNotification.vue'

const router = useRouter()

const activeForm = ref<FormSchema | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const submitting = ref(false)
const submissionSuccess = ref(false)
const submissionError = ref(false)
const submissionErrorMessage = ref('')

const formData = ref<Record<string, any>>({})
const uploadedFiles = ref<Record<string, File[]>>({})
const validationErrors = ref<Record<string, string>>({})

const toast = ref({
  show: false,
  message: '',
  type: 'success' as 'success' | 'error' | 'warning'
})

const sortedFields = computed(() => {
  if (!activeForm.value) return []
  return [...activeForm.value.fields].sort((a, b) => {
    const orderA = a.displayOrder ?? 999
    const orderB = b.displayOrder ?? 999
    return orderA - orderB
  })
})

function showToast(message: string, type: 'success' | 'error' | 'warning' = 'success') {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 5000)
}

async function loadActiveForm() {
  loading.value = true
  error.value = null
  
  try {
    const form = await registryService.getActiveForm()
    if (!form) {
      activeForm.value = null
    } else {
      activeForm.value = form
      // Initialize form data
      form.fields.forEach(field => {
        if (field.type !== 'file') {
          formData.value[field.id] = ''
        } else {
          uploadedFiles.value[field.id] = []
        }
      })
    }
  } catch (err: any) {
    console.error('Failed to load active form:', err)
    error.value = err.message || 'Failed to load form. Please try again.'
  } finally {
    loading.value = false
  }
}

function handleFileChange(event: Event, fieldId: string) {
  const target = event.target as HTMLInputElement
  const files = target.files
  
  if (!files || files.length === 0) return
  
  // Always allow multiple files - add to existing files
  const existingFiles = uploadedFiles.value[fieldId] || []
  uploadedFiles.value[fieldId] = [...existingFiles, ...Array.from(files)]
  
  // Clear validation error if exists
  if (validationErrors.value[fieldId]) {
    delete validationErrors.value[fieldId]
  }
}

function removeFile(fieldId: string, index: number) {
  if (uploadedFiles.value[fieldId]) {
    uploadedFiles.value[fieldId].splice(index, 1)
  }
}

function validateForm(): boolean {
  validationErrors.value = {}
  let isValid = true

  if (!activeForm.value) return false

  for (const field of activeForm.value.fields) {
    // Check required fields
    if (field.required) {
      if (field.type === 'file') {
        const files = uploadedFiles.value[field.id]
        if (!files || files.length === 0) {
          validationErrors.value[field.id] = `${field.label} is required`
          isValid = false
        }
      } else {
        const value = formData.value[field.id]
        if (value === '' || value === null || value === undefined) {
          validationErrors.value[field.id] = `${field.label} is required`
          isValid = false
        }
      }
    }

    // Validate field-specific rules
    if (formData.value[field.id] && field.validationRules) {
      const value = formData.value[field.id]
      const rules = field.validationRules

      if (field.type === 'text' || field.type === 'textarea') {
        if (rules.minLength && value.length < rules.minLength) {
          validationErrors.value[field.id] = `Minimum length is ${rules.minLength} characters`
          isValid = false
        }
        if (rules.maxLength && value.length > rules.maxLength) {
          validationErrors.value[field.id] = `Maximum length is ${rules.maxLength} characters`
          isValid = false
        }
        if (rules.pattern && !new RegExp(rules.pattern).test(value)) {
          validationErrors.value[field.id] = `Invalid format for ${field.label}`
          isValid = false
        }
      }

      if (field.type === 'number') {
        if (rules.minValue !== undefined && value < rules.minValue) {
          validationErrors.value[field.id] = `Minimum value is ${rules.minValue}`
          isValid = false
        }
        if (rules.maxValue !== undefined && value > rules.maxValue) {
          validationErrors.value[field.id] = `Maximum value is ${rules.maxValue}`
          isValid = false
        }
      }
    }
  }

  if (!isValid) {
    showToast('Please fix the errors in the form', 'error')
  }

  return isValid
}

async function handleSubmit() {
  if (!validateForm()) {
    return
  }

  if (!activeForm.value) {
    showToast('No active form available', 'error')
    return
  }

  submitting.value = true

  try {
    // Collect all files from all file fields
    const allFiles: File[] = []
    Object.values(uploadedFiles.value).forEach(files => {
      allFiles.push(...files)
    })

    // Submit the form
    await registryService.submitForm({
      formSchemaId: activeForm.value.id,
      formData: formData.value,
      files: allFiles
    })

    // Show success state
    submissionSuccess.value = true
    submissionError.value = false

  } catch (err: any) {
    console.error('Failed to submit form:', err)
    // Show error state
    submissionSuccess.value = false
    submissionError.value = true
    submissionErrorMessage.value = err.message || 'Failed to submit application. Please try again.'
  } finally {
    submitting.value = false
  }
}

function resetSubmission() {
  submissionError.value = false
  submissionErrorMessage.value = ''
  // Clear form and reset to initial state
  loadActiveForm()
}

onMounted(() => {
  loadActiveForm()
})
</script>

<style scoped>
.bg-bloodsa-red {
  background-color: #8B0000;
}

.text-bloodsa-red {
  color: #8B0000;
}

.focus\:ring-bloodsa-red:focus {
  --tw-ring-color: #8B0000;
}

.hover\:bg-red-700:hover {
  background-color: #7f0000;
}
</style>


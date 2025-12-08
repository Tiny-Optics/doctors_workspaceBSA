<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Back Button -->
    <button 
      @click="$router.push('/admin/registry')"
      class="mb-6 inline-flex items-center text-gray-600 hover:text-bloodsa-red transition-colors"
    >
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
      Back to Registry Settings
    </button>

    <!-- Header -->
    <div class="mb-8">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Form Schemas</h1>
          <p class="mt-2 text-gray-600">
            Manage custom forms for the African HOPeR Registry.
          </p>
        </div>
        <button
          @click="showCreateForm = true"
          class="px-6 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors flex items-center space-x-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          <span>Create Form</span>
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-bloodsa-red"></div>
    </div>

    <!-- Forms List -->
    <div v-else class="space-y-6">
      <!-- Empty State -->
      <div v-if="forms.length === 0" class="text-center py-12">
        <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No Form Schemas</h3>
        <p class="text-gray-600 mb-4">Create your first form schema to get started.</p>
        <button
          @click="showCreateForm = true"
          class="px-6 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors"
        >
          Create Form Schema
        </button>
      </div>

      <!-- Forms Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="form in forms"
          :key="form.id"
          class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow"
        >
          <div class="flex justify-between items-start mb-4">
            <div>
              <h3 class="text-lg font-semibold text-gray-900">{{ form.formName }}</h3>
              <p v-if="form.description" class="text-sm text-gray-600 mt-1">{{ form.description }}</p>
            </div>
            <div class="flex items-center space-x-2">
              <span
                :class="form.isActive ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                class="px-2 py-1 text-xs font-medium rounded-full"
              >
                {{ form.isActive ? 'Active' : 'Inactive' }}
              </span>
            </div>
          </div>

          <div class="space-y-2 mb-4">
            <div class="flex items-center text-sm text-gray-600">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              {{ form.fields.length }} fields
            </div>
            <div class="flex items-center text-sm text-gray-600">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Created {{ formatDate(form.createdAt) }}
            </div>
          </div>

          <div class="flex justify-between items-center">
            <div class="flex space-x-2">
              <button
                @click="editForm(form)"
                class="px-3 py-1 text-sm text-bloodsa-red hover:text-red-700 hover:bg-red-50 rounded transition-colors"
              >
                Edit
              </button>
              <button
                @click="previewForm(form)"
                class="px-3 py-1 text-sm text-blue-600 hover:text-blue-700 hover:bg-blue-50 rounded transition-colors"
              >
                Preview
              </button>
            </div>
            <div class="flex space-x-1">
              <button
                @click="toggleFormStatus(form)"
                :class="form.isActive ? 'text-orange-600 hover:text-orange-700' : 'text-green-600 hover:text-green-700'"
                class="p-1 hover:bg-gray-100 rounded transition-colors"
                :title="form.isActive ? 'Deactivate' : 'Activate'"
              >
                <svg v-if="form.isActive" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h8m-9-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </button>
              <button
                @click="deleteForm(form)"
                class="p-1 text-red-600 hover:text-red-700 hover:bg-red-50 rounded transition-colors"
                title="Delete"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalForms > itemsPerPage" class="bg-white rounded-lg shadow-sm border border-gray-200 px-4 py-3 flex items-center justify-between sm:px-6 mt-6">
        <div class="flex-1 flex justify-between sm:hidden">
          <button
            @click="previousPage"
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <button
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700">
              Showing
              <span class="font-medium">{{ (currentPage - 1) * itemsPerPage + 1 }}</span>
              to
              <span class="font-medium">{{ Math.min(currentPage * itemsPerPage, totalForms) }}</span>
              of
              <span class="font-medium">{{ totalForms }}</span>
              results
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
              <!-- Skip to first page -->
              <button
                @click="goToFirstPage"
                :disabled="currentPage === 1"
                class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                title="Go to first page"
              >
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M15.707 15.707a1 1 0 01-1.414 0l-5-5a1 1 0 010-1.414l5-5a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 010 1.414zm-6 0a1 1 0 01-1.414 0l-5-5a1 1 0 010-1.414l5-5a1 1 0 011.414 1.414L5.414 10l4.293 4.293a1 1 0 010 1.414z" clip-rule="evenodd" />
                </svg>
              </button>
              <!-- Previous page -->
              <button
                @click="previousPage"
                :disabled="currentPage === 1"
                class="relative inline-flex items-center px-2 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                title="Previous page"
              >
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
              <!-- Page numbers -->
              <button
                v-for="page in visiblePages"
                :key="page"
                @click="goToPage(page)"
                class="relative inline-flex items-center px-4 py-2 border text-sm font-medium"
                :class="page === currentPage 
                  ? 'z-10 bg-bloodsa-red border-bloodsa-red text-white' 
                  : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'"
              >
                {{ page }}
              </button>
              <!-- Next page -->
              <button
                @click="nextPage"
                :disabled="currentPage === totalPages"
                class="relative inline-flex items-center px-2 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                title="Next page"
              >
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
              <!-- Skip to last page -->
              <button
                @click="goToLastPage"
                :disabled="currentPage === totalPages"
                class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                title="Go to last page"
              >
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10.293 15.707a1 1 0 010-1.414L14.586 10l-4.293-4.293a1 1 0 111.414-1.414l5 5a1 1 0 010 1.414l-5 5a1 1 0 01-1.414 0zm-4 0a1 1 0 010-1.414L9.586 10l-3.293-3.293a1 1 0 111.414-1.414l5 5a1 1 0 010 1.414l-5 5a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Form Modal -->
    <div 
      v-if="showCreateForm || editingForm" 
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
    >
      <div 
        @click.stop
        class="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto"
      >
        <div class="p-6">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-2xl font-bold text-gray-900">
              {{ editingForm ? 'Edit Form Schema' : 'Create Form Schema' }}
            </h2>
            <button
              @click="closeFormModal"
              class="text-gray-400 hover:text-gray-600"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Form Builder -->
          <div class="space-y-6">
            <!-- Basic Information -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Form Name</label>
                <input
                  v-model="formData.formName"
                  type="text"
                  placeholder="Enter form name"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
                <select
                  v-model="formData.isActive"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                >
                  <option :value="true">Active</option>
                  <option :value="false">Inactive</option>
                </select>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
              <textarea
                v-model="formData.description"
                rows="3"
                placeholder="Enter form description"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
              ></textarea>
            </div>

            <!-- Form Fields -->
            <div>
              <div class="flex justify-between items-center mb-4">
                <h3 class="text-lg font-semibold text-gray-900">Form Fields</h3>
                <button
                  @click="addField"
                  class="px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors flex items-center space-x-2"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                  </svg>
                  <span>Add Field</span>
                </button>
              </div>

              <div class="space-y-4">
                <div
                  v-for="(field, index) in formData.fields"
                  :key="index"
                  class="border border-gray-200 rounded-lg p-4"
                >
                  <div class="flex justify-between items-start mb-4">
                    <h4 class="font-medium text-gray-900">Field {{ index + 1 }}</h4>
                    <button
                      @click="removeField(index)"
                      class="text-red-600 hover:text-red-700"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  </div>

                  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">Field ID</label>
                      <input
                        v-model="field.id"
                        type="text"
                        placeholder="e.g., patient_name"
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                      />
                      <p class="mt-1 text-xs text-gray-500">Alphanumeric and underscores only</p>
                    </div>
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">Field Label</label>
                      <input
                        v-model="field.label"
                        type="text"
                        placeholder="e.g., Patient Name"
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">Field Type</label>
                      <select
                        v-model="field.type"
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                      >
                        <option value="text">Text</option>
                        <option value="textarea">Textarea</option>
                        <option value="email">Email</option>
                        <option value="number">Number</option>
                        <option value="date">Date</option>
                        <option value="select">Select</option>
                        <option value="file">File Upload</option>
                      </select>
                    </div>
                    <div>
                      <label class="block text-sm font-medium text-gray-700 mb-2">Required</label>
                      <select
                        v-model="field.required"
                        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                      >
                        <option :value="true">Yes</option>
                        <option :value="false">No</option>
                      </select>
                    </div>
                  </div>

                  <div v-if="field.type === 'select'" class="mt-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">Options (one per line)</label>
                    <textarea
                      v-model="field.options"
                      rows="3"
                      placeholder="Option 1&#10;Option 2&#10;Option 3"
                      class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    ></textarea>
                  </div>

                  <div class="mt-4">
                    <label class="block text-sm font-medium text-gray-700 mb-2">Placeholder</label>
                    <input
                      v-model="field.placeholder"
                      type="text"
                      placeholder="Enter placeholder text"
                      class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex justify-end space-x-4 pt-6 border-t border-gray-200">
              <button
                @click="closeFormModal"
                class="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                @click="saveForm"
                :disabled="saving"
                class="px-6 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
              >
                <svg v-if="saving" class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                <span>{{ saving ? 'Saving...' : 'Save Form' }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Preview Modal -->
    <div
      v-if="showPreview"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
    >
      <div
        @click.stop
        class="bg-white rounded-lg shadow-xl max-w-3xl w-full max-h-[90vh] overflow-y-auto"
      >
        <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex justify-between items-center">
          <h2 class="text-2xl font-bold text-gray-900">Form Preview: {{ previewForm_data?.formName }}</h2>
          <button
            @click="showPreview = false"
            class="text-gray-400 hover:text-gray-600"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-6">
          <p v-if="previewForm_data?.description" class="text-gray-600 mb-6">
            {{ previewForm_data.description }}
          </p>

          <div class="space-y-4">
            <div
              v-for="field in previewForm_data?.fields"
              :key="field.id"
              class="border border-gray-200 rounded-lg p-4"
            >
              <div class="flex items-start justify-between mb-2">
                <label class="block text-sm font-medium text-gray-700">
                  {{ field.label }}
                  <span v-if="field.required" class="text-red-600">*</span>
                </label>
                <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded">
                  {{ field.type }}
                </span>
              </div>

              <!-- Text input preview -->
              <input
                v-if="field.type === 'text' || field.type === 'email'"
                type="text"
                :placeholder="field.placeholder"
                disabled
                class="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50"
              />

              <!-- Textarea preview -->
              <textarea
                v-else-if="field.type === 'textarea'"
                :placeholder="field.placeholder"
                disabled
                rows="3"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50"
              ></textarea>

              <!-- Select preview -->
              <select
                v-else-if="field.type === 'select'"
                disabled
                class="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50"
              >
                <option value="">Select an option...</option>
                <option v-for="option in field.options" :key="option" :value="option">
                  {{ option }}
                </option>
              </select>

              <!-- Radio preview -->
              <div v-else-if="field.type === 'radio'" class="space-y-2">
                <label
                  v-for="option in field.options"
                  :key="option"
                  class="flex items-center space-x-2"
                >
                  <input type="radio" :name="field.id" disabled class="text-bloodsa-red" />
                  <span class="text-sm text-gray-700">{{ option }}</span>
                </label>
              </div>

              <!-- Number preview -->
              <input
                v-else-if="field.type === 'number'"
                type="number"
                :placeholder="field.placeholder"
                disabled
                class="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50"
              />

              <!-- Date preview -->
              <input
                v-else-if="field.type === 'date'"
                type="date"
                disabled
                class="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50"
              />

              <!-- File preview -->
              <div v-else-if="field.type === 'file'" class="border-2 border-dashed border-gray-300 rounded-lg p-4 text-center bg-gray-50">
                <svg class="w-8 h-8 text-gray-400 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                </svg>
                <p class="text-sm text-gray-500">
                  {{ field.allowMultiple ? 'Upload multiple files' : 'Upload file' }}
                </p>
              </div>

              <!-- Help text -->
              <p v-if="field.helpText" class="mt-2 text-xs text-gray-500">
                {{ field.helpText }}
              </p>
            </div>
          </div>

          <div class="mt-6 flex justify-end">
            <button
              @click="showPreview = false"
              class="px-6 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
            >
              Close Preview
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div 
      v-if="showDeleteModal" 
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
    >
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <!-- Warning Icon -->
          <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 mb-4">
            <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.964-1.333-3.732 0L3.732 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          
          <!-- Modal Content -->
          <div class="text-center">
            <h3 class="text-lg font-medium text-gray-900 mb-2">Delete Form Schema</h3>
            <div class="mt-2 px-7 py-3">
              <p class="text-sm text-gray-500 mb-4">
                Are you sure you want to delete this form? This action cannot be undone.
              </p>
              
              <!-- Form Info -->
              <div v-if="formToDelete" class="bg-gray-50 rounded-lg p-4 mb-4">
                <div class="text-left">
                  <p class="text-sm font-medium text-gray-900 mb-1">
                    {{ formToDelete.formName }}
                  </p>
                  <p class="text-xs text-gray-500">
                    {{ formToDelete.fields.length }} fields • 
                    {{ formToDelete.isActive ? 'Active' : 'Inactive' }}
                  </p>
                  <p v-if="formToDelete.description" class="text-xs text-gray-600 mt-2">
                    {{ formToDelete.description }}
                  </p>
                </div>
              </div>

              <div v-if="formToDelete?.isActive" class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 mb-4">
                <p class="text-xs text-yellow-800">
                  ⚠️ This form is currently active. Deleting it will prevent users from submitting new entries.
                </p>
              </div>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex justify-center space-x-3 mt-6">
              <button
                @click="cancelDelete"
                :disabled="loading"
                class="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Cancel
              </button>
              <button
                @click="confirmDelete"
                :disabled="loading"
                class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
              >
                <svg v-if="loading" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ loading ? 'Deleting...' : 'Delete Form' }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { registryService } from '@/services/registryService'
import type { FormSchema, FormField } from '@/services/registryService'
import { useToast } from '@/composables/useToast'

const loading = ref(false)
const saving = ref(false)
const showCreateForm = ref(false)
const showPreview = ref(false)
const showDeleteModal = ref(false)
const editingForm = ref<FormSchema | null>(null)
const previewForm_data = ref<FormSchema | null>(null)
const formToDelete = ref<FormSchema | null>(null)
const toast = useToast()

const forms = ref<FormSchema[]>([])
const totalForms = ref(0)
const currentPage = ref(1)
const itemsPerPage = ref(9) // 3x3 grid

// Computed properties
const totalPages = computed(() => Math.ceil(totalForms.value / itemsPerPage.value))

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, start + 4)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

const formData = ref<Partial<FormSchema>>({
  formName: '',
  description: '',
  fields: [],
  isActive: false
})

async function loadForms() {
  loading.value = true
  try {
    const response = await registryService.listFormSchemas({
      page: currentPage.value,
      limit: itemsPerPage.value
    })
    console.log('Forms response:', response)
    forms.value = Array.isArray(response.schemas) ? response.schemas : []
    totalForms.value = response.total || 0
    console.log('Forms array:', forms.value)
    console.log('Total forms:', totalForms.value)
  } catch (error) {
    console.error('Failed to load forms:', error)
    // If it's a 404, it means no forms exist yet - that's okay
    forms.value = []
    totalForms.value = 0
    // Only show error toast if it's not a 404
    if (error instanceof Error && !error.message.includes('404')) {
      toast.error('Failed to load forms. Please try again.')
    }
  } finally {
    loading.value = false
  }
}

async function saveForm() {
  saving.value = true
  try {
    if (editingForm.value) {
      await registryService.updateFormSchema(editingForm.value.id, formData.value)
      toast.success('Form updated successfully!')
    } else {
      await registryService.createFormSchema(formData.value)
      toast.success('Form created successfully!')
    }
    await loadForms()
    closeFormModal()
  } catch (error) {
    console.error('Failed to save form:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    toast.error(`Failed to save form: ${errorMessage}`)
  } finally {
    saving.value = false
  }
}

function editForm(form: FormSchema) {
  editingForm.value = form
  formData.value = { ...form }
  showCreateForm.value = true
}

function previewForm(form: FormSchema) {
  previewForm_data.value = form
  showPreview.value = true
}

async function toggleFormStatus(form: FormSchema) {
  try {
    const newStatus = !form.isActive
    console.log('Toggling form:', form.id, 'to', newStatus)
    await registryService.updateFormSchema(form.id, { isActive: newStatus })
    toast.success(newStatus ? 'Form activated successfully!' : 'Form deactivated successfully!')
    await loadForms()
  } catch (error) {
    console.error('Failed to toggle form status:', error)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    toast.error(`Failed to update form status: ${errorMessage}`)
  }
}

function deleteForm(form: FormSchema) {
  formToDelete.value = form
  showDeleteModal.value = true
}

async function confirmDelete() {
  if (!formToDelete.value) return
  
  console.log('Deleting form:', formToDelete.value.id, 'Name:', formToDelete.value.formName)
  
  try {
    await registryService.deleteFormSchema(formToDelete.value.id)
    toast.success(`"${formToDelete.value.formName}" deleted successfully!`)
    showDeleteModal.value = false
    formToDelete.value = null
    await loadForms()
  } catch (error) {
    console.error('Failed to delete form:', error)
    console.error('Form ID was:', formToDelete.value?.id)
    const errorMessage = error instanceof Error ? error.message : 'Unknown error'
    toast.error(`Failed to delete form: ${errorMessage}`)
  }
}

function cancelDelete() {
  showDeleteModal.value = false
  formToDelete.value = null
}

function addField() {
  const fieldCount = (formData.value.fields?.length || 0) + 1
  formData.value.fields?.push({
    id: '',
    label: '',
    type: 'text',
    required: false,
    placeholder: '',
    displayOrder: fieldCount
  })
}

function removeField(index: number) {
  formData.value.fields?.splice(index, 1)
}

function closeFormModal() {
  showCreateForm.value = false
  editingForm.value = null
  formData.value = {
    formName: '',
    description: '',
    fields: [],
    isActive: false
  }
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString()
}

// Pagination functions
const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadForms()
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadForms()
  }
}

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadForms()
  }
}

const goToFirstPage = () => {
  if (currentPage.value !== 1) {
    currentPage.value = 1
    loadForms()
  }
}

const goToLastPage = () => {
  if (currentPage.value !== totalPages.value) {
    currentPage.value = totalPages.value
    loadForms()
  }
}

onMounted(() => {
  loadForms()
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

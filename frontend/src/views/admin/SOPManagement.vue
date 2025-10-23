<template>
  <div class="p-6">
    <!-- Page Header -->
    <div class="mb-8 flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">SOP Management</h1>
        <p class="text-gray-600 mt-2">Manage Standard Operating Procedure categories</p>
      </div>
      <div class="flex gap-3">
        <button
          @click="seedCategories"
          v-if="categories.length === 0"
          class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="loading || !dropboxConfigured"
          :title="!dropboxConfigured ? 'Configure Dropbox first in System Settings' : ''"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Seed Initial Categories
        </button>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="!dropboxConfigured"
          :title="!dropboxConfigured ? 'Configure Dropbox first in System Settings' : ''"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Create Category
        </button>
      </div>
    </div>

    <!-- Dropbox Not Configured Warning -->
    <div v-if="!checkingDropbox && !dropboxConfigured" class="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6">
      <div class="flex items-start">
        <svg class="w-5 h-5 text-yellow-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div class="ml-3 flex-1">
          <h3 class="text-sm font-medium text-yellow-800">Dropbox Not Configured</h3>
          <p class="text-sm text-yellow-700 mt-1">
            Dropbox must be configured before creating or managing SOP categories. 
            Please configure Dropbox in System Settings first.
          </p>
          <router-link 
            to="/admin/settings" 
            class="inline-flex items-center mt-3 px-4 py-2 bg-yellow-600 text-white text-sm rounded-md hover:bg-yellow-700 transition-colors"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            Go to System Settings
          </router-link>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading && categories.length === 0" class="text-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-bloodsa-red mx-auto"></div>
      <p class="text-gray-600 mt-4">Loading categories...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="!loading && categories.length === 0" class="bg-white rounded-xl shadow-lg p-12 text-center">
      <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <h3 class="text-xl font-semibold text-gray-900 mb-2">No SOP Categories Yet</h3>
      <p v-if="dropboxConfigured" class="text-gray-600 mb-2">Get started by seeding the 4 default categories or create your own.</p>
      <p v-else class="text-gray-600 mb-2">Configure Dropbox first, then create your SOP categories.</p>
      <p v-if="dropboxConfigured" class="text-sm text-gray-500 mb-6">
        <strong>Seed Categories</strong> will automatically create: Anemia, Lymphoma, Myeloma, and General Business
      </p>
      <div class="flex gap-3 justify-center">
        <button
          v-if="dropboxConfigured"
          @click="seedCategories"
          class="px-6 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors font-medium"
        >
          <svg class="w-5 h-5 inline-block mr-2 -mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Seed Initial Categories
        </button>
        <button
          v-if="dropboxConfigured"
          @click="openCreateModal"
          class="px-6 py-3 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors font-medium"
        >
          <svg class="w-5 h-5 inline-block mr-2 -mt-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Create Custom Category
        </button>
        <router-link
          v-else
          to="/admin/settings"
          class="px-6 py-3 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 transition-colors font-medium inline-flex items-center"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          Configure Dropbox First
        </router-link>
      </div>
    </div>

    <!-- Categories Table -->
    <div v-else class="bg-white rounded-xl shadow-lg overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Image</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Slug</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Order</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="category in sortedCategories" :key="category.id">
            <td class="px-6 py-4 whitespace-nowrap">
              <img 
                v-if="category.imagePath" 
                :src="`${category.imagePath}`" 
                :alt="category.name"
                class="h-12 w-12 rounded object-cover"
              />
              <div v-else class="h-12 w-12 rounded bg-gray-200 flex items-center justify-center">
                <svg class="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
            </td>
            <td class="px-6 py-4">
              <div class="text-sm font-medium text-gray-900">{{ category.name }}</div>
              <div class="text-sm text-gray-500">{{ category.description || 'No description' }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ category.slug }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span 
                :class="category.isActive ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                class="px-2 py-1 text-xs font-semibold rounded-full"
              >
                {{ category.isActive ? 'Active' : 'Inactive' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ category.displayOrder }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
              <button
                @click="viewFiles(category)"
                class="text-blue-600 hover:text-blue-900"
                title="View Files"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                </svg>
              </button>
              <button
                @click="openEditModal(category)"
                class="text-indigo-600 hover:text-indigo-900 disabled:opacity-50 disabled:cursor-not-allowed"
                :disabled="!dropboxConfigured"
                :title="!dropboxConfigured ? 'Configure Dropbox first in System Settings' : 'Edit'"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </button>
              <button
                @click="toggleActive(category)"
                class="text-yellow-600 hover:text-yellow-900"
                :title="category.isActive ? 'Deactivate' : 'Activate'"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </button>
              <button
                @click="confirmDelete(category)"
                class="text-red-600 hover:text-red-900"
                title="Delete"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="p-6">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-2xl font-bold text-gray-900">
              {{ isEditMode ? 'Edit Category' : 'Create Category' }}
            </h2>
            <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form @submit.prevent="submitForm" class="space-y-4">
            <!-- Name -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Name <span class="text-red-500">*</span>
              </label>
              <input
                v-model="formData.name"
                type="text"
                required
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                placeholder="e.g., Anemia"
              />
            </div>

            <!-- Description -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
              <textarea
                v-model="formData.description"
                rows="3"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                placeholder="Brief description of this category..."
              ></textarea>
            </div>

            <!-- Display Order -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Display Order <span class="text-red-500">*</span>
              </label>
              <input
                v-model.number="formData.displayOrder"
                type="number"
                required
                min="0"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
              />
            </div>

            <!-- Image Upload -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Category Image</label>
              <div class="flex items-center gap-4">
                <div 
                  v-if="imagePreview || (isEditMode && editingCategory?.imagePath)"
                  class="relative"
                >
                  <img
                    :src="imagePreview || `${editingCategory?.imagePath}`"
                    alt="Preview"
                    class="h-24 w-24 rounded object-cover"
                  />
                  <button
                    @click="clearImage"
                    type="button"
                    class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full p-1"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
                <label class="flex-1 cursor-pointer">
                  <div class="border-2 border-dashed border-gray-300 rounded-lg p-4 text-center hover:border-bloodsa-red transition-colors">
                    <svg class="w-8 h-8 text-gray-400 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                    </svg>
                    <p class="text-sm text-gray-600">Click to upload image</p>
                    <p class="text-xs text-gray-500 mt-1">JPG, PNG, WEBP (max 5MB)</p>
                  </div>
                  <input
                    type="file"
                    @change="handleImageUpload"
                    accept="image/jpeg,image/jpg,image/png,image/webp"
                    class="hidden"
                  />
                </label>
              </div>
            </div>

            <!-- Active Status (Edit mode only) -->
            <div v-if="isEditMode">
              <label class="flex items-center gap-2">
                <input
                  v-model="formData.isActive"
                  type="checkbox"
                  class="w-4 h-4 text-bloodsa-red border-gray-300 rounded focus:ring-bloodsa-red"
                />
                <span class="text-sm font-medium text-gray-700">Active</span>
              </label>
            </div>

            <!-- Actions -->
            <div class="flex justify-end gap-3 pt-4">
              <button
                type="button"
                @click="closeModal"
                class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="uploading"
                class="px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50"
              >
                {{ uploading ? 'Saving...' : isEditMode ? 'Update' : 'Create' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Files Modal -->
    <div v-if="showFilesModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-2xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
        <div class="p-6">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-2xl font-bold text-gray-900">
              Files in "{{ selectedCategory?.name }}"
            </h2>
            <button @click="showFilesModal = false" class="text-gray-400 hover:text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Files List -->
          <div v-if="loadingFiles" class="text-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-bloodsa-red mx-auto"></div>
          </div>
          <div v-else-if="files.length === 0" class="text-center py-8">
            <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
            <p class="text-gray-600">No files in this category</p>
            <p class="text-sm text-gray-500 mt-2">Upload files manually to Dropbox folder:</p>
            <code class="text-xs bg-gray-100 px-2 py-1 rounded mt-2 inline-block">{{ selectedCategoryDropboxPath }}</code>
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="file in files"
              :key="file.path"
              class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50"
            >
              <div class="flex items-center gap-3">
                <svg class="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
                <div>
                  <p class="text-sm font-medium text-gray-900">{{ file.name }}</p>
                  <p class="text-xs text-gray-500">
                    {{ formatFileSize(file.size) }} • {{ formatDate(file.modifiedTime) }}
                  </p>
                </div>
              </div>
              <button
                @click="downloadFile(file)"
                class="px-3 py-1 bg-bloodsa-red text-white text-sm rounded hover:bg-red-700 transition-colors"
              >
                Download
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteConfirm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl shadow-2xl max-w-md w-full p-6">
        <h3 class="text-xl font-bold text-gray-900 mb-4">Confirm Delete</h3>
        <p class="text-gray-600 mb-6">
          Are you sure you want to delete "{{ categoryToDelete?.name }}"? 
          The Dropbox folder will remain intact.
        </p>
        <div class="flex justify-end gap-3">
          <button
            @click="showDeleteConfirm = false"
            class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
          <button
            @click="deleteCategory"
            class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { sopService } from '@/services/sopService'
import type { SOPCategory, SOPFile, CreateCategoryRequest, UpdateCategoryRequest } from '@/types/sop'
import { useToast } from '@/composables/useToast'
import { dropboxAdminService } from '@/services/dropboxAdminService'

const toast = useToast()

const categories = ref<SOPCategory[]>([])
const loading = ref(false)
const dropboxConfigured = ref(false)
const checkingDropbox = ref(true)
const showModal = ref(false)
const showFilesModal = ref(false)
const showDeleteConfirm = ref(false)
const isEditMode = ref(false)
const editingCategory = ref<SOPCategory | null>(null)
const categoryToDelete = ref<SOPCategory | null>(null)
const selectedCategory = ref<SOPCategory | null>(null)
const files = ref<SOPFile[]>([])
const loadingFiles = ref(false)
const uploading = ref(false)
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)

const formData = ref<CreateCategoryRequest & { isActive?: boolean }>({
  name: '',
  description: '',
  imagePath: '',
  displayOrder: 0,
  isActive: true
})

const sortedCategories = computed(() => {
  return [...categories.value].sort((a, b) => a.displayOrder - b.displayOrder)
})

const selectedCategoryDropboxPath = computed(() => {
  const path = selectedCategory.value?.dropboxPath || ''
  // Decode URL-encoded path for user-friendly display
  try {
    return decodeURIComponent(path)
  } catch {
    return path
  }
})

onMounted(async () => {
  await checkDropboxStatus()
  loadCategories()
})

async function checkDropboxStatus() {
  checkingDropbox.value = true
  try {
    const status = await dropboxAdminService.getStatus()
    dropboxConfigured.value = status.isConnected || false
  } catch (error: any) {
    console.error('Failed to check Dropbox status:', error)
    dropboxConfigured.value = false
  } finally {
    checkingDropbox.value = false
  }
}

async function loadCategories() {
  loading.value = true
  try {
    console.log('Loading categories...')
    const response = await sopService.listCategories({ limit: 100 })
    console.log('Categories response:', response)
    categories.value = response.categories || []
  } catch (error: any) {
    console.error('Failed to load categories:', error)
    toast.error(`Failed to load categories: ${error.message}`)
    categories.value = []
  } finally {
    loading.value = false
    console.log('Loading complete, categories:', categories.value.length)
  }
}

async function seedCategories() {
  if (!dropboxConfigured.value) {
    toast.error('Please configure Dropbox in System Settings before creating categories.')
    return
  }
  
  if (!confirm('This will create 4 initial categories (Anemia, Lymphoma, Myeloma, General Business). Continue?')) {
    return
  }
  
  loading.value = true
  try {
    console.log('Seeding categories...')
    const result = await sopService.seedCategories()
    console.log('Seed result:', result)
    await loadCategories()
    toast.success(`Successfully created ${result.count} categories!`)
  } catch (error: any) {
    console.error('Failed to seed categories:', error)
    toast.error(`Failed to seed categories: ${error.message}`)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  if (!dropboxConfigured.value) {
    toast.error('Please configure Dropbox in System Settings before creating categories.')
    return
  }
  
  isEditMode.value = false
  editingCategory.value = null
  formData.value = {
    name: '',
    description: '',
    imagePath: '',
    displayOrder: categories.value.length + 1,
    isActive: true
  }
  imageFile.value = null
  imagePreview.value = null
  showModal.value = true
}

function openEditModal(category: SOPCategory) {
  if (!dropboxConfigured.value) {
    toast.error('Please configure Dropbox in System Settings before editing categories.')
    return
  }
  
  isEditMode.value = true
  editingCategory.value = category
  formData.value = {
    name: category.name,
    description: category.description || '',
    imagePath: category.imagePath || '',
    displayOrder: category.displayOrder,
    isActive: category.isActive
  }
  imageFile.value = null
  imagePreview.value = null
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  isEditMode.value = false
  editingCategory.value = null
  imageFile.value = null
  imagePreview.value = null
}

function handleImageUpload(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    imageFile.value = file
    imagePreview.value = URL.createObjectURL(file)
  }
}

function clearImage() {
  imageFile.value = null
  imagePreview.value = null
  formData.value.imagePath = ''
}

async function submitForm() {
  uploading.value = true
  try {
    // Upload image if selected
    if (imageFile.value) {
      const imagePath = await sopService.uploadImage(imageFile.value)
      formData.value.imagePath = imagePath
    }

    if (isEditMode.value && editingCategory.value) {
      // Update
      const updateData: UpdateCategoryRequest = {}
      if (formData.value.name !== editingCategory.value.name) updateData.name = formData.value.name
      if (formData.value.description !== editingCategory.value.description) updateData.description = formData.value.description
      if (formData.value.imagePath !== editingCategory.value.imagePath) updateData.imagePath = formData.value.imagePath
      if (formData.value.displayOrder !== editingCategory.value.displayOrder) updateData.displayOrder = formData.value.displayOrder
      if (formData.value.isActive !== editingCategory.value.isActive) updateData.isActive = formData.value.isActive

      await sopService.updateCategory(editingCategory.value.id, updateData)
      toast.success(`Successfully updated "${formData.value.name}" category!`)
    } else {
      // Create
      await sopService.createCategory(formData.value)
      toast.success(`Successfully created "${formData.value.name}" category!`)
    }

    await loadCategories()
    closeModal()
  } catch (error: any) {
    console.error('Failed to save category:', error)
    toast.error(error.message || 'Failed to save category')
  } finally {
    uploading.value = false
  }
}

async function toggleActive(category: SOPCategory) {
  const action = category.isActive ? 'deactivated' : 'activated'
  try {
    await sopService.updateCategory(category.id, { isActive: !category.isActive })
    await loadCategories()
    toast.success(`Category "${category.name}" has been ${action}`)
  } catch (error: any) {
    console.error('Failed to toggle status:', error)
    toast.error(error.message || 'Failed to toggle status')
  }
}

function confirmDelete(category: SOPCategory) {
  categoryToDelete.value = category
  showDeleteConfirm.value = true
}

async function deleteCategory() {
  if (!categoryToDelete.value) return

  const categoryName = categoryToDelete.value.name
  try {
    await sopService.deleteCategory(categoryToDelete.value.id)
    await loadCategories()
    showDeleteConfirm.value = false
    categoryToDelete.value = null
    toast.success(`Category "${categoryName}" has been deleted`)
  } catch (error: any) {
    console.error('Failed to delete category:', error)
    toast.error(error.message || 'Failed to delete category')
  }
}

async function viewFiles(category: SOPCategory) {
  selectedCategory.value = category
  showFilesModal.value = true
  loadingFiles.value = true
  try {
    const fileData = await sopService.getCategoryFiles(category.id)
    // Handle null, undefined, or empty responses
    files.value = Array.isArray(fileData) ? fileData : []
    if (files.value.length === 0) {
      toast.info(`No files found in "${category.name}". Upload files to Dropbox folder: ${category.dropboxPath}`)
    }
  } catch (error: any) {
    console.error('Failed to load files:', error)
    toast.error(error.message || 'Failed to load files. Ensure Dropbox is configured in System Settings.')
    files.value = []
  } finally {
    loadingFiles.value = false
  }
}

async function downloadFile(file: SOPFile) {
  if (!selectedCategory.value) return

  try {
    const downloadLink = await sopService.getFileDownloadLink(selectedCategory.value.id, file.name)
    window.open(downloadLink, '_blank')
    toast.success(`Downloading "${file.name}"`)
  } catch (error: any) {
    console.error('Failed to get download link:', error)
    toast.error(error.message || 'Failed to get download link')
  }
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}
</script>


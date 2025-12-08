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
      <h1 class="text-3xl font-bold text-gray-900">Registry Submissions</h1>
      <p class="mt-2 text-gray-600">
        Review and manage all registry form submissions.
      </p>
    </div>

    <!-- Filters and Search -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
          <select
            v-model="filters.status"
            @change="loadSubmissions"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
          >
            <option value="">All Statuses</option>
            <option value="submitted">Submitted</option>
            <option value="pending">Pending</option>
            <option value="approved">Approved</option>
            <option value="rejected">Rejected</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Search User</label>
          <input
            v-model="filters.userSearch"
            type="text"
            placeholder="Search by name or email..."
            @input="debouncedLoadSubmissions"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Date Range</label>
          <div class="flex items-center gap-2">
            <input
              v-model="filters.dateFrom"
              type="date"
              @change="loadSubmissions"
              placeholder="Start date"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent text-sm"
            />
            <span class="text-gray-500">to</span>
            <input
              v-model="filters.dateTo"
              type="date"
              @change="loadSubmissions"
              placeholder="End date"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent text-sm"
            />
            <button
              v-if="filters.dateFrom || filters.dateTo"
              @click="clearDateFilter"
              class="p-2 text-gray-400 hover:text-gray-600 transition-colors"
              title="Clear date filter"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-4 border-bloodsa-red"></div>
    </div>

    <!-- Submissions Table -->
    <div v-else class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                User
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Form
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Submitted
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Documents
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="submission in submissions" :key="submission.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="w-10 h-10 bg-bloodsa-red rounded-full flex items-center justify-center text-white text-sm font-medium">
                    {{ getUserInitials(submission.userName || submission.userId) }}
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">{{ submission.userName || submission.userId }}</div>
                    <div class="text-sm text-gray-500">{{ submission.userEmail || 'User' }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ submission.formName || submission.formSchemaId }}</div>
                <div class="text-sm text-gray-500">Form</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  :class="getStatusClass(submission.status)"
                  class="px-2 py-1 text-xs font-medium rounded-full"
                >
                  {{ submission.status }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ formatDate(submission.createdAt) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-sm text-gray-900">{{ submission.uploadedDocuments.length }} files</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <div class="flex items-center space-x-2">
                  <button
                    @click="viewSubmission(submission)"
                    class="text-blue-600 hover:text-blue-700"
                  >
                    View
                  </button>
                  <div class="relative" v-if="submission.status === 'pending'">
                    <button
                      @click="showStatusMenu = submission.id"
                      class="text-gray-600 hover:text-gray-700"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
                      </svg>
                    </button>
                    <div
                      v-if="showStatusMenu === submission.id"
                      class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-10 border border-gray-200"
                    >
                      <div class="py-1">
                        <button
                          @click="updateStatus(submission, 'approved')"
                          class="block w-full text-left px-4 py-2 text-sm text-green-700 hover:bg-green-50"
                        >
                          Approve
                        </button>
                        <button
                          @click="updateStatus(submission, 'rejected')"
                          class="block w-full text-left px-4 py-2 text-sm text-red-700 hover:bg-red-50"
                        >
                          Reject
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
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
              <span class="font-medium">{{ (currentPage - 1) * pageSize + 1 }}</span>
              to
              <span class="font-medium">{{ Math.min(currentPage * pageSize, total) }}</span>
              of
              <span class="font-medium">{{ total }}</span>
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

    <!-- Submission Details Modal -->
    <div v-if="selectedSubmission" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div @click.stop class="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div class="p-6">
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-2xl font-bold text-gray-900">Submission Details</h2>
            <button
              @click="selectedSubmission = null"
              class="text-gray-400 hover:text-gray-600"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="space-y-6">
            <!-- Submission Info -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700">User ID</label>
                <p class="mt-1 text-sm text-gray-900">{{ selectedSubmission.userId }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Status</label>
                <span
                  :class="getStatusClass(selectedSubmission.status)"
                  class="mt-1 inline-block px-2 py-1 text-xs font-medium rounded-full"
                >
                  {{ selectedSubmission.status }}
                </span>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Submitted</label>
                <p class="mt-1 text-sm text-gray-900">{{ formatDate(selectedSubmission.createdAt) }}</p>
              </div>
            </div>

            <!-- Form Data -->
            <div>
              <h3 class="text-lg font-semibold text-gray-900 mb-4">Form Data</h3>
              <div class="bg-gray-50 rounded-lg p-4">
                <div v-if="Object.keys(selectedSubmission.formData || {}).length > 0" class="space-y-3">
                  <div v-for="(value, key) in selectedSubmission.formData" :key="key" class="border-b border-gray-200 pb-2 last:border-0">
                    <dt class="text-sm font-medium text-gray-700">{{ formatFieldName(key) }}</dt>
                    <dd class="mt-1 text-sm text-gray-900">{{ formatFieldValue(value) }}</dd>
                  </div>
                </div>
                <p v-else class="text-sm text-gray-500 italic">No form data submitted</p>
              </div>
            </div>

            <!-- Documents -->
            <div v-if="selectedSubmission.uploadedDocuments.length > 0">
              <h3 class="text-lg font-semibold text-gray-900 mb-4">Uploaded Documents</h3>
              <div class="space-y-2">
                <div
                  v-for="(document, index) in selectedSubmission.uploadedDocuments"
                  :key="index"
                  class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
                >
                  <div class="flex items-center space-x-3">
                    <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <span class="text-sm text-gray-900">{{ document }}</span>
                  </div>
                  <button
                    @click="downloadDocument(document)"
                    class="text-blue-600 hover:text-blue-700 text-sm"
                  >
                    Download
                  </button>
                </div>
              </div>
            </div>

            <!-- Status Change Actions -->
            <div class="border-t pt-6">
              <h3 class="text-lg font-semibold text-gray-900 mb-4">Change Status</h3>
              <div class="flex flex-wrap gap-3">
                <button
                  v-if="selectedSubmission.status !== 'submitted'"
                  @click="changeStatus('submitted')"
                  :disabled="updatingStatus"
                  class="px-4 py-2 border border-blue-300 text-blue-700 rounded-lg hover:bg-blue-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  <svg v-if="updatingStatus" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Mark as Submitted
                </button>
                <button
                  v-if="selectedSubmission.status !== 'pending'"
                  @click="changeStatus('pending')"
                  :disabled="updatingStatus"
                  class="px-4 py-2 border border-yellow-300 text-yellow-700 rounded-lg hover:bg-yellow-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  <svg v-if="updatingStatus" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Mark as Pending
                </button>
                <button
                  v-if="selectedSubmission.status !== 'approved'"
                  @click="changeStatus('approved')"
                  :disabled="updatingStatus"
                  class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  <svg v-if="updatingStatus" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Approve
                </button>
                <button
                  v-if="selectedSubmission.status !== 'rejected'"
                  @click="openRejectModal()"
                  :disabled="updatingStatus"
                  class="px-4 py-2 border border-red-300 text-red-700 rounded-lg hover:bg-red-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  <svg v-if="updatingStatus" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Reject
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Reject Reason Modal -->
    <div v-if="showRejectModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div @click.stop class="bg-white rounded-lg shadow-xl max-w-md w-full">
        <div class="p-6">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-xl font-bold text-gray-900">Reject Submission</h3>
            <button
              @click="closeRejectModal"
              class="text-gray-400 hover:text-gray-600"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Rejection Reason <span class="text-red-500">*</span>
            </label>
            <textarea
              v-model="rejectionReason"
              rows="4"
              placeholder="Please provide a reason for rejection. This will be sent to the user via email."
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
              :class="{ 'border-red-500': rejectionError }"
            ></textarea>
            <p v-if="rejectionError" class="text-sm text-red-500 mt-1">{{ rejectionError }}</p>
            <p class="text-sm text-gray-500 mt-1">The user will receive an email with this reason.</p>
          </div>

          <div class="flex justify-end space-x-3">
            <button
              @click="closeRejectModal"
              :disabled="updatingStatus"
              class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
            <button
              @click="confirmReject"
              :disabled="updatingStatus"
              class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              <svg v-if="updatingStatus" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ updatingStatus ? 'Rejecting...' : 'Confirm Rejection' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { registryService } from '@/services/registryService'
import type { Submission } from '@/services/registryService'
import { useToast } from '@/composables/useToast'

const loading = ref(false)
const submissions = ref<Submission[]>([])
const selectedSubmission = ref<Submission | null>(null)
const showStatusMenu = ref<string | null>(null)
const showRejectModal = ref(false)
const rejectionReason = ref('')
const rejectionError = ref('')
const updatingStatus = ref(false)
const toast = useToast()

const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const totalPages = ref(0)

const filters = ref({
  status: '',
  userSearch: '',
  dateFrom: '',
  dateTo: ''
})

let debounceTimer: ReturnType<typeof setTimeout>

async function loadSubmissions() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      ...(filters.value.status && { status: filters.value.status }),
      ...(filters.value.userSearch && { userSearch: filters.value.userSearch }),
      ...(filters.value.dateFrom && { dateFrom: filters.value.dateFrom }),
      ...(filters.value.dateTo && { dateTo: filters.value.dateTo })
    }

    const response = await registryService.getAllSubmissions(params)
    submissions.value = Array.isArray(response.submissions) ? response.submissions : []
    total.value = response.total || 0
    totalPages.value = Math.ceil((response.total || 0) / pageSize.value)
  } catch (error) {
    console.error('Failed to load submissions:', error)
    // Set empty state instead of showing error for 404
    submissions.value = []
    total.value = 0
    totalPages.value = 0
    // Only show error toast if it's not a 404
    if (error instanceof Error && !error.message.includes('404')) {
      toast.error('Failed to load submissions. Please try again.')
    }
  } finally {
    loading.value = false
  }
}

function debouncedLoadSubmissions() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    loadSubmissions()
  }, 500)
}

function clearDateFilter() {
  filters.value.dateFrom = ''
  filters.value.dateTo = ''
  loadSubmissions()
}

async function changeStatus(status: 'submitted' | 'pending' | 'approved' | 'rejected') {
  if (!selectedSubmission.value) return
  
  // For rejection, we need a reason - handled by openRejectModal
  if (status === 'rejected') {
    openRejectModal()
    return
  }
  
  updatingStatus.value = true
  
  try {
    await registryService.updateSubmissionStatus(selectedSubmission.value.id, status)
    
    // Update local state
    selectedSubmission.value.status = status
    const submissionIndex = submissions.value.findIndex(s => s.id === selectedSubmission.value!.id)
    if (submissionIndex !== -1 && submissions.value[submissionIndex]) {
      submissions.value[submissionIndex].status = status
    }
    
    const statusText = status === 'approved' ? 'approved' : 
                      status === 'pending' ? 'marked as pending' : 
                      'marked as submitted'
    toast.success(`Submission ${statusText} successfully! User has been notified via email.`)
    
    // Close the details modal
    selectedSubmission.value = null
  } catch (error: any) {
    console.error('Failed to update status:', error)
    toast.error(error.message || 'Failed to update submission status. Please try again.')
  } finally {
    updatingStatus.value = false
  }
}

function openRejectModal() {
  showRejectModal.value = true
  rejectionReason.value = ''
  rejectionError.value = ''
}

function closeRejectModal() {
  showRejectModal.value = false
  rejectionReason.value = ''
  rejectionError.value = ''
}

async function confirmReject() {
  if (!selectedSubmission.value) return
  
  // Validate rejection reason
  if (!rejectionReason.value.trim()) {
    rejectionError.value = 'Rejection reason is required'
    return
  }
  
  updatingStatus.value = true
  rejectionError.value = ''
  
  try {
    // Update status with rejection reason
    await registryService.updateSubmissionStatus(selectedSubmission.value.id, 'rejected', rejectionReason.value)
    
    // Update local state
    selectedSubmission.value.status = 'rejected'
    const submissionIndex = submissions.value.findIndex(s => s.id === selectedSubmission.value!.id)
    if (submissionIndex !== -1 && submissions.value[submissionIndex]) {
      submissions.value[submissionIndex].status = 'rejected'
    }
    
    toast.success('Submission rejected. User has been notified via email with the reason.')
    
    // Close both modals
    showRejectModal.value = false
    selectedSubmission.value = null
  } catch (error: any) {
    console.error('Failed to reject submission:', error)
    toast.error(error.message || 'Failed to reject submission. Please try again.')
  } finally {
    updatingStatus.value = false
  }
}

async function updateStatus(submission: Submission, status: 'approved' | 'rejected') {
  try {
    await registryService.updateSubmissionStatus(submission.id, status)
    submission.status = status
    showStatusMenu.value = null
    if (selectedSubmission.value?.id === submission.id) {
      selectedSubmission.value.status = status
    }
    toast.success(`Submission ${status} successfully!`)
  } catch (error) {
    console.error('Failed to update status:', error)
    toast.error('Failed to update submission status. Please try again.')
  }
}

function viewSubmission(submission: Submission) {
  selectedSubmission.value = submission
}

function viewDocuments(submission: Submission) {
  // TODO: Implement document viewer
  console.log('View documents for submission:', submission.id)
}

async function downloadDocument(documentName: string) {
  if (!selectedSubmission.value) return
  
  try {
    // Construct the full path to the document
    const fullPath = `${selectedSubmission.value.documentsPath}/${documentName}`
    
    // Use the existing Dropbox service to get download link
    // We'll need to create a method in the registry service for this
    const response = await fetch(`/api/registry/document-download?path=${encodeURIComponent(fullPath)}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to get download link')
    }
    
    const data = await response.json()
    // Open the download link in a new tab
    window.open(data.link, '_blank')
  } catch (err: any) {
    console.error('Failed to download document:', err)
    toast.error('Failed to download document. Please try again.')
  }
}

function formatFieldName(fieldKey: string): string {
  // Convert camelCase or snake_case to Title Case
  return fieldKey
    .replace(/([A-Z])/g, ' $1')
    .replace(/_/g, ' ')
    .replace(/^./, str => str.toUpperCase())
    .trim()
}

function formatFieldValue(value: any): string {
  if (value === null || value === undefined) return 'N/A'
  if (typeof value === 'boolean') return value ? 'Yes' : 'No'
  if (typeof value === 'object') return JSON.stringify(value, null, 2)
  return String(value)
}

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, start + 4)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

function previousPage() {
  if (currentPage.value > 1) {
    currentPage.value--
    loadSubmissions()
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadSubmissions()
  }
}

function goToPage(page: number) {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadSubmissions()
  }
}

function goToFirstPage() {
  if (currentPage.value !== 1) {
    currentPage.value = 1
    loadSubmissions()
  }
}

function goToLastPage() {
  if (currentPage.value !== totalPages.value) {
    currentPage.value = totalPages.value
    loadSubmissions()
  }
}

function getStatusClass(status: string) {
  switch (status) {
    case 'submitted':
      return 'bg-blue-100 text-blue-800'
    case 'pending':
      return 'bg-yellow-100 text-yellow-800'
    case 'approved':
      return 'bg-green-100 text-green-800'
    case 'rejected':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

function getUserInitials(userNameOrId: string) {
  if (!userNameOrId) return '??'
  
  // If it looks like a name (contains space), get initials from first and last name
  if (userNameOrId.includes(' ')) {
    const names = userNameOrId.split(' ')
    return (names[0]?.[0] || '') + (names[names.length - 1]?.[0] || '')
  }
  
  // Otherwise, use first two characters (for user IDs)
  return userNameOrId.substring(0, 2).toUpperCase()
}

function getFormName(formSchemaId: string) {
  // This is now handled by the backend populating formName
  return formSchemaId
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString()
}

// Close status menu when clicking outside
document.addEventListener('click', (e) => {
  if (!(e.target as Element).closest('.relative')) {
    showStatusMenu.value = null
  }
})

onMounted(() => {
  loadSubmissions()
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

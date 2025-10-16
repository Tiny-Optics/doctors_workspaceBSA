<template>
  <div class="p-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Institution Management</h1>
        <p class="text-gray-600 mt-2">Manage institutions, organizations, and facilities</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="bg-bloodsa-red text-white px-6 py-3 rounded-lg hover:bg-opacity-90 transition-colors font-medium flex items-center space-x-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
        </svg>
        <span>Add Institution</span>
      </button>
    </div>

    <!-- Filters and Search -->
    <div class="bg-white rounded-xl shadow-lg p-6 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Search -->
        <div class="md:col-span-2">
          <label class="block text-sm font-medium text-gray-700 mb-2">Search Institutions</label>
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search by name, city, or province..."
              class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
            <svg class="absolute left-3 top-2.5 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>

        <!-- Type Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Filter by Type</label>
          <select
            v-model="typeFilter"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
          >
            <option value="">All Types</option>
            <option value="university">University</option>
            <option value="hospital">Hospital</option>
            <option value="laboratory">Lab</option>
            <option value="research_center">Research</option>
            <option value="government">Govt</option>
            <option value="private_practice">Private</option>
            <option value="ngo">NGO</option>
            <option value="other">Other</option>
          </select>
        </div>

        <!-- Status Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Filter by Status</label>
          <select
            v-model="statusFilter"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
          >
            <option value="">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error</h3>
          <div class="mt-2 text-sm text-red-700">
            {{ error }}
          </div>
        </div>
        <div class="ml-auto pl-3">
          <button @click="error = null" class="text-red-400 hover:text-red-600">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Institutions Table -->
    <div class="bg-white rounded-xl shadow-lg overflow-hidden">
      <!-- Table Content with Loading Overlay -->
      <div class="overflow-x-auto relative min-h-[400px]">
        <!-- Loading Overlay -->
        <transition name="fade">
          <div 
            v-if="loading" 
            class="absolute inset-0 bg-white bg-opacity-90 flex items-center justify-center z-10 backdrop-blur-sm"
          >
            <div class="inline-flex items-center space-x-2">
              <svg class="animate-spin h-5 w-5 text-bloodsa-red" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span class="text-gray-600">Loading...</span>
            </div>
          </div>
        </transition>
        
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Institution
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Type
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Location
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Contact
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <!-- Empty State -->
            <tr v-if="!loading && displayInstitutions.length === 0">
              <td colspan="6" class="px-6 py-12 text-center">
                <div class="flex flex-col items-center justify-center space-y-3">
                  <svg class="w-16 h-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                  </svg>
                  <div class="text-gray-500">
                    <p class="text-lg font-medium">No institutions found</p>
                    <p class="text-sm mt-1">
                      <span v-if="searchQuery || typeFilter || statusFilter">
                        Try adjusting your search or filters
                      </span>
                      <span v-else>
                        No institutions exist in the system yet
                      </span>
                    </p>
                  </div>
                </div>
              </td>
            </tr>
            
            <!-- Institution Rows -->
            <tr v-for="institution in displayInstitutions" :key="institution.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="w-10 h-10 bg-bloodsa-red rounded-full flex items-center justify-center text-white font-semibold text-xs">
                    {{ institution.shortName || institution.name.substring(0, 2).toUpperCase() }}
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">
                      {{ institution.name }}
                    </div>
                    <div class="text-sm text-gray-500">{{ institution.shortName }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="getTypeBadgeClass(institution.type)"
                >
                  {{ getTypeDisplayName(institution.type) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                <div>{{ institution.city }}</div>
                <div v-if="institution.province" class="text-xs text-gray-500">{{ institution.province }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <div v-if="institution.phone">{{ institution.phone }}</div>
                <div v-if="institution.email" class="text-xs">{{ institution.email }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="institution.isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                >
                  {{ institution.isActive ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div class="flex items-center justify-end space-x-2">
                  <button
                    @click="viewInstitution(institution)"
                    class="text-blue-600 hover:text-blue-900 p-1 rounded"
                    title="View Details"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </button>
                  <button
                    @click="editInstitution(institution)"
                    class="text-gray-600 hover:text-gray-900 p-1 rounded"
                    title="Edit Institution"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button
                    @click="toggleInstitutionStatus(institution)"
                    :disabled="loading"
                    class="p-1 rounded disabled:opacity-50 disabled:cursor-not-allowed"
                    :class="institution.isActive ? 'text-red-600 hover:text-red-900' : 'text-green-600 hover:text-green-900'"
                    :title="institution.isActive ? 'Deactivate Institution' : 'Activate Institution'"
                  >
                    <svg v-if="institution.isActive" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728L5.636 5.636m12.728 12.728L18.364 5.636M5.636 18.364l12.728-12.728" />
                    </svg>
                    <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </button>
                  <button
                    @click="deleteInstitution(institution)"
                    :disabled="loading"
                    class="text-red-600 hover:text-red-900 p-1 rounded disabled:opacity-50 disabled:cursor-not-allowed"
                    title="Delete Institution"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
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
              <span class="font-medium">{{ (currentPage - 1) * itemsPerPage + 1 }}</span>
              to
              <span class="font-medium">{{ Math.min(currentPage * itemsPerPage, totalInstitutions) }}</span>
              of
              <span class="font-medium">{{ totalInstitutions }}</span>
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

    <!-- Create Institution Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-10 mx-auto p-6 border w-full max-w-3xl shadow-lg rounded-lg bg-white my-10">
        <div>
          <!-- Modal Header -->
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-2xl font-bold text-gray-900">Create New Institution</h3>
            <button
              @click="showCreateModal = false"
              class="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form @submit.prevent="createInstitution">
            <div class="space-y-6">
              <!-- Basic Information -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Basic Information
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Full Name <span class="text-red-500">*</span></label>
                    <input
                      v-model="newInstitution.name"
                      type="text"
                      required
                      placeholder="e.g., University of Cape Town"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Short Name / Abbreviation</label>
                    <input
                      v-model="newInstitution.shortName"
                      type="text"
                      placeholder="e.g., UCT"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Type <span class="text-red-500">*</span></label>
                    <select
                      v-model="newInstitution.type"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    >
                      <option value="">Select Type</option>
                      <option value="university">University</option>
                      <option value="hospital">Hospital</option>
                      <option value="laboratory">Lab</option>
                      <option value="research_center">Research</option>
                      <option value="government">Govt</option>
                      <option value="private_practice">Private</option>
                      <option value="ngo">NGO</option>
                      <option value="other">Other</option>
                    </select>
                  </div>
                </div>
              </div>

              <!-- Location Information -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  Location Information
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">City <span class="text-red-500">*</span></label>
                    <input
                      v-model="newInstitution.city"
                      type="text"
                      required
                      placeholder="e.g., Cape Town"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Province / State</label>
                    <input
                      v-model="newInstitution.province"
                      type="text"
                      placeholder="e.g., Western Cape"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Country <span class="text-red-500">*</span></label>
                    <input
                      v-model="newInstitution.country"
                      type="text"
                      required
                      placeholder="e.g., South Africa"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Postal Code</label>
                    <input
                      v-model="newInstitution.postalCode"
                      type="text"
                      placeholder="e.g., 7701"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div class="md:col-span-2">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Address</label>
                    <input
                      v-model="newInstitution.address"
                      type="text"
                      placeholder="e.g., Private Bag X3, Rondebosch"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                </div>
              </div>

              <!-- Contact Information -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                  Contact Information
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
                    <input
                      v-model="newInstitution.phone"
                      type="tel"
                      placeholder="e.g., +27 21 650 9111"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
                    <input
                      v-model="newInstitution.email"
                      type="email"
                      placeholder="e.g., info@uct.ac.za"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div class="md:col-span-2">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Website</label>
                    <input
                      v-model="newInstitution.website"
                      type="url"
                      placeholder="e.g., https://www.uct.ac.za"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Form Actions -->
            <div class="flex justify-end space-x-3 mt-8 pt-6 border-t border-gray-200">
              <button
                type="button"
                @click="showCreateModal = false"
                class="px-6 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="loading"
                class="px-6 py-2 bg-bloodsa-red text-white rounded-md hover:bg-opacity-90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
              >
                <svg v-if="loading" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ loading ? 'Creating Institution...' : 'Create Institution' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- View Institution Details Modal -->
    <div v-if="showViewModal && institutionToView" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-10 mx-auto p-6 border w-full max-w-2xl shadow-lg rounded-lg bg-white">
        <!-- Modal Header -->
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center space-x-4">
            <div class="w-16 h-16 bg-bloodsa-red rounded-full flex items-center justify-center text-white text-xl font-bold">
              {{ institutionToView.shortName || institutionToView.name.substring(0, 2).toUpperCase() }}
            </div>
            <div>
              <h3 class="text-2xl font-bold text-gray-900">
                {{ institutionToView.name }}
              </h3>
              <p class="text-sm text-gray-500 flex items-center space-x-2 mt-1">
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="getTypeBadgeClass(institutionToView.type)"
                >
                  {{ getTypeDisplayName(institutionToView.type) }}
                </span>
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="institutionToView.isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                >
                  {{ institutionToView.isActive ? 'Active' : 'Inactive' }}
                </span>
              </p>
            </div>
          </div>
          <button
            @click="closeViewModal"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Modal Content -->
        <div class="space-y-6">
          <!-- Location Information -->
          <div>
            <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
              <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              Location Information
            </h4>
            <div class="bg-gray-50 rounded-lg p-4 space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">City</span>
                <span class="text-sm text-gray-900">{{ institutionToView.city }}</span>
              </div>
              <div v-if="institutionToView.province" class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Province</span>
                <span class="text-sm text-gray-900">{{ institutionToView.province }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Country</span>
                <span class="text-sm text-gray-900">{{ institutionToView.country }}</span>
              </div>
              <div v-if="institutionToView.address" class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Address</span>
                <span class="text-sm text-gray-900 text-right">{{ institutionToView.address }}</span>
              </div>
              <div v-if="institutionToView.postalCode" class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Postal Code</span>
                <span class="text-sm text-gray-900">{{ institutionToView.postalCode }}</span>
              </div>
            </div>
          </div>

          <!-- Contact Information -->
          <div>
            <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
              <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
              Contact Information
            </h4>
            <div class="bg-gray-50 rounded-lg p-4 space-y-3">
              <div v-if="institutionToView.phone" class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Phone</span>
                <span class="text-sm text-gray-900">{{ institutionToView.phone }}</span>
              </div>
              <div v-if="institutionToView.email" class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Email</span>
                <span class="text-sm text-gray-900">{{ institutionToView.email }}</span>
              </div>
              <div v-if="institutionToView.website" class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Website</span>
                <a :href="institutionToView.website" target="_blank" class="text-sm text-bloodsa-red hover:underline">
                  {{ institutionToView.website }}
                </a>
              </div>
            </div>
          </div>

          <!-- Institution Details -->
          <div>
            <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
              <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Institution Details
            </h4>
            <div class="bg-gray-50 rounded-lg p-4 space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Institution ID</span>
                <span class="text-sm text-gray-900 font-mono">{{ institutionToView.id }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Status</span>
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="institutionToView.isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                >
                  {{ institutionToView.isActive ? 'Active' : 'Inactive' }}
                </span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Created</span>
                <span class="text-sm text-gray-900">{{ new Date(institutionToView.createdAt).toLocaleDateString() }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Last Updated</span>
                <span class="text-sm text-gray-900">{{ new Date(institutionToView.updatedAt).toLocaleDateString() }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="flex justify-end space-x-3 mt-6 pt-4 border-t border-gray-200">
          <button
            @click="closeViewModal"
            class="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors"
          >
            Close
          </button>
          <button
            @click="editInstitution(institutionToView)"
            class="px-4 py-2 bg-bloodsa-red text-white rounded-md hover:bg-opacity-90 transition-colors flex items-center space-x-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            <span>Edit Institution</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Institution Modal -->
    <div v-if="showEditModal && institutionToEdit" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-10 mx-auto p-6 border w-full max-w-3xl shadow-lg rounded-lg bg-white my-10">
        <div>
          <!-- Modal Header -->
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-2xl font-bold text-gray-900">Edit Institution</h3>
            <button
              @click="showEditModal = false; institutionToEdit = null"
              class="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form @submit.prevent="updateInstitution">
            <div class="space-y-6">
              <!-- Basic Information -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Basic Information
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Full Name <span class="text-red-500">*</span></label>
                    <input
                      v-model="institutionToEdit.name"
                      type="text"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Short Name</label>
                    <input
                      v-model="institutionToEdit.shortName"
                      type="text"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Type <span class="text-red-500">*</span></label>
                    <select
                      v-model="institutionToEdit.type"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    >
                      <option value="university">University</option>
                      <option value="hospital">Hospital</option>
                      <option value="laboratory">Lab</option>
                      <option value="research_center">Research</option>
                      <option value="government">Govt</option>
                      <option value="private_practice">Private</option>
                      <option value="ngo">NGO</option>
                      <option value="other">Other</option>
                    </select>
                  </div>
                </div>
              </div>

              <!-- Location Information -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  Location Information
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">City <span class="text-red-500">*</span></label>
                    <input
                      v-model="institutionToEdit.city"
                      type="text"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Province</label>
                    <input
                      v-model="institutionToEdit.province"
                      type="text"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Country <span class="text-red-500">*</span></label>
                    <input
                      v-model="institutionToEdit.country"
                      type="text"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Postal Code</label>
                    <input
                      v-model="institutionToEdit.postalCode"
                      type="text"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div class="md:col-span-2">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Address</label>
                    <input
                      v-model="institutionToEdit.address"
                      type="text"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                </div>
              </div>

              <!-- Contact Information -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                  Contact Information
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Phone</label>
                    <input
                      v-model="institutionToEdit.phone"
                      type="tel"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
                    <input
                      v-model="institutionToEdit.email"
                      type="email"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div class="md:col-span-2">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Website</label>
                    <input
                      v-model="institutionToEdit.website"
                      type="url"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Form Actions -->
            <div class="flex justify-end space-x-3 mt-8 pt-6 border-t border-gray-200">
              <button
                type="button"
                @click="showEditModal = false; institutionToEdit = null"
                class="px-6 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="loading"
                class="px-6 py-2 bg-bloodsa-red text-white rounded-md hover:bg-opacity-90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
              >
                <svg v-if="loading" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ loading ? 'Updating...' : 'Update Institution' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <!-- Warning Icon -->
          <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 mb-4">
            <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          
          <!-- Modal Content -->
          <div class="text-center">
            <h3 class="text-lg font-medium text-gray-900 mb-2">Delete Institution</h3>
            <div class="mt-2 px-7 py-3">
              <p class="text-sm text-gray-500 mb-4">
                Are you sure you want to delete this institution? This action cannot be undone.
              </p>
              
              <!-- Institution Info -->
              <div v-if="institutionToDelete" class="bg-gray-50 rounded-lg p-4 mb-4">
                <div class="flex items-center space-x-3">
                  <div class="w-10 h-10 bg-bloodsa-red rounded-full flex items-center justify-center text-white font-semibold text-xs">
                    {{ institutionToDelete.shortName || institutionToDelete.name.substring(0, 2).toUpperCase() }}
                  </div>
                  <div class="text-left">
                    <p class="text-sm font-medium text-gray-900">
                      {{ institutionToDelete.name }}
                    </p>
                    <p class="text-sm text-gray-500">{{ institutionToDelete.city }}, {{ institutionToDelete.province }}</p>
                    <p class="text-xs text-gray-400">{{ getTypeDisplayName(institutionToDelete.type) }}</p>
                  </div>
                </div>
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
                <span>{{ loading ? 'Deleting...' : 'Delete Institution' }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useInstitutionsStore } from '@/stores/institutions'
import { getInstitutionTypeDisplayName, type Institution } from '@/types/institution'
import { useToast } from '@/composables/useToast'

const institutionsStore = useInstitutionsStore()
const toast = useToast()

// Reactive data
const searchQuery = ref('')
const typeFilter = ref('')
const statusFilter = ref('')
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const showViewModal = ref(false)
const institutionToDelete = ref<Institution | null>(null)
const institutionToEdit = ref<Institution | null>(null)
const institutionToView = ref<Institution | null>(null)
const currentPage = ref(1)
const itemsPerPage = 10

// New institution form data
const newInstitution = ref({
  name: '',
  shortName: '',
  type: '',
  country: 'South Africa',
  province: '',
  city: '',
  address: '',
  postalCode: '',
  phone: '',
  email: '',
  website: ''
})

// Real institutions data from API
const institutions = ref<Institution[]>([])
const totalInstitutions = ref(0)
const loading = ref(false)
const error = ref<string | null>(null)

// Computed properties for server-side pagination
const totalPages = computed(() => Math.ceil(totalInstitutions.value / itemsPerPage))

// Display institutions directly (search is handled by backend)
const displayInstitutions = computed(() => institutions.value)

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, start + 4)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// Methods
const getTypeDisplayName = (type: string) => {
  return getInstitutionTypeDisplayName(type as any)
}

const getTypeBadgeClass = (type: string) => {
  const classes: Record<string, string> = {
    university: 'bg-blue-100 text-blue-800',
    hospital: 'bg-green-100 text-green-800',
    laboratory: 'bg-purple-100 text-purple-800',
    research_center: 'bg-indigo-100 text-indigo-800',
    government: 'bg-yellow-100 text-yellow-800',
    private_practice: 'bg-pink-100 text-pink-800',
    ngo: 'bg-orange-100 text-orange-800',
    other: 'bg-gray-100 text-gray-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const viewInstitution = (institution: Institution) => {
  institutionToView.value = institution
  showViewModal.value = true
}

const closeViewModal = () => {
  showViewModal.value = false
  institutionToView.value = null
}

const editInstitution = (institution: Institution) => {
  // Create a copy to avoid mutating the original
  institutionToEdit.value = { ...institution }
  showViewModal.value = false // Close view modal if open
  showEditModal.value = true
}

const updateInstitution = async () => {
  if (!institutionToEdit.value) return
  
  const institutionName = institutionToEdit.value.name
  const institutionId = institutionToEdit.value.id
  
  try {
    loading.value = true
    error.value = null
    
    const updateData = {
      name: institutionToEdit.value.name,
      shortName: institutionToEdit.value.shortName || undefined,
      type: institutionToEdit.value.type,
      country: institutionToEdit.value.country,
      province: institutionToEdit.value.province || undefined,
      city: institutionToEdit.value.city,
      address: institutionToEdit.value.address || undefined,
      postalCode: institutionToEdit.value.postalCode || undefined,
      phone: institutionToEdit.value.phone || undefined,
      email: institutionToEdit.value.email || undefined,
      website: institutionToEdit.value.website || undefined
    }
    
    await institutionsStore.updateInstitution(institutionId, updateData)
    await loadInstitutions()
    
    showEditModal.value = false
    institutionToEdit.value = null
    
    toast.success(`${institutionName} has been successfully updated`)
  } catch (err) {
    console.error('Failed to update institution:', err)
    error.value = err instanceof Error ? err.message : 'Failed to update institution'
    toast.error(`Failed to update ${institutionName}. ${err instanceof Error ? err.message : 'Please try again.'}`)
  } finally {
    loading.value = false
  }
}

const toggleInstitutionStatus = async (institution: Institution) => {
  const institutionName = institution.name
  const action = institution.isActive ? 'deactivated' : 'activated'
  
  try {
    loading.value = true
    if (institution.isActive) {
      await institutionsStore.deactivateInstitution(institution.id)
    } else {
      await institutionsStore.activateInstitution(institution.id)
    }
    // Refresh institutions list
    await loadInstitutions()
    // Show success toast
    toast.success(`${institutionName} has been ${action}`)
  } catch (err) {
    console.error('Failed to toggle institution status:', err)
    error.value = err instanceof Error ? err.message : 'Failed to update institution status'
    toast.error(`Failed to ${action} ${institutionName}. Please try again.`)
  } finally {
    loading.value = false
  }
}

const deleteInstitution = (institution: Institution) => {
  institutionToDelete.value = institution
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (!institutionToDelete.value) return
  
  const institutionName = institutionToDelete.value.name
  const institutionId = institutionToDelete.value.id
  
  try {
    loading.value = true
    await institutionsStore.deleteInstitution(institutionId)
    // Close modal
    showDeleteModal.value = false
    institutionToDelete.value = null
    // Reload list
    await loadInstitutions()
    // Show success toast
    toast.success(`${institutionName} has been successfully deleted`)
  } catch (err) {
    console.error('Failed to delete institution:', err)
    error.value = err instanceof Error ? err.message : 'Failed to delete institution'
    toast.error(`Failed to delete ${institutionName}. Please try again.`)
  } finally {
    loading.value = false
  }
}

const cancelDelete = () => {
  showDeleteModal.value = false
  institutionToDelete.value = null
}

const createInstitution = async () => {
  const institutionName = newInstitution.value.name
  
  try {
    loading.value = true
    error.value = null
    
    const institutionData = {
      name: newInstitution.value.name,
      shortName: newInstitution.value.shortName || undefined,
      type: newInstitution.value.type as any,
      country: newInstitution.value.country,
      province: newInstitution.value.province || undefined,
      city: newInstitution.value.city,
      address: newInstitution.value.address || undefined,
      postalCode: newInstitution.value.postalCode || undefined,
      phone: newInstitution.value.phone || undefined,
      email: newInstitution.value.email || undefined,
      website: newInstitution.value.website || undefined
    }
    
    await institutionsStore.createInstitution(institutionData)
    await loadInstitutions()
    
    showCreateModal.value = false
    
    // Reset form
    newInstitution.value = {
      name: '',
      shortName: '',
      type: '',
      country: 'South Africa',
      province: '',
      city: '',
      address: '',
      postalCode: '',
      phone: '',
      email: '',
      website: ''
    }
    
    toast.success(`${institutionName} has been successfully created`)
  } catch (err) {
    console.error('Failed to create institution:', err)
    error.value = err instanceof Error ? err.message : 'Failed to create institution'
    toast.error(`Failed to create ${institutionName}. ${err instanceof Error ? err.message : 'Please try again.'}`)
  } finally {
    loading.value = false
  }
}

// Reset pagination when filters change
const resetPagination = () => {
  currentPage.value = 1
}

// Watch for filter changes and reload data
watch([typeFilter, statusFilter], () => {
  resetPagination()
  loadInstitutions()
})

// Watch for search query changes with debouncing
let searchTimeout: ReturnType<typeof setTimeout> | null = null
watch(searchQuery, () => {
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  
  searchTimeout = setTimeout(() => {
    resetPagination()
    loadInstitutions()
  }, 300)
})

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadInstitutions()
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadInstitutions()
  }
}

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadInstitutions()
  }
}

const goToFirstPage = () => {
  if (currentPage.value !== 1) {
    currentPage.value = 1
    loadInstitutions()
  }
}

const goToLastPage = () => {
  if (currentPage.value !== totalPages.value) {
    currentPage.value = totalPages.value
    loadInstitutions()
  }
}

// Ensure current page is valid when total pages change
watch(totalPages, (newTotalPages) => {
  if (currentPage.value > newTotalPages && newTotalPages > 0) {
    currentPage.value = newTotalPages
  }
})

// Load institutions from API with server-side pagination
const loadInstitutions = async () => {
  try {
    loading.value = true
    error.value = null
    
    const skip = (currentPage.value - 1) * itemsPerPage
    
    await institutionsStore.fetchInstitutions({ 
      limit: itemsPerPage,
      skip: skip,
      type: typeFilter.value as any || undefined,
      isActive: statusFilter.value === 'active' ? true : statusFilter.value === 'inactive' ? false : undefined,
      search: searchQuery.value || undefined
    })
    
    institutions.value = institutionsStore.institutions || []
    totalInstitutions.value = institutionsStore.total || 0
  } catch (err) {
    console.error('Failed to load institutions:', err)
    error.value = err instanceof Error ? err.message : 'Failed to load institutions'
    institutions.value = []
    totalInstitutions.value = 0
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadInstitutions()
})
</script>

<style scoped>
/* Fade transition for loading overlay */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>


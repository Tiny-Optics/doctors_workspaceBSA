<template>
  <div class="p-6">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">User Management</h1>
        <p class="text-gray-600 mt-2">Manage users, roles, and permissions</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="bg-bloodsa-red text-white px-6 py-3 rounded-lg hover:bg-opacity-90 transition-colors font-medium flex items-center space-x-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
        </svg>
        <span>Add User</span>
      </button>
    </div>

    <!-- Filters and Search -->
    <div class="bg-white rounded-xl shadow-lg p-6 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
        <!-- Search -->
        <div class="md:col-span-2">
          <label class="block text-sm font-medium text-gray-700 mb-2">Search Users</label>
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search by name, email, or role..."
              class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
            />
            <svg class="absolute left-3 top-2.5 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
        </div>

        <!-- Institution Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Filter by Institution</label>
          <select
            v-model="institutionFilter"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
          >
            <option value="">All Institutions</option>
            <option 
              v-for="institution in institutionsStore.institutions" 
              :key="institution.id" 
              :value="institution.id"
              :title="institution.name"
            >
              {{ institution.shortName || institution.name.substring(0, 20) }}
            </option>
          </select>
        </div>

        <!-- Role Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Filter by Role</label>
          <select
            v-model="roleFilter"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
          >
            <option value="">All Roles</option>
            <option value="haematologist">Haematologist</option>
            <option value="physician">Physician</option>
            <option value="data_capturer">Data Capturer</option>
            <option value="admin">Admin</option>
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

    <!-- Users Table -->
    <div class="bg-white rounded-xl shadow-lg overflow-hidden">
      <!-- Users Table Content with Loading Overlay -->
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
                User
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Role
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Institution
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Last Login
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <!-- Empty State -->
              <tr v-if="!loading && displayUsers.length === 0">
                <td colspan="6" class="px-6 py-12 text-center">
                  <div class="flex flex-col items-center justify-center space-y-3">
                    <svg class="w-16 h-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
                    </svg>
                    <div class="text-gray-500">
                      <p class="text-lg font-medium">No users found</p>
                      <p class="text-sm mt-1">
                        <span v-if="searchQuery || roleFilter || statusFilter">
                          Try adjusting your search or filters
                        </span>
                        <span v-else>
                          No users exist in the system yet
                        </span>
                      </p>
                    </div>
                  </div>
                </td>
              </tr>
              
              <!-- User Rows -->
              <tr v-for="user in displayUsers" :key="user.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="w-10 h-10 bg-bloodsa-red rounded-full flex items-center justify-center text-white font-semibold">
                    {{ getUserInitials(user) }}
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">
                      {{ user.profile.firstName }} {{ user.profile.lastName }}
                    </div>
                    <div class="text-sm text-gray-500">{{ user.email }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="getRoleBadgeClass(user.role)"
                >
                  {{ getRoleDisplayName(user.role) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ getInstitutionDisplay(user.profile.institutionId) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="user.isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                >
                  {{ user.isActive ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatLastLogin(user.lastLoginAt) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div class="flex items-center justify-end space-x-2">
                  <button
                    @click="viewUser(user)"
                    class="text-blue-600 hover:text-blue-900 p-1 rounded"
                    title="View Details"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </button>
                  <button
                    @click="editUser(user)"
                    class="text-gray-600 hover:text-gray-900 p-1 rounded"
                    title="Edit User"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button
                    @click="toggleUserStatus(user)"
                    :disabled="loading"
                    class="p-1 rounded disabled:opacity-50 disabled:cursor-not-allowed"
                    :class="user.isActive ? 'text-red-600 hover:text-red-900' : 'text-green-600 hover:text-green-900'"
                    :title="user.isActive ? 'Deactivate User' : 'Activate User'"
                  >
                    <svg v-if="user.isActive" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728L5.636 5.636m12.728 12.728L18.364 5.636M5.636 18.364l12.728-12.728" />
                    </svg>
                    <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </button>
                  <button
                    @click="deleteUser(user)"
                    :disabled="loading"
                    class="text-red-600 hover:text-red-900 p-1 rounded disabled:opacity-50 disabled:cursor-not-allowed"
                    title="Delete User"
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
              <span class="font-medium">{{ Math.min(currentPage * itemsPerPage, totalUsers) }}</span>
              of
              <span class="font-medium">{{ totalUsers }}</span>
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

    <!-- Create User Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="showCreateModal = false">
      <div class="relative top-10 mx-auto p-6 border w-full max-w-3xl shadow-lg rounded-lg bg-white my-10" @click.stop>
        <div>
          <!-- Modal Header -->
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-2xl font-bold text-gray-900">Create New User</h3>
            <button
              @click="showCreateModal = false"
              class="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form @submit.prevent="createUser">
            <div class="space-y-6">
              <!-- Account Information -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  Account Information
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Username <span class="text-red-500">*</span></label>
                    <input
                      v-model="newUser.username"
                      type="text"
                      required
                      placeholder="e.g., john.doe"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                    <p class="text-xs text-gray-500 mt-1">3-50 characters, alphanumeric and underscores only</p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Email <span class="text-red-500">*</span></label>
                    <input
                      v-model="newUser.email"
                      type="email"
                      required
                      placeholder="e.g., john.doe@hospital.com"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Password <span class="text-red-500">*</span></label>
                    <input
                      v-model="newUser.password"
                      type="password"
                      required
                      placeholder="Enter secure password"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                    <p class="text-xs text-gray-500 mt-1">Min 8 chars with uppercase, lowercase, number & special char</p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Confirm Password <span class="text-red-500">*</span></label>
                    <input
                      v-model="newUser.confirmPassword"
                      type="password"
                      required
                      placeholder="Re-enter password"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                </div>
              </div>

              <!-- Personal Information -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Personal Information
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">First Name <span class="text-red-500">*</span></label>
                    <input
                      v-model="newUser.firstName"
                      type="text"
                      required
                      placeholder="e.g., John"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Last Name <span class="text-red-500">*</span></label>
                    <input
                      v-model="newUser.lastName"
                      type="text"
                      required
                      placeholder="e.g., Doe"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
                    <input
                      v-model="newUser.phoneNumber"
                      type="tel"
                      placeholder="e.g., +27 11 123 4567"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                </div>
              </div>

              <!-- Role & Permissions -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                  </svg>
                  Role & Permissions
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Role <span class="text-red-500">*</span></label>
                    <select
                      v-model="newUser.role"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    >
                      <option value="">Select Role</option>
                      <option value="haematologist">Haematologist</option>
                      <option value="physician">Physician</option>
                      <option value="data_capturer">Data Capturer</option>
                      <option value="admin">Admin</option>
                    </select>
                  </div>
                  <div v-if="newUser.role === 'admin'">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Admin Level <span class="text-red-500">*</span></label>
                    <select
                      v-model="newUser.adminLevel"
                      :required="newUser.role === 'admin'"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    >
                      <option value="">Select Admin Level</option>
                      <option value="user_manager">User Manager</option>
                      <option value="super_admin">Super Admin</option>
                    </select>
                  </div>
                </div>
              </div>

              <!-- Professional Information -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                  Professional Information
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div class="md:col-span-2">
                    <label class="block text-sm font-medium text-gray-700 mb-1">Institution <span class="text-red-500">*</span></label>
                    <select
                      v-model="newUser.institutionId"
                      required
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    >
                      <option value="">Select Institution</option>
                      <option 
                        v-for="institution in institutionsStore.institutions" 
                        :key="institution.id" 
                        :value="institution.id"
                        :title="`${institution.name} - ${institution.city}, ${institution.province || ''}`"
                      >
                        {{ institution.shortName || institution.name.substring(0, 30) }} - {{ institution.city }}
                      </option>
                    </select>
                    <p class="text-xs text-gray-500 mt-1">Select the user's affiliated institution (hover for full name)</p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Specialty</label>
                    <input
                      v-model="newUser.specialty"
                      type="text"
                      placeholder="e.g., Haematology, Internal Medicine"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Registration Number</label>
                    <input
                      v-model="newUser.registrationNumber"
                      type="text"
                      placeholder="e.g., HPCSA123456"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    />
                    <p class="text-xs text-gray-500 mt-1">HPCSA or other professional registration number</p>
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
                <span>{{ loading ? 'Creating User...' : 'Create User' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="cancelDelete">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <!-- Warning Icon -->
          <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 mb-4">
            <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          
          <!-- Modal Content -->
          <div class="text-center">
            <h3 class="text-lg font-medium text-gray-900 mb-2">Delete User</h3>
            <div class="mt-2 px-7 py-3">
              <p class="text-sm text-gray-500 mb-4">
                Are you sure you want to delete this user? This action cannot be undone.
              </p>
              
              <!-- User Info -->
              <div v-if="userToDelete" class="bg-gray-50 rounded-lg p-4 mb-4">
                <div class="flex items-center space-x-3">
                  <div class="w-10 h-10 bg-bloodsa-red rounded-full flex items-center justify-center text-white font-semibold">
                    {{ getUserInitials(userToDelete) }}
                  </div>
                  <div class="text-left">
                    <p class="text-sm font-medium text-gray-900">
                      {{ userToDelete.profile.firstName }} {{ userToDelete.profile.lastName }}
                    </p>
                    <p class="text-sm text-gray-500">{{ userToDelete.email }}</p>
                    <p class="text-xs text-gray-400">{{ getRoleDisplayName(userToDelete.role) }} • {{ getInstitutionName(userToDelete.profile.institutionId) }}</p>
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
                <span>{{ loading ? 'Deleting...' : 'Delete User' }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- View User Details Modal -->
    <div v-if="showViewModal && userToView" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click="closeViewModal">
      <div class="relative top-10 mx-auto p-6 border w-full max-w-2xl shadow-lg rounded-lg bg-white" @click.stop>
        <!-- Modal Header -->
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center space-x-4">
            <div class="w-16 h-16 bg-bloodsa-red rounded-full flex items-center justify-center text-white text-2xl font-bold">
              {{ getUserInitials(userToView) }}
            </div>
            <div>
              <h3 class="text-2xl font-bold text-gray-900">
                {{ userToView.profile.firstName }} {{ userToView.profile.lastName }}
              </h3>
              <p class="text-sm text-gray-500 flex items-center space-x-2 mt-1">
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="getRoleBadgeClass(userToView.role)"
                >
                  {{ getRoleDisplayName(userToView.role) }}
                </span>
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="userToView.isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                >
                  {{ userToView.isActive ? 'Active' : 'Inactive' }}
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
          <!-- Contact Information -->
          <div>
            <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
              <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
              Contact Information
            </h4>
            <div class="bg-gray-50 rounded-lg p-4 space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Email</span>
                <span class="text-sm text-gray-900">{{ userToView.email }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Username</span>
                <span class="text-sm text-gray-900">{{ userToView.username }}</span>
              </div>
              <div v-if="userToView.profile.phoneNumber" class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Phone</span>
                <span class="text-sm text-gray-900">{{ userToView.profile.phoneNumber }}</span>
              </div>
            </div>
          </div>

          <!-- Professional Information -->
          <div>
            <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
              <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
              Professional Information
            </h4>
            <div class="bg-gray-50 rounded-lg p-4 space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Institution</span>
                <span class="text-sm text-gray-900">{{ getInstitutionName(userToView.profile.institutionId) }}</span>
              </div>
              <div v-if="userToView.profile.institutionId" class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Location</span>
                <span class="text-sm text-gray-900">{{ getInstitutionLocation(userToView.profile.institutionId) }}</span>
              </div>
              <div v-if="userToView.profile.specialty" class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Specialty</span>
                <span class="text-sm text-gray-900">{{ userToView.profile.specialty }}</span>
              </div>
              <div v-if="userToView.profile.registrationNumber" class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Registration Number</span>
                <span class="text-sm text-gray-900">{{ userToView.profile.registrationNumber }}</span>
              </div>
            </div>
          </div>

          <!-- Account Information -->
          <div>
            <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
              <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Account Information
            </h4>
            <div class="bg-gray-50 rounded-lg p-4 space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">User ID</span>
                <span class="text-sm text-gray-900 font-mono">{{ userToView.id }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Account Status</span>
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="userToView.isActive ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                >
                  {{ userToView.isActive ? 'Active' : 'Inactive' }}
                </span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Last Login</span>
                <span class="text-sm text-gray-900">{{ formatLastLogin(userToView.lastLoginAt) }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Account Created</span>
                <span class="text-sm text-gray-900">{{ new Date(userToView.createdAt).toLocaleDateString() }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Last Updated</span>
                <span class="text-sm text-gray-900">{{ new Date(userToView.updatedAt).toLocaleDateString() }}</span>
              </div>
            </div>
          </div>

          <!-- Admin Information (if applicable) -->
          <div v-if="userToView.role === 'admin' && userToView.adminLevel">
            <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
              <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              Admin Permissions
            </h4>
            <div class="bg-orange-50 rounded-lg p-4">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-600">Admin Level</span>
                <span class="px-3 py-1 text-sm font-medium rounded-full bg-orange-100 text-orange-800">
                  {{ userToView.adminLevel.replace('_', ' ').toUpperCase() }}
                </span>
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
            @click="editUser(userToView)"
            class="px-4 py-2 bg-bloodsa-red text-white rounded-md hover:bg-opacity-90 transition-colors flex items-center space-x-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            <span>Edit User</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useUsersStore } from '@/stores/users'
import { useInstitutionsStore } from '@/stores/institutions'
import { getUserRoleDisplayName, type User } from '@/types/user'
import type { Institution } from '@/types/institution'
import { useToast } from '@/composables/useToast'

const usersStore = useUsersStore()
const institutionsStore = useInstitutionsStore()
const toast = useToast()

// Reactive data
const searchQuery = ref('')
const roleFilter = ref('')
const statusFilter = ref('')
const institutionFilter = ref('')
const showCreateModal = ref(false)
const showDeleteModal = ref(false)
const showViewModal = ref(false)
const userToDelete = ref<User | null>(null)
const userToView = ref<User | null>(null)
const currentPage = ref(1)
const itemsPerPage = 10

// New user form data
const newUser = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  firstName: '',
  lastName: '',
  role: '',
  adminLevel: '',
  institutionId: '',
  specialty: '',
  registrationNumber: '',
  phoneNumber: ''
})

// Real users data from API
const users = ref<User[]>([])
const totalUsers = ref(0)
const loading = ref(false)
const error = ref<string | null>(null)

// Computed properties for server-side pagination
const totalPages = computed(() => Math.ceil(totalUsers.value / itemsPerPage))

// For server-side pagination, we display the users with client-side institution filtering
// (Backend doesn't support institution filtering yet)
const displayUsers = computed(() => {
  let filtered = users.value

  // Apply client-side institution filter
  if (institutionFilter.value) {
    filtered = filtered.filter(user => user.profile.institutionId === institutionFilter.value)
  }

  return filtered
})

// Helper to get institution name by ID
const getInstitutionName = (institutionId?: string): string => {
  if (!institutionId) return 'N/A'
  const institution = institutionsStore.institutions.find(i => i.id === institutionId)
  return institution ? `${institution.name}` : 'Unknown Institution'
}

// Helper to get institution display with location
const getInstitutionDisplay = (institutionId?: string): string => {
  if (!institutionId) return 'N/A'
  const institution = institutionsStore.institutions.find(i => i.id === institutionId)
  if (!institution) return 'Unknown Institution'
  return `${institution.shortName || institution.name} - ${institution.city}`
}

// Helper to get institution location
const getInstitutionLocation = (institutionId?: string): string => {
  if (!institutionId) return 'N/A'
  const institution = institutionsStore.institutions.find(i => i.id === institutionId)
  if (!institution) return 'Unknown'
  let location = institution.city
  if (institution.province) location += ', ' + institution.province
  if (institution.country) location += ', ' + institution.country
  return location
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

// Methods
const getUserInitials = (user: any) => {
  const firstName = user.profile.firstName || ''
  const lastName = user.profile.lastName || ''
  return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase()
}

// Reset pagination when filters change
const resetPagination = () => {
  currentPage.value = 1
}

// Watch for filter changes and reload data
watch([roleFilter, statusFilter], () => {
  resetPagination()
  loadUsers() // Reload with new filters
})

// Watch for institution filter (client-side only for now)
watch(institutionFilter, () => {
  // No need to reload from server, handled client-side
  // Institution filtering is client-side since backend doesn't support it yet
})

// Watch for search query changes with debouncing (backend search)
let searchTimeout: ReturnType<typeof setTimeout> | null = null
watch(searchQuery, () => {
  // Clear previous timeout
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  
  // Debounce search - wait 300ms after user stops typing
  searchTimeout = setTimeout(() => {
    resetPagination()
    loadUsers() // Reload with search query from backend
  }, 300)
})

const getRoleDisplayName = (role: string) => {
  return getUserRoleDisplayName(role as any)
}

const getRoleBadgeClass = (role: string) => {
  const classes = {
    haematologist: 'bg-red-100 text-red-800',
    physician: 'bg-green-100 text-green-800',
    data_capturer: 'bg-purple-100 text-purple-800',
    admin: 'bg-orange-100 text-orange-800'
  }
  return classes[role as keyof typeof classes] || 'bg-gray-100 text-gray-800'
}

const formatLastLogin = (lastLogin?: string) => {
  if (!lastLogin) return 'Never'
  const date = new Date(lastLogin)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  
  if (diffMins < 60) return `${diffMins}m ago`
  const diffHours = Math.floor(diffMins / 60)
  if (diffHours < 24) return `${diffHours}h ago`
  const diffDays = Math.floor(diffHours / 24)
  return `${diffDays}d ago`
}

const viewUser = (user: any) => {
  userToView.value = user
  showViewModal.value = true
}

const closeViewModal = () => {
  showViewModal.value = false
  userToView.value = null
}

const editUser = (user: any) => {
  console.log('Edit user:', user)
  // TODO: Implement edit user modal
}

const toggleUserStatus = async (user: any) => {
  const userName = `${user.profile.firstName} ${user.profile.lastName}`
  const action = user.isActive ? 'deactivated' : 'activated'
  
  try {
    loading.value = true
    if (user.isActive) {
      await usersStore.deactivateUser(user.id)
    } else {
      await usersStore.activateUser(user.id)
    }
    // Update local state
    user.isActive = !user.isActive
    // Refresh users list to get updated data
    await loadUsers()
    // Show success toast
    toast.success(`${userName} has been ${action}`)
  } catch (err) {
    console.error('Failed to toggle user status:', err)
    error.value = err instanceof Error ? err.message : 'Failed to update user status'
    toast.error(`Failed to ${user.isActive ? 'deactivate' : 'activate'} ${userName}. Please try again.`)
  } finally {
    loading.value = false
  }
}

const deleteUser = (user: any) => {
  userToDelete.value = user
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (!userToDelete.value) return
  
  const userName = `${userToDelete.value.profile.firstName} ${userToDelete.value.profile.lastName}`
  const userId = userToDelete.value.id
  
  try {
    loading.value = true
    await usersStore.deleteUser(userId)
    // Remove from local state
    const index = users.value.findIndex(u => u.id === userId)
    if (index > -1) {
      users.value.splice(index, 1)
    }
    // Close modal
    showDeleteModal.value = false
    userToDelete.value = null
    // Show success toast
    toast.success(`${userName} has been successfully deleted`)
  } catch (err) {
    console.error('Failed to delete user:', err)
    error.value = err instanceof Error ? err.message : 'Failed to delete user'
    toast.error(`Failed to delete ${userName}. Please try again.`)
  } finally {
    loading.value = false
  }
}

const cancelDelete = () => {
  showDeleteModal.value = false
  userToDelete.value = null
}

const createUser = async () => {
  const userName = `${newUser.value.firstName} ${newUser.value.lastName}`
  
  // Validate password match
  if (newUser.value.password !== newUser.value.confirmPassword) {
    error.value = 'Passwords do not match'
    toast.error('Passwords do not match. Please check and try again.')
    return
  }
  
  // Validate admin level for admin role
  if (newUser.value.role === 'admin' && !newUser.value.adminLevel) {
    error.value = 'Admin level is required for admin role'
    toast.error('Please select an admin level for admin users.')
    return
  }
  
  try {
    loading.value = true
    error.value = null
    
    const userData = {
      username: newUser.value.username,
      email: newUser.value.email,
      password: newUser.value.password,
      role: newUser.value.role as any,
      adminLevel: newUser.value.role === 'admin' ? newUser.value.adminLevel as any : undefined,
      firstName: newUser.value.firstName,
      lastName: newUser.value.lastName,
      institutionId: newUser.value.institutionId,
      specialty: newUser.value.specialty || undefined,
      registrationNumber: newUser.value.registrationNumber || undefined,
      phoneNumber: newUser.value.phoneNumber || undefined
    }
    
    const createdUser = await usersStore.createUser(userData)
    // Refresh users list to include the new user
    await loadUsers()
    
    showCreateModal.value = false
    
    // Reset form
    newUser.value = {
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      firstName: '',
      lastName: '',
      role: '',
      adminLevel: '',
      institutionId: '',
      specialty: '',
      registrationNumber: '',
      phoneNumber: ''
    }
    
    // Show success toast
    toast.success(`${userName} has been successfully created with username: ${userData.username}`)
  } catch (err) {
    console.error('Failed to create user:', err)
    error.value = err instanceof Error ? err.message : 'Failed to create user'
    toast.error(`Failed to create ${userName}. ${err instanceof Error ? err.message : 'Please try again.'}`)
  } finally {
    loading.value = false
  }
}

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadUsers() // Load previous page from server
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadUsers() // Load next page from server
  }
}

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadUsers() // Load new page from server
  }
}

const goToFirstPage = () => {
  if (currentPage.value !== 1) {
    currentPage.value = 1
    loadUsers() // Load first page from server
  }
}

const goToLastPage = () => {
  if (currentPage.value !== totalPages.value) {
    currentPage.value = totalPages.value
    loadUsers() // Load last page from server
  }
}

// Ensure current page is valid when total pages change
watch(totalPages, (newTotalPages) => {
  if (currentPage.value > newTotalPages && newTotalPages > 0) {
    currentPage.value = newTotalPages
  }
})

// Load users from API with server-side pagination
const loadUsers = async () => {
  try {
    loading.value = true
    error.value = null
    console.log('Loading users from API...', { page: currentPage.value, limit: itemsPerPage })
    
    // Calculate skip value for server-side pagination
    const skip = (currentPage.value - 1) * itemsPerPage
    
    // Fetch users with pagination and search parameters
    await usersStore.fetchUsers({ 
      limit: itemsPerPage,
      skip: skip,
      role: roleFilter.value as any || undefined,
      isActive: statusFilter.value === 'active' ? true : statusFilter.value === 'inactive' ? false : undefined,
      search: searchQuery.value || undefined
    })
    
    console.log('Users loaded from store:', usersStore.users)
    console.log('Total users:', usersStore.total)
    
    // Get users and total from store state after fetch
    users.value = usersStore.users || []
    totalUsers.value = usersStore.total || 0
    
    console.log('Users set in component:', users.value)
    console.log('Total users set:', totalUsers.value)
  } catch (err) {
    console.error('Failed to load users:', err)
    error.value = err instanceof Error ? err.message : 'Failed to load users'
    // Set empty array on error to prevent undefined issues
    users.value = []
    totalUsers.value = 0
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  // Load institutions first for dropdowns
  await institutionsStore.fetchInstitutions({ isActive: true, limit: 1000 })
  // Then load users
  loadUsers()
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

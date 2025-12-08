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
          <p class="text-lg text-gray-900">{{ user?.role ? getRoleDisplayName(user.role as any) : '' }}</p>
        </div>
        <div class="bg-gray-50 rounded-lg p-4">
          <label class="block text-sm font-medium text-gray-600 mb-1">Institution</label>
          <div class="flex items-center justify-between">
            <p class="text-lg text-gray-900">{{ institutionName }}</p>
            <div class="flex items-center gap-2">
              <button
                v-if="canEditInstitution"
                @click="openEditInstitutionModal"
                class="px-3 py-1.5 text-sm font-medium text-blue-600 border border-blue-600 rounded-lg hover:bg-blue-600 hover:text-white transition-colors"
              >
                Edit
              </button>
              <button
                @click="showInstitutionModal = true"
                class="px-3 py-1.5 text-sm font-medium text-bloodsa-red border border-bloodsa-red rounded-lg hover:bg-bloodsa-red hover:text-white transition-colors"
              >
                Change
              </button>
            </div>
          </div>
          <p v-if="canEditInstitution" class="text-xs text-gray-500 mt-2">
            ℹ️ You can edit this institution because you created it. Only the creator can edit an institution.
          </p>
        </div>
      </div>
      <p class="text-sm text-gray-500 mt-3">
        ℹ️ Contact an administrator to change your username, email, or role.
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

    <!-- Change Institution Modal -->
    <div v-if="showInstitutionModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-10 mx-auto p-6 border w-full max-w-3xl shadow-lg rounded-lg bg-white my-10" @click.stop>
        <div>
          <!-- Modal Header -->
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-2xl font-bold text-gray-900">Change Institution</h3>
            <button
              @click="closeInstitutionModal"
              class="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Select or Add New -->
          <div v-if="!showCreateForm" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Select Institution <span class="text-red-500">*</span>
              </label>
              <select
                v-model="selectedInstitutionId"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-bloodsa-red focus:border-bloodsa-red transition-colors"
                :disabled="changingInstitution"
              >
                <option value="">{{ institutionsStore.isLoading ? 'Loading institutions...' : 'Select an institution' }}</option>
                <option v-for="institution in availableInstitutions" :key="institution.id" :value="institution.id">
                  {{ institution.name }} {{ !institution.isActive ? '(Pending Approval)' : '' }}
                </option>
              </select>
            </div>

            <div class="flex items-center gap-4">
              <div class="flex-1 border-t border-gray-300"></div>
              <span class="text-sm text-gray-500">OR</span>
              <div class="flex-1 border-t border-gray-300"></div>
            </div>

            <button
              @click="showCreateForm = true"
              type="button"
              class="w-full px-4 py-3 border-2 border-dashed border-gray-300 rounded-lg hover:border-bloodsa-red hover:bg-red-50 transition-colors flex items-center justify-center gap-2"
              :disabled="changingInstitution"
            >
              <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              <span class="text-sm font-medium text-gray-700">Add New Institution</span>
            </button>

            <div class="bg-blue-50 border-l-4 border-blue-400 p-4 rounded">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                  </svg>
                </div>
                <div class="ml-3">
                  <p class="text-sm text-blue-700">
                    <strong>Note:</strong> If you create a new institution, it will be created as active and you will be automatically switched to it.
                  </p>
                </div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex justify-end gap-3 pt-4 border-t border-gray-200">
              <button
                @click="closeInstitutionModal"
                type="button"
                :disabled="changingInstitution"
                class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
              >
                Cancel
              </button>
              <button
                @click="changeInstitution"
                type="button"
                :disabled="!selectedInstitutionId || changingInstitution"
                class="px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
              >
                <svg v-if="changingInstitution" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ changingInstitution ? 'Changing...' : 'Change Institution' }}
              </button>
            </div>
          </div>

          <!-- Create New Institution Form -->
          <div v-else class="space-y-6">
            <div class="bg-green-50 border-l-4 border-green-400 p-4 rounded">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                  </svg>
                </div>
                <div class="ml-3">
                  <p class="text-sm text-green-700">
                    <strong>Note:</strong> Your new institution will be created as active and you will be automatically switched to it once created.
                  </p>
                </div>
              </div>
            </div>

            <form @submit.prevent="createAndSwitchInstitution" class="space-y-4">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Name -->
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    Institution Name <span class="text-red-500">*</span>
                  </label>
                  <input
                    v-model="newInstitution.name"
                    type="text"
                    required
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    placeholder="e.g., University of Cape Town"
                  />
                </div>

                <!-- Short Name -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Short Name / Abbreviation</label>
                  <input
                    v-model="newInstitution.shortName"
                    type="text"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    placeholder="e.g., UCT"
                  />
                </div>

                <!-- Type -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    Type <span class="text-red-500">*</span>
                  </label>
                  <select
                    v-model="newInstitution.type"
                    required
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  >
                    <option value="">Select type</option>
                    <option value="university">University</option>
                    <option value="hospital">Hospital</option>
                    <option value="laboratory">Laboratory</option>
                    <option value="research_center">Research Center</option>
                    <option value="government">Government</option>
                    <option value="private_practice">Private Practice</option>
                    <option value="ngo">NGO</option>
                    <option value="other">Other</option>
                  </select>
                </div>

                <!-- Country -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    Country <span class="text-red-500">*</span>
                  </label>
                  <input
                    v-model="newInstitution.country"
                    type="text"
                    required
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    placeholder="e.g., South Africa"
                  />
                </div>

                <!-- Province -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Province</label>
                  <input
                    v-model="newInstitution.province"
                    type="text"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    placeholder="e.g., Western Cape"
                  />
                </div>

                <!-- City -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    City <span class="text-red-500">*</span>
                  </label>
                  <input
                    v-model="newInstitution.city"
                    type="text"
                    required
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    placeholder="e.g., Cape Town"
                  />
                </div>

                <!-- Address -->
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium text-gray-700 mb-1">Address</label>
                  <input
                    v-model="newInstitution.address"
                    type="text"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    placeholder="Street address"
                  />
                </div>

                <!-- Postal Code -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Postal Code</label>
                  <input
                    v-model="newInstitution.postalCode"
                    type="text"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    placeholder="e.g., 7700"
                  />
                </div>

                <!-- Phone -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Phone</label>
                  <input
                    v-model="newInstitution.phone"
                    type="tel"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    placeholder="+27 21 123 4567"
                  />
                </div>

                <!-- Email -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
                  <input
                    v-model="newInstitution.email"
                    type="email"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    placeholder="info@institution.co.za"
                  />
                </div>

                <!-- Website -->
                <div class="md:col-span-2">
                  <label class="block text-sm font-medium text-gray-700 mb-1">Website</label>
                  <input
                    v-model="newInstitution.website"
                    type="url"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                    placeholder="https://www.institution.co.za"
                  />
                </div>
              </div>

              <!-- Logo Upload -->
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                  <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  Institution Logo (Optional)
                </h4>
                <div class="flex items-center gap-4">
                  <div 
                    v-if="institutionImagePreview"
                    class="relative"
                  >
                    <img
                      :src="institutionImagePreview"
                      alt="Preview"
                      class="h-24 w-24 rounded object-cover border-2 border-gray-200"
                    />
                    <button
                      @click="clearInstitutionImage"
                      type="button"
                      class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full p-1 hover:bg-red-600 transition-colors"
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
                      <p class="text-sm text-gray-600">Click to upload logo</p>
                      <p class="text-xs text-gray-500 mt-1">JPG, PNG, WEBP (max 5MB)</p>
                    </div>
                    <input
                      type="file"
                      @change="handleInstitutionImageUpload"
                      accept="image/jpeg,image/jpg,image/png,image/webp"
                      class="hidden"
                      :disabled="creatingInstitution"
                    />
                  </label>
                </div>
              </div>

              <!-- Action Buttons -->
              <div class="flex justify-end gap-3 pt-4 border-t border-gray-200">
                <button
                  @click="showCreateForm = false"
                  type="button"
                  :disabled="creatingInstitution"
                  class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
                >
                  Back
                </button>
                <button
                  type="submit"
                  :disabled="creatingInstitution"
                  class="px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  <svg v-if="creatingInstitution" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ creatingInstitution ? 'Creating...' : 'Create & Switch' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Institution Modal -->
    <div v-if="showEditInstitutionModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-10 mx-auto p-6 border w-full max-w-3xl shadow-lg rounded-lg bg-white my-10" @click.stop>
        <div>
          <!-- Modal Header -->
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-2xl font-bold text-gray-900">Edit Institution</h3>
            <button
              @click="closeEditInstitutionModal"
              class="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Info Note -->
          <div class="bg-blue-50 border-l-4 border-blue-400 p-4 rounded mb-6">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm text-blue-700">
                  <strong>Note:</strong> Only the creator of an institution can edit it. You can edit this institution because you created it.
                </p>
              </div>
            </div>
          </div>

          <form @submit.prevent="updateInstitution" class="space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <!-- Name -->
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Institution Name <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="editInstitution.name"
                  type="text"
                  required
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  placeholder="e.g., University of Cape Town"
                />
              </div>

              <!-- Short Name -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Short Name / Abbreviation</label>
                <input
                  v-model="editInstitution.shortName"
                  type="text"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  placeholder="e.g., UCT"
                />
              </div>

              <!-- Type -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Type <span class="text-red-500">*</span>
                </label>
                <select
                  v-model="editInstitution.type"
                  required
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                >
                  <option value="">Select type</option>
                  <option value="university">University</option>
                  <option value="hospital">Hospital</option>
                  <option value="laboratory">Laboratory</option>
                  <option value="research_center">Research Center</option>
                  <option value="government">Government</option>
                  <option value="private_practice">Private Practice</option>
                  <option value="ngo">NGO</option>
                  <option value="other">Other</option>
                </select>
              </div>

              <!-- Country -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Country <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="editInstitution.country"
                  type="text"
                  required
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  placeholder="e.g., South Africa"
                />
              </div>

              <!-- Province -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Province</label>
                <input
                  v-model="editInstitution.province"
                  type="text"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  placeholder="e.g., Western Cape"
                />
              </div>

              <!-- City -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  City <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="editInstitution.city"
                  type="text"
                  required
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  placeholder="e.g., Cape Town"
                />
              </div>

              <!-- Address -->
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">Address</label>
                <input
                  v-model="editInstitution.address"
                  type="text"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  placeholder="Street address"
                />
              </div>

              <!-- Postal Code -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Postal Code</label>
                <input
                  v-model="editInstitution.postalCode"
                  type="text"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  placeholder="e.g., 7700"
                />
              </div>

              <!-- Phone -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Phone</label>
                <input
                  v-model="editInstitution.phone"
                  type="tel"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  placeholder="+27 21 123 4567"
                />
              </div>

              <!-- Email -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
                <input
                  v-model="editInstitution.email"
                  type="email"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  placeholder="info@institution.co.za"
                />
              </div>

              <!-- Website -->
              <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">Website</label>
                <input
                  v-model="editInstitution.website"
                  type="url"
                  class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
                  placeholder="https://www.institution.co.za"
                />
              </div>
            </div>

            <!-- Logo Upload -->
            <div>
              <h4 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                <svg class="w-5 h-5 mr-2 text-bloodsa-red" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                Institution Logo (Optional)
              </h4>
              <div class="flex items-center gap-4">
                <div 
                  v-if="editInstitutionImagePreview"
                  class="relative"
                >
                  <img
                    :src="editInstitutionImagePreview"
                    alt="Preview"
                    class="h-24 w-24 rounded object-cover border-2 border-gray-200"
                  />
                  <button
                    @click="clearEditInstitutionImage"
                    type="button"
                    class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full p-1 hover:bg-red-600 transition-colors"
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
                    <p class="text-sm text-gray-600">Click to upload logo</p>
                    <p class="text-xs text-gray-500 mt-1">JPG, PNG, WEBP (max 5MB)</p>
                  </div>
                  <input
                    type="file"
                    @change="handleEditInstitutionImageUpload"
                    accept="image/jpeg,image/jpg,image/png,image/webp"
                    class="hidden"
                    :disabled="editingInstitution"
                  />
                </label>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex justify-end gap-3 pt-4 border-t border-gray-200">
              <button
                @click="closeEditInstitutionModal"
                type="button"
                :disabled="editingInstitution"
                class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="editingInstitution"
                class="px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
              >
                <svg v-if="editingInstitution" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ editingInstitution ? 'Updating...' : 'Update Institution' }}
              </button>
            </div>
          </form>
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

// Institution modal state
const showInstitutionModal = ref(false)
const showCreateForm = ref(false)
const selectedInstitutionId = ref('')
const changingInstitution = ref(false)
const creatingInstitution = ref(false)

// Edit institution modal state
const showEditInstitutionModal = ref(false)
const editingInstitution = ref(false)
const editInstitutionImageFile = ref<File | null>(null)
const editInstitutionImagePreview = ref<string | null>(null)

// Institution image upload state
const institutionImageFile = ref<File | null>(null)
const institutionImagePreview = ref<string | null>(null)

// New institution form data
const newInstitution = ref({
  name: '',
  shortName: '',
  type: '' as any,
  country: 'South Africa',
  province: '',
  city: '',
  address: '',
  postalCode: '',
  phone: '',
  email: '',
  website: '',
  imagePath: ''
})

// Edit institution form data
const editInstitution = ref({
  name: '',
  shortName: '',
  type: '' as any,
  country: '',
  province: '',
  city: '',
  address: '',
  postalCode: '',
  phone: '',
  email: '',
  website: '',
  imagePath: ''
})

// Load institutions on mount - fetch all (we'll filter in computed)
onMounted(async () => {
  await institutionsStore.fetchInstitutions({ limit: 1000 })
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
  const institutionId = user.value.profile.institutionId
  if (!institutionId) return 'Not Set'
  const institution = institutionsStore.institutions.find(i => i.id === institutionId)
  return institution ? institution.name : 'Unknown Institution'
})

// Available institutions: active ones + ones created by current user (even if inactive)
const availableInstitutions = computed(() => {
  if (!user.value) return []
  const userId = user.value.id
  return institutionsStore.institutions.filter(inst => 
    inst.isActive || (inst.createdBy === userId)
  )
})

// Check if user can edit their current institution (if they created it)
const canEditInstitution = computed(() => {
  if (!user.value || !user.value.profile.institutionId) return false
  const userId = user.value.id
  const institutionId = user.value.profile.institutionId
  if (!institutionId) return false
  const institution = institutionsStore.institutions.find(i => i.id === institutionId)
  return institution?.createdBy === userId
})

// Get current user's institution
const currentUserInstitution = computed(() => {
  if (!user.value || !user.value.profile.institutionId) return null
  const institutionId = user.value.profile.institutionId
  if (!institutionId) return null
  return institutionsStore.institutions.find(i => i.id === institutionId) || null
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
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(updateData)
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to update profile')
    }
    
    // Fetch updated user data from server to ensure we have the latest
    await authStore.fetchCurrentUser()
    
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

// Institution image handlers
function handleInstitutionImageUpload(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    institutionImageFile.value = file
    institutionImagePreview.value = URL.createObjectURL(file)
  }
}

function clearInstitutionImage() {
  institutionImageFile.value = null
  institutionImagePreview.value = null
  newInstitution.value.imagePath = ''
}

// Edit institution image handlers
function handleEditInstitutionImageUpload(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    editInstitutionImageFile.value = file
    editInstitutionImagePreview.value = URL.createObjectURL(file)
  }
}

function clearEditInstitutionImage() {
  if (editInstitutionImagePreview.value) {
    URL.revokeObjectURL(editInstitutionImagePreview.value)
  }
  editInstitutionImageFile.value = null
  editInstitutionImagePreview.value = null
  editInstitution.value.imagePath = ''
}

// Open edit institution modal
function openEditInstitutionModal() {
  if (!currentUserInstitution.value) return
  
  const inst = currentUserInstitution.value
  editInstitution.value = {
    name: inst.name,
    shortName: inst.shortName || '',
    type: inst.type,
    country: inst.country,
    province: inst.province || '',
    city: inst.city,
    address: inst.address || '',
    postalCode: inst.postalCode || '',
    phone: inst.phone || '',
    email: inst.email || '',
    website: inst.website || '',
    imagePath: inst.imagePath || ''
  }
  
  // Set image preview if exists (use full URL for existing images)
  if (inst.imagePath) {
    editInstitutionImagePreview.value = inst.imagePath
  } else {
    editInstitutionImagePreview.value = null
  }
  
  editInstitutionImageFile.value = null
  showEditInstitutionModal.value = true
}

// Close edit institution modal
function closeEditInstitutionModal() {
  showEditInstitutionModal.value = false
  clearEditInstitutionImage()
  editInstitution.value = {
    name: '',
    shortName: '',
    type: '' as any,
    country: '',
    province: '',
    city: '',
    address: '',
    postalCode: '',
    phone: '',
    email: '',
    website: '',
    imagePath: ''
  }
}

// Update institution
async function updateInstitution() {
  if (!currentUserInstitution.value || !user.value) return
  
  editingInstitution.value = true
  
  try {
    // Upload image if a new one was selected
    let imagePath: string | undefined = undefined
    if (editInstitutionImageFile.value) {
      imagePath = await institutionsStore.uploadUserImage(editInstitutionImageFile.value)
      editInstitution.value.imagePath = imagePath
    }
    
    // Update the institution
    const updatedInstitution = await institutionsStore.updateUserInstitution(currentUserInstitution.value.id, {
      name: editInstitution.value.name,
      shortName: editInstitution.value.shortName || undefined,
      type: editInstitution.value.type,
      country: editInstitution.value.country,
      province: editInstitution.value.province || undefined,
      city: editInstitution.value.city,
      address: editInstitution.value.address || undefined,
      postalCode: editInstitution.value.postalCode || undefined,
      phone: editInstitution.value.phone || undefined,
      email: editInstitution.value.email || undefined,
      website: editInstitution.value.website || undefined,
      imagePath: editInstitution.value.imagePath || undefined
    })
    
    if (!updatedInstitution) {
      throw new Error('Failed to update institution')
    }
    
    // Refresh institutions list
    await institutionsStore.fetchInstitutions({ limit: 1000 })
    
    closeEditInstitutionModal()
    toast.success('Institution updated successfully!')
  } catch (error: any) {
    console.error('Failed to update institution:', error)
    toast.error(error.message || 'Failed to update institution')
  } finally {
    editingInstitution.value = false
  }
}

// Institution modal functions
function closeInstitutionModal() {
  showInstitutionModal.value = false
  showCreateForm.value = false
  selectedInstitutionId.value = ''
  // Clear image state
  if (institutionImagePreview.value) {
    URL.revokeObjectURL(institutionImagePreview.value)
  }
  institutionImageFile.value = null
  institutionImagePreview.value = null
  newInstitution.value = {
    name: '',
    shortName: '',
    type: '' as any,
    country: 'South Africa',
    province: '',
    city: '',
    address: '',
    postalCode: '',
    phone: '',
    email: '',
    website: '',
    imagePath: ''
  }
}

async function changeInstitution() {
  if (!selectedInstitutionId.value || !user.value) return
  
  changingInstitution.value = true
  
  try {
    const response = await fetch(`/api/users/${user.value.id}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        institutionId: selectedInstitutionId.value
      })
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to change institution')
    }
    
    // Fetch updated user data from server to ensure we have the latest
    await authStore.fetchCurrentUser()
    
    // Refresh institutions list
    await institutionsStore.fetchInstitutions({ limit: 1000 })
    
    closeInstitutionModal()
    toast.success('Institution changed successfully!')
  } catch (error: any) {
    console.error('Failed to change institution:', error)
    toast.error(error.message || 'Failed to change institution')
  } finally {
    changingInstitution.value = false
  }
}

async function createAndSwitchInstitution() {
  if (!user.value) return
  
  creatingInstitution.value = true
  
  try {
    // Upload image if selected (use user upload endpoint)
    let imagePath: string | undefined = undefined
    if (institutionImageFile.value) {
      imagePath = await institutionsStore.uploadUserImage(institutionImageFile.value)
    }
    
    // Create the institution
    const institution = await institutionsStore.createUserInstitution({
      name: newInstitution.value.name,
      shortName: newInstitution.value.shortName || undefined,
      type: newInstitution.value.type,
      country: newInstitution.value.country,
      province: newInstitution.value.province || undefined,
      city: newInstitution.value.city,
      address: newInstitution.value.address || undefined,
      postalCode: newInstitution.value.postalCode || undefined,
      phone: newInstitution.value.phone || undefined,
      email: newInstitution.value.email || undefined,
      website: newInstitution.value.website || undefined,
      imagePath: imagePath
    })
    
    if (!institution) {
      throw new Error('Failed to create institution')
    }
    
    // Automatically switch to the new institution
    const response = await fetch(`/api/users/${user.value.id}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${authStore.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        institutionId: institution.id
      })
    })
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to switch to new institution')
    }
    
    // Fetch updated user data from server to ensure we have the latest
    await authStore.fetchCurrentUser()
    
    // Refresh institutions list
    await institutionsStore.fetchInstitutions({ limit: 1000 })
    
    closeInstitutionModal()
    toast.success(`Institution "${institution.name}" created successfully! You have been automatically switched to it.`)
  } catch (error: any) {
    console.error('Failed to create institution:', error)
    toast.error(error.message || 'Failed to create institution')
  } finally {
    creatingInstitution.value = false
  }
}
</script>


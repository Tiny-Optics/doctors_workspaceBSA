import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Institution, CreateInstitutionRequest, UpdateInstitutionRequest, InstitutionsListResponse, InstitutionType } from '@/types/institution'
import { useAuthStore } from './auth'

export const useInstitutionsStore = defineStore('institutions', () => {
  // State
  const institutions = ref<Institution[]>([])
  const currentInstitution = ref<Institution | null>(null)
  const total = ref(0)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Get auth store
  const authStore = useAuthStore()

  // Helper to get auth headers
  function getAuthHeaders(): HeadersInit {
    return {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${authStore.token}`
    }
  }

  // Actions
  async function fetchInstitutions(options?: {
    type?: InstitutionType
    isActive?: boolean
    search?: string
    limit?: number
    skip?: number
  }): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const params = new URLSearchParams()
      if (options?.type) params.append('type', options.type)
      if (options?.isActive !== undefined) params.append('is_active', String(options.isActive))
      if (options?.search) params.append('search', options.search)
      if (options?.limit) params.append('limit', String(options.limit))
      if (options?.skip) params.append('skip', String(options.skip))

      const response = await fetch(`/api/institutions?${params.toString()}`, {
        headers: getAuthHeaders()
      })

      if (!response.ok) {
        let errorMessage = 'Failed to fetch institutions'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          // Use status-based error message if JSON parsing fails
          errorMessage = response.status === 401 ? 'Session expired' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const data: InstitutionsListResponse = await response.json()
      institutions.value = data.institutions || []
      total.value = data.total || 0
      isLoading.value = false
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch institutions'
      institutions.value = []
      total.value = 0
      isLoading.value = false
      throw err
    }
  }

  // Public method to fetch institutions without authentication (for registration)
  async function fetchPublicInstitutions(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch('/api/institutions/public', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        let errorMessage = 'Failed to fetch institutions'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = 'Failed to fetch institutions'
        }
        throw new Error(errorMessage)
      }

      const data: InstitutionsListResponse = await response.json()
      institutions.value = data.institutions || []
      total.value = data.total || 0
      isLoading.value = false
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch institutions'
      institutions.value = []
      total.value = 0
      isLoading.value = false
      throw err
    }
  }

  async function fetchInstitution(id: string): Promise<Institution | null> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/institutions/${id}`, {
        headers: getAuthHeaders()
      })

      if (!response.ok) {
        let errorMessage = 'Failed to fetch institution'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 404 ? 'Institution not found' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const institution: Institution = await response.json()
      currentInstitution.value = institution
      isLoading.value = false
      return institution
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch institution'
      currentInstitution.value = null
      isLoading.value = false
      throw err
    }
  }

  async function createInstitution(data: CreateInstitutionRequest): Promise<Institution | null> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch('/api/institutions', {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(data)
      })

      if (!response.ok) {
        let errorMessage = 'Failed to create institution'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 409 ? 'Institution already exists' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const institution: Institution = await response.json()
      institutions.value.push(institution)
      isLoading.value = false
      return institution
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create institution'
      isLoading.value = false
      throw err
    }
  }

  // Create institution as a regular user
  async function createUserInstitution(data: CreateInstitutionRequest): Promise<Institution | null> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch('/api/institutions/user/create', {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(data)
      })

      if (!response.ok) {
        let errorMessage = 'Failed to create institution'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 409 ? 'Institution already exists' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const institution: Institution = await response.json()
      institutions.value.push(institution)
      isLoading.value = false
      return institution
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create institution'
      isLoading.value = false
      throw err
    }
  }

  async function updateInstitution(id: string, data: UpdateInstitutionRequest): Promise<Institution | null> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/institutions/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify(data)
      })

      if (!response.ok) {
        let errorMessage = 'Failed to update institution'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 404 ? 'Institution not found' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const institution: Institution = await response.json()
      const index = institutions.value.findIndex(i => i.id === id)
      if (index !== -1 && institutions.value[index]) {
        institutions.value[index] = institution
      }
      isLoading.value = false
      return institution
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update institution'
      isLoading.value = false
      throw err
    }
  }

  // Update institution as a regular user (only if they created it)
  async function updateUserInstitution(id: string, data: UpdateInstitutionRequest): Promise<Institution | null> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/institutions/user/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify(data)
      })

      if (!response.ok) {
        let errorMessage = 'Failed to update institution'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 404 ? 'Institution not found' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const institution: Institution = await response.json()
      const index = institutions.value.findIndex(i => i.id === id)
      if (index !== -1 && institutions.value[index]) {
        institutions.value[index] = institution
      }
      isLoading.value = false
      return institution
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update institution'
      isLoading.value = false
      throw err
    }
  }

  async function deleteInstitution(id: string): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/institutions/${id}`, {
        method: 'DELETE',
        headers: getAuthHeaders()
      })

      if (!response.ok) {
        let errorMessage = 'Failed to delete institution'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 404 ? 'Institution not found' : errorMessage
        }
        throw new Error(errorMessage)
      }

      institutions.value = institutions.value.filter(i => i.id !== id)
      isLoading.value = false
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete institution'
      isLoading.value = false
      throw err
    }
  }

  async function activateInstitution(id: string): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/institutions/${id}/activate`, {
        method: 'POST',
        headers: getAuthHeaders()
      })

      if (!response.ok) {
        let errorMessage = 'Failed to activate institution'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 404 ? 'Institution not found' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const index = institutions.value.findIndex(i => i.id === id)
      if (index !== -1 && institutions.value[index]) {
        institutions.value[index]!.isActive = true
      }
      isLoading.value = false
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to activate institution'
      isLoading.value = false
      throw err
    }
  }

  async function deactivateInstitution(id: string): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`/api/institutions/${id}/deactivate`, {
        method: 'POST',
        headers: getAuthHeaders()
      })

      if (!response.ok) {
        let errorMessage = 'Failed to deactivate institution'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 404 ? 'Institution not found' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const index = institutions.value.findIndex(i => i.id === id)
      if (index !== -1 && institutions.value[index]) {
        institutions.value[index]!.isActive = false
      }
      isLoading.value = false
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to deactivate institution'
      isLoading.value = false
      throw err
    }
  }

  async function uploadImage(file: File): Promise<string> {
    isLoading.value = true
    error.value = null

    try {
      const token = authStore.token
      const formData = new FormData()
      formData.append('image', file)

      const response = await fetch('/api/institutions/images/upload', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
          // Note: Don't set Content-Type, browser will set it with boundary for multipart
        },
        body: formData
      })

      if (!response.ok) {
        let errorMessage = 'Failed to upload image'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = 'Failed to upload image'
        }
        throw new Error(errorMessage)
      }

      const data = await response.json()
      isLoading.value = false
      return data.imagePath
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to upload image'
      isLoading.value = false
      throw err
    }
  }

  // Upload image for user-created institutions (no admin permission required)
  async function uploadUserImage(file: File): Promise<string> {
    isLoading.value = true
    error.value = null

    try {
      const token = authStore.token
      const formData = new FormData()
      formData.append('image', file)

      const response = await fetch('/api/institutions/user/images/upload', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
          // Note: Don't set Content-Type, browser will set it with boundary for multipart
        },
        body: formData
      })

      if (!response.ok) {
        let errorMessage = 'Failed to upload image'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = 'Failed to upload image'
        }
        throw new Error(errorMessage)
      }

      const data = await response.json()
      isLoading.value = false
      return data.imagePath
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to upload image'
      isLoading.value = false
      throw err
    }
  }

  return {
    // State
    institutions,
    currentInstitution,
    total,
    isLoading,
    error,
    // Actions
    fetchInstitutions,
    fetchPublicInstitutions,
    fetchInstitution,
    createInstitution,
    createUserInstitution,
    updateInstitution,
    updateUserInstitution,
    deleteInstitution,
    activateInstitution,
    deactivateInstitution,
    uploadImage,
    uploadUserImage
  }
})


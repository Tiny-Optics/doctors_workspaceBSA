import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User, CreateUserRequest, UpdateUserRequest, UsersListResponse, UserRole } from '@/types/user'
import { useAuthStore } from './auth'

export const useUsersStore = defineStore('users', () => {
  // State
  const users = ref<User[]>([])
  const currentUser = ref<User | null>(null)
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
  async function fetchUsers(options?: {
    role?: UserRole
    isActive?: boolean
    limit?: number
    skip?: number
  }): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const params = new URLSearchParams()
      if (options?.role) params.append('role', options.role)
      if (options?.isActive !== undefined) params.append('is_active', String(options.isActive))
      if (options?.limit) params.append('limit', String(options.limit))
      if (options?.skip) params.append('skip', String(options.skip))

      const response = await fetch(`http://localhost:8080/api/users?${params.toString()}`, {
        headers: getAuthHeaders()
      })

      if (!response.ok) {
        let errorMessage = 'Failed to fetch users'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          // Use status-based error message if JSON parsing fails
          errorMessage = response.status === 401 ? 'Session expired' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const data: UsersListResponse = await response.json()
      users.value = data.users
      total.value = data.total
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch users'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function fetchUser(id: string): Promise<User> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:8080/api/users/${id}`, {
        headers: getAuthHeaders()
      })

      if (!response.ok) {
        let errorMessage = 'Failed to fetch user'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 404 ? 'User not found' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const user: User = await response.json()
      currentUser.value = user
      return user
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch user'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function createUser(request: CreateUserRequest): Promise<User> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch('http://localhost:8080/api/users', {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(request)
      })

      if (!response.ok) {
        let errorMessage = 'Failed to create user'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          if (response.status === 403) {
            errorMessage = 'Insufficient permissions to create users'
          } else if (response.status === 409) {
            errorMessage = 'Email or username already exists'
          }
        }
        throw new Error(errorMessage)
      }

      const user: User = await response.json()
      
      // Add to local list if it's loaded
      if (users.value.length > 0) {
        users.value.unshift(user)
        total.value++
      }

      return user
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create user'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function updateUser(id: string, request: UpdateUserRequest): Promise<User> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:8080/api/users/${id}`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify(request)
      })

      if (!response.ok) {
        let errorMessage = 'Failed to update user'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 403 ? 'Insufficient permissions' : errorMessage
        }
        throw new Error(errorMessage)
      }

      const user: User = await response.json()
      
      // Update in local list
      const index = users.value.findIndex(u => u.id === id)
      if (index !== -1) {
        users.value[index] = user
      }

      if (currentUser.value?.id === id) {
        currentUser.value = user
      }

      return user
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update user'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function deleteUser(id: string): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:8080/api/users/${id}`, {
        method: 'DELETE',
        headers: getAuthHeaders()
      })

      if (!response.ok) {
        let errorMessage = 'Failed to delete user'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 403 ? 'Insufficient permissions' : errorMessage
        }
        throw new Error(errorMessage)
      }

      // Remove from local list
      users.value = users.value.filter(u => u.id !== id)
      total.value--

      if (currentUser.value?.id === id) {
        currentUser.value = null
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete user'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function deactivateUser(id: string): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:8080/api/users/${id}/deactivate`, {
        method: 'POST',
        headers: getAuthHeaders()
      })

      if (!response.ok) {
        let errorMessage = 'Failed to deactivate user'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 403 ? 'Insufficient permissions' : errorMessage
        }
        throw new Error(errorMessage)
      }

      // Update in local list
      const index = users.value.findIndex(u => u.id === id)
      if (index !== -1) {
        const user = users.value[index]
        if (user) {
          user.isActive = false
        }
      }

      if (currentUser.value?.id === id) {
        currentUser.value.isActive = false
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to deactivate user'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function activateUser(id: string): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:8080/api/users/${id}/activate`, {
        method: 'POST',
        headers: getAuthHeaders()
      })

      if (!response.ok) {
        let errorMessage = 'Failed to activate user'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = response.status === 403 ? 'Insufficient permissions' : errorMessage
        }
        throw new Error(errorMessage)
      }

      // Update in local list
      const index = users.value.findIndex(u => u.id === id)
      if (index !== -1) {
        const user = users.value[index]
        if (user) {
          user.isActive = true
        }
      }

      if (currentUser.value?.id === id) {
        currentUser.value.isActive = true
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to activate user'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  function clearError(): void {
    error.value = null
  }

  function clearCurrentUser(): void {
    currentUser.value = null
  }

  return {
    // State
    users,
    currentUser,
    total,
    isLoading,
    error,
    
    // Actions
    fetchUsers,
    fetchUser,
    createUser,
    updateUser,
    deleteUser,
    deactivateUser,
    activateUser,
    clearError,
    clearCurrentUser
  }
})


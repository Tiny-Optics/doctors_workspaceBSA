import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginRequest, LoginResponse } from '@/types/user'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const isAuthenticated = computed(() => !!user.value && !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isSuperAdmin = computed(() => user.value?.adminLevel === 'super_admin')
  const isUserManager = computed(() => user.value?.adminLevel === 'user_manager')

  // Actions
  async function login(credentials: LoginRequest): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(credentials)
      })

      if (!response.ok) {
        let errorMessage = 'Login failed. Please try again.'
        
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          // If JSON parsing fails, use status-based messages
          if (response.status === 401) {
            errorMessage = 'Invalid email or password. Please try again.'
          } else if (response.status === 403) {
            errorMessage = 'Your account has been locked. Please contact support.'
          } else if (response.status === 500) {
            errorMessage = 'Server error. Please try again later.'
          } else if (response.status === 0 || !response.status) {
            errorMessage = 'Cannot connect to server. Please check your connection.'
          }
        }
        
        throw new Error(errorMessage)
      }

      const data: LoginResponse = await response.json()
      
      user.value = data.user
      token.value = data.token
      refreshToken.value = data.refreshToken

      // Store in localStorage for persistence
      localStorage.setItem('token', data.token)
      localStorage.setItem('refreshToken', data.refreshToken)
      localStorage.setItem('user', JSON.stringify(data.user))
    } catch (err) {
      if (err instanceof Error) {
        error.value = err.message
      } else if (typeof err === 'string') {
        error.value = err
      } else {
        error.value = 'An unexpected error occurred during login'
      }
      throw new Error(error.value)
    } finally {
      isLoading.value = false
    }
  }

  async function logout(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      if (token.value) {
        await fetch('/api/auth/logout', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token.value}`
          }
        })
      }
    } catch (err) {
      console.error('Logout error:', err)
    } finally {
      // Clear state regardless of API call success
      user.value = null
      token.value = null
      refreshToken.value = null
      
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('user')
      
      isLoading.value = false
    }
  }

  async function refreshAccessToken(): Promise<void> {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await fetch('/api/auth/refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ refreshToken: refreshToken.value })
      })

      if (!response.ok) {
        let errorMessage = 'Failed to refresh token'
        
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          errorMessage = 'Session expired. Please login again.'
        }
        
        throw new Error(errorMessage)
      }

      const data: LoginResponse = await response.json()
      
      user.value = data.user
      token.value = data.token
      refreshToken.value = data.refreshToken

      localStorage.setItem('token', data.token)
      localStorage.setItem('refreshToken', data.refreshToken)
      localStorage.setItem('user', JSON.stringify(data.user))
    } catch (err) {
      // If refresh fails, log out
      await logout()
      throw err
    }
  }

  async function fetchCurrentUser(): Promise<void> {
    if (!token.value) {
      return
    }

    isLoading.value = true
    error.value = null

    try {
      const response = await fetch('/api/auth/me', {
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      })

      if (!response.ok) {
        let errorMessage = 'Failed to fetch user information'
        
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          if (response.status === 401) {
            errorMessage = 'Session expired. Please login again.'
          }
        }
        
        throw new Error(errorMessage)
      }

      const data: User = await response.json()
      user.value = data
      localStorage.setItem('user', JSON.stringify(data))
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch user'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function changePassword(oldPassword: string, newPassword: string): Promise<void> {
    if (!token.value) {
      throw new Error('Not authenticated')
    }

    isLoading.value = true
    error.value = null

    try {
      const response = await fetch('/api/auth/change-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token.value}`
        },
        body: JSON.stringify({ oldPassword, newPassword })
      })

      if (!response.ok) {
        let errorMessage = 'Failed to change password'
        
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          if (response.status === 401) {
            errorMessage = 'Current password is incorrect'
          } else if (response.status === 400) {
            errorMessage = 'New password does not meet requirements'
          }
        }
        
        throw new Error(errorMessage)
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to change password'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  function initializeFromStorage(): void {
    const storedToken = localStorage.getItem('token')
    const storedRefreshToken = localStorage.getItem('refreshToken')
    const storedUser = localStorage.getItem('user')

    if (storedToken && storedRefreshToken && storedUser) {
      token.value = storedToken
      refreshToken.value = storedRefreshToken
      try {
        user.value = JSON.parse(storedUser)
      } catch (err) {
        console.error('Failed to parse stored user:', err)
        localStorage.removeItem('user')
      }
    }
  }

  function clearError(): void {
    error.value = null
  }

  return {
    // State
    user,
    token,
    refreshToken,
    isLoading,
    error,
    
    // Computed
    isAuthenticated,
    isAdmin,
    isSuperAdmin,
    isUserManager,
    
    // Actions
    login,
    logout,
    refreshAccessToken,
    fetchCurrentUser,
    changePassword,
    initializeFromStorage,
    clearError
  }
})


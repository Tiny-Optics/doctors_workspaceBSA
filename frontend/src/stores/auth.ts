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
        const errorData = await response.json()
        throw new Error(errorData.error || 'Login failed')
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
      error.value = err instanceof Error ? err.message : 'An error occurred during login'
      throw err
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
        throw new Error('Failed to refresh token')
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
        throw new Error('Failed to fetch user')
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
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to change password')
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


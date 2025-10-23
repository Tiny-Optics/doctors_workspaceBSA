import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { 
  ForgotPasswordRequest, 
  ForgotPasswordResponse,
  ValidateResetCodeRequest,
  ValidateResetCodeResponse,
  ResetPasswordRequest,
  ResetPasswordResponse
} from '@/types/user'

export const useForgotPasswordStore = defineStore('forgotPassword', () => {
  // State
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  const success = ref<string | null>(null)
  const currentStep = ref<'email' | 'code' | 'password' | 'success'>('email')
  const resetToken = ref<string | null>(null)

  // Actions
  async function requestPasswordReset(email: string): Promise<void> {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const response = await fetch('/api/auth/forgot-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email })
      })

      if (!response.ok) {
        let errorMessage = 'Failed to send password reset code'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          if (response.status === 404) {
            errorMessage = 'No account found with this email address'
          } else if (response.status === 429) {
            errorMessage = 'Too many requests. Please wait before trying again'
          } else if (response.status === 403) {
            errorMessage = 'Account is inactive. Please contact an administrator'
          } else if (response.status === 0 || !response.status) {
            errorMessage = 'Cannot connect to server. Please check your connection'
          }
        }
        throw new Error(errorMessage)
      }

      const data: ForgotPasswordResponse = await response.json()
      success.value = data.message
      currentStep.value = 'code'
    } catch (err) {
      if (err instanceof Error) {
        error.value = err.message
      } else if (typeof err === 'string') {
        error.value = err
      } else {
        error.value = 'An unexpected error occurred'
      }
      throw new Error(error.value)
    } finally {
      isLoading.value = false
    }
  }

  async function validateResetCode(code: string): Promise<void> {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const response = await fetch('/api/auth/validate-reset-code', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ code })
      })

      if (!response.ok) {
        let errorMessage = 'Invalid or expired reset code'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          if (response.status === 404) {
            errorMessage = 'Invalid or expired reset code'
          } else if (response.status === 0 || !response.status) {
            errorMessage = 'Cannot connect to server. Please check your connection'
          }
        }
        throw new Error(errorMessage)
      }

      const data: ValidateResetCodeResponse = await response.json()
      resetToken.value = data.token
      success.value = data.message
      currentStep.value = 'password'
    } catch (err) {
      if (err instanceof Error) {
        error.value = err.message
      } else if (typeof err === 'string') {
        error.value = err
      } else {
        error.value = 'An unexpected error occurred'
      }
      throw new Error(error.value)
    } finally {
      isLoading.value = false
    }
  }

  async function resetPassword(newPassword: string): Promise<void> {
    if (!resetToken.value) {
      throw new Error('No reset token available')
    }

    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const response = await fetch('/api/auth/reset-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          token: resetToken.value,
          newPassword
        })
      })

      if (!response.ok) {
        let errorMessage = 'Failed to reset password'
        try {
          const errorData = await response.json()
          errorMessage = errorData.error || errorMessage
        } catch (parseError) {
          if (response.status === 404) {
            errorMessage = 'Invalid or expired reset token'
          } else if (response.status === 0 || !response.status) {
            errorMessage = 'Cannot connect to server. Please check your connection'
          }
        }
        throw new Error(errorMessage)
      }

      const data: ResetPasswordResponse = await response.json()
      success.value = data.message
      currentStep.value = 'success'
    } catch (err) {
      if (err instanceof Error) {
        error.value = err.message
      } else if (typeof err === 'string') {
        error.value = err
      } else {
        error.value = 'An unexpected error occurred'
      }
      throw new Error(error.value)
    } finally {
      isLoading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  function clearSuccess() {
    success.value = null
  }

  function reset() {
    isLoading.value = false
    error.value = null
    success.value = null
    currentStep.value = 'email'
    resetToken.value = null
  }

  return {
    // State
    isLoading,
    error,
    success,
    currentStep,
    resetToken,
    // Actions
    requestPasswordReset,
    validateResetCode,
    resetPassword,
    clearError,
    clearSuccess,
    reset
  }
})

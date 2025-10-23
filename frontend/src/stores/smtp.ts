import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useSMTPStore = defineStore('smtp', () => {
  // State
  const isConfigured = ref(false)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Actions
  async function checkSMTPConfiguration(): Promise<void> {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch('/api/smtp/status', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        // If we can't access the endpoint or it returns an error, assume SMTP is not configured
        isConfigured.value = false
        return
      }

      const result = await response.json()
      isConfigured.value = result.isConfigured || false
    } catch (err) {
      // If there's any error (network, CORS, etc.), assume SMTP is not configured
      isConfigured.value = false
      error.value = err instanceof Error ? err.message : 'Failed to check SMTP configuration'
    } finally {
      isLoading.value = false
    }
  }

  function reset() {
    isConfigured.value = false
    isLoading.value = false
    error.value = null
  }

  return {
    // State
    isConfigured,
    isLoading,
    error,
    // Actions
    checkSMTPConfiguration,
    reset
  }
})

import type {
  DropboxStatus,
  InitiateAuthRequest,
  InitiateAuthResponse,
  CompleteAuthRequest,
  CompleteAuthResponse,
  TestConnectionResponse,
  RefreshTokenResponse,
  DeleteConfigResponse
} from '@/types/dropbox'

const API_URL = '/api'

// Helper function to get auth headers
function getAuthHeaders(): HeadersInit {
  const token = localStorage.getItem('token')
  return {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
}

// Helper function to handle errors
async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }))
    throw new Error(error.error || `HTTP error ${response.status}`)
  }
  return response.json()
}

export const dropboxAdminService = {
  /**
   * Get Dropbox connection status
   */
  async getStatus(): Promise<DropboxStatus> {
    const response = await fetch(`${API_URL}/admin/dropbox/status`, {
      headers: getAuthHeaders()
    })
    return handleResponse<DropboxStatus>(response)
  },

  /**
   * Initiate OAuth authorization flow
   */
  async initiateAuth(request: InitiateAuthRequest): Promise<InitiateAuthResponse> {
    const response = await fetch(`${API_URL}/admin/dropbox/authorize`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(request)
    })
    return handleResponse<InitiateAuthResponse>(response)
  },

  /**
   * Complete OAuth authorization with code
   */
  async completeAuth(request: CompleteAuthRequest): Promise<CompleteAuthResponse> {
    const response = await fetch(`${API_URL}/admin/dropbox/callback`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(request)
    })
    return handleResponse<CompleteAuthResponse>(response)
  },

  /**
   * Force refresh access token
   */
  async forceRefresh(): Promise<RefreshTokenResponse> {
    const response = await fetch(`${API_URL}/admin/dropbox/refresh`, {
      method: 'POST',
      headers: getAuthHeaders()
    })
    return handleResponse<RefreshTokenResponse>(response)
  },

  /**
   * Test Dropbox connection
   */
  async testConnection(): Promise<TestConnectionResponse> {
    const response = await fetch(`${API_URL}/admin/dropbox/test`, {
      method: 'POST',
      headers: getAuthHeaders()
    })
    return handleResponse<TestConnectionResponse>(response)
  },

  /**
   * Delete Dropbox configuration
   */
  async deleteConfiguration(): Promise<DeleteConfigResponse> {
    const response = await fetch(`${API_URL}/admin/dropbox/configuration`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })
    return handleResponse<DeleteConfigResponse>(response)
  }
}


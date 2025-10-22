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

export interface ReferralConfig {
  id?: string
  redcapUrl: string
  isEnabled: boolean
  createdAt?: string
  updatedAt?: string
  updatedBy?: string
}

export interface UpdateReferralConfigRequest {
  redcapUrl?: string
  isEnabled?: boolean
}

export interface ReferralConfigResponse {
  isConfigured: boolean
  isEnabled: boolean
  redcapUrl: string
}

export interface ReferralAccessResponse {
  redirectUrl: string
}

class ReferralService {
  /**
   * Get referral configuration for regular users
   */
  async getReferralConfig(): Promise<ReferralConfigResponse> {
    const response = await fetch(`${API_URL}/referrals/config`, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<ReferralConfigResponse>(response)
  }

  /**
   * Get full referral configuration for admin users
   */
  async getAdminReferralConfig(): Promise<ReferralConfig> {
    const response = await fetch(`${API_URL}/admin/referrals/config`, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<ReferralConfig>(response)
  }

  /**
   * Update referral configuration (admin only)
   */
  async updateReferralConfig(config: UpdateReferralConfigRequest): Promise<ReferralConfig> {
    const response = await fetch(`${API_URL}/admin/referrals/config`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(config)
    })
    return handleResponse<ReferralConfig>(response)
  }

  /**
   * Log referral access and get redirect URL
   */
  async logReferralAccess(): Promise<ReferralAccessResponse> {
    const response = await fetch(`${API_URL}/referrals/access`, {
      method: 'POST',
      headers: getAuthHeaders()
    })
    return handleResponse<ReferralAccessResponse>(response)
  }
}

export const referralService = new ReferralService()

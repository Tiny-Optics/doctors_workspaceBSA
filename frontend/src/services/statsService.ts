// Admin statistics service

export interface AdminStats {
  totalUsers: number
  activeUsers: number
  newUsersThisMonth: number
  newUsersThisWeek: number
  newUsersToday: number
  totalInstitutions: number
  totalSOPs: number
  roleDistribution: RoleDistribution[]
}

export interface RoleDistribution {
  role: string
  count: number
}

const API_URL = 'http://localhost:8080/api'

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

export const statsService = {
  /**
   * Get admin dashboard statistics
   */
  async getAdminStats(): Promise<AdminStats> {
    const response = await fetch(`${API_URL}/stats/admin`, {
      headers: getAuthHeaders()
    })
    return handleResponse<AdminStats>(response)
  }
}


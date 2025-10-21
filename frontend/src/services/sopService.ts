import type {
  SOPCategory,
  SOPFile,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  CategoryListResponse,
  FileListResponse,
  DownloadLinkResponse,
  ImageUploadResponse,
  SeedResponse
} from '@/types/sop'

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

export const sopService = {
  // ===== CATEGORIES =====
  
  /**
   * List all SOP categories
   */
  async listCategories(params?: { page?: number; limit?: number; search?: string }): Promise<CategoryListResponse> {
    const queryParams = new URLSearchParams()
    if (params?.page) queryParams.append('page', params.page.toString())
    if (params?.limit) queryParams.append('limit', params.limit.toString())
    if (params?.search) queryParams.append('search', params.search)
    
    const url = `${API_URL}/sops/categories${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    
    const response = await fetch(url, {
      headers: getAuthHeaders()
    })
    
    return handleResponse<CategoryListResponse>(response)
  },
  
  /**
   * Get a single category by ID
   */
  async getCategory(id: string): Promise<SOPCategory> {
    const response = await fetch(`${API_URL}/sops/categories/${id}`, {
      headers: getAuthHeaders()
    })
    
    return handleResponse<SOPCategory>(response)
  },
  
  /**
   * Get a category by slug
   */
  async getCategoryBySlug(slug: string): Promise<SOPCategory | null> {
    const response = await this.listCategories({ search: slug })
    const category = response.categories.find(c => c.slug === slug)
    return category || null
  },
  
  /**
   * Create a new SOP category
   */
  async createCategory(data: CreateCategoryRequest): Promise<SOPCategory> {
    const response = await fetch(`${API_URL}/sops/categories`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data)
    })
    
    return handleResponse<SOPCategory>(response)
  },
  
  /**
   * Update an existing category
   */
  async updateCategory(id: string, data: UpdateCategoryRequest): Promise<SOPCategory> {
    const response = await fetch(`${API_URL}/sops/categories/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data)
    })
    
    return handleResponse<SOPCategory>(response)
  },
  
  /**
   * Delete a category
   */
  async deleteCategory(id: string): Promise<void> {
    const response = await fetch(`${API_URL}/sops/categories/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })
    
    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Delete failed' }))
      throw new Error(error.error || `HTTP error ${response.status}`)
    }
  },
  
  // ===== FILES =====
  
  /**
   * Get all files in a category from Dropbox
   */
  async getCategoryFiles(categoryId: string): Promise<SOPFile[]> {
    const response = await fetch(`${API_URL}/sops/categories/${categoryId}/files`, {
      headers: getAuthHeaders()
    })
    
    const data = await handleResponse<FileListResponse>(response)
    return data.files
  },
  
  /**
   * Get a temporary download link for a file
   */
  async getFileDownloadLink(categoryId: string, filePath: string): Promise<string> {
    const params = new URLSearchParams({ path: filePath })
    
    const response = await fetch(
      `${API_URL}/sops/categories/${categoryId}/files/download?${params.toString()}`,
      { headers: getAuthHeaders() }
    )
    
    const data = await handleResponse<DownloadLinkResponse>(response)
    return data.downloadLink
  },
  
  // ===== IMAGES =====
  
  /**
   * Upload a category image
   */
  async uploadImage(file: File): Promise<string> {
    const token = localStorage.getItem('token')
    const formData = new FormData()
    formData.append('image', file)
    
    const response = await fetch(`${API_URL}/sops/images/upload`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
        // Note: Don't set Content-Type, browser will set it with boundary for multipart
      },
      body: formData
    })
    
    const data = await handleResponse<ImageUploadResponse>(response)
    return data.imagePath
  },
  
  // ===== SEEDING =====
  
  /**
   * Seed initial categories (admin only)
   */
  async seedCategories(): Promise<SeedResponse> {
    const response = await fetch(`${API_URL}/sops/seed`, {
      method: 'POST',
      headers: getAuthHeaders()
    })
    
    return handleResponse<SeedResponse>(response)
  }
}


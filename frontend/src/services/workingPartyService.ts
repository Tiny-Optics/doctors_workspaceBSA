import type {
  WorkingPartyCategory,
  WorkingPartyFile,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  CategoryListResponse,
  FileListResponse,
  DownloadLinkResponse,
  ImageUploadResponse
} from '@/types/workingParty'

const API_URL = '/api'

function getAuthHeaders(): HeadersInit {
  const token = localStorage.getItem('token')
  return {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }))
    throw new Error(error.error || `HTTP error ${response.status}`)
  }
  return response.json()
}

export const workingPartyService = {
  async listCategories(params?: { page?: number; limit?: number; search?: string }): Promise<CategoryListResponse> {
    const queryParams = new URLSearchParams()
    if (params?.page) queryParams.append('page', params.page.toString())
    if (params?.limit) queryParams.append('limit', params.limit.toString())
    if (params?.search) queryParams.append('search', params.search)

    const url = `${API_URL}/working-parties/categories${queryParams.toString() ? '?' + queryParams.toString() : ''}`

    const response = await fetch(url, {
      headers: getAuthHeaders()
    })

    return handleResponse<CategoryListResponse>(response)
  },

  async getCategory(id: string): Promise<WorkingPartyCategory> {
    const response = await fetch(`${API_URL}/working-parties/categories/${id}`, {
      headers: getAuthHeaders()
    })

    return handleResponse<WorkingPartyCategory>(response)
  },

  async getCategoryBySlug(slug: string): Promise<WorkingPartyCategory | null> {
    const response = await this.listCategories({ search: slug })
    const category = response.categories?.find(c => c.slug === slug)
    return category || null
  },

  async createCategory(data: CreateCategoryRequest): Promise<WorkingPartyCategory> {
    const response = await fetch(`${API_URL}/working-parties/categories`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(data)
    })

    return handleResponse<WorkingPartyCategory>(response)
  },

  async updateCategory(id: string, data: UpdateCategoryRequest): Promise<WorkingPartyCategory> {
    const response = await fetch(`${API_URL}/working-parties/categories/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(data)
    })

    return handleResponse<WorkingPartyCategory>(response)
  },

  async deleteCategory(id: string): Promise<void> {
    const response = await fetch(`${API_URL}/working-parties/categories/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Delete failed' }))
      throw new Error(error.error || `HTTP error ${response.status}`)
    }
  },

  async getCategoryFiles(categoryId: string): Promise<WorkingPartyFile[]> {
    const response = await fetch(`${API_URL}/working-parties/categories/${categoryId}/files`, {
      headers: getAuthHeaders()
    })

    const data = await handleResponse<FileListResponse>(response)
    return data.files
  },

  async getFileDownloadLink(categoryId: string, filePath: string): Promise<string> {
    const params = new URLSearchParams({ path: filePath })

    const response = await fetch(
      `${API_URL}/working-parties/categories/${categoryId}/files/download?${params.toString()}`,
      { headers: getAuthHeaders() }
    )

    const data = await handleResponse<DownloadLinkResponse>(response)
    return data.downloadLink
  },

  async uploadImage(file: File): Promise<string> {
    const token = localStorage.getItem('token')
    const formData = new FormData()
    formData.append('image', file)

    const response = await fetch(`${API_URL}/working-parties/images/upload`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    })

    const data = await handleResponse<ImageUploadResponse>(response)
    return data.imagePath
  }
}

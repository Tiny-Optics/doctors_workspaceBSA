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
  
  // Handle 204 No Content responses (DELETE operations)
  if (response.status === 204) {
    return undefined as T
  }
  
  return response.json()
}

export interface RegistryConfig {
  id: string
  videoUrl?: string
  documentsPath?: string
  notificationEmails: string[]
  smtpConfig: {
    host: string
    port: number
    username: string
    password: string
    fromEmail: string
    fromName: string
  }
  createdAt: string
  updatedAt: string
}

export interface FormField {
  id: string
  label: string
  type: 'text' | 'textarea' | 'select' | 'radio' | 'date' | 'number' | 'email' | 'file'
  required: boolean
  placeholder?: string
  helpText?: string
  options?: string[]
  allowMultiple?: boolean
  validationRules?: {
    minLength?: number
    maxLength?: number
    minValue?: number
    maxValue?: number
    pattern?: string
  }
  displayOrder?: number
}

export interface FormSchema {
  id: string
  formName: string
  description?: string
  fields: FormField[]
  documentFieldId?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface Submission {
  id: string
  userId: string
  formSchemaId: string
  formData: Record<string, any>
  documentsPath?: string
  uploadedDocuments: string[]
  status: 'submitted' | 'pending' | 'approved' | 'rejected'
  createdAt: string
  updatedAt: string
  userName?: string
  userEmail?: string
  formName?: string
}

class RegistryService {
  // Admin endpoint - requires super admin permission
  async getConfiguration(): Promise<RegistryConfig> {
    const response = await fetch(`${API_URL}/admin/registry/config`, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<RegistryConfig>(response)
  }

  async updateConfiguration(config: Partial<RegistryConfig>): Promise<RegistryConfig> {
    // Transform the config to match the backend's UpdateRegistryConfigRequest format
    const requestBody: any = {
      videoUrl: config.videoUrl,
      documentsPath: config.documentsPath,
      notificationEmails: config.notificationEmails
    }
    
    // Flatten SMTP config fields to match backend expectations
    if (config.smtpConfig) {
      requestBody.smtpHost = config.smtpConfig.host
      requestBody.smtpPort = config.smtpConfig.port
      requestBody.smtpUsername = config.smtpConfig.username
      requestBody.smtpPassword = config.smtpConfig.password
      requestBody.smtpFromEmail = config.smtpConfig.fromEmail
      requestBody.smtpFromName = config.smtpConfig.fromName
    }
    
    const response = await fetch(`${API_URL}/admin/registry/config`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(requestBody)
    })
    return handleResponse<RegistryConfig>(response)
  }

  // Public endpoint - for regular users viewing the registry
  async getPublicConfiguration(): Promise<RegistryConfig> {
    const response = await fetch(`${API_URL}/registry/config`, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<RegistryConfig>(response)
  }

  // Get active form for users
  async getActiveForm(): Promise<FormSchema | null> {
    const response = await fetch(`${API_URL}/registry/form-schema`, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<FormSchema | null>(response)
  }

  // Admin endpoints for form management
  async createFormSchema(schema: Partial<FormSchema>): Promise<FormSchema> {
    const response = await fetch(`${API_URL}/admin/registry/form-schema`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(schema)
    })
    return handleResponse<FormSchema>(response)
  }

  async updateFormSchema(id: string, schema: Partial<FormSchema>): Promise<FormSchema> {
    const response = await fetch(`${API_URL}/admin/registry/form-schema/${id}`, {
      method: 'PUT',
      headers: getAuthHeaders(),
      body: JSON.stringify(schema)
    })
    return handleResponse<FormSchema>(response)
  }

  async deleteFormSchema(id: string): Promise<void> {
    const response = await fetch(`${API_URL}/admin/registry/form-schema/${id}`, {
      method: 'DELETE',
      headers: getAuthHeaders()
    })
    await handleResponse<void>(response)
  }

  async listFormSchemas(params?: { page?: number, limit?: number }): Promise<{ schemas: FormSchema[], total: number, page: number, limit: number }> {
    const url = new URL(`${API_URL}/admin/registry/form-schemas`, window.location.origin)
    if (params?.page) {
      url.searchParams.append('page', params.page.toString())
    }
    if (params?.limit) {
      url.searchParams.append('limit', params.limit.toString())
    }
    
    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<{ schemas: FormSchema[], total: number, page: number, limit: number }>(response)
  }

  // User endpoints for submissions
  async submitForm(submission: {
    formSchemaId: string
    formData: Record<string, any>
    files: File[]
  }): Promise<Submission> {
    const formData = new FormData()
    formData.append('formSchemaId', submission.formSchemaId)
    formData.append('formData', JSON.stringify(submission.formData))
    
    submission.files.forEach((file, index) => {
      formData.append('files', file)
    })

    const token = localStorage.getItem('token')
    const response = await fetch(`${API_URL}/registry/submit`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    })
    return handleResponse<Submission>(response)
  }

  async getUserSubmissions(): Promise<Submission[]> {
    const response = await fetch(`${API_URL}/registry/submissions`, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<Submission[]>(response)
  }

  async getSubmission(id: string): Promise<Submission> {
    const response = await fetch(`${API_URL}/registry/submissions/${id}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<Submission>(response)
  }

  // Admin endpoints for submissions
  async getAllSubmissions(params?: {
    page?: number
    limit?: number
    status?: string
    userId?: string
  }): Promise<{ submissions: Submission[], total: number, page: number, limit: number }> {
    const url = new URL(`${API_URL}/admin/registry/submissions`, window.location.origin)
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          url.searchParams.append(key, value.toString())
        }
      })
    }
    
    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<{ submissions: Submission[], total: number, page: number, limit: number }>(response)
  }

  async updateSubmissionStatus(id: string, status: 'submitted' | 'pending' | 'approved' | 'rejected', reviewNotes?: string): Promise<Submission> {
    const body: any = { status }
    if (reviewNotes) {
      body.reviewNotes = reviewNotes
    }
    
    const response = await fetch(`${API_URL}/admin/registry/submissions/${id}/status`, {
      method: 'PATCH',
      headers: getAuthHeaders(),
      body: JSON.stringify(body)
    })
    return handleResponse<Submission>(response)
  }

  // Send test email
  async sendTestEmail(email: string): Promise<{ message: string }> {
    const response = await fetch(`${API_URL}/admin/registry/test-email`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({ email })
    })
    return handleResponse<{ message: string }>(response)
  }

  // Get example documents
  async getExampleDocuments(): Promise<{ files: any[] }> {
    const response = await fetch(`${API_URL}/registry/example-documents`, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<{ files: any[] }>(response)
  }

  // Get download link for example document
  async getDocumentDownloadLink(filePath: string): Promise<{ link: string }> {
    const params = new URLSearchParams({ path: filePath })
    const response = await fetch(`${API_URL}/registry/example-documents/download?${params.toString()}`, {
      method: 'GET',
      headers: getAuthHeaders()
    })
    return handleResponse<{ link: string }>(response)
  }
}

export const registryService = new RegistryService()

// Common API response types

export interface ApiError {
  error: string
}

export interface ApiSuccess {
  message: string
}

export interface PaginationParams {
  limit?: number
  skip?: number
}

export interface ApiResponse<T> {
  data?: T
  error?: string
  status: number
}


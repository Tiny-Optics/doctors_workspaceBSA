import { useAuthStore } from '@/stores/auth'
import router from '@/router'

/**
 * Decodes a JWT token and returns the payload
 * JWT format: header.payload.signature
 * Uses base64url encoding
 */
function decodeJWT(token: string): any {
  try {
    const parts = token.split('.')
    if (parts.length !== 3 || !parts[1]) {
      return null
    }
    
    // Decode base64url (JWT uses base64url, not standard base64)
    const payload = parts[1]
    // Convert base64url to base64
    let base64 = payload.replace(/-/g, '+').replace(/_/g, '/')
    // Add padding if needed
    while (base64.length % 4) {
      base64 += '='
    }
    
    // Decode base64
    const decoded = atob(base64)
    // Parse JSON
    return JSON.parse(decoded)
  } catch (error) {
    console.error('Failed to decode JWT:', error)
    return null
  }
}

/**
 * Checks if a JWT token is expired
 */
export function isTokenExpired(token: string | null): boolean {
  if (!token) {
    return true
  }

  const decoded = decodeJWT(token)
  if (!decoded || !decoded.exp) {
    return true
  }

  // exp is in seconds, Date.now() is in milliseconds
  const expirationTime = decoded.exp * 1000
  const currentTime = Date.now()
  
  // Add 5 minute buffer to account for clock skew and network delays
  const buffer = 5 * 60 * 1000
  
  return currentTime >= (expirationTime - buffer)
}

/**
 * Centralized API service that handles authentication and error responses
 */
export class ApiService {
  /**
   * Makes an authenticated API request
   */
  static async request(
    url: string,
    options: RequestInit = {}
  ): Promise<Response> {
    const authStore = useAuthStore()
    
    // Check if token is expired before making the request
    if (isTokenExpired(authStore.token)) {
      console.warn('Token expired, logging out...')
      await authStore.logout()
      router.push({ name: 'login' })
      throw new Error('Session expired. Please login again.')
    }

    // Add Authorization header if token exists
    const headers = new Headers(options.headers)
    if (authStore.token) {
      headers.set('Authorization', `Bearer ${authStore.token}`)
    }
    
    // Ensure Content-Type is set for POST/PUT/PATCH requests
    if (options.method && ['POST', 'PUT', 'PATCH'].includes(options.method)) {
      if (!headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json')
      }
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    // Handle 401 Unauthorized - token expired or invalid
    if (response.status === 401) {
      console.warn('Received 401 Unauthorized, logging out...')
      
      // Only logout if we have a token (to avoid infinite loops)
      if (authStore.token) {
        await authStore.logout()
        router.push({ name: 'login' })
      }
      
      throw new Error('Session expired. Please login again.')
    }

    return response
  }

  /**
   * Convenience method for GET requests
   */
  static async get(url: string, options: RequestInit = {}): Promise<Response> {
    return this.request(url, { ...options, method: 'GET' })
  }

  /**
   * Convenience method for POST requests
   */
  static async post(url: string, body?: any, options: RequestInit = {}): Promise<Response> {
    return this.request(url, {
      ...options,
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  /**
   * Convenience method for PUT requests
   */
  static async put(url: string, body?: any, options: RequestInit = {}): Promise<Response> {
    return this.request(url, {
      ...options,
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  /**
   * Convenience method for PATCH requests
   */
  static async patch(url: string, body?: any, options: RequestInit = {}): Promise<Response> {
    return this.request(url, {
      ...options,
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  /**
   * Convenience method for DELETE requests
   */
  static async delete(url: string, options: RequestInit = {}): Promise<Response> {
    return this.request(url, { ...options, method: 'DELETE' })
  }
}


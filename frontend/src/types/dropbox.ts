// Dropbox configuration types

export interface DropboxStatus {
  isConnected: boolean
  tokenExpiry: string
  lastRefreshSuccess: string
  lastRefreshAttempt: string
  consecutiveFailures: number
  lastError: string
  needsReconnection: boolean
  parentFolder: string
  configured?: boolean
  message?: string
}

export interface InitiateAuthRequest {
  appKey: string
  appSecret: string
  parentFolder: string
  redirectUri?: string
}

export interface InitiateAuthResponse {
  authUrl: string
  message: string
  instructions?: string
}

export interface CompleteAuthRequest {
  code: string
  appKey: string
  appSecret: string
  parentFolder: string
  redirectUri?: string
}

export interface CompleteAuthResponse {
  message: string
  status: DropboxStatus
}

export interface TestConnectionResponse {
  success: boolean
  message: string
}

export interface RefreshTokenResponse {
  message: string
}

export interface DeleteConfigResponse {
  message: string
}


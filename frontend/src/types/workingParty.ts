export interface WorkingPartyCategory {
  id: string
  name: string
  slug: string
  description?: string
  imagePath?: string
  dropboxPath: string
  displayOrder: number
  isActive: boolean
  createdAt: string
  updatedAt: string
  createdBy?: string
}

export interface WorkingPartyFile {
  name: string
  path: string
  size: number
  modifiedTime: string
  isFolder: boolean
  children?: WorkingPartyFile[]
}

export interface CreateCategoryRequest {
  name: string
  description?: string
  imagePath?: string
  displayOrder: number
}

export interface UpdateCategoryRequest {
  name?: string
  description?: string
  imagePath?: string
  displayOrder?: number
  isActive?: boolean
}

export interface CategoryListResponse {
  categories: WorkingPartyCategory[]
  total: number
  page: number
  limit: number
}

export interface FileListResponse {
  files: WorkingPartyFile[]
}

export interface DownloadLinkResponse {
  downloadLink: string
}

export interface ImageUploadResponse {
  imagePath: string
  message: string
}

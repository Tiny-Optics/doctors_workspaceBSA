export type InstitutionType = 
  | 'university'
  | 'hospital'
  | 'laboratory'
  | 'research_center'
  | 'government'
  | 'private_practice'
  | 'ngo'
  | 'other'

export interface Institution {
  id: string
  name: string
  shortName?: string
  type: InstitutionType
  country: string
  province?: string
  city: string
  address?: string
  postalCode?: string
  phone?: string
  email?: string
  website?: string
  imagePath?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
  createdBy?: string
}

export interface CreateInstitutionRequest {
  name: string
  shortName?: string
  type: InstitutionType
  country: string
  province?: string
  city: string
  address?: string
  postalCode?: string
  phone?: string
  email?: string
  website?: string
  imagePath?: string
}

export interface UpdateInstitutionRequest {
  name?: string
  shortName?: string
  type?: InstitutionType
  country?: string
  province?: string
  city?: string
  address?: string
  postalCode?: string
  phone?: string
  email?: string
  website?: string
  imagePath?: string
  isActive?: bool
}

export interface InstitutionsListResponse {
  institutions: Institution[]
  total: number
  limit: number
  skip: number
}

export function getInstitutionTypeDisplayName(type: InstitutionType): string {
  const displayNames: Record<InstitutionType, string> = {
    university: 'University',
    hospital: 'Hospital',
    laboratory: 'Laboratory',
    research_center: 'Research Center',
    government: 'Government',
    private_practice: 'Private Practice',
    ngo: 'NGO',
    other: 'Other'
  }
  return displayNames[type] || type
}


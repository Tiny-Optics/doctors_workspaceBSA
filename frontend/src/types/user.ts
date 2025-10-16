export interface User {
  id: string
  username: string
  email: string
  role: UserRole
  adminLevel?: AdminLevel
  isActive: boolean
  profile: UserProfile
  createdAt: string
  updatedAt: string
  lastLoginAt?: string
  createdBy?: string
}

export interface UserProfile {
  firstName: string
  lastName: string
  institutionId?: string
  specialty?: string
  registrationNumber?: string
  phoneNumber?: string
}

export type UserRole = 'haematologist' | 'physician' | 'data_capturer' | 'admin'
export type AdminLevel = 'user_manager' | 'super_admin'

export type Permission =
  | 'view_sops'
  | 'download_sops'
  | 'access_referrals'
  | 'view_registry'
  | 'upload_ethics_approval'
  | 'manage_users'
  | 'assign_roles'
  | 'view_audit_logs'
  | 'manage_system'
  | 'delete_users'

export interface CreateUserRequest {
  username: string
  email: string
  password: string
  role: UserRole
  adminLevel?: AdminLevel
  firstName: string
  lastName: string
  institutionId: string
  specialty?: string
  registrationNumber?: string
  phoneNumber?: string
}

export interface UpdateUserRequest {
  firstName?: string
  lastName?: string
  institutionId?: string
  specialty?: string
  registrationNumber?: string
  phoneNumber?: string
  role?: UserRole
  adminLevel?: AdminLevel
  isActive?: boolean
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  refreshToken: string
  user: User
  expiresAt: string
}

export interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
}

export interface UsersListResponse {
  users: User[]
  total: number
  limit: number
  skip: number
}

// Helper function to get full name from user
export function getUserFullName(user: User): string {
  return `${user.profile.firstName} ${user.profile.lastName}`
}

// Helper function to get user role display name
export function getUserRoleDisplayName(role: UserRole): string {
  const roleMap: Record<UserRole, string> = {
    haematologist: 'Haematologist',
    physician: 'Physician',
    data_capturer: 'Data Capturer',
    admin: 'Administrator'
  }
  return roleMap[role]
}

// Helper function to get admin level display name
export function getAdminLevelDisplayName(adminLevel?: AdminLevel): string {
  if (!adminLevel) return ''
  const levelMap: Record<AdminLevel, string> = {
    user_manager: 'User Manager',
    super_admin: 'Super Admin'
  }
  return levelMap[adminLevel]
}

// Permission matrix helper
export function getPermissionsForRole(role: UserRole, adminLevel?: AdminLevel): Permission[] {
  const clinicalPermissions: Permission[] = [
    'view_sops',
    'download_sops',
    'access_referrals',
    'view_registry',
    'upload_ethics_approval'
  ]

  if (role !== 'admin') {
    return clinicalPermissions
  }

  switch (adminLevel) {
    case 'user_manager':
      return [...clinicalPermissions, 'manage_users', 'assign_roles']
    case 'super_admin':
      return [
        ...clinicalPermissions,
        'manage_users',
        'assign_roles',
        'view_audit_logs',
        'manage_system',
        'delete_users'
      ]
    default:
      return clinicalPermissions
  }
}

// Check if user has a specific permission
export function hasPermission(user: User, permission: Permission): boolean {
  const permissions = getPermissionsForRole(user.role, user.adminLevel)
  return permissions.includes(permission)
}

// Check if user can manage another user
export function canManageUser(currentUser: User, targetUser: User): boolean {
  if (!hasPermission(currentUser, 'manage_users')) {
    return false
  }

  if (currentUser.adminLevel === 'super_admin') {
    return true
  }

  if (currentUser.adminLevel === 'user_manager') {
    return targetUser.role !== 'admin'
  }

  return false
}


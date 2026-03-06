// ── Project Types ────────────────────────────────────────────────────────────

export type ProjectRole = 'admin' | 'member' | 'viewer'

export interface Project {
    id: string
    organizationId: string
    name: string
    slug: string
    description: string
    isActive: boolean
    createdBy: string
    maxApiKeys: number
    maxMembers: number
    memberCount: number
    createdAt: string
}

export interface ProjectWithMembers extends Project {
    members: ProjectMember[]
}

export interface ProjectMember {
    id: string
    userId: string
    email: string
    fullName: string
    role: ProjectRole
    isActive: boolean
    joinedAt: string
    removedAt?: string | null
}

// ── API Key Types ────────────────────────────────────────────────────────────

export type ApiKeyMode = 'dev' | 'live'

export interface ApiKey {
    id: string
    projectId: string
    name: string
    keyPrefix: string
    mode: ApiKeyMode
    createdBy: string
    expiresAt?: string | null
    lastUsedAt?: string | null
    isRevoked: boolean
    revokedAt?: string | null
    requestLimit: number
    rateLimit: number
    createdAt: string
}

export interface ApiKeyCreated extends ApiKey {
    rawKey: string
}

export interface OrganizationMember {
    id: string
    userId: string
    email: string
    fullName: string
    isActive: boolean
    joinedAt: string
    invitedBy?: string | null
    invitedAt?: string | null
    deactivatedBy?: string | null
    deactivatedAt?: string | null
}

/** OrganizationDTO — returned by list endpoint */
export interface Organization {
    id: string
    name: string
    slug: string
    description: string
    ownerId: string
    isActive: boolean
    maxProjects: number
    maxMembers: number
    memberCount: number
    deactivatedBy?: string | null
    deactivatedAt?: string | null
    createdAt: string
}

/** OrganizationResponseDTO — returned by create / get-by-id */
export interface OrganizationWithMembers extends Organization {
    members: OrganizationMember[]
}

export interface OrganizationMemberPermission {
    id: string
    code: string
    resource: string
    action: string
    group: string
    isSystem: boolean
}

export interface OrganizationMemberRole {
    id: string,
    name: string,
    description: string,
    isSystemRole: true,
    isDefault: false
}

/** OrganizationMemberWithRBACDTO */
export interface OrganizationMemberWithRBAC {
    id: string
    userId: string
    email: string
    fullName: string
    isActive: boolean
    roles: OrganizationMemberRole[]
    permissions: OrganizationMemberPermission[]
}

/** OrganizationMemberInvitationDTO */
export interface OrganizationMemberInvitation {
    id: string
    organizationId: string
    invitedBy: string
    email: string
    roleIds: string[]
    invitedAt: string
    isAccepted: boolean
    acceptedAt?: string | null
    expiredAt: string
    isPending: boolean
}

/** OrganizationRole service struct */
export interface OrganizationRole {
    id: string
    name: string
    description: string
    isSystem: boolean
    isDefault: boolean
}

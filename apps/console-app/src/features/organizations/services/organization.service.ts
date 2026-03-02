import { api } from '@/lib/api'
import type {
    Organization,
    OrganizationWithMembers,
    OrganizationMemberWithRBAC,
    OrganizationMemberInvitation,
    OrganizationRole,
    OrganizationMember,
} from '../types'

// ─── Payloads ─────────────────────────────────────────────────────────────────

export interface CreateOrganizationPayload {
    name: string
    description?: string
}

export interface UpdateOrganizationPayload {
    name?: string
    description?: string
    slug?: string
}

export interface CreateInvitationPayload {
    email: string
    roleIds?: string[]
}

export interface ChangeRolePayload {
    roleId: string
}

export interface TransferOwnershipPayload {
    newOwnerMemberId: string
}

export interface DeleteOrganizationPayload {
    confirmName: string
}

// ─── Organization API ─────────────────────────────────────────────────────────

export async function getOrganizations(): Promise<Organization[]> {
    return api.get('/organizations')
}

export async function getOrganizationBySlug(slug: string): Promise<OrganizationWithMembers> {
    return api.get(`/organizations/slug/${slug}`)
}

export async function getOrganizationById(id: string): Promise<OrganizationWithMembers> {
    return api.get(`/organizations/${id}`)
}

export async function createOrganization(payload: CreateOrganizationPayload): Promise<OrganizationWithMembers> {
    return api.post('/organizations', payload)
}

export async function onboardOrganization(payload: CreateOrganizationPayload): Promise<OrganizationWithMembers> {
    return api.post('/organizations/onboarding', payload)
}

export async function updateOrganization(
    orgId: string,
    payload: UpdateOrganizationPayload,
): Promise<OrganizationWithMembers> {
    return api.patch(`/organizations/${orgId}`, payload)
}

export async function archiveOrganization(orgId: string): Promise<void> {
    return api.post(`/organizations/${orgId}/archive`)
}

export async function restoreOrganization(orgId: string): Promise<void> {
    return api.post(`/organizations/${orgId}/restore`)
}

export async function deleteOrganization(orgId: string, payload: DeleteOrganizationPayload): Promise<void> {
    return api.delete(`/organizations/${orgId}`, payload)
}

// ─── Member API ───────────────────────────────────────────────────────────────

export async function getMyMembershipInOrg(orgId: string): Promise<OrganizationMemberWithRBAC> {
    return api.get(`/organizations/${orgId}/members/me`)
}

export async function getMembers(orgId: string): Promise<OrganizationMember[]> {
    return api.get(`/organizations/${orgId}/members`)
}

export async function removeMember(orgId: string, memberId: string): Promise<void> {
    return api.delete(`/organizations/${orgId}/members/${memberId}`)
}

export async function changeMemberRole(orgId: string, memberId: string, payload: ChangeRolePayload): Promise<void> {
    return api.patch(`/organizations/${orgId}/members/${memberId}/role`, payload)
}

export async function leaveOrganization(orgId: string): Promise<void> {
    return api.post(`/organizations/${orgId}/members/leave`)
}

export async function transferOwnership(orgId: string, payload: TransferOwnershipPayload): Promise<void> {
    return api.post(`/organizations/${orgId}/members/transfer-ownership`, payload)
}

// ─── Invitation API ───────────────────────────────────────────────────────────

export async function getInvitations(orgId: string): Promise<OrganizationMemberInvitation[]> {
    return api.get(`/organizations/${orgId}/invitations`)
}

export async function createInvitation(
    orgId: string,
    payload: CreateInvitationPayload,
): Promise<OrganizationMemberInvitation> {
    return api.post(`/organizations/${orgId}/invitations`, payload)
}

export async function resendInvitation(orgId: string, invitationId: string): Promise<OrganizationMemberInvitation> {
    return api.post(`/organizations/${orgId}/invitations/${invitationId}/resend`)
}

export async function cancelInvitation(orgId: string, invitationId: string): Promise<void> {
    return api.delete(`/organizations/${orgId}/invitations/${invitationId}`)
}

export async function acceptInvitation(token: string): Promise<void> {
    return api.get(`/organizations/invitations/accept?token=${token}`)
}

// ─── RBAC API ─────────────────────────────────────────────────────────────────

export async function getRoles(orgId: string): Promise<OrganizationRole[]> {
    return api.get(`/organizations/${orgId}/rbac/roles`)
}

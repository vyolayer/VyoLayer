import { api } from '@/lib/api'
import type { ProjectMember, ProjectRole } from '../types'

// ─── Payloads ─────────────────────────────────────────────────────────────────

export interface AddProjectMemberPayload {
    userId: string
    role: ProjectRole
}

export interface ChangeProjectMemberRolePayload {
    role: ProjectRole
}

// ─── Project Member API ───────────────────────────────────────────────────────

const base = (orgId: string, projectId: string) => `/organizations/${orgId}/projects/${projectId}/members`

export async function listProjectMembers(orgId: string, projectId: string): Promise<ProjectMember[]> {
    return api.get(base(orgId, projectId))
}

export async function getCurrentProjectMember(orgId: string, projectId: string): Promise<ProjectMember> {
    return api.get(`${base(orgId, projectId)}/me`)
}

export async function addProjectMember(
    orgId: string,
    projectId: string,
    payload: AddProjectMemberPayload,
): Promise<ProjectMember> {
    return api.post(base(orgId, projectId), payload)
}

export async function changeProjectMemberRole(
    orgId: string,
    projectId: string,
    memberId: string,
    payload: ChangeProjectMemberRolePayload,
): Promise<void> {
    return api.patch(`${base(orgId, projectId)}/${memberId}/role`, payload)
}

export async function removeProjectMember(orgId: string, projectId: string, memberId: string): Promise<void> {
    return api.delete(`${base(orgId, projectId)}/${memberId}`)
}

export async function leaveProject(orgId: string, projectId: string): Promise<void> {
    return api.post(`${base(orgId, projectId)}/leave`)
}

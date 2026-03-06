import { api } from '@/lib/api'
import type { Project, ProjectWithMembers } from '../types'

// ─── Payloads ─────────────────────────────────────────────────────────────────

export interface CreateProjectPayload {
    name: string
    description?: string
}

export interface UpdateProjectPayload extends Partial<CreateProjectPayload> {}

export interface DeleteProjectPayload {
    confirmName: string
}

// ─── Project API ──────────────────────────────────────────────────────────────

export async function listProjects(orgId: string): Promise<Project[]> {
    return api.get(`/organizations/${orgId}/projects`)
}

export async function getProjectById(orgId: string, projectId: string): Promise<ProjectWithMembers> {
    return api.get(`/organizations/${orgId}/projects/${projectId}`)
}

export async function createProject(orgId: string, payload: CreateProjectPayload): Promise<ProjectWithMembers> {
    return api.post(`/organizations/${orgId}/projects`, payload)
}

export async function updateProject(
    orgId: string,
    projectId: string,
    payload: UpdateProjectPayload,
): Promise<ProjectWithMembers> {
    return api.patch(`/organizations/${orgId}/projects/${projectId}`, payload)
}

export async function archiveProject(orgId: string, projectId: string): Promise<void> {
    return api.post(`/organizations/${orgId}/projects/${projectId}/archive`)
}

export async function restoreProject(orgId: string, projectId: string): Promise<void> {
    return api.post(`/organizations/${orgId}/projects/${projectId}/restore`)
}

export async function deleteProject(orgId: string, projectId: string): Promise<void> {
    return api.delete(`/organizations/${orgId}/projects/${projectId}`)
}

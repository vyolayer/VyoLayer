import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
    listProjects,
    getProjectById,
    createProject,
    updateProject,
    archiveProject,
    restoreProject,
    deleteProject,
    type CreateProjectPayload,
    type UpdateProjectPayload,
} from '../services/project.service'

// ─── Query Keys ──────────────────────────────────────────────────────────────

export const projectKeys = {
    all: (orgId: string) => ['projects', orgId] as const,
    detail: (orgId: string, projectId: string) => ['projects', orgId, projectId] as const,
}

// ─── Queries ─────────────────────────────────────────────────────────────────

export function useProjects(orgId: string) {
    return useQuery({
        queryKey: projectKeys.all(orgId),
        queryFn: () => listProjects(orgId),
        enabled: !!orgId,
    })
}

export function useProject(orgId: string, projectId: string) {
    return useQuery({
        queryKey: projectKeys.detail(orgId, projectId),
        queryFn: () => getProjectById(orgId, projectId),
        enabled: !!orgId && !!projectId,
    })
}

// ─── Mutations ───────────────────────────────────────────────────────────────

export function useCreateProject(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: CreateProjectPayload) => createProject(orgId, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: projectKeys.all(orgId) })
        },
    })
}

export function useUpdateProject(orgId: string, projectId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: UpdateProjectPayload) => updateProject(orgId, projectId, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: projectKeys.all(orgId) })
            qc.invalidateQueries({ queryKey: projectKeys.detail(orgId, projectId) })
        },
    })
}

export function useArchiveProject(orgId: string, projectId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: () => archiveProject(orgId, projectId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: projectKeys.all(orgId) })
            qc.invalidateQueries({ queryKey: projectKeys.detail(orgId, projectId) })
        },
    })
}

export function useRestoreProject(orgId: string, projectId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: () => restoreProject(orgId, projectId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: projectKeys.all(orgId) })
            qc.invalidateQueries({ queryKey: projectKeys.detail(orgId, projectId) })
        },
    })
}

export function useDeleteProject(orgId: string, projectId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: () => deleteProject(orgId, projectId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: projectKeys.all(orgId) })
        },
    })
}

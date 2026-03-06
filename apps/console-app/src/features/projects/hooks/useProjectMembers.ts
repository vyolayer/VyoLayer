import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import {
    listProjectMembers,
    getCurrentProjectMember,
    addProjectMember,
    changeProjectMemberRole,
    removeProjectMember,
    leaveProject,
    type AddProjectMemberPayload,
    type ChangeProjectMemberRolePayload,
} from '../services/project-member.service'
import { projectKeys } from './useProjects'

// ─── Query Keys ──────────────────────────────────────────────────────────────

export const projectMemberKeys = {
    all: (orgId: string, projectId: string) => ['project-members', orgId, projectId] as const,
    me: (orgId: string, projectId: string) => ['project-members', orgId, projectId, 'me'] as const,
}

// ─── Queries ─────────────────────────────────────────────────────────────────

export function useProjectMembers(orgId: string, projectId: string) {
    return useQuery({
        queryKey: projectMemberKeys.all(orgId, projectId),
        queryFn: () => listProjectMembers(orgId, projectId),
        enabled: !!orgId && !!projectId,
    })
}

export function useCurrentProjectMember(orgId: string, projectId: string) {
    return useQuery({
        queryKey: projectMemberKeys.me(orgId, projectId),
        queryFn: () => getCurrentProjectMember(orgId, projectId),
        enabled: !!orgId && !!projectId,
    })
}

// ─── Mutations ───────────────────────────────────────────────────────────────

export function useAddProjectMember(orgId: string, projectId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: AddProjectMemberPayload) => addProjectMember(orgId, projectId, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: projectMemberKeys.all(orgId, projectId) })
            qc.invalidateQueries({ queryKey: projectKeys.detail(orgId, projectId) })
        },
    })
}

export function useChangeProjectMemberRole(orgId: string, projectId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: ({ memberId, payload }: { memberId: string; payload: ChangeProjectMemberRolePayload }) =>
            changeProjectMemberRole(orgId, projectId, memberId, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: projectMemberKeys.all(orgId, projectId) })
        },
    })
}

export function useRemoveProjectMember(orgId: string, projectId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (memberId: string) => removeProjectMember(orgId, projectId, memberId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: projectMemberKeys.all(orgId, projectId) })
            qc.invalidateQueries({ queryKey: projectKeys.detail(orgId, projectId) })
        },
    })
}

export function useLeaveProject(orgId: string, projectId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: () => leaveProject(orgId, projectId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: projectKeys.all(orgId) })
        },
    })
}

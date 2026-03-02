import { useMutation, useQueryClient } from '@tanstack/react-query'
import {
    updateOrganization,
    archiveOrganization,
    restoreOrganization,
    deleteOrganization,
    createInvitation,
    resendInvitation,
    cancelInvitation,
    removeMember,
    changeMemberRole,
    leaveOrganization,
    transferOwnership,
    type UpdateOrganizationPayload,
    type DeleteOrganizationPayload,
    type CreateInvitationPayload,
    type ChangeRolePayload,
    type TransferOwnershipPayload,
} from '../services/organization.service'

// ─── Organization Mutations ──────────────────────────────────────────────────

export function useUpdateOrganization(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: UpdateOrganizationPayload) => updateOrganization(orgId, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['organizations'] })
        },
    })
}

export function useArchiveOrganization(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: () => archiveOrganization(orgId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['organizations'] })
        },
    })
}

export function useRestoreOrganization(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: () => restoreOrganization(orgId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['organizations'] })
        },
    })
}

export function useDeleteOrganization(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: DeleteOrganizationPayload) => deleteOrganization(orgId, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['organizations'] })
        },
    })
}

// ─── Member Mutations ────────────────────────────────────────────────────────

export function useRemoveMember(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (memberId: string) => removeMember(orgId, memberId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['org-members', orgId] })
            qc.invalidateQueries({ queryKey: ['organizations'] })
        },
    })
}

export function useChangeMemberRole(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: ({ memberId, payload }: { memberId: string; payload: ChangeRolePayload }) =>
            changeMemberRole(orgId, memberId, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['org-members', orgId] })
            qc.invalidateQueries({ queryKey: ['organizations'] })
        },
    })
}

export function useLeaveOrganization(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: () => leaveOrganization(orgId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['organizations'] })
        },
    })
}

export function useTransferOwnership(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: TransferOwnershipPayload) => transferOwnership(orgId, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['organizations'] })
            qc.invalidateQueries({ queryKey: ['org-members', orgId] })
        },
    })
}

// ─── Invitation Mutations ────────────────────────────────────────────────────

export function useInviteMember(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (payload: CreateInvitationPayload) => createInvitation(orgId, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['org-invitations', orgId] })
        },
    })
}

export function useResendInvitation(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (invitationId: string) => resendInvitation(orgId, invitationId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['org-invitations', orgId] })
        },
    })
}

export function useCancelInvitation(orgId: string) {
    const qc = useQueryClient()
    return useMutation({
        mutationFn: (invitationId: string) => cancelInvitation(orgId, invitationId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['org-invitations', orgId] })
        },
    })
}

import { useCallback, useState } from 'react'

import {
    OrganizationMemberList,
    OrganizationMembersInvite,
    OrganizationInvitationPending,
} from '@/org/components/members'
import { Separator } from '@/components/ui/separator'
import { ConfirmDialog } from '@/components/organization/ConfirmDialog'

import { useOrganization } from '@/contexts/OrganizationContext'
import { useRemoveMember } from '@/org/hooks/useOrganizationMutations'
import { useOrganizationMembers } from '@/org/hooks/useOrganizationMembers'

export function OrganizationMembersPage() {
    const { org, me, isAdmin } = useOrganization()
    const { data: members = org.members } = useOrganizationMembers(org.id)

    const [removeMemberId, setRemoveMemberId] = useState<string | null>(null)
    const removeMemberName = members.find((m) => m.id === removeMemberId)?.fullName ?? ''
    const removeMutation = useRemoveMember(org.id)

    const handleRemoveMember = useCallback(
        (memberId: string | null) => {
            if (!memberId) return

            removeMutation.mutate(memberId, {
                onSuccess: () => setRemoveMemberId(null),
            })
        },
        [removeMutation],
    )

    return (
        <div className="space-y-8">
            {/* ── Invite Member ── */}
            {isAdmin && <OrganizationMembersInvite orgId={org.id} isAdmin={isAdmin} />}
            {/* ── Pending Invitations ── */}
            {isAdmin && <OrganizationInvitationPending orgId={org.id} isAdmin={isAdmin} />}
            <Separator />
            {/* ── Member List ── */}
            <OrganizationMemberList
                members={members}
                org={org}
                isAdmin={isAdmin}
                currentUserId={me?.userId}
                onMemberRemove={setRemoveMemberId}
            />
            {/* ── Remove Member Dialog ── */}
            <ConfirmDialog
                open={!!removeMemberId}
                onOpenChange={() => setRemoveMemberId(null)}
                title="Remove Member"
                description={`Are you sure you want to remove ${removeMemberName} from this organization? They will lose all access immediately.`}
                confirmLabel="Remove"
                isLoading={removeMutation.isPending}
                onConfirm={() => handleRemoveMember(removeMemberId)}
            />
        </div>
    )
}

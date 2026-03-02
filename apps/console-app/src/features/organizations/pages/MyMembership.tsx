import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'

import { Separator } from '@/components/ui/separator'

import { useOrganization } from '@/contexts/OrganizationContext'
import { useLeaveOrganization } from '@/org/hooks/useOrganizationMutations'
import {
    OrganizationLeaveCard,
    OrganizationLeaveDialog,
    OrganizationMemberProfileAsMeCard,
} from '../components/members'

export function OrganizationMyMembershipPage() {
    const { org, me, isOwner } = useOrganization()
    const navigate = useNavigate()
    const [showLeaveDialog, setShowLeaveDialog] = useState(false)
    const leaveMutation = useLeaveOrganization(org.id)

    const handleLeave = () => {
        leaveMutation.mutate(undefined, {
            onSuccess: () => {
                setShowLeaveDialog(false)
                navigate({ to: '/' })
            },
        })
    }

    return (
        <div className="space-y-8">
            {/* ── My Profile in this Org ── */}
            <OrganizationMemberProfileAsMeCard me={me} organizationName={org.name} />

            <Separator />

            {/* ── Leave Organization ── */}
            <OrganizationLeaveCard
                organizationName={org.name}
                isOwner={isOwner}
                onLeave={() => setShowLeaveDialog(true)}
            />

            <OrganizationLeaveDialog
                open={showLeaveDialog}
                onOpenChange={setShowLeaveDialog}
                organizationName={org.name}
                isPending={leaveMutation.isPending}
                onLeave={handleLeave}
            />
        </div>
    )
}

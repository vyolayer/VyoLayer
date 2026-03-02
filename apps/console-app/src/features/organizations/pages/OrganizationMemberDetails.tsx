import { Link, useParams } from '@tanstack/react-router'
import { useOrganization } from '@/contexts/OrganizationContext'
import { useOrganizationMembers } from '@/features/organizations/hooks/useOrganizationMembers'
import { useOrganizationRoles } from '@/features/organizations/hooks/useOrganizationRoles'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { ArrowLeft } from 'lucide-react'
import {
    OrganizationMemberProfile,
    OrganizationMemberRemove,
    OrganizationMemberProfileRoleManagement,
} from '../components/members'

export const OrganizationMemberDetailPage = () => {
    const { id: memberId } = useParams({ from: '/_authenticated/org/$slug/_layout/members/$id' })
    const { org, me, isAdmin, isOwner } = useOrganization()
    const { data: members = org.members } = useOrganizationMembers(org.id)
    const { data: roles = [] } = useOrganizationRoles(isAdmin ? org.id : undefined)

    const member = members.find((m) => m.id === memberId || m.id === `org_member_${memberId}`)
    const isMemberOwner = member?.userId === org.ownerId
    const isSelf = member?.userId === me?.userId
    const canChangeRole = isAdmin && !isSelf && !(isMemberOwner && !isOwner)
    const canRemove = isAdmin && !isMemberOwner && !isSelf

    if (!member) return <NoMemberFound slug={org.slug} />

    const isCanRoleManagement = canChangeRole && roles.length > 0

    return (
        <div className="space-y-6">
            <Button variant="ghost" size="sm" asChild>
                <Link to="/org/$slug/members" params={{ slug: org.slug }}>
                    <ArrowLeft className="h-4 w-4 mr-2" />
                    All Members
                </Link>
            </Button>

            {/* ── Member Profile ── */}
            <OrganizationMemberProfile member={member} isMemberOwner={isMemberOwner} isSelf={isSelf} />

            {/* ── Role Management ── */}
            {isCanRoleManagement && (
                <OrganizationMemberProfileRoleManagement
                    member={member}
                    roles={roles}
                    isOwner={isOwner}
                    orgId={org.id}
                />
            )}

            {/* ── Remove Member ── */}
            {canRemove && (
                <>
                    <Separator />
                    <OrganizationMemberRemove orgId={org.id} member={member} />
                </>
            )}
        </div>
    )
}

const NoMemberFound: React.FC<{ slug: string }> = ({ slug }) => {
    return (
        <div className="py-20 text-center space-y-4">
            <p className="text-muted-foreground">Member not found.</p>
            <Button variant="outline" asChild>
                <Link to="/org/$slug/members" params={{ slug }}>
                    <ArrowLeft className="h-4 w-4 mr-2" />
                    Back to members
                </Link>
            </Button>
        </div>
    )
}

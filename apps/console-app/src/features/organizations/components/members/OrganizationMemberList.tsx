import * as React from 'react'
import type { Organization, OrganizationMember } from '@/org/types'

import { formatDate } from '@/lib/utils'
import { UserMinus } from 'lucide-react'

import { Link } from '@tanstack/react-router'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { RoleBadge } from '@/org/components/RoleBadge'
import { Card, CardContent } from '@/components/ui/card'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'

export const OrganizationMemberList: React.FC<{
    members: OrganizationMember[]
    isAdmin: boolean
    org: Organization
    currentUserId: string
    onMemberRemove: (memberId: string) => void
}> = ({ members, org, isAdmin, currentUserId, onMemberRemove }) => {
    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <h2 className="text-lg font-semibold">
                    Members <span className="text-muted-foreground font-normal">({members.length})</span>
                </h2>
            </div>

            <Card>
                <CardContent className="p-0">
                    <ul className="divide-y divide-border">
                        {members.map((member) => (
                            <MemberDetails
                                key={member.id}
                                member={member}
                                org={org}
                                isAdmin={isAdmin}
                                currentUserId={currentUserId}
                                onMemberRemove={onMemberRemove}
                            />
                        ))}
                    </ul>
                </CardContent>
            </Card>
        </div>
    )
}

const MemberDetails = ({
    member,
    org,
    isAdmin,
    currentUserId,
    onMemberRemove,
}: {
    member: OrganizationMember
    org: Organization
    isAdmin: boolean
    currentUserId: string
    onMemberRemove: (memberId: string) => void
}) => {
    const isMemberOwner = member.userId === org.ownerId
    const isSelf = member.userId === currentUserId
    const canRemove = isAdmin && !isMemberOwner && !isSelf
    const memberId = member.id.split('org_member_')[1]
    const memberFallbackChar = member.fullName?.charAt(0)?.toUpperCase() ?? '?'

    return (
        <li key={member.id} className="flex items-center justify-between px-4 py-3 gap-3">
            <Link
                to="/org/$slug/members/$id"
                params={{ slug: org.slug, id: memberId }}
                className="flex items-center gap-3 min-w-0 flex-1 group"
            >
                <Avatar className="h-8 w-8 shrink-0">
                    <AvatarFallback className="text-xs font-medium bg-primary/10 text-primary">
                        {memberFallbackChar}
                    </AvatarFallback>
                </Avatar>
                <div className="min-w-0">
                    <div className="flex items-center gap-2">
                        <p className="text-sm font-medium truncate group-hover:text-primary transition-colors">
                            {member.fullName}
                        </p>
                        {isMemberOwner && <RoleBadge name="Owner" />}
                        {isSelf && (
                            <Badge variant="outline" className="text-xs">
                                You
                            </Badge>
                        )}
                    </div>
                    <p className="text-xs text-muted-foreground truncate">{member.email}</p>
                </div>
            </Link>

            <div className="flex items-center gap-2 shrink-0">
                <span className="text-xs text-muted-foreground hidden sm:inline">{formatDate(member.joinedAt)}</span>
                {canRemove && (
                    <Button
                        variant="ghost"
                        size="sm"
                        className="text-destructive hover:text-destructive"
                        onClick={() => onMemberRemove(member.id)}
                    >
                        <UserMinus className="w-3.5 h-3.5" />
                    </Button>
                )}
            </div>
        </li>
    )
}

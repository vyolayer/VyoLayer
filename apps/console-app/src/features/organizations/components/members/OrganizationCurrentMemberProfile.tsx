import * as React from 'react'
import { Mail, Shield } from 'lucide-react'

import { RoleBadge } from '@/org/components/RoleBadge'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import type { OrganizationMemberPermission, OrganizationMemberWithRBAC } from '@/org/types'
import { Badge } from '@/components/ui/badge'

export const OrganizationMemberProfileAsMeCard: React.FC<{
    me: OrganizationMemberWithRBAC
    organizationName: string
}> = ({ me, organizationName }) => {
    return (
        <Card>
            <CardHeader>
                <CardTitle>My Membership</CardTitle>
                <CardDescription>Your profile and status within {organizationName}.</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="flex items-start gap-4">
                    <Avatar className="h-16 w-16 shrink-0">
                        <AvatarFallback className="text-lg font-semibold bg-primary/10 text-primary">
                            {me?.fullName?.charAt(0)?.toUpperCase() ?? '?'}
                        </AvatarFallback>
                    </Avatar>
                    <div className="flex-1 space-y-3">
                        <div>
                            <h3 className="text-lg font-semibold">{me?.fullName}</h3>
                            <div className="flex items-center gap-2 text-sm text-muted-foreground mt-1">
                                <Mail className="h-4 w-4" />
                                {me?.email}
                            </div>
                        </div>

                        {/* Roles */}
                        <div className="space-y-2">
                            <p className="text-sm font-medium">Roles</p>
                            <div className="flex flex-wrap gap-2">
                                {me?.roles?.map((role) => (
                                    <RoleBadge key={role.id} name={role.name} size="md" />
                                ))}
                            </div>
                        </div>

                        {/* Permissions */}
                        {me?.permissions && me.permissions.length > 0 && (
                            <OrganizationMemberPermission permissions={me.permissions} />
                        )}

                        {/* Dates */}
                        {/* {myMember && (
                            <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                <Calendar className="h-4 w-4" />
                                Joined {formatDate(myMember.joinedAt)}
                            </div>
                        )} */}
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}

const OrganizationMemberPermission: React.FC<{
    permissions: OrganizationMemberPermission[]
}> = ({ permissions }) => {
    return (
        <div className="space-y-2">
            <p className="text-sm font-medium">Permissions</p>
            <div className="flex flex-wrap gap-1.5">
                {permissions.map((perm) => (
                    <Badge variant="outline" key={perm.id}>
                        <Shield className="size-2.5" />
                        <span>{perm.code}</span>
                    </Badge>
                ))}
            </div>
        </div>
    )
}

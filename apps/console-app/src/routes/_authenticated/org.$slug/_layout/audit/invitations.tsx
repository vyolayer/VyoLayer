import { createFileRoute } from '@tanstack/react-router'
import { useOrganization } from '@/contexts/OrganizationContext'
import { useOrganizationInvitations } from '@/features/organizations/hooks/useOrganizationInvitations'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Mail, CheckCircle2, Clock, XCircle } from 'lucide-react'
import { formatDate } from '@/lib/utils'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/audit/invitations')({
    component: AuditInvitationsPage,
})

function AuditInvitationsPage() {
    const { org } = useOrganization()
    const { data: invitations = [], isLoading } = useOrganizationInvitations(org.id)

    if (isLoading) {
        return (
            <Card>
                <CardContent className="py-10">
                    <div className="space-y-3">
                        {[1, 2, 3].map((n) => (
                            <div key={n} className="h-12 animate-pulse bg-muted rounded-lg" />
                        ))}
                    </div>
                </CardContent>
            </Card>
        )
    }

    if (invitations.length === 0) {
        return (
            <Card>
                <CardContent className="py-16 text-center space-y-3">
                    <Mail className="h-10 w-10 text-muted-foreground mx-auto" />
                    <div>
                        <p className="text-sm font-medium">No Invitations</p>
                        <p className="text-sm text-muted-foreground mt-1">
                            No invitations have been sent for this organization yet.
                        </p>
                    </div>
                </CardContent>
            </Card>
        )
    }

    const pending = invitations.filter((i) => i.isPending)
    const accepted = invitations.filter((i) => i.isAccepted)
    const expired = invitations.filter((i) => !i.isPending && !i.isAccepted)

    return (
        <div className="space-y-6">
            {pending.length > 0 && (
                <Card>
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm flex items-center gap-2">
                            <Clock className="h-4 w-4 text-amber-500" />
                            Pending ({pending.length})
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="p-0">
                        <ul className="divide-y divide-border">
                            {pending.map((inv) => (
                                <li key={inv.id} className="flex items-center justify-between px-4 py-3">
                                    <div className="flex items-center gap-3">
                                        <Mail className="w-4 h-4 text-muted-foreground" />
                                        <div>
                                            <p className="text-sm font-medium">{inv.email}</p>
                                            <p className="text-xs text-muted-foreground">
                                                Invited {formatDate(inv.invitedAt)} · Expires{' '}
                                                {formatDate(inv.expiredAt)}
                                            </p>
                                        </div>
                                    </div>
                                    <Badge variant="outline" className="text-amber-600 border-amber-500/30">
                                        Pending
                                    </Badge>
                                </li>
                            ))}
                        </ul>
                    </CardContent>
                </Card>
            )}

            {accepted.length > 0 && (
                <Card>
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm flex items-center gap-2">
                            <CheckCircle2 className="h-4 w-4 text-emerald-500" />
                            Accepted ({accepted.length})
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="p-0">
                        <ul className="divide-y divide-border">
                            {accepted.map((inv) => (
                                <li key={inv.id} className="flex items-center justify-between px-4 py-3">
                                    <div className="flex items-center gap-3">
                                        <Mail className="w-4 h-4 text-muted-foreground" />
                                        <div>
                                            <p className="text-sm font-medium">{inv.email}</p>
                                            <p className="text-xs text-muted-foreground">
                                                Accepted {inv.acceptedAt ? formatDate(inv.acceptedAt) : 'N/A'}
                                            </p>
                                        </div>
                                    </div>
                                    <Badge variant="outline" className="text-emerald-600 border-emerald-500/30">
                                        Accepted
                                    </Badge>
                                </li>
                            ))}
                        </ul>
                    </CardContent>
                </Card>
            )}

            {expired.length > 0 && (
                <Card>
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm flex items-center gap-2">
                            <XCircle className="h-4 w-4 text-muted-foreground" />
                            Expired / Cancelled ({expired.length})
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="p-0">
                        <ul className="divide-y divide-border">
                            {expired.map((inv) => (
                                <li key={inv.id} className="flex items-center justify-between px-4 py-3">
                                    <div className="flex items-center gap-3">
                                        <Mail className="w-4 h-4 text-muted-foreground" />
                                        <div>
                                            <p className="text-sm font-medium text-muted-foreground">{inv.email}</p>
                                            <p className="text-xs text-muted-foreground">
                                                Expired {formatDate(inv.expiredAt)}
                                            </p>
                                        </div>
                                    </div>
                                    <Badge variant="outline" className="text-muted-foreground">
                                        Expired
                                    </Badge>
                                </li>
                            ))}
                        </ul>
                    </CardContent>
                </Card>
            )}
        </div>
    )
}

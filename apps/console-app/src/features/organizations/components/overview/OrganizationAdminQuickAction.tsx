import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { useOrganization } from '@/contexts/OrganizationContext'
import { Link } from '@tanstack/react-router'

export const OrganizationAdminQuickAction = () => {
    const {
        org: { slug },
        isAdmin,
    } = useOrganization()

    if (!isAdmin) return null

    return (
        <>
            <Separator />
            <div className="space-y-3">
                <h2 className="text-lg font-semibold">Quick Actions</h2>
                <div className="flex flex-wrap gap-2">
                    <Button variant="outline" size="sm" asChild>
                        <Link to={`/org/$slug/members`} params={{ slug }}>
                            Invite Members
                        </Link>
                    </Button>
                    <Button variant="outline" size="sm" asChild>
                        <Link to={`/org/$slug/settings`} params={{ slug }}>
                            Organization Settings
                        </Link>
                    </Button>
                </div>
            </div>
        </>
    )
}

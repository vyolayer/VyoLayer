import { LogOutIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { ConfirmDialog } from '@/components/organization/ConfirmDialog'

const Descriptions = {
    owner: 'As the owner, you must transfer ownership before you can leave.',
    member: 'You will lose all access to this organization and its resources.',
}

export const OrganizationLeaveCard: React.FC<{
    organizationName: string
    isOwner: boolean
    onLeave: () => void
}> = ({ organizationName, isOwner, onLeave }) => {
    const access = isOwner ? 'owner' : 'member'

    return (
        <Card className="border-destructive/30">
            <CardHeader>
                <CardTitle className="text-base text-destructive">Leave Organization</CardTitle>
                <CardDescription>{Descriptions[access]}</CardDescription>
            </CardHeader>
            <CardContent>
                <Button variant="destructive" size="sm" onClick={onLeave} disabled={isOwner}>
                    <LogOutIcon className="w-4 h-4 mr-2" />
                    Leave {organizationName}
                </Button>
            </CardContent>
        </Card>
    )
}

export const OrganizationLeaveDialog: React.FC<{
    open: boolean
    onOpenChange: (open: boolean) => void
    organizationName: string
    onLeave: () => void
    isPending: boolean
}> = ({ onOpenChange, open, organizationName, onLeave, isPending }) => {
    return (
        <ConfirmDialog
            open={open}
            onOpenChange={onOpenChange}
            title="Leave Organization"
            description={`Are you sure you want to leave ${organizationName}? You will lose all access to this organization and its resources.`}
            confirmLabel="Leave Organization"
            isLoading={isPending}
            onConfirm={onLeave}
        />
    )
}

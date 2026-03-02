import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useMutation } from '@tanstack/react-query'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { CheckCircle2, XCircle, Loader2 } from 'lucide-react'
import { useEffect, useState } from 'react'
import { acceptInvitation } from '@/features/organizations/services/organization.service'

export const Route = createFileRoute('/invite/$token')({
    component: AcceptInvitePage,
})

function AcceptInvitePage() {
    const { token } = Route.useParams()
    const navigate = useNavigate()
    const [autoAccepted, setAutoAccepted] = useState(false)

    const acceptMutation = useMutation({
        mutationFn: () => acceptInvitation(token),
    })

    // Auto-accept on mount
    useEffect(() => {
        if (token && !autoAccepted) {
            setAutoAccepted(true)
            acceptMutation.mutate()
        }
    }, [token])

    return (
        <div className="flex min-h-[calc(100vh-3.5rem)] items-center justify-center px-4">
            <Card className="w-full max-w-md">
                <CardHeader className="text-center">
                    <CardTitle>Organization Invitation</CardTitle>
                    <CardDescription>
                        {acceptMutation.isPending && 'Accepting your invitation...'}
                        {acceptMutation.isSuccess && 'You have been added to the organization!'}
                        {acceptMutation.isError && 'Failed to accept invitation.'}
                    </CardDescription>
                </CardHeader>
                <CardContent className="flex flex-col items-center gap-4">
                    {acceptMutation.isPending && <Loader2 className="w-10 h-10 animate-spin text-primary" />}
                    {acceptMutation.isSuccess && (
                        <>
                            <CheckCircle2 className="w-10 h-10 text-emerald-500" />
                            <p className="text-sm text-muted-foreground text-center">
                                You can now access the organization. Head to your dashboard to get started.
                            </p>
                            <Button onClick={() => navigate({ to: '/' })}>Go to Dashboard</Button>
                        </>
                    )}
                    {acceptMutation.isError && (
                        <>
                            <XCircle className="w-10 h-10 text-destructive" />
                            <p className="text-sm text-destructive text-center">
                                {acceptMutation.error instanceof Error
                                    ? acceptMutation.error.message
                                    : 'The invitation may have expired or already been accepted.'}
                            </p>
                            <Button variant="outline" onClick={() => navigate({ to: '/' })}>
                                Go to Dashboard
                            </Button>
                        </>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}

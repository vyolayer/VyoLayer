import { Button } from '@/components/ui/button'
import { AlertCircle } from 'lucide-react'
import { useNavigate } from '@tanstack/react-router'

export const OrganizationDashboardError = ({ error }: { error: Error }) => {
    const navigate = useNavigate()

    return (
        <div className="mx-auto max-w-xl py-20 px-4 text-center space-y-4">
            <AlertCircle className="mx-auto h-10 w-10 text-destructive" />
            <h2 className="text-lg font-semibold">Could not load organizations</h2>
            <p className="text-sm text-muted-foreground">
                {error instanceof Error ? error.message : 'An unexpected error occurred.'}
            </p>
            <Button variant="outline" onClick={() => navigate({ to: '/' })}>
                Try again
            </Button>
        </div>
    )
}

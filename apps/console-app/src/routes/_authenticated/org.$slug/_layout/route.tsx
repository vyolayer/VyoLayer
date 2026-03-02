import { OrganizationHeader } from '@/org/components/OrganizationHeader'
import { OrganizationProvider } from '@/contexts/OrganizationContext'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout')({
    component: () => (
        <OrganizationProvider slug={Route.useParams().slug}>
            <div className="mx-auto py-5 max-w-7xl px-4 space-y-8">
                <OrganizationHeader />
                <Outlet />
            </div>
        </OrganizationProvider>
    ),
})

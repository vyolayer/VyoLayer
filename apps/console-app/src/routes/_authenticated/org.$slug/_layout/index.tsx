import { createFileRoute } from '@tanstack/react-router'
import { OrganizationOverviewPage } from '@/features/organizations/pages'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/')({
    component: OrganizationOverviewPage,
})

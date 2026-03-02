import { createFileRoute } from '@tanstack/react-router'
import { OrganizationArchivePage } from '@/features/organizations/pages'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/archive')({
    component: OrganizationArchivePage,
})

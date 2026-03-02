import { createFileRoute } from '@tanstack/react-router'
import { OrganizationMembersPage } from '@/features/organizations/pages'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/members/')({
    component: OrganizationMembersPage,
})

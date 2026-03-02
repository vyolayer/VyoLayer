import { OrganizationMyMembershipPage } from '@/org/pages'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/me')({
    component: OrganizationMyMembershipPage,
})

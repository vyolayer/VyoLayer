import { createFileRoute } from '@tanstack/react-router'
import { OrganizationDashboardPage } from '@/features/organizations/pages'

export const Route = createFileRoute('/_authenticated/dashboard')({
    component: OrganizationDashboardPage,
})

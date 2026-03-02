import { OrganizationDangerZone } from '@/features/organizations/pages/OrganizationDangerZone'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/danger')({
    component: OrganizationDangerZone,
})

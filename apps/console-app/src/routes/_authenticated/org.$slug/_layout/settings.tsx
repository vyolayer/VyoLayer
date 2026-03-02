import { createFileRoute } from '@tanstack/react-router'
import { OrganizationSettingsPage } from '@/org/pages'

export const Route = createFileRoute('/_authenticated/org/$slug/_layout/settings')({
    component: OrganizationSettingsPage,
})

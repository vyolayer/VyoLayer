import { createFileRoute } from '@tanstack/react-router'
import SettingsPage from '@/pages/user/Settings'

export const Route = createFileRoute('/_authenticated/settings')({
    component: SettingsPage,
})

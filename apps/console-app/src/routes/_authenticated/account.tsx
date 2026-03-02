import AccountPage from '@/pages/user/Account'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/account')({
    loader(ctx) {
        const search = ctx.location.search
        if (search && (search as any).tab === 'change-password') {
            return { tab: 'change-password' }
        }
        return { tab: null }
    },
    component: AccountPage,
})

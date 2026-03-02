import { Badge } from '@/components/ui/badge'
import { Crown, Shield, User } from 'lucide-react'

const ROLE_CONFIG: Record<
    string,
    { icon: typeof Crown; variant: 'default' | 'secondary' | 'outline'; className: string }
> = {
    Owner: {
        icon: Crown,
        variant: 'default',
        className: 'bg-amber-500/15 text-amber-700 dark:text-amber-400 border-amber-500/20',
    },
    Admin: {
        icon: Shield,
        variant: 'default',
        className: 'bg-blue-500/15 text-blue-700 dark:text-blue-400 border-blue-500/20',
    },
    Member: { icon: User, variant: 'secondary', className: '' },
}

interface RoleBadgeProps {
    name: string
    size?: 'sm' | 'md'
}

export function RoleBadge({ name, size = 'sm' }: RoleBadgeProps) {
    const config = ROLE_CONFIG[name] ?? ROLE_CONFIG.Member
    const Icon = config.icon

    return (
        <Badge variant="outline" className={`gap-1 ${size === 'sm' ? 'text-xs' : 'text-sm'} ${config.className}`}>
            <Icon className={size === 'sm' ? 'h-2.5 w-2.5' : 'h-3 w-3'} />
            {name}
        </Badge>
    )
}

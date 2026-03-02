import * as React from 'react'
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar'
import { useAuth } from '@/contexts/AuthContext'
import { useAvatarRandomBg } from '@/hooks/useAvatarRandomBg'

export const UserAvatar: React.FC<{
    src?: string
    alt?: string
    fallback?: string
}> = ({ src, alt, fallback }) => {
    return (
        <Avatar>
            <AvatarImage src={src} alt={alt} />
            <AvatarFallback>{fallback}</AvatarFallback>
        </Avatar>
    )
}

const CurrentUserAvatarMemo: React.FC<{ className?: string }> = ({ className }) => {
    const { user } = useAuth()
    const fallback = user?.fullName.charAt(0).toUpperCase()
    const randomColor = useAvatarRandomBg()

    return (
        <Avatar className={className}>
            <AvatarImage src={undefined} alt={user?.fullName} />
            <AvatarFallback className={`${randomColor} font-medium text-white`}>{fallback}</AvatarFallback>
        </Avatar>
    )
}

export const CurrentUserAvatar = React.memo(CurrentUserAvatarMemo)

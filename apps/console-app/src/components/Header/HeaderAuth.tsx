import { Link } from '@tanstack/react-router'
import type { User } from '@/lib/services/auth.service'
import { LogIn, LogOutIcon, SettingsIcon, User2Icon, UserPlus } from 'lucide-react'

import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '../ui/dropdown-menu'
import { Button } from '../ui/button'
import { CurrentUserAvatar } from '../UserAvatar'

interface HeaderAuthProps {
    isAuthenticated: boolean
    user: User | null
    onLogout: () => void
}

export function HeaderAuth({ isAuthenticated, user, onLogout }: HeaderAuthProps) {
    if (isAuthenticated && user) {
        return (
            <DropdownMenu>
                <DropdownMenuTrigger asChild>
                    <Button variant="ghost" className="relative h-9 w-9 rounded-full p-0">
                        <CurrentUserAvatar />
                    </Button>
                </DropdownMenuTrigger>

                <DropdownMenuContent className="w-56" align="end" sideOffset={8} forceMount>
                    <DropdownMenuLabel className="font-normal">
                        <div className="flex flex-col space-y-1">
                            <p className="text-sm font-medium leading-none">{user.fullName}</p>
                            <p className="text-xs leading-none text-muted-foreground">{user.email}</p>
                        </div>
                    </DropdownMenuLabel>

                    <DropdownMenuSeparator />

                    <DropdownMenuItem asChild>
                        <Link to="/account" className="flex items-center gap-2">
                            <User2Icon className="h-4 w-4" />
                            My Account
                        </Link>
                    </DropdownMenuItem>

                    <DropdownMenuItem asChild>
                        <Link to="/settings" className="flex items-center gap-2">
                            <SettingsIcon className="h-4 w-4" />
                            Settings
                        </Link>
                    </DropdownMenuItem>

                    <DropdownMenuSeparator />

                    <DropdownMenuItem
                        onClick={onLogout}
                        className="flex cursor-pointer items-center gap-2 text-destructive focus:bg-destructive/10 focus:text-destructive"
                    >
                        <LogOutIcon className="h-4 w-4" />
                        Logout
                    </DropdownMenuItem>
                </DropdownMenuContent>
            </DropdownMenu>
        )
    }

    return (
        <div className="flex items-center gap-2">
            <Button variant="ghost" size="sm" asChild>
                <Link to="/auth/login" className="flex items-center gap-2">
                    <LogIn className="h-4 w-4" />
                    <span className="hidden sm:inline">Sign in</span>
                </Link>
            </Button>

            <Button size="sm" asChild>
                <Link to="/auth/register" className="flex items-center gap-2">
                    <UserPlus className="h-4 w-4" />
                    <span className="hidden sm:inline">Sign up</span>
                </Link>
            </Button>
        </div>
    )
}

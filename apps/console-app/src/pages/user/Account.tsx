import * as React from 'react'

import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'

import { useAuth } from '@/contexts/AuthContext'
import ChangePasswordPage from './ChangePassword'
import { CurrentUserAvatar } from '@/components/UserAvatar'
import { Link, useLoaderData } from '@tanstack/react-router'
import { AuthPageLayout } from '@/components/auth/AuthPageLayout'

export default function AccountPage() {
    const { user } = useAuth()
    const { tab } = useLoaderData({ from: '/_authenticated/account' })

    const [fullName, setFullName] = React.useState(user?.fullName ?? '')

    const handleSave = () => {
        // TODO: connect to backend update profile endpoint
        console.log('Saving profile', { fullName })
    }

    if (!user) return null

    if (tab === 'change-password') {
        return (
            <AuthPageLayout>
                <ChangePasswordPage />
            </AuthPageLayout>
        )
    }

    return (
        <div className="mx-auto max-w-4xl space-y-8 py-8">
            <div>
                <h1 className="text-2xl font-semibold tracking-tight">Account</h1>
                <p className="text-sm text-muted-foreground">Manage your personal profile and security settings.</p>
            </div>

            {/* Profile Section */}
            <Card>
                <CardHeader>
                    <CardTitle>Profile</CardTitle>
                    <CardDescription>Update your personal information.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                    <div className="flex items-center gap-6">
                        <CurrentUserAvatar className="h-16 w-16 text-2xl" />

                        <div className="space-y-1">
                            <div className="flex items-center gap-2">
                                <p className="text-sm font-medium">{user.email}</p>
                                {user.isEmailVerified && <Badge variant="success">Verified</Badge>}
                                {!user.isEmailVerified && (
                                    <div className="flex items-center gap-2">
                                        <Badge variant="destructive">Not Verified</Badge>
                                        <Button variant="link">Resend Verification Email</Button>
                                    </div>
                                )}
                            </div>
                            <p className="text-xs text-muted-foreground">Email cannot be changed.</p>
                        </div>
                    </div>

                    <Separator />

                    <div className="grid gap-4 sm:grid-cols-2">
                        <div className="space-y-2">
                            <label className="text-sm font-medium">Full Name</label>
                            <Input value={fullName} onChange={(e) => setFullName(e.target.value)} />
                        </div>
                    </div>

                    <div className="flex justify-end">
                        <Button onClick={handleSave}>Save Changes</Button>
                    </div>
                </CardContent>
            </Card>

            {/* Security Section */}
            <Card>
                <CardHeader>
                    <CardTitle>Security</CardTitle>
                    <CardDescription>Manage authentication and session security.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4 flex gap-4">
                    <Link to="/account" search={{ tab: 'change-password' }}>
                        <Button variant="outline">Change Password</Button>
                    </Link>

                    <Button variant="outline" disabled>
                        Enable Two-Factor Authentication
                    </Button>
                </CardContent>
            </Card>
        </div>
    )
}

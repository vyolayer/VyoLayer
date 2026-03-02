import * as React from 'react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import { Separator } from '@/components/ui/separator'
import { useTheme } from '@/components/ThemeProvider'
import { Button } from '@/components/ui/button'

export default function SettingsPage() {
    const { theme, setTheme } = useTheme()

    const [compactMode, setCompactMode] = React.useState(false)
    const [emailNotifications, setEmailNotifications] = React.useState(true)

    return (
        <div className="mx-auto max-w-4xl space-y-8 py-8">
            <div>
                <h1 className="text-2xl font-semibold tracking-tight">Settings</h1>
                <p className="text-sm text-muted-foreground">Configure your Worklayer console preferences.</p>
            </div>

            {/* Appearance */}
            <Card>
                <CardHeader>
                    <CardTitle>Appearance</CardTitle>
                    <CardDescription>Customize your console theme and layout.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm font-medium">Dark Mode</p>
                            <p className="text-xs text-muted-foreground">Toggle between light and dark theme.</p>
                        </div>

                        <Switch
                            checked={theme === 'dark'}
                            onCheckedChange={(checked) => setTheme(checked ? 'dark' : 'light')}
                        />
                    </div>

                    <Separator />

                    <div className="flex items-center justify-between">
                        <div>
                            <p className="text-sm font-medium">Compact Mode</p>
                            <p className="text-xs text-muted-foreground">Reduce spacing in tables and panels.</p>
                        </div>

                        <Switch checked={compactMode} onCheckedChange={setCompactMode} />
                    </div>
                </CardContent>
            </Card>

            {/* Notifications */}
            <Card>
                <CardHeader>
                    <CardTitle>Notifications</CardTitle>
                    <CardDescription>Control how you receive updates.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="flex items-center justify-between">
                        <span className="text-sm font-medium">Email Notifications</span>

                        <Switch checked={emailNotifications} onCheckedChange={setEmailNotifications} />
                    </div>
                </CardContent>
            </Card>

            {/* Developer Preferences */}
            <Card>
                <CardHeader>
                    <CardTitle>Developer Preferences</CardTitle>
                    <CardDescription>Console behavior and API settings.</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4 flex gap-4">
                    <Button variant="outline" disabled>
                        Regenerate API Token
                    </Button>

                    <Button variant="destructive" disabled>
                        Reset Preferences
                    </Button>
                </CardContent>
            </Card>
        </div>
    )
}

import * as React from 'react'

export type Theme = 'light' | 'dark' | 'system'

interface ThemeContextValue {
    theme: Theme
    setTheme: (theme: Theme) => void
}

const ThemeContext = React.createContext<ThemeContextValue | null>(null)

export function ThemeProvider({
    children,
    storageKey = 'worklayer-theme',
    defaultTheme = 'system',
}: {
    children: React.ReactNode
    storageKey?: string
    defaultTheme?: Theme
}) {
    const [theme, setThemeState] = React.useState<Theme>(() => {
        if (typeof window === 'undefined') return defaultTheme
        return (localStorage.getItem(storageKey) as Theme) || defaultTheme
    })

    const applyTheme = React.useCallback((theme: Theme) => {
        const root = document.documentElement
        root.classList.remove('dark')

        if (theme === 'dark') {
            root.classList.add('dark')
        } else if (theme === 'system') {
            const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches
            if (systemDark) root.classList.add('dark')
        }
    }, [])

    React.useEffect(() => {
        applyTheme(theme)
        localStorage.setItem(storageKey, theme)
    }, [theme, applyTheme, storageKey])

    // Listen for system changes if using system theme
    React.useEffect(() => {
        if (theme !== 'system') return

        const media = window.matchMedia('(prefers-color-scheme: dark)')
        const listener = () => applyTheme('system')

        media.addEventListener('change', listener)
        return () => media.removeEventListener('change', listener)
    }, [theme, applyTheme])

    return <ThemeContext.Provider value={{ theme, setTheme: setThemeState }}>{children}</ThemeContext.Provider>
}

export function useTheme() {
    const ctx = React.useContext(ThemeContext)
    if (!ctx) throw new Error('useTheme must be used inside ThemeProvider')
    return ctx
}

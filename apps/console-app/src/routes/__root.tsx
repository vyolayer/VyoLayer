import { HeadContent, Scripts, createRootRouteWithContext } from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'

import Header from '../components/Header/Header'
import { AuthProvider } from '../contexts/AuthContext'

import TanStackQueryProvider from '../integrations/tanstack-query/root-provider'

import TanStackQueryDevtools from '../integrations/tanstack-query/devtools'

import appCss from '../styles.css?url'

import type { QueryClient } from '@tanstack/react-query'
import { ThemeProvider } from '@/components/ThemeProvider'

interface MyRouterContext {
    queryClient: QueryClient
}

function DevTools() {
    return (
        <TanStackDevtools
            config={{ position: 'bottom-right' }}
            plugins={[
                {
                    name: 'Tanstack Router',
                    render: <TanStackRouterDevtoolsPanel />,
                },
                TanStackQueryDevtools,
            ]}
        />
    )
}

function AppProviders({ children }: { children: React.ReactNode }) {
    return (
        <ThemeProvider defaultTheme="system">
            <AuthProvider>
                <TanStackQueryProvider>{children}</TanStackQueryProvider>
            </AuthProvider>
        </ThemeProvider>
    )
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
    head: () => ({
        meta: [
            {
                charSet: 'utf-8',
            },
            {
                name: 'viewport',
                content: 'width=device-width, initial-scale=1',
            },
            {
                title: 'WorkLayer Console',
            },
        ],
        links: [
            {
                rel: 'stylesheet',
                href: appCss,
            },
        ],
    }),
    shellComponent: RootDocument,
})

function RootDocument({ children }: { children: React.ReactNode }) {
    return (
        <html lang="en" suppressHydrationWarning>
            <head>
                <HeadContent />

                <script
                    dangerouslySetInnerHTML={{
                        __html: `
                            (function () {
                            try {
                                const storageKey = 'worklayer-theme';
                                const stored = localStorage.getItem(storageKey);
                                const root = document.documentElement;

                                if (stored === 'dark') {
                                root.classList.add('dark');
                                } else if (stored === 'light') {
                                root.classList.remove('dark');
                                } else {
                                const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
                                if (systemDark) root.classList.add('dark');
                                }
                            } catch (_) {}
                            })();
                            `,
                    }}
                />
            </head>
            <body suppressHydrationWarning>
                <AppProviders>
                    <Header />
                    {children}
                    <DevTools />
                </AppProviders>
                <Scripts />
            </body>
        </html>
    )
}

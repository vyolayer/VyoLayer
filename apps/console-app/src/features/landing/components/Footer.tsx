import { Separator } from '@/components/ui/separator'
import { Link } from '@tanstack/react-router'

export function Footer() {
    return (
        <footer className="border-t border-border/40 px-6 py-12">
            <div className="mx-auto max-w-6xl">
                <div className="grid gap-10 md:grid-cols-4">
                    {/* Brand */}
                    <div className="md:col-span-1">
                        <span className="text-lg font-bold tracking-tight">WorkLayer</span>
                        <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                            The open-source backend platform for building internal tools.
                        </p>
                    </div>

                    {/* Product */}
                    <div>
                        <h4 className="mb-3 text-sm font-semibold uppercase tracking-wider text-muted-foreground/70">
                            Product
                        </h4>
                        <ul className="space-y-2 text-sm">
                            <li>
                                <a
                                    href="#features"
                                    className="text-muted-foreground transition-colors hover:text-foreground"
                                >
                                    Features
                                </a>
                            </li>
                            <li>
                                <Link
                                    to="/auth/register"
                                    className="text-muted-foreground transition-colors hover:text-foreground"
                                >
                                    Get Started
                                </Link>
                            </li>
                            <li>
                                <Link
                                    to="/auth/login"
                                    className="text-muted-foreground transition-colors hover:text-foreground"
                                >
                                    Console
                                </Link>
                            </li>
                        </ul>
                    </div>

                    {/* Developers */}
                    <div>
                        <h4 className="mb-3 text-sm font-semibold uppercase tracking-wider text-muted-foreground/70">
                            Developers
                        </h4>
                        <ul className="space-y-2 text-sm">
                            <li>
                                <a
                                    href="https://github.com/worklayer"
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-muted-foreground transition-colors hover:text-foreground"
                                >
                                    GitHub
                                </a>
                            </li>
                            <li>
                                <a href="#" className="text-muted-foreground transition-colors hover:text-foreground">
                                    API Reference
                                </a>
                            </li>
                            <li>
                                <a href="#" className="text-muted-foreground transition-colors hover:text-foreground">
                                    TypeScript SDK
                                </a>
                            </li>
                        </ul>
                    </div>

                    {/* Legal */}
                    <div>
                        <h4 className="mb-3 text-sm font-semibold uppercase tracking-wider text-muted-foreground/70">
                            Legal
                        </h4>
                        <ul className="space-y-2 text-sm">
                            <li>
                                <a href="#" className="text-muted-foreground transition-colors hover:text-foreground">
                                    Privacy Policy
                                </a>
                            </li>
                            <li>
                                <a href="#" className="text-muted-foreground transition-colors hover:text-foreground">
                                    Terms of Service
                                </a>
                            </li>
                        </ul>
                    </div>
                </div>

                <Separator className="my-8 opacity-40" />

                <div className="flex flex-col items-center justify-between gap-4 text-sm text-muted-foreground sm:flex-row">
                    <span>© {new Date().getFullYear()} WorkLayer. All rights reserved.</span>
                    <span className="flex items-center gap-1.5">
                        Built with <span className="text-primary">♥</span> for developers
                    </span>
                </div>
            </div>
        </footer>
    )
}

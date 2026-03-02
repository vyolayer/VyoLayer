import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'
import { ArrowRightIcon, Code2Icon, GitBranchIcon, GlobeIcon, TerminalIcon, ZapIcon } from 'lucide-react'
import { HERO_CODE } from '../data'

export function HeroSection() {
    return (
        <section className="relative px-6 pt-20 pb-16 md:pt-32 md:pb-24">
            <div className="mx-auto max-w-6xl">
                <div className="grid items-center gap-12 lg:grid-cols-2 lg:gap-16">
                    {/* Left: Copy */}
                    <div className="space-y-8">
                        <div className="space-y-2">
                            <Badge
                                variant="outline"
                                className="gap-1.5 border-primary/30 px-3 py-1 text-xs font-medium text-primary"
                            >
                                <ZapIcon className="h-3 w-3" />
                                Open Source Backend Platform
                            </Badge>
                        </div>

                        <h1 className="text-4xl font-bold leading-tight tracking-tight md:text-5xl lg:text-6xl">
                            The backend platform for <span className="landing-gradient-text">internal tools</span>
                        </h1>

                        <p className="max-w-lg text-lg leading-relaxed text-muted-foreground md:text-xl">
                            Authentication, RBAC, audit logs, and project-scoped APIs — pre-built so your team can ship
                            internal tools <strong className="text-foreground">in hours, not weeks</strong>.
                        </p>

                        <div className="flex flex-wrap items-center gap-4">
                            <Button size="lg" asChild className="group gap-2 text-base">
                                <Link to="/auth/register">
                                    Get Started
                                    <ArrowRightIcon className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
                                </Link>
                            </Button>
                            <Button variant="outline" size="lg" asChild className="gap-2 text-base">
                                <Link to="/auth/login">
                                    <TerminalIcon className="h-4 w-4" />
                                    View Console
                                </Link>
                            </Button>
                        </div>

                        <div className="flex items-center gap-6 text-sm text-muted-foreground">
                            <span className="flex items-center gap-1.5">
                                <GitBranchIcon className="h-3.5 w-3.5" />
                                Open source
                            </span>
                            <span className="flex items-center gap-1.5">
                                <GlobeIcon className="h-3.5 w-3.5" />
                                Self-hostable
                            </span>
                            <span className="flex items-center gap-1.5">
                                <Code2Icon className="h-3.5 w-3.5" />
                                REST API
                            </span>
                        </div>
                    </div>

                    {/* Right: Hero code block */}
                    <div className="landing-hero-card relative rounded-xl border border-border/50 bg-card/60 p-1 shadow-2xl backdrop-blur-sm">
                        <div className="flex items-center gap-2 border-b border-border/50 px-4 py-3">
                            <div className="flex gap-1.5">
                                <span className="h-3 w-3 rounded-full bg-red-400/70" />
                                <span className="h-3 w-3 rounded-full bg-yellow-400/70" />
                                <span className="h-3 w-3 rounded-full bg-green-400/70" />
                            </div>
                            <span className="ml-2 text-xs text-muted-foreground font-mono">terminal</span>
                        </div>
                        <pre className="overflow-x-auto p-5 text-sm leading-relaxed bg-transparent! rounded-none!">
                            <code className="bg-transparent! p-0! text-foreground/90 font-mono">{HERO_CODE}</code>
                        </pre>
                    </div>
                </div>
            </div>
        </section>
    )
}

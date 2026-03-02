import { ChevronRightIcon, ExternalLinkIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'
import { WORKLAYER_GITHUB_REPO } from '@/lib/constant/url'

export function CTABanner() {
    return (
        <section className="relative px-6 py-20 md:py-28">
            <div className="landing-cta-card mx-auto max-w-4xl overflow-hidden rounded-2xl border border-primary/20 bg-linear-to-br from-primary/5 via-card/80 to-card/60 p-12 text-center shadow-xl backdrop-blur-sm md:p-16">
                {/* Glow */}
                <div className="pointer-events-none absolute inset-0 -z-10">
                    <div className="absolute top-0 left-1/2 h-40 w-80 -translate-x-1/2 rounded-full bg-primary/15 blur-3xl" />
                </div>

                <h2 className="mb-4 text-3xl font-bold tracking-tight md:text-4xl">Start building in minutes</h2>
                <p className="mx-auto mb-8 max-w-lg text-lg text-muted-foreground">
                    Create your free account, spin up an organization, and start integrating the API into your internal
                    tools.
                </p>
                <div className="flex flex-wrap items-center justify-center gap-4">
                    <Button size="lg" asChild className="group gap-2 text-base">
                        <Link to="/auth/register">
                            Get Started Free
                            <ChevronRightIcon className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
                        </Link>
                    </Button>
                    <Button variant="ghost" size="lg" asChild className="gap-2 text-base text-muted-foreground">
                        <a href={WORKLAYER_GITHUB_REPO} target="_blank" rel="noopener noreferrer">
                            <ExternalLinkIcon className="h-4 w-4" />
                            GitHub
                        </a>
                    </Button>
                </div>
            </div>
        </section>
    )
}

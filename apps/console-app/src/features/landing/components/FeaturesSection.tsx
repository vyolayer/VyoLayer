import { Badge } from '@/components/ui/badge'
import { FEATURES } from '../data'

export function FeaturesSection() {
    return (
        <section id="features" className="relative px-6 py-20 md:py-28">
            <div className="mx-auto max-w-6xl">
                <div className="mb-16 text-center">
                    <Badge
                        variant="outline"
                        className="mb-4 gap-1.5 border-primary/30 px-3 py-1 text-xs font-medium text-primary"
                    >
                        Features
                    </Badge>
                    <h2 className="mb-4 text-3xl font-bold tracking-tight md:text-4xl">
                        Everything you need, <span className="landing-gradient-text">nothing you don't</span>
                    </h2>
                    <p className="mx-auto max-w-2xl text-lg text-muted-foreground">
                        Stop rebuilding the same backend infrastructure. WorkLayer gives you production-ready building
                        blocks for every internal tool.
                    </p>
                </div>

                <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
                    {FEATURES.map((feature) => (
                        <div
                            key={feature.title}
                            className="landing-feature-card group relative overflow-hidden rounded-xl border border-border/50 bg-card/40 p-6 backdrop-blur-sm transition-all duration-300 hover:border-primary/30 hover:bg-card/70 hover:shadow-lg"
                        >
                            {/* Hover glow */}
                            <div className="pointer-events-none absolute inset-0 -z-10 opacity-0 transition-opacity duration-300 group-hover:opacity-100">
                                <div className="absolute -top-12 -right-12 h-32 w-32 rounded-full bg-primary/10 blur-2xl" />
                            </div>

                            <div className="mb-4 flex h-11 w-11 items-center justify-center rounded-lg border border-border/60 bg-muted/50">
                                <feature.icon className="h-5 w-5 text-primary" />
                            </div>
                            <h3 className="mb-2 text-base font-semibold">{feature.title}</h3>
                            <p className="text-sm leading-relaxed text-muted-foreground">{feature.description}</p>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    )
}

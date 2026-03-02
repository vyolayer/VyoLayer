import { Badge } from '@/components/ui/badge'
import { STEPS } from '../data'

export function HowItWorksSection() {
    return (
        <section className="relative px-6 py-20 md:py-28">
            <div className="mx-auto max-w-6xl">
                <div className="mb-16 text-center">
                    <Badge
                        variant="outline"
                        className="mb-4 gap-1.5 border-primary/30 px-3 py-1 text-xs font-medium text-primary"
                    >
                        How It Works
                    </Badge>
                    <h2 className="mb-4 text-3xl font-bold tracking-tight md:text-4xl">
                        Up and running in <span className="landing-gradient-text">Three Steps</span>
                    </h2>
                </div>

                <div className="grid gap-8 md:grid-cols-3">
                    {STEPS.map((item, index) => (
                        <div key={item.step} className="relative text-center">
                            {/* Connector line */}
                            {index < STEPS.length - 1 && (
                                <div className="absolute top-8 left-[calc(50%+3rem)] hidden h-px w-[calc(100%-6rem)] bg-linear-to-r from-primary/30 to-transparent md:block" />
                            )}

                            <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-2xl border border-border/50 bg-card/60 text-2xl font-bold text-primary shadow-lg backdrop-blur-sm">
                                {item.step}
                            </div>
                            <h3 className="mb-2 text-lg font-semibold">{item.title}</h3>
                            <p className="text-sm leading-relaxed text-muted-foreground">{item.description}</p>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    )
}

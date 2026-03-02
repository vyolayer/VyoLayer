import { Badge } from '@/components/ui/badge'
import { GitBranchIcon } from 'lucide-react'
import { TECH_STACK } from '../data'

export function TechStackSection() {
    return (
        <section className="relative px-6 py-20 md:py-28">
            <div className="mx-auto max-w-4xl text-center">
                <Badge
                    variant="outline"
                    className="mb-4 gap-1.5 border-primary/30 px-3 py-1 text-xs font-medium text-primary"
                >
                    <GitBranchIcon className="h-3 w-3" />
                    Open Source
                </Badge>
                <h2 className="mb-4 text-3xl font-bold tracking-tight md:text-4xl">
                    Built with a <span className="landing-gradient-text">Modern Stack</span>
                </h2>
                <p className="mx-auto mb-10 max-w-2xl text-lg text-muted-foreground">
                    WorkLayer is open source and self-hostable. Built with battle-tested technologies for reliability
                    and performance.
                </p>

                <div className="flex flex-wrap items-center justify-center gap-3">
                    {TECH_STACK.map((tech) => (
                        <span
                            key={tech.name}
                            className={`inline-flex items-center rounded-full border px-4 py-2 text-sm font-medium transition-transform hover:scale-105 ${tech.color}`}
                        >
                            {tech.name}
                        </span>
                    ))}
                </div>
            </div>
        </section>
    )
}

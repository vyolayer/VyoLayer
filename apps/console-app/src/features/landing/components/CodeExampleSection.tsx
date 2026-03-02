import { Badge } from '@/components/ui/badge'
import { CURL_SNIPPET, SDK_SNIPPET } from '../data'
import { Code2Icon, Terminal } from 'lucide-react'

export function CodeExampleSection() {
    return (
        <section className="relative px-6 py-20 md:py-28">
            <div className="mx-auto max-w-6xl">
                <div className="mb-16 text-center">
                    <Badge
                        variant="outline"
                        className="mb-4 gap-1.5 border-primary/30 px-3 py-1 text-xs font-medium text-primary"
                    >
                        <Terminal className="h-3 w-3" />
                        Developer Experience
                    </Badge>
                    <h2 className="mb-4 text-3xl font-bold tracking-tight md:text-4xl">
                        Ship faster with a <span className="landing-gradient-text">Clean API</span>
                    </h2>
                    <p className="mx-auto max-w-2xl text-lg text-muted-foreground">
                        Use the REST API directly or the typed TypeScript SDK — both are designed for speed and
                        simplicity.
                    </p>
                </div>

                <div className="grid gap-8 lg:grid-cols-2">
                    {/* cURL example */}
                    <div className="landing-code-card overflow-hidden rounded-xl border border-border/50 bg-card/60 shadow-lg backdrop-blur-sm">
                        <div className="flex items-center gap-2 border-b border-border/50 px-5 py-3">
                            <Terminal className="h-4 w-4 text-muted-foreground" />
                            <span className="text-sm font-medium text-muted-foreground">cURL</span>
                            <Badge
                                variant="outline"
                                className="ml-auto text-[10px] px-1.5 py-0 border-primary/30 text-primary"
                            >
                                REST
                            </Badge>
                        </div>
                        <pre className="overflow-x-auto p-5 text-sm leading-relaxed bg-transparent! rounded-none!">
                            <code className="bg-transparent! p-0! text-foreground/90 font-mono">{CURL_SNIPPET}</code>
                        </pre>
                    </div>

                    {/* TypeScript SDK example */}
                    <div className="landing-code-card overflow-hidden rounded-xl border border-border/50 bg-card/60 shadow-lg backdrop-blur-sm">
                        <div className="flex items-center gap-2 border-b border-border/50 px-5 py-3">
                            <Code2Icon className="h-4 w-4 text-muted-foreground" />
                            <span className="text-sm font-medium text-muted-foreground">TypeScript SDK</span>
                            <Badge
                                variant="outline"
                                className="ml-auto text-[10px] px-1.5 py-0 border-primary/30 text-primary"
                            >
                                SDK
                            </Badge>
                        </div>
                        <pre className="overflow-x-auto p-5 text-sm leading-relaxed bg-transparent! rounded-none!">
                            <code className="bg-transparent! p-0! text-foreground/90 font-mono">{SDK_SNIPPET}</code>
                        </pre>
                    </div>
                </div>
            </div>
        </section>
    )
}

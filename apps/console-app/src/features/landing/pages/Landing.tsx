import { CodeExampleSection } from '../components/CodeExampleSection'
import { CTABanner } from '../components/CTABanner'
import { FeaturesSection } from '../components/FeaturesSection'
import { Footer } from '../components/Footer'
import { HeroSection } from '../components/HeroSection'
import { HowItWorksSection } from '../components/HowItWorksSection'
import { TechStackSection } from '../components/TechStackSection'

export function LandingPage() {
    return (
        <>
            <HeroSection />
            <FeaturesSection />
            <CodeExampleSection />
            <HowItWorksSection />
            <TechStackSection />
            <CTABanner />
            <Footer />
        </>
    )
}

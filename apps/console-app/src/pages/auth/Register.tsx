import { Link } from '@tanstack/react-router'
import { CheckCircle } from 'lucide-react'
import { AuthPageLayout } from '@/components/auth/AuthPageLayout'
import { AuthCard } from '@/components/auth/AuthCard'
import { ServerError } from '@/components/auth/AuthFormUtils'
import { Button } from '@/components/ui/button'
import { RegisterForm } from '@/components/auth/RegisterForm'
import { useRegister } from '@/hooks/auth/useRegister'

export default function RegisterPage() {
    const { register, handleSubmit, errors, isSubmitting, onSubmit, isRegistered, serverError } = useRegister()

    // ── Success state ──────────────────────────────────────────────
    if (isRegistered) {
        return (
            <AuthPageLayout>
                <div className="text-center">
                    <div className="inline-flex items-center justify-center w-16 h-16 rounded-full border border-border bg-card mb-6">
                        <CheckCircle className="w-8 h-8 text-primary" />
                    </div>
                    <h2 className="text-2xl font-bold mb-3">Account created!</h2>
                    <p className="text-muted-foreground mb-8">
                        Your account has been created successfully. You can now sign in with your credentials.
                    </p>
                    <Button asChild size="lg">
                        <Link to="/auth/login">Go to sign in</Link>
                    </Button>
                </div>
            </AuthPageLayout>
        )
    }

    // ── Registration form ──────────────────────────────────────────
    return (
        <AuthPageLayout>
            <AuthCard
                title="Create an account"
                subtitle="Get started with WorkLayer today"
                footer={
                    <>
                        Already have an account?{' '}
                        <Link to="/auth/login" className="text-primary font-medium hover:underline">
                            Sign in
                        </Link>
                    </>
                }
            >
                <ServerError message={serverError} />

                <RegisterForm
                    register={register}
                    errors={errors}
                    isSubmitting={isSubmitting}
                    onSubmit={onSubmit}
                    handleSubmit={handleSubmit}
                />
            </AuthCard>
        </AuthPageLayout>
    )
}

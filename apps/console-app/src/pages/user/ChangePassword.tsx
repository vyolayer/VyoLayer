import { CheckCircle } from 'lucide-react'
import { Link } from '@tanstack/react-router'

import { AuthCard } from '@/components/auth/AuthCard'
import { ServerError } from '@/components/auth/AuthFormUtils'
import { Button } from '@/components/ui/button'
import { ChangePasswordForm } from '@/components/auth/ChangePasswordForm'
import { useChangePassword } from '@/hooks/auth/useChangePassword'

export default function ChangePasswordPage() {
    const {
        register,
        handleSubmit,
        formState: { errors, isSubmitting },
        serverError,
        isSuccess,
        onSubmit,
    } = useChangePassword()

    // ── Success state ──────────────────────────────────────────────
    if (isSuccess) {
        return (
            <div className="text-center py-8">
                <div className="inline-flex items-center justify-center w-16 h-16 rounded-full border border-border bg-card mb-6">
                    <CheckCircle className="w-8 h-8 text-primary" />
                </div>
                <h2 className="text-2xl font-bold mb-3">Password updated!</h2>
                <p className="text-muted-foreground mb-8">Your password has been changed successfully.</p>
                <Button asChild size="lg">
                    <Link to="/account">Back to account</Link>
                </Button>
            </div>
        )
    }

    // ── Change password form ───────────────────────────────────────
    return (
        <AuthCard
            title="Change password"
            subtitle="Enter your current password and choose a new one"
            footer={
                <>
                    Changed your mind?{' '}
                    <Link to="/account" className="text-primary font-medium hover:underline">
                        Back to account
                    </Link>
                </>
            }
        >
            <ServerError message={serverError} />

            <ChangePasswordForm
                register={register}
                errors={errors}
                isSubmitting={isSubmitting}
                onSubmit={onSubmit}
                handleSubmit={handleSubmit}
            />
        </AuthCard>
    )
}

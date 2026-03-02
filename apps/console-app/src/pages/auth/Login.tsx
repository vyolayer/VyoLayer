import { Link } from '@tanstack/react-router'

import { AuthCard } from '@/components/auth/AuthCard'
import { AuthPageLayout } from '@/components/auth/AuthPageLayout'
import { ServerError } from '@/components/auth/AuthFormUtils'

import { useLogin } from '@/hooks/auth/useLogin'
import { LoginForm } from '@/components/auth/LoginForm'

export default function LoginPage() {
    const { register, handleSubmit, errors, isSubmitting, serverError, onSubmit } = useLogin()

    return (
        <AuthPageLayout>
            <AuthCard
                title="Welcome back"
                subtitle="Sign in to your WorkLayer account"
                footer={
                    <>
                        Don&apos;t have an account?{' '}
                        <Link to="/auth/register" className="text-primary font-medium hover:underline">
                            Create one
                        </Link>
                    </>
                }
            >
                <ServerError message={serverError} />

                <LoginForm
                    onSubmit={onSubmit}
                    errors={errors}
                    isSubmitting={isSubmitting}
                    register={register}
                    handleSubmit={handleSubmit}
                />
            </AuthCard>
        </AuthPageLayout>
    )
}

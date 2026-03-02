import * as React from 'react'
import { useRouter } from '@tanstack/react-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { ServerError } from '@/components/auth/AuthFormUtils'
import { CreateOrganizationForm } from '@/features/organizations/components/dashboard/OrganizationCreateForm'
import { useCreateOrganization } from '@/org/hooks/useCreateOrganization'

export function OrganizationOnboardingPage() {
    const router = useRouter()
    const { register, handleSubmit, errors, isSubmitting, serverError, createdOrg, onSubmit } = useCreateOrganization()

    // Navigate to the new org once created
    React.useEffect(() => {
        if (createdOrg) {
            router.navigate({ to: `/org/${createdOrg.slug}`, replace: true })
        }
    }, [createdOrg, router])

    return (
        <div className="mx-auto max-w-xl py-16 space-y-8 px-4">
            <div>
                <h1 className="text-2xl font-semibold tracking-tight">Create Organization</h1>
                <p className="text-muted-foreground text-sm mt-1">
                    Set up your workspace to start managing projects and API keys.
                </p>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Organization Details</CardTitle>
                    <CardDescription>
                        You will be set as the owner and can invite team members after creation.
                    </CardDescription>
                </CardHeader>

                <CardContent className="pt-2 pb-6">
                    <ServerError message={serverError} />

                    <CreateOrganizationForm
                        register={register}
                        handleSubmit={handleSubmit}
                        errors={errors}
                        isSubmitting={isSubmitting}
                        onSubmit={onSubmit}
                    />
                </CardContent>
            </Card>
        </div>
    )
}

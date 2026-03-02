export const OrganizationDashboardLoading = () => {
    return (
        <div className="mx-auto max-w-5xl py-12 px-4 space-y-6">
            <div className="h-7 w-48 animate-pulse bg-muted rounded" />
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
                {[1, 2, 3, 4].map((n) => (
                    <div key={n} className="h-40 w-full animate-pulse bg-muted rounded-xl" />
                ))}
            </div>
        </div>
    )
}

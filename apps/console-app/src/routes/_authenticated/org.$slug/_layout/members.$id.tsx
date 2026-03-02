import { createFileRoute } from "@tanstack/react-router";
import { OrganizationMemberDetailPage } from "@/features/organizations/pages";


export const Route = createFileRoute('/_authenticated/org/$slug/_layout/members/$id')({
    component: OrganizationMemberDetailPage,
})

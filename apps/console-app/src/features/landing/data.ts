import { ShieldCheck, Users, ScrollText, KeyRound, Database, Building2 } from 'lucide-react'

export const FEATURES = [
    {
        icon: ShieldCheck,
        title: 'Authentication',
        description:
            'Cookie-based auth with JWT tokens — registration, login, session refresh, and logout out of the box.',
    },
    {
        icon: Users,
        title: 'Role-Based Access',
        description: 'Fine-grained RBAC with organization-level roles and permissions across every resource.',
    },
    {
        icon: ScrollText,
        title: 'Audit Logs',
        description:
            'Immutable audit trail for every action across your organization. Full traceability for compliance.',
    },
    {
        icon: KeyRound,
        title: 'API Keys',
        description: 'Scoped API keys for programmatic access. Rotate, revoke, and manage them from the console.',
    },
    {
        icon: Database,
        title: 'Project-Scoped APIs',
        description: 'Data APIs scoped to projects — no custom backend needed. Build internal dashboards in hours.',
    },
    {
        icon: Building2,
        title: 'Organizations',
        description: 'Multi-tenant org management with invitations, member control, and team-level governance.',
    },
]

export const STEPS = [
    {
        step: '01',
        title: 'Create an account',
        description: 'Sign up and create your first organization in under a minute.',
    },
    {
        step: '02',
        title: 'Configure access',
        description: 'Set up roles, invite team members, and generate API keys.',
    },
    {
        step: '03',
        title: 'Build your tools',
        description: 'Use the REST API or TypeScript SDK to ship internal tools fast.',
    },
]

export const TECH_STACK = [
    { name: 'Go', color: 'bg-sky-500/15 text-sky-400 border-sky-500/25' },
    { name: 'PostgreSQL', color: 'bg-blue-500/15 text-blue-400 border-blue-500/25' },
    { name: 'Fiber', color: 'bg-violet-500/15 text-violet-400 border-violet-500/25' },
    { name: 'GORM', color: 'bg-amber-500/15 text-amber-400 border-amber-500/25' },
    { name: 'React', color: 'bg-cyan-500/15 text-cyan-400 border-cyan-500/25' },
    { name: 'TanStack', color: 'bg-rose-500/15 text-rose-400 border-rose-500/25' },
    { name: 'TypeScript', color: 'bg-blue-500/15 text-blue-400 border-blue-500/25' },
]

export const CURL_SNIPPET = `# Create a new organization
curl -X POST https://api.worklayer.dev/v1/organizations \\
  -H "Authorization: Bearer $TOKEN" \\
  -H "Content-Type: application/json" \\
  -d '{
    "name": "Acme Corp",
    "description": "Internal tools platform"
  }'`

export const SDK_SNIPPET = `import { WorklayerClient } from "@worklayer/sdk";

const wl = new WorklayerClient({
  baseUrl: "https://api.worklayer.dev",
  token: process.env.WORKLAYER_TOKEN,
});

// Create an organization
const org = await wl.organizations.create({
  name: "Acme Corp",
  description: "Internal tools platform",
});

// Invite a team member
await wl.organizations.invite(org.id, {
  email: "dev@acme.com",
  role: "admin",
});`

export const HERO_CODE = `$ curl https://api.worklayer.dev/v1/health
{
  "status": "ok",
  "version": "1.0.0"
}`

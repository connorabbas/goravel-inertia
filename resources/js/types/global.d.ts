/// <reference types="vite/client" />

import '@inertiajs/core'

// Shared props injected by the backend (provider Share/ShareFunc + middleware
// ShareSession). Augmenting PageProps gives `usePage().props` proper typing.
declare module '@inertiajs/core' {
    interface PageProps {
        appName: string
        flash?: Record<string, unknown>
		auth: { user: { id: number; name: string; email: string; emailVerifiedAt: string | null } | null }
		config: { appName: string; timezone: string }
		csrfToken: string
    }
}

export { }

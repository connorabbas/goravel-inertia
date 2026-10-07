/// <reference types="vite/client" />

import '@inertiajs/core'

// Shared props injected by the backend (provider Share/ShareFunc + middleware
// ShareSession). Augmenting PageProps gives `usePage().props` proper typing.
declare module '@inertiajs/core' {
    interface PageProps {
        appName: string
    }
}

export { }

import { createApp, h, type DefineComponent } from 'vue'
import { createInertiaApp } from '@inertiajs/vue3'

const pages = import.meta.glob<DefineComponent>('./pages/**/*.vue', { import: 'default' })

createInertiaApp({
    title: (title) => (title ? `${title} · Gonertia` : 'Gonertia'),
    resolve: (name) => {
        const page = pages[`./pages/${name}.vue`]
        if (!page) {
            throw new Error(`Page not found: ${name}`)
        }
        return page()
    },
    setup({ el, App, props, plugin }) {
        createApp({ render: () => h(App, props) })
            .use(plugin)
            .mount(el)
    },
})

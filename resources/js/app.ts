import '../css/app.css'
import { createApp, h, type DefineComponent } from 'vue'
import { createInertiaApp } from '@inertiajs/vue3'
import ui from '@nuxt/ui/vue-plugin'
import AppRoot from './components/AppRoot.vue'

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
        createApp({ render: () => h(AppRoot, {}, () => h(App, props)) })
            .use(plugin)
            .use(ui)
            .mount(el)
    },
})

import { defineConfig, loadEnv, type Plugin } from 'vite'
import vue from '@vitejs/plugin-vue'
import ui from '@nuxt/ui/vite'
import { resolve } from 'node:path'
import { rmSync, writeFileSync } from 'node:fs'

// The adapter reads public/hot only while the developer runs Vite.
function goravelHot(origin: string): Plugin {
    const hotFile = resolve('public/hot')
    const clean = () => rmSync(hotFile, { force: true })

    return {
        name: 'goravel-hot',
        apply: 'serve',
        configureServer(server) {
            server.httpServer?.once('listening', () => writeFileSync(hotFile, origin))
            process.once('SIGINT', () => {
                clean()
                process.exit()
            })
            process.once('exit', clean)
        },
        closeBundle: clean,
    }
}

export default defineConfig(({ command, mode }) => {
    // Compose settings override .env, matching Goravel's configuration precedence.
    const env = { ...loadEnv(mode, process.cwd(), ''), ...process.env }
    const origin = new URL(env.VITE_ORIGIN || 'http://vite.goravel.localhost')

    return {
        base: command === 'build' ? '/build/' : '/',
        plugins: [
            vue(),
            goravelHot(origin.origin),
            ui({
                router: 'inertia',
                icon: { clientBundle: { icons: ['simple-icons:github'] } },
                ui: {
                    colors: {
                        primary: 'sky',
                        neutral: 'mist'
                    },
                    button: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    badge: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    input: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    select: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    textarea: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    selectMenu: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    inputMenu: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    inputNumber: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    inputTags: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    inputDate: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    inputTime: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    },
                    pinInput: {
                        defaultVariants: {
                            variant: 'soft'
                        }
                    }
                }
            })
        ],
        resolve: { alias: { '@': resolve('resources/js') } },
        publicDir: false,
        build: {
            outDir: 'public/build',
            manifest: true,
            rollupOptions: { input: 'resources/js/app.ts' },
        },
        server: {
            host: '0.0.0.0',
            port: 5173,
            strictPort: true,
            origin: origin.origin,
            allowedHosts: [origin.hostname],
            cors: { origin: env.APP_URL || 'http://goravel.localhost' },
            hmr: {
                host: origin.hostname,
                protocol: origin.protocol === 'https:' ? 'wss' : 'ws',
                clientPort: Number(origin.port || (origin.protocol === 'https:' ? 443 : 80)),
            },
            watch: { usePolling: env.VITE_USE_POLLING === 'true' },
        },
    }
})

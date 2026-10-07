type Notification = {
    color: 'success' | 'error' | 'warning' | 'info' | 'neutral'
    title: string
    description: string
    icon?: string
}

const tones = {
    success: { color: 'success', title: 'Success', icon: 'i-lucide-circle-check' },
    info: { color: 'info', title: 'Info', icon: 'i-lucide-info' },
    warning: { color: 'warning', title: 'Warning', icon: 'i-lucide-triangle-alert' },
    error: { color: 'error', title: 'Error', icon: 'i-lucide-circle-x' },
    neutral: { color: 'neutral', title: 'Notice', icon: 'i-lucide-megaphone' },
} as const

export function flashToasts(flash: unknown): Notification[] {
    if (!flash || typeof flash !== 'object') return []

    return Object.entries(flash).flatMap(([key, value]) => {
        if (!key.endsWith('_toast') || typeof value !== 'string' || !value.trim()) return []

        const prefix = key.split('_')[0]
        const tone = prefix === 'warn' ? 'warning' : prefix
        const presentation = Object.prototype.hasOwnProperty.call(tones, tone)
            ? tones[tone as keyof typeof tones]
            : tones.neutral

        return [{ ...presentation, description: value.trim() }]
    })
}

export function errorToast(data: unknown, status: number): Notification | undefined {
    if (!data || typeof data !== 'object' || !Number.isInteger(status) || status < 400 || status > 599) return
    const payload = data as Record<string, unknown>
    if (payload.status !== status ||
        typeof payload.errorSummary !== 'string' || !payload.errorSummary.trim() ||
        typeof payload.errorDetail !== 'string' || !payload.errorDetail.trim() ||
        typeof payload.errorIcon !== 'string' || !payload.errorIcon.trim()) return

    return {
        color: status >= 500 ? 'error' : 'warning',
        title: payload.errorSummary.trim(),
        description: payload.errorDetail.trim(),
        icon: payload.errorIcon.trim(),
    }
}

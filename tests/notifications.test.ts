import { assert, test } from 'vitest'
import { errorToast, flashToasts } from '../resources/js/utils/notifications.ts'

test('only explicit, nonempty flash toasts are presented', () => {
    assert.deepEqual(flashToasts({
        success_toast: ' Saved ', warning_toast: '', error_toast: 123,
        success: 'not a toast',
    }), [{ color: 'success', title: 'Success', description: 'Saved', icon: 'i-lucide-circle-check' }])
    assert.deepEqual(flashToasts(null), [])
    assert.deepEqual(flashToasts('invalid'), [])
})

test('flash tones match Laravel, including warn and neutral fallbacks', () => {
    for (const [key, color, title, icon] of [
        ['success_toast', 'success', 'Success', 'i-lucide-circle-check'],
        ['info_toast', 'info', 'Info', 'i-lucide-info'],
        ['warning_toast', 'warning', 'Warning', 'i-lucide-triangle-alert'],
        ['warn_toast', 'warning', 'Warning', 'i-lucide-triangle-alert'],
        ['error_toast', 'error', 'Error', 'i-lucide-circle-x'],
        ['neutral_toast', 'neutral', 'Notice', 'i-lucide-megaphone'],
        ['unexpected_toast', 'neutral', 'Notice', 'i-lucide-megaphone'],
        ['toString_toast', 'neutral', 'Notice', 'i-lucide-megaphone'],
    ]) {
        assert.deepEqual(flashToasts({ [key]: ' Message ' }), [
            { color, title, icon, description: 'Message' },
        ])
    }
})

test('structured HTTP errors must match the response status', () => {
    const payload = {
        status: 403, errorSummary: 'Forbidden - 403',
        errorDetail: 'Permission denied.', errorIcon: 'i-lucide-triangle-alert',
    }
    assert.equal(errorToast(payload, 403)?.color, 'warning')
    assert.equal(errorToast({ ...payload, status: 503 }, 503)?.color, 'error')
    assert.equal(errorToast(payload, 500), undefined)
    assert.equal(errorToast({ ...payload, errorDetail: ' ' }, 403), undefined)
    assert.equal(errorToast('<html>debug response</html>', 500), undefined)
    assert.equal(errorToast({ ...payload, status: 200 }, 200), undefined)
})

<script setup lang="ts">
import { Head } from '@inertiajs/vue3'

defineProps<{
    title: string
    detail: string
    status: number
    homepageRoute: string
}>()
</script>

<template>
    <Head :title="title">
        <meta name="description" :content="detail" />
    </Head>
    <UError
        :clear="false"
        :error="{ statusCode: status, statusMessage: title, message: detail }"
        style="--ui-header-height: 0rem"
    >
        <template #links>
            <UButton :to="homepageRoute" size="lg">Back to home</UButton>
        </template>
    </UError>
</template>

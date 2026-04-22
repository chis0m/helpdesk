<template>
  <RouterView />
</template>

<!-- Centralize app brand hydration -->
<script setup lang="ts">
import { onMounted, provide, ref, watch } from 'vue'
import { appBrandKey, brandShortFromAppName, loadAppDetail } from '@/stores/app-detail'

const appName = ref('SecWeb HelpDesk')
const brandShort = ref('SecWeb')
provide(appBrandKey, { appName, brandShort })

onMounted(async () => {
  const d = await loadAppDetail()
  appName.value = d.app_name
  brandShort.value = brandShortFromAppName(d.app_name)
})

watch(appName, (v) => {
  if (typeof document !== 'undefined') {
    const t = String(v || '').trim()
    document.title = t.length > 0 ? t : 'HelpDesk'
  }
}, { immediate: true })
</script>

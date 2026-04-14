// SEC-05: Audit logging functionality in the backend is implemented but not in the frontend.
import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import { getApiBaseUrl } from '@/api/base-url'
import { initSessionRefresh, registerSessionRefreshFailure } from '@/api/session-refresh'
import { paths } from '@/constants/routes'
import { logger } from '@/utils/logger'

const rawEnv = import.meta.env.VITE_API_BASE_URL
logger.debug('app', 'API base URL (resolved)', {
  active: getApiBaseUrl(),
  VITE_API_BASE_URL:
    typeof rawEnv === 'string' && rawEnv.trim().length > 0 ? rawEnv.trim() : '(unset — default http://localhost:8080)',
})

registerSessionRefreshFailure(() => {
  void router.replace({
    path: paths.login,
    query: { redirect: router.currentRoute.value.fullPath },
  })
})

initSessionRefresh()

createApp(App).use(router).mount('#app')

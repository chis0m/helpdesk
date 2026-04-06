// VULN-06: No security/audit logging middleware in the SPA — insufficient audit is implemented server-side (see backend routes).
import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'

createApp(App).use(router).mount('#app')

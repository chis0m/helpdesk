<template>
  <div class="relative flex min-h-screen flex-col overflow-hidden bg-[var(--surface-page)] text-[var(--text-primary)]">
    <div
      class="pointer-events-none absolute inset-0 overflow-hidden"
      aria-hidden="true"
    >
      <div
        class="auth-blob absolute -right-24 top-0 h-[380px] w-[380px] rounded-full bg-[var(--surface-mint)] opacity-60 blur-3xl"
      />
      <div
        class="auth-blob absolute bottom-0 left-0 h-[320px] w-[320px] rounded-full bg-[var(--surface-muted)] opacity-80 blur-3xl"
      />
      <div class="auth-grid absolute inset-0 opacity-[0.3]" />
    </div>

    <header
      class="relative z-50 border-b border-[var(--border-subtle)]/80 bg-white/75 px-5 py-3 backdrop-blur-xl lg:px-8"
    >
      <div class="mx-auto flex max-w-6xl items-center justify-between">
        <RouterLink
          to="/"
          class="text-xl font-semibold tracking-tight text-[var(--brand-green-dark)] transition hover:opacity-90"
        >
          {{ appName }}
        </RouterLink>
        <div class="flex items-center gap-2 sm:gap-3">
          <span
            class="hidden rounded-full bg-[var(--brand-green)] px-3 py-1.5 text-xs font-semibold text-[var(--text-on-green)] shadow-sm sm:inline"
          >Product support</span>
          <RouterLink
            v-if="route.name === 'login'"
            to="/signup"
            class="auth-header-cta rounded-full bg-[var(--surface-mint)] px-4 py-2 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--surface-mint-hover)] hover:shadow-sm"
          >
            Create account
          </RouterLink>
          <button
            v-else-if="route.name === 'change-password-required'"
            type="button"
            class="text-sm font-semibold text-[var(--text-secondary)] transition hover:text-[var(--text-primary)]"
            @click="onSignOut"
          >
            Sign out
          </button>
          <RouterLink
            v-else
            to="/login"
            class="text-sm font-semibold text-[var(--text-secondary)] transition hover:text-[var(--text-primary)]"
          >
            Sign in
          </RouterLink>
        </div>
      </div>
    </header>

    <main
      class="relative z-10 mx-auto grid w-full max-w-6xl flex-1 grid-cols-1 gap-8 px-5 py-7 lg:grid-cols-2 lg:items-center lg:gap-12 lg:px-8 lg:py-9"
    >
      <!-- Brand column (desktop) -->
      <div
        class="hidden lg:flex lg:flex-col lg:justify-center lg:pr-8"
      >
        <div class="auth-panel-in max-w-md">
          <template v-if="route.name === 'login'">
            <p
              class="inline-flex rounded-full border border-[var(--border-subtle)] bg-white/80 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)]"
            >
              Sign in
            </p>
            <h2 class="mt-5 text-3xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:text-4xl lg:leading-tight">
              Welcome back to
              <span class="text-[var(--brand-green-dark)]">{{ appName }}</span>.
            </h2>
            <p class="mt-4 text-base leading-relaxed text-[var(--text-secondary)]">
              Track your tickets and replies from the {{ brandShort }} team — the same flow as signing in to a product support portal.
            </p>
          </template>
          <template v-else-if="route.name === 'change-password-required'">
            <p
              class="inline-flex rounded-full border border-[var(--border-subtle)] bg-white/80 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)]"
            >
              Account security
            </p>
            <h2 class="mt-5 text-3xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:text-4xl lg:leading-tight">
              Update your
              <span class="text-[var(--brand-green-dark)]">password</span>.
            </h2>
            <p class="mt-4 text-base leading-relaxed text-[var(--text-secondary)]">
              For your security, choose a new password before accessing tickets and profile.
            </p>
          </template>
          <template v-else>
            <p
              class="inline-flex rounded-full border border-[var(--border-subtle)] bg-white/80 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)]"
            >
              Create account
            </p>
            <h2 class="mt-5 text-3xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:text-4xl lg:leading-tight">
              Create your
              <span class="text-[var(--brand-green-dark)]">support account</span>.
            </h2>
            <p class="mt-4 text-base leading-relaxed text-[var(--text-secondary)]">
              For people using {{ brandShort }} products — register once, then sign in to open and follow up on tickets.
            </p>
          </template>
        </div>
      </div>

      <!-- Form card -->
      <div class="flex justify-center lg:justify-end">
        <div
          class="auth-card-in w-full max-w-[440px] rounded-2xl border border-[var(--border-subtle)] bg-white/90 p-6 shadow-[0_24px_80px_-32px_rgba(0,0,0,0.14)] backdrop-blur-md sm:p-8"
        >
          <RouterView />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, provide, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { appBrandKey, brandShortFromAppName, loadAppDetail } from '@/stores/app-detail'
import { performLogout } from '@/utils/perform-logout'

const route = useRoute()
const router = useRouter()

const appName = ref('SecWeb HelpDesk')
const brandShort = ref('SecWeb')
provide(appBrandKey, { appName, brandShort })

async function hydrateBrandIfNeeded() {
  if (route.name !== 'login' && route.name !== 'signup')
    return
  const d = await loadAppDetail()
  appName.value = d.app_name
  brandShort.value = brandShortFromAppName(d.app_name)
}

onMounted(() => {
  void hydrateBrandIfNeeded()
})
watch(() => route.name, () => {
  void hydrateBrandIfNeeded()
})

function onSignOut() {
  void performLogout(router)
}
</script>

<style scoped>
.auth-grid {
  background-image: radial-gradient(circle at center, var(--border-subtle) 1px, transparent 1px);
  background-size: 24px 24px;
}

.auth-blob {
  animation: auth-blob-drift 20s ease-in-out infinite alternate;
}

.auth-blob:nth-child(2) {
  animation-duration: 24s;
  animation-delay: -5s;
}

@keyframes auth-blob-drift {
  0% {
    transform: translate(0, 0) scale(1);
  }
  100% {
    transform: translate(18px, 12px) scale(1.04);
  }
}

.auth-panel-in {
  animation: auth-fade-up 0.7s cubic-bezier(0.22, 1, 0.36, 1) both;
}

.auth-card-in {
  animation: auth-fade-up 0.7s cubic-bezier(0.22, 1, 0.36, 1) 0.08s both;
}

@keyframes auth-fade-up {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.auth-header-cta {
  transform: translateZ(0);
  transition: transform 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
}

.auth-header-cta:hover {
  transform: scale(1.02);
}
</style>

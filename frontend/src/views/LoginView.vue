<template>
  <div>
    <p
      class="inline-flex rounded-full border border-[var(--border-subtle)] bg-[var(--surface-main)] px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)] lg:hidden"
    >
      Sign in
    </p>
    <h1 class="mt-4 text-2xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:mt-0 lg:text-3xl">
      Sign in to
      <span class="text-[var(--brand-green-dark)]">SecWeb Helpdesk</span>
    </h1>
    <p class="mt-2 text-sm leading-relaxed text-[var(--text-secondary)]">
      Use the email tied to your SecWeb products to view your support requests and updates.
    </p>

    <div
      v-if="justRegistered"
      class="mt-6 rounded-2xl border border-emerald-200/80 bg-emerald-50/90 px-4 py-3 text-sm text-emerald-900"
      role="status"
    >
      You’re all set — sign in whenever you’re ready.
    </div>

    <form
      class="mt-8 space-y-5"
      @submit.prevent="onSubmit"
    >
      <div
        v-if="errorMessage"
        class="rounded-2xl border border-red-200/80 bg-red-50/90 px-4 py-3 text-sm text-red-900"
        role="alert"
      >
        {{ errorMessage }}
      </div>

      <div>
        <label
          for="login-email"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Email</label>
        <input
          id="login-email"
          v-model.trim="email"
          type="email"
          name="email"
          autocomplete="email"
          required
          class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
          placeholder="you@company.com"
        >
      </div>

      <div>
        <div class="mb-1.5 flex items-center justify-between gap-2">
          <label
            for="login-password"
            class="text-sm font-medium text-[var(--text-primary)]"
          >Password</label>
          <button
            type="button"
            class="text-xs font-semibold text-[var(--brand-green-dark)] transition hover:underline"
          >
            Forgot password?
          </button>
        </div>
        <div class="relative">
          <input
            id="login-password"
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            name="password"
            autocomplete="current-password"
            required
            minlength="8"
            maxlength="128"
            class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white py-3.5 pl-4 pr-14 text-sm text-[var(--text-primary)] shadow-sm outline-none transition-[box-shadow,border-color] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
            placeholder="Your password"
          >
          <button
            type="button"
            class="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg px-2.5 py-1 text-xs font-semibold text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)] hover:text-[var(--text-primary)]"
            :aria-pressed="showPassword"
            :aria-label="showPassword ? 'Hide password' : 'Show password'"
            @click="showPassword = !showPassword"
          >
            {{ showPassword ? 'Hide' : 'Show' }}
          </button>
        </div>
        <p class="mt-1.5 text-xs text-[var(--text-muted)]">
          At least 8 characters (required by the server).
        </p>
      </div>

      <button
        type="submit"
        :disabled="submitting"
        class="auth-submit mt-2 w-full rounded-full bg-[var(--brand-green)] py-3.5 text-sm font-semibold text-[var(--text-on-green)] shadow-md transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {{ submitting ? 'Signing in…' : 'Continue' }}
      </button>
    </form>

    <p class="mt-8 border-t border-[var(--border-subtle)] pt-8 text-center text-sm text-[var(--text-secondary)]">
      No account yet?
      <RouterLink
        to="/signup"
        class="font-semibold text-[var(--brand-green-dark)] transition hover:underline"
      >
        Create one
      </RouterLink>
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchPublicCsrfToken, loginRequest } from '@/api/auth'
import { setAuthSessionFromLogin } from '@/stores/auth-session'

const router = useRouter()
const route = useRoute()
const justRegistered = computed(() => route.query.registered === '1')

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const submitting = ref(false)
const errorMessage = ref('')

async function onSubmit() {
  errorMessage.value = ''
  if (password.value.length < 8) {
    errorMessage.value = 'Password must be at least 8 characters.'
    return
  }

  submitting.value = true
  try {
    const csrf = await fetchPublicCsrfToken()
    if (!csrf.ok) {
      errorMessage.value = csrf.message
      return
    }

    const result = await loginRequest(email.value, password.value, csrf.token)
    if (!result.ok) {
      errorMessage.value = result.message
      return
    }

    setAuthSessionFromLogin(result.data)

    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : ''
    const target = redirect && redirect.startsWith('/') ? redirect : '/dashboard'
    await router.replace(target)
  }
  catch (e) {
    const msg = e instanceof Error ? e.message : String(e)
    errorMessage.value = msg.includes('fetch')
      ? 'Could not reach the server. Check that the API is running and VITE_API_BASE_URL matches it (and CORS FRONTEND_URL matches this app).'
      : 'Something went wrong. Please try again.'
  }
  finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.auth-submit {
  transform: translateZ(0);
  transition:
    transform 0.2s ease,
    filter 0.2s ease;
}

.auth-submit:hover {
  transform: scale(1.02);
}

.auth-submit:active {
  transform: scale(0.98);
}
</style>

<template>
  <div>
    <p
      class="inline-flex rounded-full border border-[var(--border-subtle)] bg-[var(--surface-main)] px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)] lg:hidden"
    >
      Staff invite
    </p>
    <h1 class="mt-4 text-2xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:mt-0 lg:text-3xl">
      Accept your
      <span class="text-[var(--brand-green-dark)]">invite</span>
    </h1>
    <p class="mt-2 text-sm leading-relaxed text-[var(--text-secondary)]">
      Set a password to activate your {{ appName }} account. Use the full link from your invitation email.
    </p>

    <div
      v-if="loadingVerify"
      class="mt-8 text-sm text-[var(--text-secondary)]"
    >
      Checking invite…
    </div>

    <div
      v-else-if="verifyError"
      class="mt-8 rounded-2xl border border-red-200/80 bg-red-50/90 px-4 py-3 text-sm text-red-900"
      role="alert"
    >
      {{ verifyError }}
    </div>

    <div
      v-else-if="verified && verified.valid === false"
      class="mt-8 rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-950"
      role="status"
    >
      This invite link is invalid, expired, or already used. Ask an admin to send a new invite.
    </div>

    <form
      v-else-if="verified && verified.valid === true"
      class="mt-8 space-y-5"
      @submit.prevent="onSubmit"
    >
      <!-- SEC-04: Accept uses public CSRF + POST; Public CSRF is validated by backend CSRF middleware. -->
      <div
        v-if="errorMessage"
        class="rounded-2xl border border-red-200/80 bg-red-50/90 px-4 py-3 text-sm text-red-900"
        role="alert"
      >
        {{ errorMessage }}
      </div>

      <div class="rounded-2xl border border-[var(--border-subtle)] bg-[var(--surface-main)] px-4 py-3 text-sm text-[var(--text-secondary)]">
        <p><span class="font-medium text-[var(--text-primary)]">{{ verified.email }}</span></p>
        <p class="mt-1">
          {{ verified.first_name }} {{ verified.last_name }}
        </p>
      </div>

      <div>
        <label
          for="ai-password"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Password</label>
        <input
          id="ai-password"
          v-model="password"
          type="password"
          name="new-password"
          autocomplete="new-password"
          required
          minlength="8"
          maxlength="128"
          class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
        >
      </div>

      <div>
        <label
          for="ai-password2"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Confirm password</label>
        <input
          id="ai-password2"
          v-model="password2"
          type="password"
          autocomplete="new-password"
          required
          minlength="8"
          maxlength="128"
          class="auth-input w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3.5 text-sm text-[var(--text-primary)] shadow-sm outline-none focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
        >
      </div>

      <button
        type="submit"
        :disabled="submitting"
        class="auth-submit mt-2 w-full rounded-full bg-[var(--brand-green)] py-3.5 text-sm font-semibold text-[var(--text-on-green)] shadow-md transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {{ submitting ? 'Creating account…' : 'Activate account' }}
      </button>
    </form>

    <p class="mt-8 border-t border-[var(--border-subtle)] pt-8 text-center text-sm text-[var(--text-secondary)]">
      <RouterLink
        :to="paths.login"
        class="font-semibold text-[var(--brand-green-dark)] transition hover:underline"
      >
        Back to sign in
      </RouterLink>
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchPublicCsrfToken } from '@/api/auth'
import { acceptInviteRequest, fetchInviteVerify, type InviteVerifyValid } from '@/api/invites'
import { paths } from '@/constants/routes'
import { appBrandKey } from '@/stores/app-detail'

const route = useRoute()
const router = useRouter()

const brand = inject(appBrandKey, null)
const appName = computed(() => brand?.appName.value ?? 'SecWeb HelpDesk')

const token = computed(() => {
  const q = route.query.token
  return typeof q === 'string' ? q.trim() : ''
})

const loadingVerify = ref(true)
const verifyError = ref('')
const verified = ref<InviteVerifyValid | { valid: false } | null>(null)

const password = ref('')
const password2 = ref('')
const submitting = ref(false)
const errorMessage = ref('')

onMounted(async () => {
  if (!token.value) {
    loadingVerify.value = false
    verifyError.value = 'Missing invite token. Open the full link from your invitation email.'
    return
  }
  const res = await fetchInviteVerify(token.value)
  loadingVerify.value = false
  if (!res.ok) {
    verifyError.value = res.message
    verified.value = null
    return
  }
  verified.value = res.data
})

async function onSubmit() {
  errorMessage.value = ''
  if (password.value !== password2.value) {
    errorMessage.value = 'Passwords do not match.'
    return
  }
  if (password.value.length < 8) {
    errorMessage.value = 'Password must be at least 8 characters.'
    return
  }
  if (!token.value) {
    errorMessage.value = 'Missing invite token.'
    return
  }

  submitting.value = true
  try {
    const csrf = await fetchPublicCsrfToken()
    if (!csrf.ok) {
      errorMessage.value = csrf.message
      return
    }
    const result = await acceptInviteRequest(token.value, password.value, csrf.token)
    if (!result.ok) {
      errorMessage.value = result.message
      return
    }
    const target = result.data.redirect_to.startsWith('/') ? result.data.redirect_to : paths.login
    await router.replace(target)
  }
  finally {
    submitting.value = false
  }
}
</script>

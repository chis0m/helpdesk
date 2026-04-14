<template>
  <!-- VULN-02: Profile UI is keyed by route `userId`; no check that it matches the signed-in user (backend IDOR). -->
  <div class="mx-auto max-w-xl space-y-6">
    <header>
      <h2 class="text-lg font-semibold text-[var(--text-primary)]">
        Your profile
      </h2>
      <p class="mt-1 text-sm text-[var(--text-secondary)]">
        User ID {{ userId }} — update the details on your SecWeb Helpdesk account.
      </p>
    </header>

    <div
      v-if="loadError"
      class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-900"
      role="alert"
    >
      {{ loadError }}
    </div>

    <div
      v-else-if="loading"
      class="rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-10 text-center text-sm text-[var(--text-secondary)]"
    >
      Loading profile…
    </div>

    <form
      v-else
      class="space-y-5 rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-sm"
      @submit.prevent="onSave"
    >
      <div
        v-if="savedBanner"
        class="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-900"
        role="status"
      >
        Your profile was updated.
      </div>

      <div
        v-if="saveError"
        class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-900"
        role="alert"
      >
        {{ saveError }}
      </div>

      <div>
        <label
          for="pf-email"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Email</label>
        <input
          id="pf-email"
          v-model="form.email"
          type="email"
          autocomplete="email"
          class="w-full rounded-2xl border border-[var(--border-subtle)] bg-[var(--surface-main)] px-4 py-3 text-sm text-[var(--text-primary)]"
        >
      </div>

      <div class="grid gap-4 sm:grid-cols-2">
        <div>
          <label
            for="pf-first"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >First name</label>
          <input
            id="pf-first"
            v-model="form.firstName"
            type="text"
            class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3 text-sm"
          >
        </div>
        <div>
          <label
            for="pf-last"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Last name</label>
          <input
            id="pf-last"
            v-model="form.lastName"
            type="text"
            class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3 text-sm"
          >
        </div>
      </div>

      <div>
        <label
          for="pf-middle"
          class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
        >Middle name <span class="font-normal text-[var(--text-muted)]">(optional)</span></label>
        <input
          id="pf-middle"
          v-model="form.middleName"
          type="text"
          maxlength="100"
          class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3 text-sm"
          placeholder="Optional"
        >
      </div>

      <button
        type="submit"
        class="rounded-full bg-[var(--brand-green)] px-6 py-3 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95 disabled:opacity-50"
        :disabled="saving"
      >
        {{ saving ? 'Saving…' : 'Save changes' }}
      </button>
    </form>

    <!-- VULN-01: Change password POST uses session cookies. VULN-04: `X-CSRF-Token` on POST; weak verification is backend CSRF middleware. -->
    <section
      v-if="!loading && !loadError"
      class="space-y-5 rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-sm"
    >
      <h3 class="text-base font-semibold text-[var(--text-primary)]">
        Change password
      </h3>
      <p class="text-sm text-[var(--text-secondary)]">
        Updates the password for the account loaded in this view (same session as the profile above).
      </p>
      <div
        v-if="pwSavedBanner"
        class="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-900"
        role="status"
      >
        Password updated.
      </div>
      <div
        v-if="pwError"
        class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-900"
        role="alert"
      >
        {{ pwError }}
      </div>
      <form
        class="space-y-4"
        @submit.prevent="onChangePassword"
      >
        <div>
          <label
            for="pw-current"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Current password</label>
          <div class="relative">
            <input
              id="pw-current"
              v-model="pwCurrent"
              :type="showPwCurrent ? 'text' : 'password'"
              autocomplete="current-password"
              class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white py-3 pl-4 pr-14 text-sm text-[var(--text-primary)] outline-none focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
            >
            <button
              type="button"
              class="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg px-2.5 py-1 text-xs font-semibold text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)] hover:text-[var(--text-primary)]"
              :aria-pressed="showPwCurrent"
              :aria-label="showPwCurrent ? 'Hide current password' : 'Show current password'"
              @click="showPwCurrent = !showPwCurrent"
            >
              {{ showPwCurrent ? 'Hide' : 'Show' }}
            </button>
          </div>
        </div>
        <div>
          <label
            for="pw-new"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >New password</label>
          <div class="relative">
            <input
              id="pw-new"
              v-model="pwNew"
              :type="showPwNew ? 'text' : 'password'"
              autocomplete="new-password"
              minlength="8"
              class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white py-3 pl-4 pr-14 text-sm text-[var(--text-primary)] outline-none focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
            >
            <button
              type="button"
              class="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg px-2.5 py-1 text-xs font-semibold text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)] hover:text-[var(--text-primary)]"
              :aria-pressed="showPwNew"
              :aria-label="showPwNew ? 'Hide new password' : 'Show new password'"
              @click="showPwNew = !showPwNew"
            >
              {{ showPwNew ? 'Hide' : 'Show' }}
            </button>
          </div>
        </div>
        <div>
          <label
            for="pw-new2"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Confirm new password</label>
          <div class="relative">
            <input
              id="pw-new2"
              v-model="pwNew2"
              :type="showPwNew2 ? 'text' : 'password'"
              autocomplete="new-password"
              minlength="8"
              class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white py-3 pl-4 pr-14 text-sm text-[var(--text-primary)] outline-none focus:border-transparent focus:ring-2 focus:ring-[var(--brand-green)]"
            >
            <button
              type="button"
              class="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg px-2.5 py-1 text-xs font-semibold text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)] hover:text-[var(--text-primary)]"
              :aria-pressed="showPwNew2"
              :aria-label="showPwNew2 ? 'Hide confirm password' : 'Show confirm password'"
              @click="showPwNew2 = !showPwNew2"
            >
              {{ showPwNew2 ? 'Hide' : 'Show' }}
            </button>
          </div>
        </div>
        <button
          type="submit"
          class="rounded-full bg-[var(--brand-green)] px-6 py-3 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95 disabled:opacity-50"
          :disabled="pwSaving"
        >
          {{ pwSaving ? 'Updating…' : 'Change password' }}
        </button>
      </form>
    </section>
  </div>
</template>

<script setup lang="ts">
// VULN-02: fetch/patch user by numeric id from route — pairs with backend GET/PATCH /users/:id IDOR.
import { computed, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { changePasswordRequest } from '@/api/auth'
import { fetchUser, patchUser } from '@/api/users'
import { getSessionCsrfToken } from '@/stores/auth-session'

const route = useRoute()

const userId = computed(() => {
  const raw = route.params.userId
  const n = typeof raw === 'string' ? Number.parseInt(raw, 10) : Number.NaN
  return Number.isFinite(n) && n > 0 ? n : null
})

const loading = ref(true)
const loadError = ref('')
const saveError = ref('')
const saving = ref(false)
const savedBanner = ref(false)

const showPwCurrent = ref(false)
const showPwNew = ref(false)
const showPwNew2 = ref(false)

const pwCurrent = ref('')
const pwNew = ref('')
const pwNew2 = ref('')
const pwError = ref('')
const pwSaving = ref(false)
const pwSavedBanner = ref(false)

const form = reactive({
  email: '',
  firstName: '',
  lastName: '',
  middleName: '',
})

async function loadProfile() {
  loadError.value = ''
  loading.value = true
  const id = userId.value
  if (id === null) {
    loadError.value = 'Invalid profile link.'
    loading.value = false
    return
  }

  const result = await fetchUser(id)
  loading.value = false
  if (!result.ok) {
    loadError.value = result.message
    return
  }
  form.email = result.data.email
  form.firstName = result.data.first_name
  form.lastName = result.data.last_name
  form.middleName = result.data.middle_name ?? ''
}

watch(
  () => route.params.userId,
  () => {
    void loadProfile()
  },
  { immediate: true },
)

async function onSave() {
  saveError.value = ''
  savedBanner.value = false
  const id = userId.value
  if (id === null) {
    saveError.value = 'Invalid profile.'
    return
  }
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    saveError.value = 'Your session expired. Sign in again.'
    return
  }

  saving.value = true
  const result = await patchUser(
    id,
    {
      email: form.email.trim(),
      first_name: form.firstName.trim(),
      last_name: form.lastName.trim(),
      middle_name: form.middleName.trim() === '' ? null : form.middleName.trim(),
    },
    csrf,
  )
  saving.value = false
  if (!result.ok) {
    saveError.value = result.message
    return
  }
  form.email = result.data.email
  form.firstName = result.data.first_name
  form.lastName = result.data.last_name
  form.middleName = result.data.middle_name ?? ''
  savedBanner.value = true
  window.setTimeout(() => {
    savedBanner.value = false
  }, 3000)
}

async function onChangePassword() {
  pwError.value = ''
  pwSavedBanner.value = false
  if (pwNew.value !== pwNew2.value) {
    pwError.value = 'New passwords do not match.'
    return
  }
  if (pwNew.value.length < 8) {
    pwError.value = 'New password must be at least 8 characters.'
    return
  }
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    pwError.value = 'Your session expired. Sign in again.'
    return
  }
  pwSaving.value = true
  const result = await changePasswordRequest(
    {
      current_password: pwCurrent.value,
      new_password: pwNew.value,
    },
    csrf,
  )
  pwSaving.value = false
  if (!result.ok) {
    pwError.value = result.message
    return
  }
  pwCurrent.value = ''
  pwNew.value = ''
  pwNew2.value = ''
  pwSavedBanner.value = true
  window.setTimeout(() => {
    pwSavedBanner.value = false
  }, 4000)
}
</script>

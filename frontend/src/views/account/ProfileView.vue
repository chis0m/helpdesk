<template>
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
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
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
</script>

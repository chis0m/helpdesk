<template>
  <div class="space-y-6">
    <AdminSubnav />

    <div class="mx-auto max-w-lg space-y-6">
      <header>
        <p class="text-xs font-semibold uppercase tracking-wider text-[var(--brand-green-dark)]">
          Admin
        </p>
        <h2 class="mt-2 text-lg font-semibold text-[var(--text-primary)]">
          Create staff account
        </h2>
        <p class="mt-2 text-sm leading-relaxed text-[var(--text-secondary)]">
          Send an email invite (they choose a password when they accept), or create an account now with a temporary password you share out of band.
        </p>
      </header>

      <div
        class="flex gap-2 rounded-2xl border border-[var(--border-subtle)] bg-[var(--surface-muted)]/30 p-1"
        role="tablist"
      >
        <button
          type="button"
          role="tab"
          :aria-selected="mode === 'invite'"
          class="flex-1 rounded-xl px-3 py-2.5 text-sm font-semibold transition"
          :class="mode === 'invite'
            ? 'bg-white text-[var(--text-primary)] shadow-sm ring-1 ring-[var(--border-subtle)]'
            : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)]'"
          @click="mode = 'invite'"
        >
          Email invite
        </button>
        <button
          type="button"
          role="tab"
          :aria-selected="mode === 'direct'"
          class="flex-1 rounded-xl px-3 py-2.5 text-sm font-semibold transition"
          :class="mode === 'direct'
            ? 'bg-white text-[var(--text-primary)] shadow-sm ring-1 ring-[var(--border-subtle)]'
            : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)]'"
          @click="mode = 'direct'"
        >
          Account + password
        </button>
      </div>

      <form
        class="space-y-5 rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-sm"
        @submit.prevent="onSubmit"
      >
        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-900"
          role="alert"
        >
          {{ errorMessage }}
        </div>

        <div
          v-if="banner"
          class="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-900"
          role="status"
        >
          {{ banner }}
        </div>

        <div
          v-if="inviteNotice"
          class="space-y-3 rounded-2xl border border-emerald-200 bg-emerald-50/80 px-4 py-3 text-sm text-emerald-950"
        >
          <p class="whitespace-pre-wrap">
            {{ inviteNotice }}
          </p>
          <div v-if="inviteUrl">
            <label
              for="invite-url-copy"
              class="mb-1 block text-xs font-semibold uppercase tracking-wide text-emerald-900/80"
            >Invite link (copy for the new staff member)</label>
            <div class="flex gap-2">
              <input
                id="invite-url-copy"
                :value="inviteUrl"
                readonly
                class="min-w-0 flex-1 rounded-xl border border-emerald-200/80 bg-white px-3 py-2 font-mono text-xs text-emerald-950"
              >
              <button
                type="button"
                class="shrink-0 rounded-xl border border-emerald-300 bg-white px-3 py-2 text-xs font-semibold text-emerald-900 transition hover:bg-emerald-100"
                @click="copyInviteUrl"
              >
                {{ copyFeedback ? 'Copied' : 'Copy' }}
              </button>
            </div>
          </div>
        </div>

        <div>
          <label
            for="st-email"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Work email</label>
          <input
            id="st-email"
            v-model="form.email"
            type="email"
            required
            autocomplete="email"
            placeholder="new.staff@example.internal"
            class="w-full rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-3 text-sm outline-none ring-[var(--brand-green)] focus:border-transparent focus:ring-2"
          >
        </div>

        <div class="grid gap-4 sm:grid-cols-2">
          <div>
            <label
              for="st-first"
              class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
            >First name</label>
            <input
              id="st-first"
              v-model="form.firstName"
              type="text"
              required
              minlength="2"
              autocomplete="section-staff given-name"
              class="w-full rounded-2xl border border-[var(--border-subtle)] px-4 py-3 text-sm"
            >
          </div>
          <div>
            <label
              for="st-last"
              class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
            >Last name</label>
            <input
              id="st-last"
              v-model="form.lastName"
              type="text"
              required
              minlength="2"
              autocomplete="section-staff family-name"
              class="w-full rounded-2xl border border-[var(--border-subtle)] px-4 py-3 text-sm"
            >
          </div>
        </div>

        <div class="flex items-start gap-3 rounded-2xl border border-[var(--border-subtle)] bg-[var(--surface-muted)]/40 px-4 py-3">
          <input
            id="st-admin"
            v-model="form.isAdmin"
            type="checkbox"
            class="mt-1 h-4 w-4 rounded border-[var(--border-strong)] text-[var(--brand-green)]"
          >
          <div>
            <label
              for="st-admin"
              class="text-sm font-medium text-[var(--text-primary)]"
            >Administrator</label>
            <p class="mt-0.5 text-xs text-[var(--text-muted)]">
              Sends <code class="rounded bg-black/5 px-1 py-0.5 text-[11px]">role: admin</code>.
              Only an <strong>administrator</strong> or <strong>super administrator</strong> may choose this; staff and customers cannot.
            </p>
          </div>
        </div>

        <div v-if="mode === 'direct'">
          <label
            for="st-temp"
            class="mb-1.5 block text-sm font-medium text-[var(--text-primary)]"
          >Temporary password</label>
          <input
            id="st-temp"
            v-model="form.password"
            type="password"
            autocomplete="new-password"
            placeholder="At least 8 characters"
            :required="mode === 'direct'"
            minlength="8"
            class="w-full rounded-2xl border border-[var(--border-subtle)] px-4 py-3 text-sm"
          >
          <p class="mt-1.5 text-xs text-[var(--text-muted)]">
            No email is sent — share this password securely. For an email link instead, use <strong>Email invite</strong>.
          </p>
        </div>

        <p
          v-else
          class="text-xs text-[var(--text-muted)]"
        >
          No password here — the invitee sets one when they accept the invite. If a link appears in the success message, copy it for them.
        </p>

        <button
          type="submit"
          :disabled="submitting"
          class="w-full rounded-full bg-[var(--brand-green)] py-3.5 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {{ submitLabel }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import AdminSubnav from '@/components/admin/AdminSubnav.vue'
import { createStaffInvite, createStaffUser } from '@/api/admin'
import { getSessionCsrfToken } from '@/stores/auth-session'

const mode = ref<'invite' | 'direct'>('invite')

const form = reactive({
  email: '',
  firstName: '',
  lastName: '',
  isAdmin: false,
  password: '',
})

const banner = ref('')
const errorMessage = ref('')
const submitting = ref(false)
const inviteNotice = ref('')
const inviteUrl = ref('')
const copyFeedback = ref(false)

const submitLabel = computed(() => {
  if (submitting.value)
    return mode.value === 'invite' ? 'Sending…' : 'Creating…'
  return mode.value === 'invite' ? 'Send invite' : 'Create account'
})

watch(mode, () => {
  errorMessage.value = ''
  inviteNotice.value = ''
  inviteUrl.value = ''
})

async function copyInviteUrl() {
  if (!inviteUrl.value)
    return
  try {
    await navigator.clipboard.writeText(inviteUrl.value)
    copyFeedback.value = true
    window.setTimeout(() => {
      copyFeedback.value = false
    }, 2000)
  }
  catch {
    /* ignore */
  }
}

async function onSubmit() {
  errorMessage.value = ''
  banner.value = ''
  inviteNotice.value = ''
  inviteUrl.value = ''

  const csrf = getSessionCsrfToken()
  if (!csrf) {
    errorMessage.value = 'Your session is missing a security token. Sign out and sign in again, then retry.'
    return
  }

  submitting.value = true
  try {
    if (mode.value === 'invite') {
      const res = await createStaffInvite(
        {
          email: form.email.trim(),
          first_name: form.firstName.trim(),
          last_name: form.lastName.trim(),
          role: form.isAdmin ? 'admin' : 'staff',
        },
        csrf,
      )
      if (!res.ok) {
        errorMessage.value = res.message
        return
      }
      const d = res.data
      inviteNotice.value = `Invitation for ${d.email} (${d.target_role}). Valid until ${formatExpiry(d.expires_at_utc)}.`
      inviteUrl.value = typeof d.invite_url === 'string' ? d.invite_url : ''
      form.email = ''
      form.firstName = ''
      form.lastName = ''
      form.isAdmin = false
      return
    }

    if (form.password.length < 8) {
      errorMessage.value = 'Password must be at least 8 characters.'
      return
    }

    const res = await createStaffUser(
      {
        email: form.email.trim(),
        password: form.password,
        first_name: form.firstName.trim(),
        last_name: form.lastName.trim(),
        role: form.isAdmin ? 'admin' : 'staff',
      },
      csrf,
    )
    if (!res.ok) {
      errorMessage.value = res.message
      return
    }
    const note =
      res.data.role === 'admin'
        ? ' Account role is administrator.'
        : ''
    banner.value = `Account created for ${res.data.email} (id ${res.data.user_id}, role ${res.data.role}).${note}`
    form.email = ''
    form.firstName = ''
    form.lastName = ''
    form.password = ''
    form.isAdmin = false
    window.setTimeout(() => {
      banner.value = ''
    }, 8000)
  }
  finally {
    submitting.value = false
  }
}

function formatExpiry(iso: string): string {
  try {
    const d = new Date(iso)
    if (Number.isNaN(d.getTime()))
      return iso
    return d.toLocaleString()
  }
  catch {
    return iso
  }
}
</script>

<template>
  <div class="space-y-6">
    <header>
      <p class="text-xs font-semibold uppercase tracking-wider text-[var(--brand-green-dark)]">
        Admin
      </p>
      <h1 class="mt-2 text-xl font-semibold tracking-tight text-[var(--text-primary)]">
        Audit logs
      </h1>
      <p class="mt-2 max-w-2xl text-sm leading-relaxed text-[var(--text-secondary)]">
        Append-only security and activity records. Only administrators can open this page.
      </p>
    </header>

    <AdminSubnav />

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
      Loading audit logs…
    </div>

    <div
      v-else
      class="space-y-4"
    >
      <div
        v-if="total > 0"
        class="flex flex-wrap items-center justify-between gap-3 text-sm text-[var(--text-secondary)]"
      >
        <span>
          Page {{ page }} of {{ totalPages }} · {{ total }} entries
        </span>
        <div class="flex gap-2">
          <button
            type="button"
            class="rounded-full border border-[var(--border-subtle)] bg-white px-4 py-2 font-medium text-[var(--text-primary)] transition hover:bg-[var(--surface-hover)] disabled:opacity-40"
            :disabled="page <= 1 || loading"
            @click="goPrev"
          >
            Previous
          </button>
          <button
            type="button"
            class="rounded-full border border-[var(--border-subtle)] bg-white px-4 py-2 font-medium text-[var(--text-primary)] transition hover:bg-[var(--surface-hover)] disabled:opacity-40"
            :disabled="page >= totalPages || loading"
            @click="goNext"
          >
            Next
          </button>
        </div>
      </div>

      <p
        v-if="items.length === 0"
        class="rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-10 text-center text-sm text-[var(--text-secondary)]"
      >
        {{ total > 0 ? 'No entries on this page.' : 'No audit entries yet.' }}
      </p>

      <div
        v-else
        class="overflow-x-auto rounded-2xl border border-[var(--border-subtle)] bg-white"
      >
        <table class="min-w-full divide-y divide-[var(--border-subtle)] text-left text-sm">
          <thead class="bg-[var(--surface-muted)] text-xs font-semibold uppercase tracking-wide text-[var(--text-secondary)]">
            <tr>
              <th class="whitespace-nowrap px-4 py-3">
                Time
              </th>
              <th class="whitespace-nowrap px-4 py-3">
                Action
              </th>
              <th class="whitespace-nowrap px-4 py-3">
                OK
              </th>
              <th class="whitespace-nowrap px-4 py-3">
                Method / path
              </th>
              <th class="whitespace-nowrap px-4 py-3">
                Resource
              </th>
              <th class="min-w-[12rem] px-4 py-3">
                Actor
              </th>
              <th class="min-w-[14rem] px-4 py-3">
                Metadata
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[var(--border-subtle)] text-[var(--text-primary)]">
            <tr
              v-for="row in items"
              :key="row.id"
              class="align-top"
            >
              <td class="whitespace-nowrap px-4 py-3 font-mono text-xs text-[var(--text-secondary)]">
                {{ formatTime(row.created_at) }}
              </td>
              <td class="max-w-[14rem] px-4 py-3 font-mono text-xs break-all">
                {{ row.action }}
              </td>
              <td class="whitespace-nowrap px-4 py-3">
                <span
                  :class="row.success ? 'text-emerald-700' : 'text-red-700'"
                  class="font-medium"
                >{{ row.success ? 'yes' : 'no' }}</span>
                <span
                  v-if="row.error_code"
                  class="mt-1 block font-mono text-xs text-red-800"
                >{{ row.error_code }}</span>
              </td>
              <td class="max-w-[16rem] px-4 py-3 font-mono text-xs break-all text-[var(--text-secondary)]">
                {{ row.http_method }} {{ row.path }}
              </td>
              <td class="max-w-[12rem] px-4 py-3 font-mono text-xs break-all text-[var(--text-secondary)]">
                <template v-if="row.resource_type || row.resource_id != null">
                  {{ row.resource_type ?? '—' }}
                  <span v-if="row.resource_id != null">#{{ row.resource_id }}</span>
                </template>
                <template v-else>
                  —
                </template>
              </td>
              <td class="max-w-[14rem] px-4 py-3 font-mono text-xs break-all text-[var(--text-secondary)]">
                {{ row.actor_user_uuid ?? '—' }}
              </td>
              <td class="max-w-[20rem] px-4 py-3 font-mono text-xs break-all text-[var(--text-secondary)]">
                {{ formatMetadata(row.metadata) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { fetchAuditLogs, type AuditLogListItem } from '@/api/admin'
import AdminSubnav from '@/components/admin/AdminSubnav.vue'

const limit = 25
const page = ref(1)
const total = ref(0)
const items = ref<AuditLogListItem[]>([])
const loading = ref(true)
const loadError = ref('')

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit)))

function formatTime(iso: string) {
  try {
    const d = new Date(iso)
    return Number.isNaN(d.getTime()) ? iso : d.toLocaleString()
  }
  catch {
    return iso
  }
}

function formatMetadata(meta: unknown) {
  if (meta == null)
    return '—'
  if (typeof meta === 'string')
    return meta.length > 400 ? `${meta.slice(0, 400)}…` : meta
  try {
    const s = JSON.stringify(meta)
    return s.length > 400 ? `${s.slice(0, 400)}…` : s
  }
  catch {
    return String(meta)
  }
}

async function load() {
  loading.value = true
  loadError.value = ''
  const result = await fetchAuditLogs({ page: page.value, limit })
  loading.value = false
  if (!result.ok) {
    loadError.value = result.message
    return
  }
  items.value = result.data.items
  total.value = result.data.pagination.total
}

function goPrev() {
  if (page.value <= 1)
    return
  page.value -= 1
  void load()
}

function goNext() {
  if (page.value >= totalPages.value)
    return
  page.value += 1
  void load()
}

onMounted(() => {
  void load()
})
</script>

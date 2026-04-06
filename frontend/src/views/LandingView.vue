<template>
  <div class="relative min-h-screen overflow-hidden bg-[var(--surface-page)] text-[var(--text-primary)]">
    <!-- Ambient background -->
    <div
      class="pointer-events-none absolute inset-0 overflow-hidden"
      aria-hidden="true"
    >
      <div
        class="landing-blob absolute -right-32 -top-40 h-[520px] w-[520px] rounded-full bg-[var(--surface-mint)] opacity-70 blur-3xl"
      />
      <div
        class="landing-blob absolute -bottom-48 left-1/4 h-[420px] w-[420px] rounded-full bg-[var(--surface-muted)] opacity-90 blur-3xl"
      />
      <div
        class="landing-grid absolute inset-0 opacity-[0.35]"
      />
    </div>

    <header
      class="landing-nav sticky top-0 z-50 border-b border-[var(--border-subtle)]/80 bg-white/75 px-6 py-4 backdrop-blur-xl lg:px-10"
    >
      <div class="mx-auto flex max-w-6xl items-center justify-between">
        <RouterLink
          to="/"
          class="text-xl font-semibold tracking-tight text-[var(--brand-green-dark)] transition hover:opacity-90"
        >
          SecWeb Helpdesk
        </RouterLink>
        <div class="flex items-center gap-2 sm:gap-3">
          <span
            class="hidden rounded-full bg-[var(--brand-green)] px-3 py-1.5 text-xs font-semibold text-[var(--text-on-green)] shadow-sm sm:inline"
          >Product support</span>
          <RouterLink
            to="/login"
            class="text-sm font-semibold text-[var(--text-secondary)] transition hover:text-[var(--text-primary)]"
          >
            Sign in
          </RouterLink>
          <RouterLink
            to="/signup"
            class="rounded-full bg-[var(--surface-mint)] px-4 py-2 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--surface-mint-hover)] hover:shadow-sm"
          >
            Create account
          </RouterLink>
        </div>
      </div>
    </header>

    <main class="relative mx-auto max-w-6xl px-6 pb-28 pt-10 lg:px-10 lg:pt-14">
      <!-- Hero -->
      <div
        class="grid items-center gap-12 lg:grid-cols-2 lg:gap-16"
      >
        <div class="land-in max-w-xl">
          <p
            class="inline-flex rounded-full border border-[var(--border-subtle)] bg-white/90 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-[var(--text-secondary)] shadow-sm"
          >
            SecWeb product support
          </p>
          <h1
            class="mt-5 text-4xl font-semibold tracking-[-0.02em] text-[var(--text-primary)] lg:text-[2.75rem] lg:leading-[1.12]"
          >
            Get help with
            <span class="text-[var(--brand-green-dark)]">SecWeb</span>
            products — report issues and track fixes.
          </h1>
          <p class="mt-5 text-lg leading-relaxed text-[var(--text-secondary)]">
            Create an account, sign in, and open a ticket when something breaks — the same idea as vendor portals like Sisense or other B2B product support sites.
          </p>
          <div class="mt-10 flex flex-wrap items-center gap-3">
            <RouterLink
              to="/signup"
              class="landing-cta-primary inline-flex items-center justify-center rounded-full bg-[var(--brand-green)] px-6 py-3 text-sm font-semibold text-[var(--text-on-green)] shadow-md transition hover:brightness-95"
            >
              Get started
            </RouterLink>
            <RouterLink
              to="/login"
              class="landing-cta-secondary inline-flex items-center justify-center rounded-full bg-[var(--surface-mint)] px-6 py-3 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--surface-mint-hover)]"
            >
              Sign in
            </RouterLink>
            <RouterLink
              to="/dashboard"
              class="landing-cta-secondary inline-flex items-center justify-center rounded-full border border-[var(--border-strong)] bg-white/90 px-6 py-3 text-sm font-semibold text-[var(--text-primary)] shadow-sm transition hover:bg-[var(--surface-hover)]"
            >
              View support home (demo)
            </RouterLink>
          </div>
        </div>

        <!-- Product preview -->
        <div
          class="land-in land-delay-2 lg:justify-self-end"
        >
          <div
            class="w-full max-w-[440px] rounded-2xl border border-[var(--border-subtle)] bg-white shadow-[0_24px_80px_-20px_rgba(0,0,0,0.12)]"
          >
            <div
              class="flex items-center gap-2 border-b border-[var(--border-subtle)] bg-[var(--surface-muted)]/80 px-4 py-3"
            >
              <span class="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
              <span class="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
              <span class="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
              <div
                class="ml-3 flex-1 rounded-lg border border-[var(--border-subtle)] bg-white px-3 py-1.5 text-center text-xs text-[var(--text-muted)]"
              >
                support.secweb · My requests
              </div>
            </div>
            <div class="space-y-4 p-5">
              <div class="flex items-end justify-between">
                <div>
                  <p class="text-xs font-medium text-[var(--text-muted)]">
                    Your open requests
                  </p>
                  <p class="mt-0.5 text-3xl font-semibold tracking-tight tabular-nums">
                    12
                  </p>
                </div>
                <span
                  class="rounded-full bg-[var(--brand-green)] px-3 py-1 text-xs font-semibold text-[var(--text-on-green)]"
                >Live</span>
              </div>
              <div class="space-y-2">
                <div
                  v-for="row in previewRows"
                  :key="row.title"
                  class="flex items-center gap-3 rounded-xl border border-[var(--border-subtle)] bg-[var(--surface-main)] px-3 py-2.5"
                >
                  <span
                    class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-white text-xs font-medium text-[var(--text-muted)]"
                  >#</span>
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-sm font-medium text-[var(--text-primary)]">
                      {{ row.title }}
                    </p>
                    <p class="text-xs text-[var(--text-muted)]">
                      {{ row.meta }}
                    </p>
                  </div>
                  <span
                    class="shrink-0 rounded-full px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide"
                    :class="row.badgeClass"
                  >{{ row.badge }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Stats -->
      <section
        class="land-in land-delay-3 mt-20 rounded-2xl border border-[var(--border-subtle)] bg-white/80 p-8 shadow-[0_1px_3px_rgba(0,0,0,0.06)] backdrop-blur-sm sm:p-10"
      >
        <div class="grid gap-10 sm:grid-cols-3 sm:gap-8">
          <div
            v-for="s in stats"
            :key="s.label"
            class="text-center sm:text-left"
          >
            <p class="text-3xl font-semibold tracking-tight text-[var(--text-primary)] tabular-nums lg:text-4xl">
              {{ s.value }}
            </p>
            <p class="mt-1 text-sm font-medium text-[var(--text-secondary)]">
              {{ s.label }}
            </p>
            <p class="mt-1 text-xs text-[var(--text-muted)]">
              {{ s.hint }}
            </p>
          </div>
        </div>
      </section>

      <!-- Bento -->
      <section class="mt-14">
        <h2 class="land-in land-delay-4 text-lg font-semibold text-[var(--text-primary)]">
          How SecWeb Helpdesk fits your workflow
        </h2>
        <p class="land-in land-delay-4 mt-1 text-sm text-[var(--text-secondary)]">
          Simple steps from account to ticket — staff triage and resolve on their side.
        </p>
        <div
          class="mt-8 grid gap-4 md:grid-cols-3"
        >
          <article
            class="land-in land-delay-5 landing-card group rounded-2xl border border-[var(--border-subtle)] bg-white p-7 shadow-sm transition duration-300 md:col-span-2"
          >
            <p class="text-base font-semibold text-[var(--text-primary)]">
              {{ featured.title }}
            </p>
            <p class="mt-3 max-w-xl text-sm leading-relaxed text-[var(--text-secondary)]">
              {{ featured.body }}
            </p>
            <div class="mt-6 flex flex-wrap gap-2">
              <span
                class="rounded-full bg-[var(--surface-mint)] px-3 py-1 text-xs font-medium text-[var(--brand-green-dark)]"
              >Your tickets</span>
              <span
                class="rounded-full bg-[var(--surface-muted)] px-3 py-1 text-xs font-medium text-[var(--text-secondary)]"
              >Staff triage</span>
              <span
                class="rounded-full bg-[var(--surface-muted)] px-3 py-1 text-xs font-medium text-[var(--text-secondary)]"
              >Product issues</span>
            </div>
          </article>
          <article
            class="land-in land-delay-6 landing-card group rounded-2xl border border-[var(--border-subtle)] bg-[var(--surface-muted)]/60 p-7 shadow-sm transition duration-300"
          >
            <p class="text-base font-semibold text-[var(--text-primary)]">
              {{ sideCard.title }}
            </p>
            <p class="mt-3 text-sm leading-relaxed text-[var(--text-secondary)]">
              {{ sideCard.body }}
            </p>
          </article>
          <article
            v-for="(card, i) in bottomCards"
            :key="card.title"
            class="landing-card group rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-sm transition duration-300"
            :class="landDelayClass(7 + i)"
          >
            <p class="text-sm font-semibold text-[var(--text-primary)]">
              {{ card.title }}
            </p>
            <p class="mt-2 text-sm leading-relaxed text-[var(--text-secondary)]">
              {{ card.body }}
            </p>
          </article>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
const previewRows = [
  {
    title: 'VPN drops on Wi‑Fi handoff',
    meta: 'Technical · updated 2h ago',
    badge: 'Open',
    badgeClass: 'bg-neutral-200 text-neutral-800',
  },
  {
    title: 'Invoice PDF not attaching',
    meta: 'Billing',
    badge: 'Progress',
    badgeClass: 'bg-amber-100 text-amber-900',
  },
]

const stats = [
  { value: '2h 10m', label: 'Typical first response', hint: 'Sample metric — not live data' },
  { value: '99.2%', label: 'Requests closed on time', hint: 'Illustrative SLA-style figure' },
  { value: '28', label: 'Issues resolved this week', hint: 'Static preview only' },
]

const featured = {
  title: 'One place for your support requests',
  body:
    'You see the tickets you opened and their status; SecWeb staff see the queue on their side. Same calm layout whether you’re reporting a bug or asking for help.',
}

const sideCard = {
  title: 'Open a ticket in minutes',
  body: 'Describe the issue, pick a category if needed, and submit. No clutter — just what product support portals are meant to do.',
}

const bottomCards = [
  {
    title: 'Account required',
    body: 'Every reporter signs in — so we know who to reply to and can keep your history in one place.',
  },
  {
    title: 'Staff triage',
    body: 'Internal roles assign and resolve tickets; you stay informed as status changes.',
  },
  {
    title: 'Built for SecWeb CA',
    body: 'UI preview for Secure Web Development — wire to the real API when you connect the backend.',
  },
]

function landDelayClass(n: number) {
  return `land-in land-delay-${n}`
}
</script>

<style scoped>
.landing-grid {
  background-image: radial-gradient(circle at center, var(--border-subtle) 1px, transparent 1px);
  background-size: 24px 24px;
}

.landing-blob {
  animation: blob-drift 18s ease-in-out infinite alternate;
}

.landing-blob:nth-child(2) {
  animation-duration: 22s;
  animation-delay: -4s;
}

@keyframes blob-drift {
  0% {
    transform: translate(0, 0) scale(1);
  }
  100% {
    transform: translate(24px, 16px) scale(1.05);
  }
}

.land-in {
  animation: land-in 0.75s cubic-bezier(0.22, 1, 0.36, 1) both;
}

.land-delay-2 {
  animation-delay: 0.1s;
}
.land-delay-3 {
  animation-delay: 0.18s;
}
.land-delay-4 {
  animation-delay: 0.24s;
}
.land-delay-5 {
  animation-delay: 0.3s;
}
.land-delay-6 {
  animation-delay: 0.36s;
}
.land-delay-7 {
  animation-delay: 0.42s;
}
.land-delay-8 {
  animation-delay: 0.48s;
}
.land-delay-9 {
  animation-delay: 0.54s;
}

@keyframes land-in {
  from {
    opacity: 0;
    transform: translateY(18px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.landing-cta-primary,
.landing-cta-secondary {
  transform: translateZ(0);
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease,
    filter 0.2s ease,
    background-color 0.2s ease;
}

.landing-cta-primary:hover,
.landing-cta-secondary:hover {
  transform: scale(1.02);
}

.landing-cta-primary:active,
.landing-cta-secondary:active {
  transform: scale(0.98);
}

.landing-card {
  transform: translateZ(0);
}

.landing-card:hover {
  border-color: color-mix(in srgb, var(--brand-green) 28%, var(--border-subtle));
  box-shadow:
    0 12px 40px -24px rgba(0, 0, 0, 0.18),
    0 0 0 1px color-mix(in srgb, var(--brand-green) 12%, transparent);
}
</style>

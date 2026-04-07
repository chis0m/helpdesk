/** App routes — landing, auth, and dashboard. */
// VULN-02: `profile/:userId` — any numeric id can be opened in the SPA (backend IDOR completes the issue).
// VULN-04: `tickets/:id` — ticket detail/comments keyed only by id in the URL (backend IDOR completes the issue).
import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'
import { getAuthUserSnapshot } from '@/stores/auth-session'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: () => import('@/views/LandingView.vue'),
    },
    {
      path: '/login',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'login',
          component: () => import('@/views/LoginView.vue'),
        },
      ],
    },
    {
      path: '/signup',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'signup',
          component: () => import('@/views/SignupView.vue'),
        },
      ],
    },
    {
      path: '/accept-invite',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'accept-invite',
          meta: { title: 'Accept invite' },
          component: () => import('@/views/auth/AcceptInviteView.vue'),
        },
      ],
    },
    {
      path: '/forgot-password',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'forgot-password',
          meta: { title: 'Forgot password' },
          component: () => import('@/views/auth/ForgotPasswordView.vue'),
        },
      ],
    },
    {
      path: '/reset-password',
      component: () => import('@/layouts/AuthLayout.vue'),
      children: [
        {
          path: '',
          name: 'reset-password',
          meta: { title: 'Reset password' },
          component: () => import('@/views/auth/ResetPasswordView.vue'),
        },
      ],
    },
    {
      path: '/dashboard',
      component: AppLayout,
      children: [
        {
          path: 'tickets/new',
          name: 'ticket-new',
          meta: { title: 'New ticket' },
          component: () => import('@/views/tickets/TicketCreateView.vue'),
        },
        {
          path: 'tickets/:id',
          name: 'ticket-detail',
          meta: { title: 'Ticket detail' },
          component: () => import('@/views/tickets/TicketDetailView.vue'),
        },
        {
          path: 'tickets',
          name: 'dashboard-tickets',
          meta: { title: 'Tickets & search' },
          component: () => import('@/views/tickets/TicketsListView.vue'),
        },
        {
          path: 'profile',
          redirect() {
            const u = getAuthUserSnapshot()
            if (!u)
              return { path: '/login', query: { redirect: '/dashboard/profile' } }
            return { name: 'dashboard-profile', params: { userId: String(u.user_id) } }
          },
        },
        {
          path: 'profile/:userId',
          name: 'dashboard-profile',
          meta: { title: 'Your profile' },
          component: () => import('@/views/account/ProfileView.vue'),
        },
        {
          path: 'sessions',
          name: 'dashboard-sessions',
          meta: { title: 'Sessions' },
          component: () => import('@/views/account/SessionsView.vue'),
        },
        {
          path: 'admin/users',
          name: 'admin-users',
          meta: { title: 'Users' },
          component: () => import('@/views/admin/AdminUsersListView.vue'),
        },
        {
          path: 'admin/staff/new',
          name: 'admin-staff-new',
          meta: { title: 'Create staff' },
          component: () => import('@/views/admin/AdminStaffView.vue'),
        },
        {
          path: 'admin/staff',
          name: 'admin-staff',
          meta: { title: 'Staff' },
          component: () => import('@/views/admin/AdminStaffListView.vue'),
        },
        {
          path: '',
          name: 'dashboard',
          meta: { title: 'Support home' },
          component: () => import('@/views/DashboardView.vue'),
        },
      ],
    },
  ],
})

export default router

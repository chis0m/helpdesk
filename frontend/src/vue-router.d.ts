/** Module augmentation — `import 'vue-router'` is required so we extend the real package instead of replacing it. */
import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    /** When true, only `admin` or `super_admin` may open the route (see `adminRoutesGuard`). */
    requiresAdmin?: boolean
  }
}

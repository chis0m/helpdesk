/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_PORT?: string
  /** Backend origin, e.g. http://localhost:8080 */
  readonly VITE_API_BASE_URL?: string
  /** When "true", logger.debug runs even for production-mode preview builds */
  readonly VITE_DEBUG?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

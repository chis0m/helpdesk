import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import checker from 'vite-plugin-checker'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [
      vue(),
      tailwindcss(),
      // TypeScript checking during development
      checker({
        typescript: true,
        vueTsc: {
          tsconfigPath: 'tsconfig.app.json',
        },
        // Show errors in browser overlay and terminal
        overlay: {
          initialIsOpen: false,
          position: 'tl',
          badgeStyle: 'position: fixed; top: 20px; right: 20px; z-index: 9999;',
        },
      }),
    ],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    server: {
      port: Number.parseInt(env.VITE_PORT || '3000', 10),
      host: '0.0.0.0',
      open: false,
      allowedHosts: ['secweb.local'],
    },
  }
})

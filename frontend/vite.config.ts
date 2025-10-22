import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    watch: {
      usePolling: true,
    },
    hmr: {
      host: 'localhost',
    },
    // Only add proxy in development mode
    ...(process.env.NODE_ENV === 'development' && {
      proxy: {
        '/api': {
          // In Docker dev environment, use service name; otherwise use localhost
          target: process.env.DOCKER_ENV ? 'http://backend:8080' : 'http://localhost:8080',
          changeOrigin: true,
          secure: false,
          configure: (proxy, options) => {
            console.log('🔧 Vite Proxy Configured:')
            console.log('  - Target:', options.target)
            console.log('  - NODE_ENV:', process.env.NODE_ENV)
            console.log('  - VITE_API_URL:', process.env.VITE_API_URL)
            console.log('  - DOCKER_ENV:', process.env.DOCKER_ENV)
          }
        },
      },
    }),
  },
})

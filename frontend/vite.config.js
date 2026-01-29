import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    host: '0.0.0.0',
    port: 3000,
    watch: {
      usePolling: true
    },
    // Allow requests from production domain
    allowedHosts: [
      'juliepogue.anhelm.com',
      'juliepogue.com',
      'localhost'
    ]
  },
  envPrefix: ['VITE_', 'SITE_URL']
})

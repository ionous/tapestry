const { resolve } = require('path')
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  // exposes the specified object as a global to the app
  define: {
    appcfg: JSON.stringify({
      mosaic: 'http://localhost:8080',
      live: 'http://localhost:8088',
    })
  },
  build: {
    rollupOptions: {
      input: {
        live: resolve(__dirname, 'live/index.html'),
        mosaic: resolve(__dirname, 'mosaic/index.html')
      }
    }
  }
})

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  define: {
    appcfg: JSON.stringify({
      mosaic: 'http://localhost:8080',
      live: 'http://localhost:8088',
    })
  }
})

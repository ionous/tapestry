import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  assetsInclude: [ "**/*.gob", "**/*.wasm"],
  // exposes the specified object as a global to the app
    define: {
      appcfg: JSON.stringify({
        shuttle: "wasm"
      }),
    },
})

const { resolve } = require("path");
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  // exposes the specified object as a global to the app
  define: {
    appcfg: JSON.stringify({
      // run by cmd/serve. somebody has to fix these names.
      live: "http://localhost:8088",
      // run by cmd/mosaic.
      mosaic: "http://localhost:8080",
    }),
  },
  build: {
    chunkSizeWarningLimit: 825, // in kb; see notes on blockly below.
    rollupOptions: {
      input: {
        // a terrible name for playing the game.
        live: resolve(__dirname, "live/index.html"),
        // the editor.
        mosaic: resolve(__dirname, "mosaic/index.html"),
      },
      output: {
        // this splits the individual node modules into separate outputs
        // blockly itself is 805.07 KiB / gzip: 179.20 KiB
        manualChunks(id) {
          if (id.includes("node_modules/blockly")) {
            return id
              .toString()
              .split("node_modules/")[1]
              .split("/")[0]
              .toString();
          }
        },
      },
    },
  },
});

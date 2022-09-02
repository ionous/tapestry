const { resolve } = require("path");
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig(({ command /*, mode, ssrBuild*/ }) => {
  const build = command === "build";
  return {
    plugins: [vue()],
    // exposes the specified object as a global to the app
    define: {
      appcfg: JSON.stringify({
        // run by cmd/serve. somebody has to fix these names.
        // plus that person should probably make these the same port.
        live: build ? "http://wails.localhost" : "http://localhost:8088",
        // run by cmd/mosaic.
        mosaic: build ? "http://wails.localhost" : "http://localhost:8080",
      }),
    },
    server: {
      port: 3000, //  fix: default is now 5173
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
  };
});

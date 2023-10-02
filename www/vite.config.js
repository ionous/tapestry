const { resolve } = require("path");
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig(({ command /*, mode, ssrBuild*/ }) => {
  const build = command === "build";
  const backend = build ? "wails://wails" : "http://localhost:8080";
  return {
    plugins: [vue()],
    // exposes the specified object as a global to the app
    define: {
      appcfg: JSON.stringify({
        shuttle: backend + "/shuttle/",
        mosaic: backend,
      }),
    },
    server: {
      port: 3000, //  fix: default is now 5173
    },
    build: {
      chunkSizeWarningLimit: 825, // in kb; see notes on blockly below.
      rollupOptions: {
        input: {
          // backend for playing the game.
          play: resolve(__dirname, "play/index.html"),
          // the editor.
          mosaic: resolve(__dirname, "index.html"),
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

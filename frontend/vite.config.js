import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueDevTools from "vite-plugin-vue-devtools";

export default defineConfig({
  // 'base' ensures that all script/style tags in index.html
  // start with /assets/ so Gin can find them easily.
  base: "/",

  plugins: [vue(), vueDevTools()],

  server: {
    host: "0.0.0.0",
    port: 5173, // This is only for your local 'npm run dev'
    proxy: {
      // During development, this redirects localhost:5173/api
      // to your Go server on port 5000.
      "/api": {
        target: "http://localhost:5000",
        changeOrigin: true,
        secure: false,
      },
      // If you are using WebSockets (as your Go code suggests)
      "/ws": {
        target: "ws://localhost:5000",
        ws: true,
        changeOrigin: true,
      },
    },
  },

  build: {
    // This ensures your output goes exactly where your go:embed expects it
    outDir: "dist",
    emptyOutDir: true,
    // This groups JS/CSS into an 'assets' folder inside dist
    assetsDir: "assets",
  },

  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
});

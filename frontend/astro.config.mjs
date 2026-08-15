// @ts-check
import { defineConfig } from 'astro/config';

// https://astro.build/config
export default defineConfig({
  vite: {
    server: {
      // Forward /api/* to the local Go backend (http://localhost:8080) so the
      // browser can call SAME-ORIGIN /api/... with no CORS. Vite's dev-server
      // proxy is the reliable way to do this (Astro runs on Vite under the hood).
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
      },
      allowedHosts: ['frontend.example.com', '.example.com', 'localhost', '127.0.0.1'] // dev / server
    },
    preview: {
      allowedHosts: ['frontend.example.com', '.example.com', 'localhost', '127.0.0.1'] // preview
    }
  }
});

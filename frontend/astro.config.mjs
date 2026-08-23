// @ts-check
import { defineConfig } from 'astro/config';

// Where /api/* is forwarded to. Switch between the local Go backend and the
// K3s ingress WITHOUT editing code:
//
//   local : npm run dev            → http://localhost:8080 (default)
//   K3s   : npm run dev:k3s        → https://frontend.example.com (ingress)
//   other : API_PROXY_TARGET=http://... npm run dev
//
// The browser still calls SAME-ORIGIN /api/... so there is no CORS in any mode.
const API_PROXY_TARGET = process.env.API_PROXY_TARGET || 'http://localhost:8080';

// https://astro.build/config
export default defineConfig({
  vite: {
    server: {
      proxy: {
        '/api': {
          target: API_PROXY_TARGET,
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

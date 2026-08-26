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

// Hosts the dev/preview servers accept (Vite host check). Inject the real
// domain at build/run time without editing code, e.g.:
//   ALLOWED_HOSTS="frontend.example.com,.example.com" docker build ...
//   docker run -e ALLOWED_HOSTS=frontend.example.com .
const ALLOWED_HOSTS = (process.env.ALLOWED_HOSTS || 'localhost,127.0.0.1')
  .split(',')
  .map((h) => h.trim())
  .filter(Boolean);

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
      allowedHosts: ALLOWED_HOSTS // dev / server
    },
    preview: {
      allowedHosts: ALLOWED_HOSTS // preview
    }
  }
});

// @ts-check
import 'dotenv/config'; // ← บรรทัดนี้ใหม่: โหลด frontend/.env → process.env
import { defineConfig } from 'astro/config';

// Where /api/* is forwarded to. Priority (highest first):
//   1. shell env   : API_PROXY_TARGET=http://... npm run dev
//   2. frontend/.env : API_PROXY_TARGET=...   (loaded via dotenv above)
//   3. default     : http://localhost:8080
//
// NOTE: astro.config.mjs is evaluated before Vite loads .env into
// import.meta.env — so we load it explicitly here with dotenv (which does
// NOT overwrite variables already set in the shell).

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

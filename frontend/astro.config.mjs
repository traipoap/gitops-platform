// @ts-check
import { defineConfig } from 'astro/config';

// https://astro.build/config
export default defineConfig({vite: {

    server: {
      allowedHosts: ['frontend.example.com', '.example.com'] // dev / server
    },
    preview: {
      allowedHosts: ['frontend.example.com', '.example.com'] // preview
    }
  }
});

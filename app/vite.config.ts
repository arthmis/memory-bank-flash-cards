import { defineConfig } from 'vite';
import solidPlugin from 'vite-plugin-solid';
import { tanstackRouter } from '@tanstack/router-plugin/vite';

export default defineConfig({
  plugins: [
    tanstackRouter({
      target: 'solid',
      autoCodeSplitting: true,
      verboseFileRoutes: false,
      // apiBase: "/dashboard/",
    }),
    solidPlugin(),
  ],
  server: {
    port: 3000,
  },
  build: {
    target: 'esnext',
  },
  // base: "/dashboard/"
});

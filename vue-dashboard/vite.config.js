import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import path from 'path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  build: {
    outDir: '../dist', // Path relative to the current directory
    emptyOutDir: true, // Ensures the output directory is cleared before build
  },
  css: {
    preprocessorOptions: {
      scss: {

        additionalData: `        
          @use "~sass/main" as *;
          @use "~sass/abstracts/_index" as *;
          `,
        api: 'modern',
      },
    }
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '~sass': path.resolve(__dirname, 'src/assets/sass'),
    },
  },
})

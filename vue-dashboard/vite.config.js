import { fileURLToPath, URL } from 'node:url';
import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueDevTools from 'vite-plugin-vue-devtools';
import path from 'path';
import dotenv from 'dotenv';

// Custom path to the `.env` file
const customEnvPath = path.resolve(__dirname, '../.env');
dotenv.config({ path: customEnvPath });

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');

  return {
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
      },
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
        '~sass': path.resolve(__dirname, 'src/assets/sass'),
      },
    },
    define: {
      'process.env': env, // Expose environment variables to your app
    },
  };
});

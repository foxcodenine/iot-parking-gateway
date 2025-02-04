### Step 1. Create a Vue Project with vite

```bash
npm init vite@latest
√ Project name: ... vue.client
√ Select a framework: » Vue
√ Select a variant: » Customize with create-vue ↗

Vue.js - The Progressive JavaScript Framework

√ Add TypeScript? ... No
√ Add JSX Support? ... No
√ Add Vue Router for Single Page Application development? ... Yes
√ Add Pinia for state management? ... Yes
√ Add Vitest for Unit Testing? ... No 
√ Add an End-to-End Testing Solution? » No
√ Add ESLint for code quality? ... No 
√ Add Vue DevTools 7 extension for debugging? (experimental) ... No 
```

### Step 2. Install sass:

```bash
npm install sass@latest sass-loader@latest
```

### Step 3. create main sass file:

```bash
code src/assets/sass/main.scss
code src/assets/sass/abstracts/_index.scss
```

### Step 4. Update vite.config.js: 

```js
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
```

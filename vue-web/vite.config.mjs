import vue from '@vitejs/plugin-vue'
import { vite as vidstack } from 'vidstack/plugins'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => tag.startsWith('media-'),
        },
      },
    }),
    vidstack({ include: /src\// }),
  ],
})

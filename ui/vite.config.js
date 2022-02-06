import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      src: path.resolve('./src')
    }
  },
  plugins: [svelte()],
  server: {
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})

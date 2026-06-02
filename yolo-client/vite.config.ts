import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 8000,
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  build: {
    outDir: '../yolo-cli/internal/frontend/client-dist',
    emptyOutDir: true,
  },
})

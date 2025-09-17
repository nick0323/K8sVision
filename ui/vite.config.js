import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    host: '0.0.0.0', // 监听所有IP地址
    strictPort: false, // 如果端口被占用，自动尝试下一个可用端口
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})

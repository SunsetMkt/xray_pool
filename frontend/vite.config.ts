import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import eslintPlugin from 'vite-plugin-eslint';
import path from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  server: {
    port: 5173,
    hmr: {
      overlay: true,
    },
  },

  plugins: [eslintPlugin(), vue({})],

  resolve: {
    extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue'],
    alias: {
      '@': path.resolve(__dirname, 'src'),
      config: path.resolve(__dirname, 'src/config'),
      src: path.resolve(__dirname, 'src'),
      app: path.resolve(__dirname, '/'),
      components: path.resolve(__dirname, 'src/components'),
      layouts: path.resolve(__dirname, 'src/layouts'),
      pages: path.resolve(__dirname, 'src/pages'),
      assets: path.resolve(__dirname, 'src/assets'),
      stores: path.resolve(__dirname, 'src/stores'),
    },
  },
});

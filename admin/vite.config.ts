import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import Components from 'unplugin-vue-components/vite';
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers';
import { resolve } from 'path';
import config from './src/configs/index.js';

export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
    Components({
      resolvers: [
        AntDesignVueResolver({
          importStyle: false,
        }),
      ],
    }),
  ],
  
  resolve: {
    alias: {
      '#': resolve(__dirname, 'src'),
      '@': resolve(__dirname, 'src'),
      '@vben/preferences': resolve(__dirname, 'src/vben-shims/preferences.ts'),
      '@vben/utils': resolve(__dirname, 'src/vben-shims/utils.ts'),
      '@vben/access': resolve(__dirname, 'src/vben-shims/access.ts'),
      '@vben/stores': resolve(__dirname, 'src/vben-shims/stores.ts'),
      '@vben/hooks': resolve(__dirname, 'src/vben-shims/hooks.ts'),
      '@vben/request': resolve(__dirname, 'src/vben-shims/request.ts'),
      '@vben/locales': resolve(__dirname, 'src/vben-shims/locales.ts'),
      '@vben/icons': resolve(__dirname, 'src/vben-shims/icons.ts'),
      '@vben/constants': resolve(__dirname, 'src/vben-shims/constants.ts'),
      '@vben/types': resolve(__dirname, 'src/vben-shims/types.ts'),
      '@vben/common-ui': resolve(__dirname, 'src/vben-shims/utils.ts'),
      '@vben/common-ui/es/loading': resolve(__dirname, 'src/vben-shims/utils.ts'),
      '@vben/common-ui/es/tippy': resolve(__dirname, 'src/vben-shims/utils.ts'),
      '@vben/plugins/motion': resolve(__dirname, 'src/vben-shims/utils.ts'),
      '@vben/plugins/vxe-table': resolve(__dirname, 'src/vben-shims/utils.ts'),
      '@vben/layouts': resolve(__dirname, 'src/vben-shims/utils.ts'),
      '@vben/styles': resolve(__dirname, 'src/vben-shims/utils.ts'),
    },
  },

  define: {
    'import.meta.env.VITE_APP_TITLE': JSON.stringify(config.APP_TITLE),
    'import.meta.env.VITE_APP_NAMESPACE': JSON.stringify(config.APP_NAMESPACE),
    'import.meta.env.VITE_APP_VERSION': JSON.stringify(config.APP_VERSION),
    'import.meta.env.VITE_GLOB_API_URL': JSON.stringify(config.API_BASE_URL),
    'import.meta.env.VITE_PORT': JSON.stringify(config.PORT),
    'import.meta.env.VITE_BASE': JSON.stringify(config.BASE),
    'import.meta.env.VITE_NITRO_MOCK': JSON.stringify(config.NITRO_MOCK),
    'import.meta.env.VITE_DEVTOOLS': JSON.stringify(config.DEVTOOLS),
    'import.meta.env.VITE_INJECT_APP_LOADING': JSON.stringify(config.INJECT_APP_LOADING),
    'import.meta.env.VITE_APP_STORE_SECURE_KEY': JSON.stringify(config.STORE_SECURE_KEY),
  },

  server: {
    port: config.PORT,
    host: '0.0.0.0',
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8989',
        changeOrigin: true,
        ws: true,
      },
    },
  },

  build: {
    target: 'es2015',
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: false,
    chunkSizeWarningLimit: 1500,
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          'antd-vendor': ['ant-design-vue', '@ant-design/icons-vue'],
        },
      },
    },
  },
});

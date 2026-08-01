import { defineConfig } from '@vben/vite-config';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      server: {
        host: '0.0.0.0',
        proxy: {
          '/api': {
            changeOrigin: true,
            target: 'http://127.0.0.1:8989',
            ws: true,
          },
        },
      },
    },
  };
});

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import zhCn from 'element-plus/es/locale/lang/zh-cn';

import App from './app.vue';
import { router } from './router';

async function bootstrap(namespace: string) {
  console.log('Bootstrap namespace:', namespace);

  const app = createApp(App);

  // 注册 Pinia
  const pinia = createPinia();
  app.use(pinia);

  // 注册 Element Plus
  app.use(ElementPlus, { locale: zhCn });

  // 配置路由
  app.use(router);

  app.mount('#app');
}

export { bootstrap };

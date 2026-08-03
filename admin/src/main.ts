import { unmountGlobalLoading } from './vben-shims/utils';

/**
 * 应用初始化完成之后再进行页面加载渲染
 */
async function initApplication() {
  const env = import.meta.env.PROD ? 'prod' : 'dev';
  const appVersion = import.meta.env.VITE_APP_VERSION || '1.0.0';
  const namespace = `${import.meta.env.VITE_APP_NAMESPACE || 'short-video-leads-admin'}-${appVersion}-${env}`;

  // 启动应用并挂载
  const { bootstrap } = await import('./bootstrap');
  await bootstrap(namespace);

  // 移除并销毁loading
  unmountGlobalLoading();
}

initApplication();

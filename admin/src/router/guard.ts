import type { Router } from 'vue-router';
import { useUserStore } from '../vben-shims/stores';

const LOGIN_PATH = '/auth/login';

/**
 * 项目守卫配置
 * @param router
 */
function createRouterGuard(router: Router) {
  router.beforeEach((to, _from, next) => {
    const userStore = useUserStore();
    const token = userStore.token;

    // 如果是认证相关页面，直接放行
    if (to.path.startsWith('/auth')) {
      next();
      return;
    }

    // 如果没有token，跳转到登录页
    if (!token) {
      next({ path: LOGIN_PATH, query: { redirect: to.fullPath } });
      return;
    }

    // 有token，放行
    next();
  });
}

export { createRouterGuard };

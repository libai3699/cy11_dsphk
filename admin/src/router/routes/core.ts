import type { RouteRecordRaw } from 'vue-router';

const LOGIN_PATH = '/auth/login';
const BasicLayout = () => import('#/layouts/basic.vue');
const AuthPageLayout = () => import('#/layouts/auth.vue');

/** 全局404页面 */
const fallbackNotFoundRoute: RouteRecordRaw = {
  component: () => import('#/views/_core/fallback/not-found.vue'),
  meta: { title: '404' },
  name: 'FallbackNotFound',
  path: '/:path(.*)*',
};

const coreRoutes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: { title: 'Root' },
    name: 'Root',
    path: '/',
    redirect: '/workspace',
    children: [],
  },
  {
    component: AuthPageLayout,
    meta: { title: 'Authentication' },
    name: 'Authentication',
    path: '/auth',
    redirect: LOGIN_PATH,
    children: [
      {
        name: 'Login',
        path: 'login',
        component: () => import('#/views/_core/authentication/login.vue'),
        meta: { title: '登录' },
      },
    ],
  },
];

export { coreRoutes, fallbackNotFoundRoute };

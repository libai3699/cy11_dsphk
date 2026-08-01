import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    name: 'WorkspacePage',
    path: '/workspace',
    component: () => import('#/views/dashboard/WorkspacePage.vue'),
    meta: {
      affixTab: true,
      icon: 'carbon:dashboard',
      order: -1,
      title: '仪表盘',
    },
  },
];

export default routes;

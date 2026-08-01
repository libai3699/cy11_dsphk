import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: { icon: 'carbon:user-multiple', order: 1, title: '商家管理' },
    name: 'UserMgmt',
    path: '/users',
    children: [
      {
        name: 'UserListPage',
        path: '/users/list',
        component: () => import('#/views/users/UserListPage.vue'),
        meta: { icon: 'carbon:user', title: '商家列表' },
      },
      {
        name: 'MerchantProcessPage',
        path: '/users/detail/:id',
        component: () => import('#/views/users/MerchantProcessPage.vue'),
        meta: {
          hideInMenu: true,
          icon: 'carbon:flow',
          title: '商家推进流程',
        },
      },
      {
        name: 'UserDevicesPage',
        path: '/users/devices',
        component: () => import('#/views/users/UserDevicesPage.vue'),
        meta: { icon: 'carbon:mobile', title: '账号授权' },
      },
      {
        name: 'AccountDiagnosisPage',
        path: '/users/account-diagnosis',
        component: () => import('#/views/users/AccountDiagnosisPage.vue'),
        meta: { icon: 'carbon:analytics', title: '账号诊断' },
      },
      {
        name: 'DurationLogPage',
        path: '/users/duration-logs',
        component: () => import('#/views/users/DurationLogPage.vue'),
        meta: { icon: 'carbon:time', title: '跟进记录' },
      },
    ],
  },
  {
    meta: { icon: 'carbon:purchase', order: 2, title: '成交与分成' },
    name: 'PlanMgmt',
    path: '/plans',
    children: [
      {
        name: 'PlanListPage',
        path: '/plans/list',
        component: () => import('#/views/plans/PlanListPage.vue'),
        meta: { icon: 'carbon:list', title: '团购套餐' },
      },
      {
        name: 'PlanOrdersPage',
        path: '/plans/orders',
        component: () => import('#/views/plans/PlanOrdersPage.vue'),
        meta: { icon: 'carbon:document', title: '分成订单' },
      },
    ],
  },
  {
    meta: { icon: 'carbon:network-4', order: 3, title: '对标分析' },
    name: 'LineMgmt',
    path: '/lines',
    children: [
      {
        name: 'LineListPage',
        path: '/lines/list',
        component: () => import('#/views/lines/LineListPage.vue'),
        meta: { icon: 'carbon:network-4', title: '对标账号库' },
      },
    ],
  },
  {
    meta: { icon: 'carbon:settings', order: 4, title: '内容生产' },
    name: 'ContentMgmt',
    path: '/content',
    children: [
      {
        name: 'ContentNoticesPage',
        path: '/content/notices',
        component: () => import('#/views/content/ContentNoticesPage.vue'),
        meta: { icon: 'carbon:notification', title: '选题中心' },
      },
      {
        name: 'ContentQuotesPage',
        path: '/content/quotes',
        component: () => import('#/views/content/ContentQuotesPage.vue'),
        meta: { icon: 'carbon:quotes', title: '文案脚本' },
      },
      {
        name: 'ContentDiscoveriesPage',
        path: '/content/discoveries',
        component: () => import('#/views/content/ContentDiscoveriesPage.vue'),
        meta: { icon: 'carbon:compass', title: '分镜脚本' },
      },
      {
        name: 'ContentPaymentsPage',
        path: '/content/payments',
        component: () => import('#/views/content/ContentPaymentsPage.vue'),
        meta: { icon: 'carbon:wallet', title: '发布排期' },
      },
      {
        name: 'ContentConfigsPage',
        path: '/content/configs',
        component: () => import('#/views/content/ContentConfigsPage.vue'),
        meta: { icon: 'carbon:settings-adjust', title: '系统规则' },
      },
    ],
  },
  {
    meta: { icon: 'carbon:document', order: 5, title: '执行复盘' },
    name: 'LogMgmt',
    path: '/logs',
    redirect: '/logs/user',
    children: [
      {
        name: 'LogUserPage',
        path: '/logs/user',
        component: () => import('#/views/logs/LogUserPage.vue'),
        meta: { icon: 'carbon:user-activity', title: '拍摄任务' },
      },
      {
        name: 'LogAdminPage',
        path: '/logs/admin',
        component: () => import('#/views/logs/LogAdminPage.vue'),
        meta: { icon: 'carbon:security', title: '数据复盘' },
      },
    ],
  },
];

export default routes;

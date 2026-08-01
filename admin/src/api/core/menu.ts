import type { RouteRecordStringComponent } from '@vben/types';

export async function getAllMenusApi(): Promise<RouteRecordStringComponent[]> {
  return Promise.resolve([
    {
      name: 'WorkspacePage',
      path: '/workspace',
      component: '/dashboard/WorkspacePage',
      meta: { affixTab: true, icon: 'carbon:dashboard', order: -1, title: '仪表盘' },
    },
    {
      meta: { icon: 'carbon:user-multiple', order: 1, title: '用户管理' },
      name: 'UserMgmt',
      path: '/users',
      children: [
        { name: 'UserListPage', path: '/users/list', component: '/users/UserListPage', meta: { icon: 'carbon:user', title: '用户列表' } },
        { name: 'UserDevicesPage', path: '/users/devices', component: '/users/UserDevicesPage', meta: { icon: 'carbon:mobile', title: '设备管理' } },
      ],
    },
    {
      meta: { icon: 'carbon:purchase', order: 2, title: '套餐管理' },
      name: 'PlanMgmt',
      path: '/plans',
      children: [
        { name: 'PlanListPage', path: '/plans/list', component: '/plans/PlanListPage', meta: { icon: 'carbon:list', title: '套餐列表' } },
        { name: 'PlanOrdersPage', path: '/plans/orders', component: '/plans/PlanOrdersPage', meta: { icon: 'carbon:document', title: '订单管理' } },
      ],
    },
    {
      meta: { icon: 'carbon:network-4', order: 3, title: '线路管理' },
      name: 'LineMgmt',
      path: '/lines',
      children: [
        { name: 'LineListPage', path: '/lines/list', component: '/lines/LineListPage', meta: { icon: 'carbon:network-4', title: '线路列表' } },
      ],
    },
    {
      meta: { icon: 'carbon:settings', order: 4, title: '内容管理' },
      name: 'ContentMgmt',
      path: '/content',
      children: [
        { name: 'ContentNoticesPage', path: '/content/notices', component: '/content/ContentNoticesPage', meta: { icon: 'carbon:notification', title: '公共通知' } },
        { name: 'ContentConfigsPage', path: '/content/configs', component: '/content/ContentConfigsPage', meta: { icon: 'carbon:settings-adjust', title: '系统配置' } },
      ],
    },
    {
      meta: { icon: 'carbon:document', order: 5, title: '日志管理' },
      name: 'LogMgmt',
      path: '/logs',
      children: [
        { name: 'LogUserPage', path: '/logs/user', component: '/logs/LogUserPage', meta: { icon: 'carbon:user-activity', title: '用户登录日志' } },
        { name: 'LogAdminPage', path: '/logs/admin', component: '/logs/LogAdminPage', meta: { icon: 'carbon:security', title: '后台登录日志' } },
      ],
    },
  ]);
}

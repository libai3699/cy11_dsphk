import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    meta: { icon: 'carbon:bookmark', order: 6, title: '运营知识库' },
    name: 'KnowledgeMgmt',
    path: '/knowledge',
    children: [
      {
        name: 'PainPointList',
        path: '/knowledge/pain-points',
        component: () => import('#/views/knowledge/PainPointList.vue'),
        meta: { icon: 'carbon:idea', title: '痛点库' },
      },
      {
        name: 'CaseStudyList',
        path: '/knowledge/case-studies',
        component: () => import('#/views/knowledge/CaseStudyList.vue'),
        meta: { icon: 'carbon:document', title: '案例库' },
      },
      {
        name: 'AccountProfileList',
        path: '/knowledge/account-profiles',
        component: () => import('#/views/knowledge/AccountProfileList.vue'),
        meta: { icon: 'carbon:user-avatar', title: '账号画像库' },
      },
      {
        name: 'PlatformRuleList',
        path: '/knowledge/platform-rules',
        component: () => import('#/views/knowledge/PlatformRuleList.vue'),
        meta: { icon: 'carbon:rule', title: '平台规则库' },
      },
      {
        name: 'ContentTemplateList',
        path: '/knowledge/content-templates',
        component: () => import('#/views/knowledge/ContentTemplateList.vue'),
        meta: { icon: 'carbon:template', title: '内容模板库' },
      },
    ],
  },
];

export default routes;

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from '../vben-shims/stores';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const collapsed = ref(false);

const menuItems = [
  {
    key: '/workspace',
    icon: '📊',
    label: '仪表盘',
  },
  {
    key: '/users',
    icon: '👥',
    label: '用户管理',
    children: [
      { key: '/users/list', label: '用户列表' },
      { key: '/users/devices', label: '设备管理' },
    ],
  },
  {
    key: '/plans',
    icon: '📦',
    label: '套餐管理',
    children: [
      { key: '/plans/list', label: '套餐列表' },
      { key: '/plans/orders', label: '订单管理' },
    ],
  },
  {
    key: '/lines',
    icon: '🌐',
    label: '线路管理',
    children: [
      { key: '/lines/list', label: '线路列表' },
    ],
  },
  {
    key: '/content',
    icon: '⚙️',
    label: '内容管理',
    children: [
      { key: '/content/notices', label: '公共通知' },
      { key: '/content/quotes', label: '精选语录' },
      { key: '/content/payments', label: '支付配置' },
      { key: '/content/configs', label: '系统配置' },
    ],
  },
  {
    key: '/logs',
    icon: '📋',
    label: '日志管理',
    children: [
      { key: '/logs/user', label: '用户登录日志' },
      { key: '/logs/admin', label: '后台登录日志' },
    ],
  },
];

// 展开的子菜单
const openKeys = ref<string[]>([]);

// 当前激活的菜单项
const selectedKey = computed(() => route.path);

// 初始化展开状态
menuItems.forEach((item) => {
  if (item.children?.some((c) => route.path.startsWith(c.key))) {
    openKeys.value.push(item.key);
  }
});

function toggleSubmenu(key: string) {
  const idx = openKeys.value.indexOf(key);
  if (idx >= 0) {
    openKeys.value.splice(idx, 1);
  } else {
    openKeys.value.push(key);
  }
}

function navigate(path: string) {
  router.push(path);
}

async function handleLogout() {
  userStore.clearToken();
  router.push('/auth/login');
}
</script>

<template>
  <div class="layout-wrap">
    <!-- 侧边栏 -->
    <aside :class="['layout-sider', { collapsed }]">
      <div class="sider-logo">
        <span v-if="!collapsed">短视频获客</span>
        <span v-else>获客</span>
      </div>
      <nav class="sider-menu">
        <template v-for="item in menuItems" :key="item.key">
          <!-- 有子菜单 -->
          <template v-if="item.children">
            <div
              :class="['menu-group-title', { open: openKeys.includes(item.key) }]"
              @click="toggleSubmenu(item.key)"
            >
              <span class="menu-icon">{{ item.icon }}</span>
              <span v-if="!collapsed" class="menu-label">{{ item.label }}</span>
              <span v-if="!collapsed" class="menu-arrow">{{ openKeys.includes(item.key) ? '▾' : '▸' }}</span>
            </div>
            <div v-if="openKeys.includes(item.key) && !collapsed" class="menu-sub">
              <div
                v-for="child in item.children"
                :key="child.key"
                :class="['menu-item', { active: selectedKey === child.key }]"
                @click="navigate(child.key)"
              >
                {{ child.label }}
              </div>
            </div>
          </template>
          <!-- 无子菜单 -->
          <template v-else>
            <div
              :class="['menu-item top-level', { active: selectedKey === item.key }]"
              @click="navigate(item.key)"
            >
              <span class="menu-icon">{{ item.icon }}</span>
              <span v-if="!collapsed" class="menu-label">{{ item.label }}</span>
            </div>
          </template>
        </template>
      </nav>
    </aside>

    <!-- 主体 -->
    <div class="layout-main">
      <!-- 顶部栏 -->
      <header class="layout-header">
        <button class="collapse-btn" @click="collapsed = !collapsed">
          {{ collapsed ? '▶' : '◀' }}
        </button>
        <div class="header-right">
          <span class="username">管理员</span>
          <button class="logout-btn" @click="handleLogout">退出登录</button>
        </div>
      </header>
      <!-- 内容区 -->
      <main class="layout-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<style scoped>
.layout-wrap {
  display: flex;
  min-height: 100vh;
  background: #f0f2f5;
}

/* 侧边栏 */
.layout-sider {
  width: 220px;
  min-height: 100vh;
  background: #001529;
  color: #fff;
  transition: width 0.2s;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
}
.layout-sider.collapsed {
  width: 64px;
}

.sider-logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  border-bottom: 1px solid rgba(255,255,255,0.1);
  white-space: nowrap;
  overflow: hidden;
}

.sider-menu {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.menu-group-title {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  color: rgba(255,255,255,0.65);
  font-size: 14px;
  transition: all 0.2s;
  white-space: nowrap;
  overflow: hidden;
}
.menu-group-title:hover,
.menu-group-title.open {
  color: #fff;
  background: rgba(255,255,255,0.05);
}

.menu-item {
  padding: 10px 16px 10px 40px;
  cursor: pointer;
  color: rgba(255,255,255,0.65);
  font-size: 14px;
  transition: all 0.2s;
  white-space: nowrap;
  overflow: hidden;
}
.menu-item.top-level {
  padding-left: 16px;
  display: flex;
  align-items: center;
}
.menu-item:hover {
  color: #fff;
  background: rgba(255,255,255,0.05);
}
.menu-item.active {
  color: #fff;
  background: #1890ff;
}

.menu-sub {
  background: rgba(0,0,0,0.2);
}

.menu-icon {
  margin-right: 10px;
  font-size: 16px;
  flex-shrink: 0;
}
.menu-label {
  flex: 1;
}
.menu-arrow {
  margin-left: auto;
  font-size: 12px;
}

/* 主体 */
.layout-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.layout-header {
  height: 64px;
  background: #fff;
  display: flex;
  align-items: center;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.08);
  justify-content: space-between;
}

.collapse-btn {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  color: #666;
  padding: 4px 8px;
}
.collapse-btn:hover { color: #1890ff; }

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}
.username { color: #666; font-size: 14px; }
.logout-btn {
  padding: 6px 16px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}
.logout-btn:hover { color: #1890ff; border-color: #1890ff; }

.layout-content {
  flex: 1;
  padding: 24px;
  overflow: auto;
}
</style>

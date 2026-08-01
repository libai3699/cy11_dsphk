<script lang="ts" setup>
interface Props {
  status?: '403' | '404' | '500';
  title?: string;
  description?: string;
}

const props = withDefaults(defineProps<Props>(), {
  status: '404',
});

const statusConfig = {
  '403': {
    title: '403',
    description: '抱歉，您无权访问此页面',
  },
  '404': {
    title: '404',
    description: '抱歉，您访问的页面不存在',
  },
  '500': {
    title: '500',
    description: '抱歉，服务器出错了',
  },
};

const config = statusConfig[props.status];
</script>

<template>
  <div class="fallback-container">
    <div class="fallback-content">
      <h1 class="fallback-status">{{ props.title || config.title }}</h1>
      <p class="fallback-description">{{ props.description || config.description }}</p>
      <router-link to="/" class="fallback-button">返回首页</router-link>
    </div>
  </div>
</template>

<style scoped>
.fallback-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.fallback-content {
  text-align: center;
  color: white;
}

.fallback-status {
  font-size: 120px;
  font-weight: bold;
  margin: 0;
  line-height: 1;
}

.fallback-description {
  font-size: 24px;
  margin: 20px 0 40px;
  opacity: 0.9;
}

.fallback-button {
  display: inline-block;
  padding: 12px 32px;
  background: white;
  color: #667eea;
  text-decoration: none;
  border-radius: 4px;
  font-weight: 600;
  transition: all 0.3s;
}

.fallback-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}
</style>

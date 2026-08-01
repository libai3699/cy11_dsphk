<script lang="ts" setup>
import { useRouter } from 'vue-router';
const router = useRouter();
defineOptions({ name: 'Fallback404' });
</script>

<template>
  <div class="not-found">
    <div class="not-found-content">
      <h1>404</h1>
      <p>抱歉，您访问的页面不存在</p>
      <button @click="router.push('/')">返回首页</button>
    </div>
  </div>
</template>

<style scoped>
.not-found {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100%;
  padding: 48px;
}
.not-found-content {
  text-align: center;
}
.not-found-content h1 {
  font-size: 96px;
  font-weight: 700;
  color: #1890ff;
  margin: 0;
  line-height: 1;
}
.not-found-content p {
  font-size: 18px;
  color: #666;
  margin: 16px 0 32px;
}
.not-found-content button {
  padding: 10px 32px;
  background: #1890ff;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  cursor: pointer;
}
.not-found-content button:hover {
  background: #40a9ff;
}
</style>

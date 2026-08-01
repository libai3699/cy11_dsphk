<script lang="ts" setup>
import { reactive } from 'vue';

import { useAuthStore } from '#/store';

defineOptions({ name: 'Login' });

const authStore = useAuthStore();

const form = reactive({
  password: 'admin123456',
  username: 'admin',
});

async function submit() {
  await authStore.authLogin({ ...form });
}
</script>

<template>
  <main class="cy-login-page">
    <section class="cy-login-card">
      <div class="cy-login-brand">
        <img alt="内容获客后台" class="cy-login-logo" src="/logo.png" />
        <div>
          <h1>本地商家抖音获客后台</h1>
          <p>后台账号由系统分配，只保留登录入口。</p>
        </div>
      </div>

      <label>
        <span>账号</span>
        <input v-model="form.username" autocomplete="username" />
      </label>

      <label>
        <span>密码</span>
        <input
          v-model="form.password"
          autocomplete="current-password"
          type="password"
          @keyup.enter="submit"
        />
      </label>

      <p class="cy-login-tip">默认超管：admin / admin123456</p>

      <button
        class="cy-login-button"
        type="button"
        :disabled="authStore.loginLoading"
        @click="submit"
      >
        {{ authStore.loginLoading ? '登录中...' : '登录后台' }}
      </button>
    </section>
  </main>
</template>

<style scoped>
.cy-login-page {
  display: flex;
  min-height: 100vh;
  align-items: center;
  justify-content: center;
  background:
    radial-gradient(circle at top left, rgba(37, 99, 235, 0.14), transparent 32%),
    linear-gradient(135deg, #f7f8fb 0%, #eef4ff 100%);
  padding: 24px;
}

.cy-login-card {
  width: 420px;
  border: 1px solid rgba(148, 163, 184, 0.24);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 20px 60px rgba(15, 23, 42, 0.12);
  padding: 32px;
  backdrop-filter: blur(10px);
}

.cy-login-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.cy-login-logo {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  object-fit: cover;
  flex-shrink: 0;
}

.cy-login-brand h1 {
  margin: 0;
  color: #0f172a;
  font-size: 26px;
}

.cy-login-brand p {
  margin: 8px 0 0;
  color: #64748b;
  font-size: 14px;
}

.cy-login-card label {
  display: block;
  margin-bottom: 16px;
}

.cy-login-card span {
  display: block;
  margin-bottom: 8px;
  color: #334155;
  font-size: 14px;
}

.cy-login-card input {
  width: 100%;
  height: 40px;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  padding: 0 12px;
  outline: none;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease;
}

.cy-login-card input:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
}

.cy-login-tip {
  margin: 0 0 16px;
  color: #64748b;
  font-size: 13px;
  line-height: 1.6;
}

.cy-login-button {
  width: 100%;
  height: 44px;
  border: 0;
  border-radius: 10px;
  background: #2563eb;
  color: white;
  cursor: pointer;
  font-weight: 600;
}

.cy-login-button:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}
</style>

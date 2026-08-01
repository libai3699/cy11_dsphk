<script lang="ts" setup>
import { reactive, ref } from 'vue';

import {
  ElButton,
  ElCard,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
} from 'element-plus';
import 'element-plus/dist/index.css';

import { useAuthStore } from '#/store';

defineOptions({ name: 'Login' });

const authStore = useAuthStore();

const form = reactive({
  captcha: '',
  password: '',
  username: '',
});

const loading = ref(false);

async function onSubmit() {
  loading.value = true;
  try {
    await authStore.authLogin({ ...form });
  } catch {
    ElMessage.error('登录失败');
  } finally {
    loading.value = false;
  }
}

</script>

<template>
  <main class="cy-login-page">
    <section class="cy-login-panel">
      <ElCard class="cy-login-card" shadow="never">
        <div class="cy-login-title">
          <h2>后台登录</h2>
          <p>请输入账号、密码和 Google Authenticator 动态验证码</p>
        </div>

        <ElForm :model="form" label-position="top" @submit.prevent>
          <ElFormItem label="用户名">
            <ElInput v-model="form.username" size="large" />
          </ElFormItem>
          <ElFormItem label="密码">
            <ElInput v-model="form.password" show-password size="large" type="password" />
          </ElFormItem>
          <ElFormItem label="Google 验证码">
            <ElInput v-model="form.captcha" maxlength="6" size="large" />
          </ElFormItem>
          <ElButton
            class="cy-login-button"
            :loading="loading || authStore.loginLoading"
            native-type="button"
            size="large"
            type="primary"
            @click="onSubmit"
          >
            登录后台
          </ElButton>
        </ElForm>
      </ElCard>
    </section>
  </main>
</template>

<style scoped>
.cy-login-page {
  display: flex;
  width: 100vw;
  min-height: 100vh;
  align-items: center;
  justify-content: center;
  background: #f8fafc;
}

.cy-login-card {
  width: 360px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
}

.cy-login-title {
  margin-bottom: 24px;
}

.cy-login-title h2 {
  margin: 0;
  color: #0f172a;
  font-size: 24px;
  font-weight: 700;
}

.cy-login-title p {
  margin: 8px 0 0;
  color: #64748b;
  font-size: 14px;
}

.cy-login-button {
  width: 100%;
}
</style>

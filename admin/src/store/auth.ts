import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { defineStore } from 'pinia';
import { useUserStore } from '../vben-shims/stores';
import { loginApi } from '#/api';

export const useAuthStore = defineStore('auth', () => {
  const userStore = useUserStore();
  const router = useRouter();
  const route = useRoute();
  const loginLoading = ref(false);

  async function authLogin(params: any, onSuccess?: () => Promise<void> | void) {
    try {
      loginLoading.value = true;
      const { accessToken } = await loginApi(params);

      if (accessToken) {
        userStore.setToken(accessToken);
        
        if (onSuccess) {
          await onSuccess();
        } else {
          // 获取重定向路径
          const redirect = route.query.redirect as string;
          await router.push(redirect || '/workspace');
        }
      }
    } finally {
      loginLoading.value = false;
    }
  }

  async function logout(redirect: boolean = true) {
    userStore.clearToken();
    if (redirect) {
      await router.replace('/auth/login');
    }
  }

  function $reset() {
    loginLoading.value = false;
  }

  return {
    $reset,
    authLogin,
    loginLoading,
    logout,
  };
});

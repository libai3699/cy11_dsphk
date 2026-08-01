// @vben/stores 替代实现

import { defineStore } from 'pinia';

export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: null as any,
    token: localStorage.getItem('admin_token') || '',
  }),
  actions: {
    setToken(token: string) {
      this.token = token;
      localStorage.setItem('admin_token', token);
    },
    clearToken() {
      this.token = '';
      localStorage.removeItem('admin_token');
    },
  },
});

export const useAccessStore = defineStore('access', {
  state: () => ({
    accessCodes: [] as string[],
  }),
});

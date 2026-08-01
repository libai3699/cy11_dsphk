import type { App } from 'vue';
import { ref } from 'vue';
import zhCn from 'ant-design-vue/es/locale/zh_CN';

// 简化的国际化配置
export function $t(key: string) {
  return key;
}

export const antdLocale = ref(zhCn);

export async function setupI18n(app: App) {
  console.log('I18n initialized');
  return Promise.resolve();
}

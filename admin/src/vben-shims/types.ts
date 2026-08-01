// @vben/types 替代实现

export interface RouteRecordStringComponent {
  name: string;
  path: string;
  component?: string;
  meta?: any;
  children?: RouteRecordStringComponent[];
}

export interface UserInfo {
  username?: string;
  realName?: string;
  avatar?: string;
  roles?: string[];
  homePath?: string;
}

export type Recordable<T = any> = Record<string, T>;

export type ComponentRecordType = Record<string, () => Promise<any>>;

export interface GenerateMenuAndRoutesOptions {
  roles?: string[];
  router?: any;
  routes?: any[];
}

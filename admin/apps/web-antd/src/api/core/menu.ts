import type { RouteRecordStringComponent } from '@vben/types';

import { requestClient } from '#/api/request';

export async function getAllMenusApi(): Promise<RouteRecordStringComponent[]> {
  return requestClient.get<RouteRecordStringComponent[]>('/menus');
}

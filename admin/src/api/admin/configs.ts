import { requestClient } from '#/api/request';

export interface AppConfig {
  id: number;
  key_name: string;
  value: string;
  label: string;
  sort_order: number;
}

export const getConfigList = () =>
  requestClient.get<AppConfig[]>('/configs');

export const createConfig = (data: {
  key_name: string;
  label: string;
  sort_order?: number;
  value?: string;
}) => requestClient.post<AppConfig>('/configs', data);

export const updateConfig = (key: string, value: string) =>
  requestClient.put(`/configs/${key}`, { value });

import { requestClient } from '#/api/request';
import type { IpDetail } from './logs';

export interface Device {
  id: number;
  device_id: string;
  display_id: string;
  user_id: number | null;
  brand: string;
  model: string;
  os_version: string;
  app_version: string;
  last_ip: string;
  last_ip_detail?: IpDetail;
  last_seen_at: string | null;
  created_at: string;
}

export interface DeviceListResult {
  list: Device[];
  total: number;
  page: number;
  size: number;
}

export const getDeviceList = (params: { page?: number; size?: number; keyword?: string }) =>
  requestClient.get<DeviceListResult>('/devices', { params });

export const getDevice = (id: number) =>
  requestClient.get<Device>(`/devices/${id}`);

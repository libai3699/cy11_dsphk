import { requestClient } from '#/api/request';
import type { IpDetail } from './logs';

export interface User {
  id: number;
  username: string;
  phone: string;
  status: number;
  device_id: string;
  display_id: string;
  last_ip: string;
  last_ip_detail?: IpDetail;
  free_used_seconds: number;
  free_limit_seconds: number;
  current_line_id: number | null;
  current_plan_id: number | null;
  plan_expired_at: string | null;
  traffic_used_bytes: number;
  traffic_limit_bytes: number | null;
  last_login_at: string | null;
  last_active_at: string | null;
  is_online: boolean;
  online_devices: number;
  created_at: string;
}

export interface UserListResult {
  list: User[];
  total: number;
  page: number;
  size: number;
}

export const getUserList = (params: { page?: number; size?: number; keyword?: string }) =>
  requestClient.get<UserListResult>('/users', { params });

export const getUser = (id: number) =>
  requestClient.get<User>(`/users/${id}`);

export const createUser = (data: { username: string; password: string; phone?: string; free_limit_seconds?: number }) =>
  requestClient.post<User>('/users', data);

export const updateUser = (id: number, data: {
  status?: number;
  phone?: string;
  password?: string;
  free_used_seconds?: number;
  free_limit_seconds?: number;
  remaining_seconds?: number;
  traffic_remaining_bytes?: number;
}) =>
  requestClient.put(`/users/${id}`, data);

export const deleteUser = (id: number) =>
  requestClient.delete(`/users/${id}`);

export const addUserDuration = (id: number, seconds: number, trafficBytes = 0) =>
  requestClient.post(`/users/${id}/add-duration`, { 
    seconds: Math.round(seconds),
    traffic_bytes: trafficBytes,
  });


import { requestClient } from '#/api/request';

export interface UserLog {
  id: number;
  user_id: number;
  device_id: string;
  ip: string;
  ip_detail?: IpDetail;
  app_version: string;
  status: number;
  status_text: string;
  fail_reason: string;
  created_at: string;
}

export interface AdminLog {
  id: number;
  username: string;
  ip: string;
  ip_detail?: IpDetail;
  user_agent: string;
  status: number;
  status_text: string;
  fail_reason: string;
  created_at: string;
}

export interface IpDetail {
  ip: string;
  is_private: boolean;
  location: string;
  type: string;
  country?: string;
  region?: string;
  city?: string;
  isp?: string;
}

export interface LogListResult<T> {
  list: T[];
  total: number;
  page: number;
  size: number;
}

export const getUserLogs = (params: { page?: number; size?: number }) =>
  requestClient.get<LogListResult<UserLog>>('/logs/user', { params });

export const getAdminLogs = (params: { page?: number; size?: number }) =>
  requestClient.get<LogListResult<AdminLog>>('/logs/admin', { params });

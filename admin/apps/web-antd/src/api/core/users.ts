import { requestClient } from '#/api/request';

export namespace UsersApi {
  /** 用户信息 */
  export interface User {
    id: number;
    username: string;
    phone?: string;
    status: number;
    device_id?: string;
    free_used_seconds: number;
    free_limit_seconds: number;
    current_plan_id?: number;
    plan_expired_at?: string;
    traffic_used_bytes: number;
    traffic_limit_bytes?: number;
    last_login_at?: string;
    created_at: string;
    updated_at: string;
  }

  /** 用户列表参数 */
  export interface UserListParams {
    page?: number;
    pageSize?: number;
    username?: string;
  }

  /** 用户列表返回 */
  export interface UserListResult {
    items: User[];
    total: number;
  }
}

/**
 * 获取用户列表
 */
export async function getUsersListApi(params: UsersApi.UserListParams) {
  return requestClient.get<UsersApi.UserListResult>('/users', { params });
}

/**
 * 获取用户详情
 */
export async function getUserDetailApi(id: number) {
  return requestClient.get<UsersApi.User>(`/users/${id}`);
}

/**
 * 更新用户
 */
export async function updateUserApi(id: number, data: Partial<UsersApi.User>) {
  return requestClient.put<UsersApi.User>(`/users/${id}`, data);
}

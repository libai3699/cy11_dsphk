import { requestClient } from '#/api/request';

export interface MerchantAccountAuth {
  id: number;
  merchantId: number;
  merchantName: string;
  platform: string;
  authMethod: string;
  accountName: string;
  accountUid: string;
  authStatus: string;
  lastLoginAt: string;
  remark: string;
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface MerchantAccountAuthPayload {
  merchantId: number;
  platform?: string;
  authMethod?: string;
  accountName?: string;
  accountUid?: string;
  authStatus?: string;
  lastLoginAt?: string;
  remark?: string;
}

export interface MerchantAccountAuthListResult {
  list: MerchantAccountAuth[];
  total: number;
  page: number;
  size: number;
}

export function getMerchantAccountAuthList(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
  status?: string;
}) {
  return requestClient.get<MerchantAccountAuthListResult>('/account-auths', {
    params,
  });
}

export function createMerchantAccountAuth(data: MerchantAccountAuthPayload) {
  return requestClient.post<MerchantAccountAuth>('/account-auths', data);
}

export function updateMerchantAccountAuth(
  id: number,
  data: MerchantAccountAuthPayload,
) {
  return requestClient.put<MerchantAccountAuth>(`/account-auths/${id}`, data);
}

export function deleteMerchantAccountAuth(id: number) {
  return requestClient.delete<{ id: number }>(`/account-auths/${id}`);
}

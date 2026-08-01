import { requestClient } from '#/api/request';

export interface Merchant {
  id: number;
  name: string;
  industry: string;
  city: string;
  address: string;
  contactName: string;
  contactPhone: string;
  douyinAccount: string;
  douyinLaikeAccount: string;
  cooperationType: string;
  commissionRate: number;
  stage: string;
  status: number;
  remark: string;
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface MerchantPayload {
  name: string;
  industry?: string;
  city?: string;
  address?: string;
  contactName?: string;
  contactPhone?: string;
  douyinAccount?: string;
  douyinLaikeAccount?: string;
  cooperationType?: string;
  commissionRate?: number;
  stage?: string;
  status?: number;
  remark?: string;
}

export interface MerchantListResult {
  list: Merchant[];
  total: number;
  page: number;
  size: number;
}

export function getMerchantList(params: {
  keyword?: string;
  page?: number;
  size?: number;
}) {
  return requestClient.get<MerchantListResult>('/merchants', { params });
}

export function createMerchant(data: MerchantPayload) {
  return requestClient.post<Merchant>('/merchants', data);
}

export function updateMerchant(id: number, data: MerchantPayload) {
  return requestClient.put<Merchant>(`/merchants/${id}`, data);
}

export function getMerchant(id: number) {
  return requestClient.get<Merchant>(`/merchants/${id}`);
}

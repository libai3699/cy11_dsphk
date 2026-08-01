import { requestClient } from '#/api/request';

export interface MerchantPackage {
  id: number;
  merchantId: number;
  merchantName: string;
  name: string;
  originalPrice: number;
  sellingPrice: number;
  costPrice: number;
  commissionRate: number;
  trafficLabel: string;
  profitGuard: string;
  grossProfit: number;
  marginRate: number;
  estimatedCommission: number;
  netAfterCommission: number;
  status: number;
  remark: string;
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface MerchantPackagePayload {
  merchantId: number;
  name: string;
  originalPrice?: number;
  sellingPrice: number;
  costPrice?: number;
  commissionRate?: number;
  trafficLabel?: string;
  profitGuard?: string;
  status?: number;
  remark?: string;
}

export interface MerchantPackageListResult {
  list: MerchantPackage[];
  total: number;
  page: number;
  size: number;
}

export function getMerchantPackageList(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
}) {
  return requestClient.get<MerchantPackageListResult>('/packages', { params });
}

export function createMerchantPackage(data: MerchantPackagePayload) {
  return requestClient.post<MerchantPackage>('/packages', data);
}

export function updateMerchantPackage(id: number, data: MerchantPackagePayload) {
  return requestClient.put<MerchantPackage>(`/packages/${id}`, data);
}

export function deleteMerchantPackage(id: number) {
  return requestClient.delete<{ id: number }>(`/packages/${id}`);
}

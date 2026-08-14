import { requestClient } from '#/api/request';

export interface SettlementOrder {
  id: number;
  merchantId: number;
  merchantName: string;
  scheduleId: number;
  videoTitle: string;
  sourceVideo: string;
  orderWindow: string;
  periodStart?: string;
  periodEnd?: string;
  redeemedAmount: number;
  commissionRate: number;
  commission: number;
  status: string;
  remark: string;
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface SettlementOrderPayload {
  merchantId: number;
  scheduleId?: number;
  videoTitle?: string;
  sourceVideo?: string;
  orderWindow?: string;
  periodStart?: string;
  periodEnd?: string;
  redeemedAmount?: number;
  commissionRate?: number;
  status?: string;
  remark?: string;
}

export interface SettlementOrderListResult {
  list: SettlementOrder[];
  total: number;
  page: number;
  size: number;
}

export interface SettlementOrderGenerateResult {
  created: number;
  list: SettlementOrder[];
}

export function getSettlementOrderList(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
  status?: string;
}) {
  return requestClient.get<SettlementOrderListResult>('/settlement-orders', {
    params,
  });
}

export function createSettlementOrder(data: SettlementOrderPayload) {
  return requestClient.post<SettlementOrder>('/settlement-orders', data);
}

export function updateSettlementOrder(id: number, data: SettlementOrderPayload) {
  return requestClient.put<SettlementOrder>(`/settlement-orders/${id}`, data);
}

export function updateSettlementOrderStatus(
  id: number,
  data: { remark?: string; status: string },
) {
  return requestClient.put<SettlementOrder>(
    `/settlement-orders/${id}/status`,
    data,
  );
}

export function deleteSettlementOrder(id: number) {
  return requestClient.delete<{ id: number }>(`/settlement-orders/${id}`);
}

export function generateSettlementOrders(data: { merchantId?: number }) {
  return requestClient.post<SettlementOrderGenerateResult>(
    '/settlement-orders/generate',
    data,
  );
}

import { requestClient } from '#/api/request';

export interface PaymentConfig {
  id: number;
  type: string;
  label: string;
  address: string;
  qr_code: string;
  is_active: number;
  sort_order: number;
  remark: string;
  created_at: string;
  updated_at: string;
}

export const getPaymentConfigList = () =>
  requestClient.get<PaymentConfig[]>('/payment-configs');

export const createPaymentConfig = (data: { type: string; label: string; address?: string; qr_code?: string; is_active?: number; sort_order?: number; remark?: string }) =>
  requestClient.post<PaymentConfig>('/payment-configs', data);

export const updatePaymentConfig = (id: number, data: Partial<PaymentConfig>) =>
  requestClient.put(`/payment-configs/${id}`, data);

export const deletePaymentConfig = (id: number) =>
  requestClient.delete(`/payment-configs/${id}`);

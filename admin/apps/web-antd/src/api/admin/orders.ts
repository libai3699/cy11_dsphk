import { requestClient } from '#/api/request';

export interface Order {
  id: number;
  user_id: number;
  plan_id: number;
  plan_name: string;
  plan_price: number;
  traffic_gb: number | null;
  duration_days: number;
  started_at: string;
  expired_at: string;
  pay_method: string;
  remark: string;
  created_at: string;
}

export interface OrderListResult {
  list: Order[];
  total: number;
  page: number;
  size: number;
}

export const getOrderList = (params: { page?: number; size?: number }) =>
  requestClient.get<OrderListResult>('/orders', { params });

export const createOrder = (data: { billing_cycle?: 'half_year' | 'month' | 'quarter' | 'year'; user_id: number; plan_id: number; remark?: string }) =>
  requestClient.post<Order>('/orders', data);

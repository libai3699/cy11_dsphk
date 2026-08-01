import { requestClient } from '#/api/request';

export interface Plan {
  discount_half_year?: number | null;
  discount_quarter?: number | null;
  discount_year?: number | null;
  id: number;
  name: string;
  price: number;
  traffic_gb: number | null;
  duration_days: number;
  sort_order: number;
  is_active: number;
  max_devices: number;
  created_at: string;
}

export const getPlanList = () =>
  requestClient.get<Plan[]>('/plans');

export const createPlan = (data: Omit<Plan, 'id' | 'created_at'>) =>
  requestClient.post<Plan>('/plans', data);

export const updatePlan = (id: number, data: Partial<Plan>) =>
  requestClient.put(`/plans/${id}`, data);

export const deletePlan = (id: number) =>
  requestClient.delete(`/plans/${id}`);

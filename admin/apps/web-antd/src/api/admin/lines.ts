import { requestClient } from '#/api/request';

export interface VpnLine {
  id: number;
  name: string;
  region: string;
  protocol: string;
  address: string;
  raw_uri: string;
  sort_order: number;
  is_default: 0 | 1;
  is_active: 0 | 1;
  description: string;
}

export const getLineList = () => requestClient.get<VpnLine[]>('/lines');

export const createLine = (data: Partial<VpnLine>) =>
  requestClient.post<VpnLine>('/lines', data);

export const updateLine = (id: number, data: Partial<VpnLine>) =>
  requestClient.put(`/lines/${id}`, data);

export const deleteLine = (id: number) => requestClient.delete(`/lines/${id}`);

export const assignUserLine = (data: {
  line_id: number;
  notice?: string;
  user_id: number;
}) => requestClient.post('/users/assign-line', data);

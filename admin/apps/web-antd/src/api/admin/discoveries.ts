import { requestClient } from '#/api/request';

export interface DiscoveryItem {
  id: number;
  name: string;
  icon_url: string;
  h5_url: string;
  is_active: number;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export type DiscoveryItemPayload = Pick<
  DiscoveryItem,
  'h5_url' | 'icon_url' | 'is_active' | 'name' | 'sort_order'
>;

export const getDiscoveryList = () =>
  requestClient.get<DiscoveryItem[]>('/discoveries');

export const createDiscovery = (data: DiscoveryItemPayload) =>
  requestClient.post<DiscoveryItem>('/discoveries', data);

export const updateDiscovery = (id: number, data: DiscoveryItemPayload) =>
  requestClient.put(`/discoveries/${id}`, data);

export const deleteDiscovery = (id: number) =>
  requestClient.delete(`/discoveries/${id}`);

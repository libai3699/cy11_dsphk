import { requestClient } from '#/api/request';

export interface Notice {
  id: number;
  target_user_id?: number | null;
  content: string;
  type: number;
  is_active: number;
  sort_order: number;
  created_at: string;
}

export const getNoticeList = () =>
  requestClient.get<Notice[]>('/notices');

export const createNotice = (data: { content: string; target_user_id?: number | null; type?: number; is_active?: number; sort_order?: number }) =>
  requestClient.post<Notice>('/notices', data);

export const updateNotice = (id: number, data: Partial<Notice>) =>
  requestClient.put(`/notices/${id}`, data);

export const deleteNotice = (id: number) =>
  requestClient.delete(`/notices/${id}`);

import { requestClient } from '#/api/request';

export interface Stats {
  today_new: number;
  today_income_cents: number;
  total_members: number;
  online_members: number;
  login_within_3_days: number;
}

export const getStats = () =>
  requestClient.get<Stats>('/stats');

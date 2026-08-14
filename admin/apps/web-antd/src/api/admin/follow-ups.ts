import { requestClient } from '#/api/request';

export interface MerchantFollowUpLog {
  id: number;
  merchantId: number;
  merchantName: string;
  stage: string;
  latestTalk: string;
  objection: string;
  nextStep: string;
  owner: string;
  followTime?: string;
  nextFollowTime?: string;
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface MerchantFollowUpPayload {
  merchantId: number;
  stage?: string;
  latestTalk: string;
  objection?: string;
  nextStep?: string;
  owner?: string;
  followTime?: string;
  nextFollowTime?: string;
}

export interface MerchantFollowUpListResult {
  list: MerchantFollowUpLog[];
  total: number;
  page: number;
  size: number;
}

export interface MerchantFollowUpSuggestion {
  talkScript: string;
  actions: string[];
}

export function getFollowUpLogList(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
  stage?: string;
}) {
  return requestClient.get<MerchantFollowUpListResult>('/follow-up-logs', {
    params,
  });
}

export function createFollowUpLog(data: MerchantFollowUpPayload) {
  return requestClient.post<MerchantFollowUpLog>('/follow-up-logs', data);
}

export function updateFollowUpLog(id: number, data: MerchantFollowUpPayload) {
  return requestClient.put<MerchantFollowUpLog>(`/follow-up-logs/${id}`, data);
}

export function deleteFollowUpLog(id: number) {
  return requestClient.delete<{ id: number }>(`/follow-up-logs/${id}`);
}

export function generateFollowUpSuggestion(id: number) {
  return requestClient.post<MerchantFollowUpSuggestion>(
    `/follow-up-logs/${id}/suggestion`,
  );
}

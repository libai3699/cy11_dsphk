import { requestClient } from '#/api/request';

export interface Merchant {
  id: number;
  name: string;
  industry: string;
  city: string;
  address: string;
  contactName: string;
  contactPhone: string;
  douyinAccount: string;
  douyinLaikeAccount: string;
  cooperationType: string;
  commissionRate: number;
  stage: string;
  status: number;
  remark: string;
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface MerchantPayload {
  name: string;
  industry?: string;
  city?: string;
  address?: string;
  contactName?: string;
  contactPhone?: string;
  douyinAccount?: string;
  douyinLaikeAccount?: string;
  cooperationType?: string;
  commissionRate?: number;
  stage?: string;
  status?: number;
  remark?: string;
}

export interface MerchantListResult {
  list: Merchant[];
  total: number;
  page: number;
  size: number;
}

export interface MerchantWorkspaceMetrics {
  acceptedTopicCount: number;
  accountAuthCount: number;
  activeAccountAuthCount: number;
  analyzedBenchmarkCount: number;
  benchmarkCount: number;
  completedDiagnosisCount: number;
  confirmedScriptCount: number;
  confirmedStoryboardCount: number;
  diagnosisCount: number;
  enabledPackageCount: number;
  packageCount: number;
  publishedScheduleCount: number;
  readyShootingTaskCount: number;
  reviewCount: number;
  scheduleCount: number;
  scriptCount: number;
  shootingTaskCount: number;
  storyboardCount: number;
  topicCount: number;
}

export interface MerchantRequirementStatus {
  done: boolean;
  key: string;
  missing: string[];
  title: string;
}

export interface MerchantWorkspace {
  completeness: number;
  merchant: Merchant;
  metrics: MerchantWorkspaceMetrics;
  requirements: MerchantRequirementStatus[];
}

export function getMerchantList(params: {
  keyword?: string;
  page?: number;
  size?: number;
}) {
  return requestClient.get<MerchantListResult>('/merchants', { params });
}

export function createMerchant(data: MerchantPayload) {
  return requestClient.post<Merchant>('/merchants', data);
}

export function updateMerchant(id: number, data: MerchantPayload) {
  return requestClient.put<Merchant>(`/merchants/${id}`, data);
}

export function deleteMerchant(id: number) {
  return requestClient.delete<{
    deleted: Record<string, number>;
    id: number;
    name: string;
  }>(`/merchants/${id}`);
}

export function getMerchant(id: number) {
  return requestClient.get<Merchant>(`/merchants/${id}`);
}

export function getMerchantWorkspace(id: number) {
  return requestClient.get<MerchantWorkspace>(`/merchants/${id}/workspace`);
}

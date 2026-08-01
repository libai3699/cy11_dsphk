import { requestClient } from '#/api/request';

export interface AccountDiagnosisResult {
  accountScore?: number;
  contentScore?: number;
  convertScore?: number;
  problems?: string[];
  nextActions?: string[];
  summary?: string;
  suggestions?: string[];
  artifacts?: Record<string, unknown>;
  [key: string]: unknown;
}

export interface AccountDiagnosisTask {
  id: number;
  merchantId: number;
  merchantName: string;
  accountAuthId: number;
  accountName: string;
  status: string;
  inputSnapshot: string;
  resultJson: string;
  errorMessage: string;
  input?: unknown;
  result?: AccountDiagnosisResult;
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface AccountDiagnosisPayload {
  merchantId: number;
  accountAuthId?: number;
  followerCount?: number;
  avgPlayCount?: number;
  bestVideoTitle?: string;
  bestVideoPlay?: number;
  recentVideoCount?: number;
  ownerAppearance?: string;
  currentProblems?: string;
  targetPackage?: string;
  operatorGoal?: string;
  recentVideos?: string[];
  remark?: string;
}

export interface AccountDiagnosisListResult {
  list: AccountDiagnosisTask[];
  total: number;
  page: number;
  size: number;
}

export function getAccountDiagnosisList(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
  status?: string;
}) {
  return requestClient.get<AccountDiagnosisListResult>('/account-diagnoses', {
    params,
  });
}

export function createAccountDiagnosis(data: AccountDiagnosisPayload) {
  return requestClient.post<AccountDiagnosisTask>('/account-diagnoses', data);
}

export function getAccountDiagnosis(id: number) {
  return requestClient.get<AccountDiagnosisTask>(`/account-diagnoses/${id}`);
}

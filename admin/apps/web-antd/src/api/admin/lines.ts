import { requestClient } from '#/api/request';

export interface BenchmarkAccount {
  id: number;
  merchantId: number;
  merchantName: string;
  accountName: string;
  platform: string;
  city: string;
  industry: string;
  accountUrl: string;
  followerCount: number;
  bestPlayCount: number;
  latestHitTitle: string;
  takeaway: string;
  risk: string;
  status: string;
  remark: string;
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface BenchmarkAccountPayload {
  merchantId: number;
  accountName: string;
  platform?: string;
  city?: string;
  industry?: string;
  accountUrl?: string;
  followerCount?: number;
  bestPlayCount?: number;
  latestHitTitle?: string;
  takeaway?: string;
  risk?: string;
  status?: string;
  remark?: string;
}

export interface BenchmarkAccountListResult {
  list: BenchmarkAccount[];
  total: number;
  page: number;
  size: number;
}

export interface BenchmarkAnalysisResult {
  summary?: string;
  suggestions?: string[];
  patterns?: string[];
  risks?: string[];
  [key: string]: unknown;
}

export interface BenchmarkAnalysisTask {
  id: number;
  merchantId: number;
  merchantName: string;
  benchmarkAccountId: number;
  benchmarkName: string;
  status: string;
  inputSnapshot: string;
  resultJson: string;
  errorMessage: string;
  input?: unknown;
  result?: BenchmarkAnalysisResult;
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface BenchmarkAnalysisListResult {
  list: BenchmarkAnalysisTask[];
  total: number;
  page: number;
  size: number;
}

export function getBenchmarkAccountList(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
  status?: string;
}) {
  return requestClient.get<BenchmarkAccountListResult>('/benchmarks', { params });
}

export function createBenchmarkAccount(data: BenchmarkAccountPayload) {
  return requestClient.post<BenchmarkAccount>('/benchmarks', data);
}

export function updateBenchmarkAccount(id: number, data: BenchmarkAccountPayload) {
  return requestClient.put<BenchmarkAccount>(`/benchmarks/${id}`, data);
}

export function deleteBenchmarkAccount(id: number) {
  return requestClient.delete<{ id: number }>(`/benchmarks/${id}`);
}

export function analyzeBenchmarkAccount(id: number) {
  return requestClient.post<BenchmarkAnalysisTask>(`/benchmarks/${id}/analyze`);
}

export function analyzeBenchmarkAccounts(data: {
  benchmarkIds?: number[];
  merchantId: number;
}) {
  return requestClient.post<BenchmarkAnalysisTask>(
    '/benchmarks/analyze-batch',
    data,
  );
}

export function getBenchmarkAnalysisList(params: {
  benchmarkAccountId?: number;
  merchantId?: number;
  page?: number;
  size?: number;
}) {
  return requestClient.get<BenchmarkAnalysisListResult>('/benchmark-analyses', {
    params,
  });
}

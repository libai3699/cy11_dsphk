import { requestClient } from '#/api/request';

export interface ContentTopic {
  id: number;
  merchantId: number;
  merchantName: string;
  benchmarkId: number;
  benchmarkName: string;
  title: string;
  hook: string;
  angle: string;
  target: string;
  riskLevel: string;
  recommendReason: string;
  tagsJson: string;
  tags: string[];
  publishWindow: string;
  status: string;
  sourceTaskId: number;
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface ContentTopicListResult {
  list: ContentTopic[];
  total: number;
  page: number;
  size: number;
}

export interface GenerateTopicsPayload {
  merchantId: number;
  benchmarkId?: number;
  benchmarkIds?: number[];
  benchmarkName?: string;
  cityHotspots?: string[];
  industryHotspots?: string[];
  nationalHotspots?: string[];
  extraRequirement?: string;
}

export interface HotspotTopicTask {
  id: number;
  merchantId: number;
  merchantName: string;
  benchmarkId: number;
  benchmarkName: string;
  status: string;
  inputSnapshot: string;
  resultJson: string;
  errorMessage: string;
  input?: unknown;
  result?: unknown;
  topics?: ContentTopic[];
  createdBy: number;
  updatedBy: number;
  createdAt: string;
  updatedAt: string;
}

export function getContentTopicList(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
  status?: string;
}) {
  return requestClient.get<ContentTopicListResult>('/topics', { params });
}

export function generateContentTopics(data: GenerateTopicsPayload) {
  return requestClient.post<HotspotTopicTask>('/topics/generate', data, {
    timeout: 120_000,
  });
}

export function updateContentTopicStatus(id: number, status: string) {
  return requestClient.put<ContentTopic>(`/topics/${id}/status`, { status });
}

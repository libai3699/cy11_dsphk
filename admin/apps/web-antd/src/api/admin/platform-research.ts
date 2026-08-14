import { requestClient } from '#/api/request';

export interface PlatformSearchResult {
  keyword: string;
  platform: string;
  query?: string;
  score?: number;
  snippet: string;
  source: string;
  title: string;
  url: string;
}

export interface PlatformCaseStudy {
  platform: string;
  reason: string;
  risks?: string[];
  takeaway: string;
  title: string;
  url?: string;
}

export interface PlatformResearchTask {
  id: number;
  merchantId: number;
  merchantName: string;
  industry: string;
  city: string;
  sources?: string[];
  keywords?: string[];
  searchResults?: PlatformSearchResult[];
  goodCases?: PlatformCaseStudy[];
  badCases?: PlatformCaseStudy[];
  insights?: string[];
  suggestions?: string[];
  summary: string;
  status: string;
  errorMessage: string;
  createdAt: string;
  updatedAt: string;
}

export interface PlatformResearchListResult {
  list: PlatformResearchTask[];
  page: number;
  size: number;
  total: number;
}

export interface GeneratePlatformResearchPayload {
  extraRequirement?: string;
  keywords?: string[];
  limit?: number;
  merchantId: number;
  sources?: string[];
}

export function getPlatformResearchList(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
}) {
  return requestClient.get<PlatformResearchListResult>('/platform-researches', {
    params,
  });
}

export function generatePlatformResearch(data: GeneratePlatformResearchPayload) {
  return requestClient.post<PlatformResearchTask>(
    '/platform-researches/generate',
    data,
    {
      timeout: 180_000,
    },
  );
}

export function getPlatformResearch(id: number) {
  return requestClient.get<PlatformResearchTask>(`/platform-researches/${id}`);
}

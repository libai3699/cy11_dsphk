import { requestClient } from '#/api/request';

export interface KnowledgeListResult<T> {
  list: T[];
  total: number;
  page: number;
  size: number;
}

/* ===================== 痛点 pain-points ===================== */
export interface PainPoint {
  id: number;
  merchantId: number;
  merchantName: string;
  source: string;
  category: string;
  content: string;
  userQuote: string;
  emotion: string;
  product: string;
  demandLevel: string;
  suggestedTopic: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface PainPointPayload {
  merchantId?: number;
  source?: string;
  category?: string;
  content: string;
  userQuote?: string;
  emotion?: string;
  product?: string;
  demandLevel?: string;
  suggestedTopic?: string;
  status?: string;
}

export function getPainPoints(params: {
  keyword?: string;
  merchantId?: number;
  platform?: string;
  status?: string;
  page?: number;
  size?: number;
}) {
  return requestClient.get<KnowledgeListResult<PainPoint>>('/knowledge/pain-points', {
    params,
  });
}

export function getPainPoint(id: number) {
  return requestClient.get<PainPoint>(`/knowledge/pain-points/${id}`);
}

export function createPainPoint(data: Partial<PainPointPayload>) {
  return requestClient.post<PainPoint>('/knowledge/pain-points', data);
}

export function updatePainPoint(id: number, data: Partial<PainPointPayload>) {
  return requestClient.put<PainPoint>(`/knowledge/pain-points/${id}`, data);
}

export function deletePainPoint(id: number) {
  return requestClient.delete(`/knowledge/pain-points/${id}`);
}

/* ===================== 案例 case-studies ===================== */
export interface CaseStudy {
  id: number;
  merchantId: number;
  merchantName: string;
  title: string;
  platform: string;
  accountName: string;
  industry: string;
  form: string;
  hookType: string;
  structure: string;
  conversionAction: string;
  reusablePoint: string;
  dataJson: string;
  lesson: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface CaseStudyPayload {
  merchantId?: number;
  title: string;
  platform?: string;
  accountName?: string;
  industry?: string;
  form?: string;
  hookType?: string;
  structure?: string;
  conversionAction?: string;
  reusablePoint?: string;
  dataJson?: string;
  lesson?: string;
  status?: string;
}

export function getCaseStudies(params: {
  keyword?: string;
  merchantId?: number;
  platform?: string;
  status?: string;
  page?: number;
  size?: number;
}) {
  return requestClient.get<KnowledgeListResult<CaseStudy>>('/knowledge/case-studies', {
    params,
  });
}

export function getCaseStudy(id: number) {
  return requestClient.get<CaseStudy>(`/knowledge/case-studies/${id}`);
}

export function createCaseStudy(data: Partial<CaseStudyPayload>) {
  return requestClient.post<CaseStudy>('/knowledge/case-studies', data);
}

export function updateCaseStudy(id: number, data: Partial<CaseStudyPayload>) {
  return requestClient.put<CaseStudy>(`/knowledge/case-studies/${id}`, data);
}

export function deleteCaseStudy(id: number) {
  return requestClient.delete(`/knowledge/case-studies/${id}`);
}

/* ===================== 账号画像 account-profiles ===================== */
export interface AccountProfile {
  id: number;
  merchantId: number;
  merchantName: string;
  accountName: string;
  platform: string;
  industry: string;
  targetAudience: string;
  product: string;
  positioning: string;
  persona: string;
  toneStyle: string;
  monetization: string;
  updateFrequency: string;
  strengths: string;
  weaknesses: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface AccountProfilePayload {
  merchantId?: number;
  accountName: string;
  platform?: string;
  industry?: string;
  targetAudience?: string;
  product?: string;
  positioning?: string;
  persona?: string;
  toneStyle?: string;
  monetization?: string;
  updateFrequency?: string;
  strengths?: string;
  weaknesses?: string;
  status?: string;
}

export function getAccountProfiles(params: {
  keyword?: string;
  merchantId?: number;
  platform?: string;
  status?: string;
  page?: number;
  size?: number;
}) {
  return requestClient.get<KnowledgeListResult<AccountProfile>>(
    '/knowledge/account-profiles',
    { params },
  );
}

export function getAccountProfile(id: number) {
  return requestClient.get<AccountProfile>(`/knowledge/account-profiles/${id}`);
}

export function createAccountProfile(data: Partial<AccountProfilePayload>) {
  return requestClient.post<AccountProfile>('/knowledge/account-profiles', data);
}

export function updateAccountProfile(id: number, data: Partial<AccountProfilePayload>) {
  return requestClient.put<AccountProfile>(`/knowledge/account-profiles/${id}`, data);
}

export function deleteAccountProfile(id: number) {
  return requestClient.delete(`/knowledge/account-profiles/${id}`);
}

/* ===================== 平台规则 platform-rules ===================== */
export interface PlatformRule {
  id: number;
  merchantId: number;
  merchantName: string;
  platform: string;
  category: string;
  title: string;
  content: string;
  riskLevel: string;
  source: string;
  effectiveDate: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface PlatformRulePayload {
  merchantId?: number;
  platform: string;
  category?: string;
  title: string;
  content?: string;
  riskLevel?: string;
  source?: string;
  effectiveDate?: string;
  status?: string;
}

export function getPlatformRules(params: {
  keyword?: string;
  platform?: string;
  status?: string;
  page?: number;
  size?: number;
}) {
  return requestClient.get<KnowledgeListResult<PlatformRule>>('/knowledge/platform-rules', {
    params,
  });
}

export function getPlatformRule(id: number) {
  return requestClient.get<PlatformRule>(`/knowledge/platform-rules/${id}`);
}

export function createPlatformRule(data: Partial<PlatformRulePayload>) {
  return requestClient.post<PlatformRule>('/knowledge/platform-rules', data);
}

export function updatePlatformRule(id: number, data: Partial<PlatformRulePayload>) {
  return requestClient.put<PlatformRule>(`/knowledge/platform-rules/${id}`, data);
}

export function deletePlatformRule(id: number) {
  return requestClient.delete(`/knowledge/platform-rules/${id}`);
}

/* ===================== 内容模板 content-templates ===================== */
export interface ContentTemplate {
  id: number;
  merchantId: number;
  merchantName: string;
  name: string;
  type: string;
  category: string;
  structureJson: string;
  content: string;
  usageNote: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface ContentTemplatePayload {
  merchantId?: number;
  name: string;
  type?: string;
  category?: string;
  structureJson?: string;
  content?: string;
  usageNote?: string;
  status?: string;
}

export function getContentTemplates(params: {
  keyword?: string;
  merchantId?: number;
  platform?: string;
  status?: string;
  page?: number;
  size?: number;
}) {
  return requestClient.get<KnowledgeListResult<ContentTemplate>>(
    '/knowledge/content-templates',
    { params },
  );
}

export function getContentTemplate(id: number) {
  return requestClient.get<ContentTemplate>(`/knowledge/content-templates/${id}`);
}

export function createContentTemplate(data: Partial<ContentTemplatePayload>) {
  return requestClient.post<ContentTemplate>('/knowledge/content-templates', data);
}

export function updateContentTemplate(id: number, data: Partial<ContentTemplatePayload>) {
  return requestClient.put<ContentTemplate>(`/knowledge/content-templates/${id}`, data);
}

export function deleteContentTemplate(id: number) {
  return requestClient.delete(`/knowledge/content-templates/${id}`);
}

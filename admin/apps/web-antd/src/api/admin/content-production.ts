import { requestClient } from '#/api/request';

export interface AgentConfig {
  agent: string;
  provider: string;
  baseUrl: string;
  model: string;
  timeoutSeconds: number;
  keyConfigured: boolean;
  enabled: boolean;
}

export interface ContentScript {
  id: number;
  merchantId: number;
  merchantName: string;
  topicId: number;
  topicTitle: string;
  title: string;
  opening: string;
  body: string;
  ending: string;
  cta: string;
  fullScript: string;
  shootingNotesJson: string;
  shootingNotes: string[];
  status: string;
  inputSnapshot: string;
  resultJson: string;
  errorMessage: string;
  createdAt: string;
  updatedAt: string;
}

export interface ContentScriptListResult {
  list: ContentScript[];
  total: number;
  page: number;
  size: number;
}

export interface StoryboardShot {
  index: number;
  duration: string;
  location: string;
  camera: string;
  content: string;
  line: string;
  note: string;
}

export interface ContentStoryboard {
  id: number;
  merchantId: number;
  merchantName: string;
  topicId: number;
  topicTitle: string;
  scriptId: number;
  scriptTitle: string;
  shotsJson: string;
  shots: StoryboardShot[];
  status: string;
  inputSnapshot: string;
  resultJson: string;
  errorMessage: string;
  createdAt: string;
  updatedAt: string;
}

export interface ContentStoryboardListResult {
  list: ContentStoryboard[];
  total: number;
  page: number;
  size: number;
}

export interface ShootingTask {
  id: number;
  merchantId: number;
  merchantName: string;
  topicId: number;
  topicTitle: string;
  scriptId: number;
  scriptTitle: string;
  storyboardId: number;
  taskTitle: string;
  shotCount: number;
  shotsJson: string;
  shots: StoryboardShot[];
  assignee: string;
  shootTime?: string;
  deadline?: string;
  status: string;
  materialUrl: string;
  remark: string;
  createdAt: string;
  updatedAt: string;
}

export interface ShootingTaskListResult {
  list: ShootingTask[];
  total: number;
  page: number;
  size: number;
}

export interface PublishSchedule {
  id: number;
  merchantId: number;
  merchantName: string;
  topicId: number;
  topicTitle: string;
  scriptId: number;
  scriptTitle: string;
  storyboardId: number;
  videoTitle: string;
  publishTime?: string;
  owner: string;
  douyinAccount: string;
  materialStatus: string;
  status: string;
  remark: string;
  createdAt: string;
  updatedAt: string;
}

export interface PublishScheduleListResult {
  list: PublishSchedule[];
  total: number;
  page: number;
  size: number;
}

export interface ContentReviewTask {
  id: number;
  merchantId: number;
  merchantName: string;
  scheduleId: number;
  videoTitle: string;
  status: string;
  resultJson: string;
  errorMessage: string;
  periodStart?: string;
  periodEnd?: string;
  playCount: number;
  likeCount: number;
  commentCount: number;
  shareCount: number;
  dealCount: number;
  writeOffAmount: number;
  result?: {
    conclusion?: string;
    nextTopics?: string[];
    optimizes?: string[];
    summary?: string;
    suggestions?: string[];
  };
  createdAt: string;
  updatedAt: string;
}

export function getContentScripts(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
  status?: string;
  topicId?: number;
}) {
  return requestClient.get<ContentScriptListResult>('/scripts', { params });
}

export function generateContentScript(data: {
  extraRequirement?: string;
  topicId: number;
}) {
  return requestClient.post<ContentScript>('/scripts/generate', data, {
    timeout: 120_000,
  });
}

export function updateContentScriptStatus(id: number, status: string) {
  return requestClient.put<ContentScript>(`/scripts/${id}/status`, { status });
}

export function getContentStoryboards(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  scriptId?: number;
  size?: number;
  status?: string;
}) {
  return requestClient.get<ContentStoryboardListResult>('/storyboards', {
    params,
  });
}

export function generateContentStoryboard(data: {
  locations?: string[];
  scriptId: number;
}) {
  return requestClient.post<ContentStoryboard>('/storyboards/generate', data, {
    timeout: 120_000,
  });
}

export function updateContentStoryboardStatus(id: number, status: string) {
  return requestClient.put<ContentStoryboard>(`/storyboards/${id}/status`, {
    status,
  });
}

export function getShootingTasks(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
  status?: string;
  storyboardId?: number;
}) {
  return requestClient.get<ShootingTaskListResult>('/shooting-tasks', {
    params,
  });
}

export function createShootingTask(data: {
  assignee?: string;
  deadline?: string;
  materialUrl?: string;
  remark?: string;
  shootTime?: string;
  status?: string;
  storyboardId: number;
  taskTitle?: string;
}) {
  return requestClient.post<ShootingTask>('/shooting-tasks', data);
}

export function updateShootingTaskStatus(
  id: number,
  data: { materialUrl?: string; remark?: string; status?: string },
) {
  return requestClient.put<ShootingTask>(`/shooting-tasks/${id}/status`, data);
}

export function getPublishSchedules(params: {
  keyword?: string;
  merchantId?: number;
  page?: number;
  size?: number;
  status?: string;
  storyboardId?: number;
}) {
  return requestClient.get<PublishScheduleListResult>('/schedules', {
    params,
  });
}

export function createPublishSchedule(data: Partial<PublishSchedule>) {
  return requestClient.post<PublishSchedule>('/schedules', data);
}

export function updatePublishSchedule(
  id: number,
  data: Partial<PublishSchedule>,
) {
  return requestClient.put<PublishSchedule>(`/schedules/${id}`, data);
}

export function updatePublishScheduleStatus(
  id: number,
  data: { materialStatus?: string; remark?: string; status?: string },
) {
  return requestClient.put<PublishSchedule>(`/schedules/${id}/status`, data);
}

export function getContentReviews(params: {
  merchantId?: number;
  page?: number;
  size?: number;
}) {
  return requestClient.get<{
    list: ContentReviewTask[];
    total: number;
    page: number;
    size: number;
  }>('/reviews', { params });
}

export function generateContentReview(data: {
  commentCount?: number;
  dealCount?: number;
  likeCount?: number;
  periodEnd?: string;
  periodStart?: string;
  playCount?: number;
  scheduleId: number;
  shareCount?: number;
  writeOffAmount?: number;
}) {
  return requestClient.post<ContentReviewTask>('/reviews/generate', data, {
    timeout: 120_000,
  });
}

export function getAgentConfigs() {
  return requestClient.get<AgentConfig[]>('/agent-configs');
}

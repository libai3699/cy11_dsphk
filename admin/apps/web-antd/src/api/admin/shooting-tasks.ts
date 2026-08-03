import { requestClient } from '#/api/request';

export interface ShootingTaskShot {
  index: number;
  duration: string;
  location: string;
  camera: string;
  content: string;
  line: string;
  note: string;
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
  shots: ShootingTaskShot[];
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

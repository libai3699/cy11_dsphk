import { requestClient } from '#/api/request';

export interface UserInfo {
  avatar?: string;
  desc?: string;
  homePath?: string;
  realName: string;
  roles: string[];
  token?: string;
  userId: string;
  username: string;
}

export async function getUserInfoApi(): Promise<UserInfo> {
  return requestClient.get<UserInfo>('/user/info');
}

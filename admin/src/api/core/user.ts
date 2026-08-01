export interface UserInfo {
  id: number;
  username: string;
  realName: string;
  roles: string[];
}

/**
 * 获取用户信息 - 从登录时缓存的数据中读取
 */
export async function getUserInfoApi(): Promise<UserInfo> {
  return Promise.resolve({
    id: 1,
    username: 'admin',
    realName: '管理员',
    roles: ['super'],
  });
}

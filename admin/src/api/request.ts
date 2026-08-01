import axios, { type AxiosInstance } from 'axios';
import { useUserStore } from '../vben-shims/stores';

// 开发环境使用相对路径（走vite代理），生产环境使用配置的完整URL
const apiURL = import.meta.env.DEV 
  ? '/api/admin' 
  : import.meta.env.VITE_GLOB_API_URL;

console.log('API URL:', apiURL, 'ENV:', import.meta.env.MODE);

class RequestClient {
  private instance: AxiosInstance;

  constructor(baseURL: string) {
    this.instance = axios.create({
      baseURL,
      timeout: 10000,
    });

    // 请求拦截器
    this.instance.interceptors.request.use(
      (config) => {
        const userStore = useUserStore();
        if (userStore.token) {
          config.headers.Authorization = `Bearer ${userStore.token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // 响应拦截器
    this.instance.interceptors.response.use(
      (response) => {
        const { data } = response;
        // 后端返回格式：{ code: 0, data: {...} }
        if (data.code === 0) {
          return data.data;
        }
        // 如果code不是0，抛出错误
        const errorMsg = data.msg || data.message || '请求失败';
        console.error('API Error:', errorMsg);
        return Promise.reject(new Error(errorMsg));
      },
      (error) => {
        console.error('Request error:', error);
        const errorMsg = error.response?.data?.msg || error.response?.data?.message || error.message || '网络错误';
        return Promise.reject(new Error(errorMsg));
      }
    );
  }

  async get<T = any>(url: string, config?: any): Promise<T> {
    return this.instance.get(url, config);
  }

  async post<T = any>(url: string, data?: any, config?: any): Promise<T> {
    return this.instance.post(url, data, config);
  }

  async put<T = any>(url: string, data?: any, config?: any): Promise<T> {
    return this.instance.put(url, data, config);
  }

  async delete<T = any>(url: string, config?: any): Promise<T> {
    return this.instance.delete(url, config);
  }
}

export const requestClient = new RequestClient(apiURL);
export const baseRequestClient = requestClient;

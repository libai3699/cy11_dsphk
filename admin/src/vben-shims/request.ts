// @vben/request 替代实现

export interface RequestOptions {
  url: string;
  method?: string;
  data?: any;
  params?: any;
  headers?: Record<string, string>;
}

export class RequestClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  async request(options: RequestOptions) {
    const url = this.baseURL + options.url;
    const response = await fetch(url, {
      method: options.method || 'GET',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      body: options.data ? JSON.stringify(options.data) : undefined,
    });
    return response.json();
  }
}

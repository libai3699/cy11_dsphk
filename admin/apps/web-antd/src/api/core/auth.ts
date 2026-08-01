import { requestClient } from '#/api/request';

export namespace AuthApi {
  export interface LoginParams {
    captcha?: string;
    password?: string;
    username?: string;
  }

  export interface LoginResult {
    accessToken: string;
  }

  export interface RefreshTokenResult {
    data: string;
    status: number;
  }

  export interface CaptchaResult {
    account: string;
    issuer: string;
    otpauth: string;
    period: number;
    secret: string;
    type: string;
  }
}

export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>('/auth/login', data);
}

export async function captchaApi() {
  return requestClient.get<AuthApi.CaptchaResult>('/auth/captcha');
}

export async function refreshTokenApi() {
  return requestClient.post<AuthApi.RefreshTokenResult>('/auth/refresh');
}

export async function logoutApi() {
  return requestClient.post('/auth/logout');
}

export async function getAccessCodesApi() {
  return requestClient.get<string[]>('/auth/codes');
}

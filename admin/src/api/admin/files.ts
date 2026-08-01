import { requestClient } from '#/api/request';

export interface UploadedFile {
  url: string;
}

export const uploadPaymentImage = (data: FormData) =>
  requestClient.post<UploadedFile>('/files/payment-image', data, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });

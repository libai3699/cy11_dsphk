import { requestClient } from '#/api/request';

export interface Quote {
  id: number;
  content: string;
  author: string;
  is_active: number;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export const getQuoteList = () =>
  requestClient.get<Quote[]>('/quotes');

export const createQuote = (data: { content: string; author?: string; is_active?: number; sort_order?: number }) =>
  requestClient.post<Quote>('/quotes', data);

export const updateQuote = (id: number, data: Partial<Quote>) =>
  requestClient.put(`/quotes/${id}`, data);

export const deleteQuote = (id: number) =>
  requestClient.delete(`/quotes/${id}`);

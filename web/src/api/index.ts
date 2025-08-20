import request from '../utils/request';

// 登录接口
export const login = (data: { username: string; password: string }) => {
  return request.post('/login', data);
};

// 创建短链接
export const createShortLink = (data: { link: string; expire: number }) => {
  return request.post('/short-link/create', data);
};

// 获取短链接列表
export const getShortLinks = (params: {
  page?: number;
  pageSize?: number;
  shortCode?: string;
  originalUrl?: string;
  status?: 'active' | 'expired';
}) => {
  return request.get('/short-link/list', { params });
};

// 获取历史短链接列表
export const getHistoryLinks = (params: {
  month?: string;
  page?: number;
  pageSize?: number;
  shortCode?: string;
  originalUrl?: string;
}) => {
  return request.get('/short-link/history', { params });
};

// 删除短链接
export const deleteShortLink = (id: number) => {
  return request.delete(`/short-link/${id}`);
};

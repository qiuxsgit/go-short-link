import axios, { AxiosResponse, AxiosError, InternalAxiosRequestConfig } from 'axios';
import { message } from 'antd';
import { getToken, removeToken } from './auth';

// 创建axios实例
const request = axios.create({
  baseURL: '/api', // 使用相对路径，结合package.json中的proxy配置
  timeout: 10000,
});

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 添加token到请求头
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    // 打印响应数据，用于调试
    console.log('API响应数据:', response.data);
    
    // 检查响应中是否包含错误信息
    if (response.data && response.data.error) {
      message.error(response.data.error);
      return Promise.reject(new Error(response.data.error));
    }
    
    return response.data;
  },
  (error: AxiosError) => {
    console.error('API请求错误:', error);
    
    if (error.response) {
      const { status } = error.response;
      console.error('错误状态码:', status);
      console.error('错误响应数据:', error.response.data);
      
      // 处理401未授权错误
      if (status === 401) {
        message.error('登录已过期，请重新登录');
        removeToken();
        // 重定向到登录页
        window.location.href = '/login';
      } else {
        // 处理其他错误
        const errorMessage = (error.response.data as any)?.error || '请求失败';
        message.error(errorMessage);
      }
    } else if (error.request) {
      console.error('请求未收到响应:', error.request);
      message.error('服务器未响应，请检查服务器状态');
    } else {
      console.error('请求配置错误:', error.message);
      message.error('网络错误，请检查您的网络连接');
    }
    
    return Promise.reject(error);
  }
);

export default request;
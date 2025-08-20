// 令牌存储键名
const TOKEN_KEY = 'go_short_link_token';
const USER_INFO_KEY = 'go_short_link_user_info';

// 保存令牌到本地存储
export const setToken = (token: string): void => {
  localStorage.setItem(TOKEN_KEY, token);
};

// 从本地存储获取令牌
export const getToken = (): string | null => {
  return localStorage.getItem(TOKEN_KEY);
};

// 从本地存储移除令牌
export const removeToken = (): void => {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_INFO_KEY);
};

// 保存用户信息到本地存储
export const setUserInfo = (userInfo: any): void => {
  localStorage.setItem(USER_INFO_KEY, JSON.stringify(userInfo));
};

// 从本地存储获取用户信息
export const getUserInfo = (): any => {
  const userInfo = localStorage.getItem(USER_INFO_KEY);
  return userInfo ? JSON.parse(userInfo) : null;
};
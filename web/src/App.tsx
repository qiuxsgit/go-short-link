import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { message } from 'antd';
import Login from './pages/Login';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import ShortLinks from './pages/ShortLinks';
import HistoryLinks from './pages/HistoryLinks';
import { getToken, removeToken } from './utils/auth';

const App: React.FC = () => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(!!getToken());
  
  useEffect(() => {
    // 检查令牌是否有效
    const token = getToken();
    if (token) {
      setIsAuthenticated(true);
    }
  }, []);

  // 登录成功处理
  const handleLoginSuccess = () => {
    setIsAuthenticated(true);
    message.success('登录成功');
  };

  // 登出处理
  const handleLogout = () => {
    removeToken();
    setIsAuthenticated(false);
    message.success('已退出登录');
  };

  return (
    <Router>
      <Routes>
        {/* 公共路由 */}
        <Route path="/login" element={
          !isAuthenticated ? <Login onLoginSuccess={handleLoginSuccess} /> : <Navigate to="/" />
        } />
        
        {/* 受保护的路由 */}
        <Route path="/" element={
          isAuthenticated ? <Layout onLogout={handleLogout} /> : <Navigate to="/login" />
        }>
          <Route index element={<Dashboard />} />
          <Route path="short-links" element={<ShortLinks />} />
          <Route path="history-links" element={<HistoryLinks />} />
        </Route>
        
        {/* 默认重定向 */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
};

export default App;
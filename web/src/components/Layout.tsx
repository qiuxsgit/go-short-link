import React, { useState } from 'react';
import { Layout as AntLayout, Menu, theme, Button } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  LinkOutlined,
  HistoryOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { getUserInfo } from '../utils/auth';

const { Header, Content, Footer } = AntLayout;

interface LayoutProps {
  onLogout: () => void;
}

const Layout: React.FC<LayoutProps> = ({ onLogout }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const [current, setCurrent] = useState(location.pathname);
  const userInfo = getUserInfo();
  
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  const handleMenuClick = (e: { key: string }) => {
    setCurrent(e.key);
    navigate(e.key);
  };

  const handleLogout = () => {
    onLogout();
  };

  return (
    <AntLayout className="layout" style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center' }}>
        <div className="logo" />
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[current]}
          onClick={handleMenuClick}
          items={[
            {
              key: '/',
              icon: <DashboardOutlined />,
              label: '仪表盘',
            },
            {
              key: '/short-links',
              icon: <LinkOutlined />,
              label: '短链接管理',
            },
            {
              key: '/history-links',
              icon: <HistoryOutlined />,
              label: '历史短链接',
            },
          ]}
          style={{ flex: 1 }}
        />
        <div style={{ color: 'white', marginRight: '20px' }}>
          欢迎，{userInfo?.username || '管理员'}
        </div>
        <Button
          type="primary"
          danger
          icon={<LogoutOutlined />}
          onClick={handleLogout}
        >
          退出
        </Button>
      </Header>
      <Content className="content-container">
        <div className="inner-content" style={{ background: colorBgContainer }}>
          <Outlet />
        </div>
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        短链接管理系统 ©{new Date().getFullYear()} 版权所有
      </Footer>
    </AntLayout>
  );
};

export default Layout;
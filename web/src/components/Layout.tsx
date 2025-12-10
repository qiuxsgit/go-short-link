import React, { useState } from 'react';
import { Layout as AntLayout, Menu, theme, Button, Space } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  LinkOutlined,
  HistoryOutlined,
  LogoutOutlined,
  KeyOutlined,
} from '@ant-design/icons';
import { getUserInfo } from '../utils/auth';
import ChangePasswordModal from './ChangePasswordModal';

const { Header, Content, Footer } = AntLayout;

interface LayoutProps {
  onLogout: () => void;
}

const Layout: React.FC<LayoutProps> = ({ onLogout }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const [current, setCurrent] = useState(location.pathname);
  const userInfo = getUserInfo();
  const [changePasswordVisible, setChangePasswordVisible] = useState(false);
  
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
    <AntLayout className="layout" style={{ minHeight: '100vh', background: 'transparent' }}>
      <Header style={{ 
        display: 'flex', 
        alignItems: 'center',
        padding: '0 48px',
        height: 64,
        lineHeight: '64px',
        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
        position: 'sticky',
        top: 0,
        zIndex: 1000,
      }}>
        <div className="logo" style={{ 
          background: 'rgba(255, 255, 255, 0.2)',
          backdropFilter: 'blur(10px)',
          width: 160,
          height: 48,
          margin: '8px 32px 8px 0',
        }}>
          短链管理
        </div>
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
          style={{ 
            flex: 1,
            background: 'transparent',
            borderBottom: 'none',
            minWidth: 0,
          }}
        />
        <div style={{ 
          color: 'rgba(255, 255, 255, 0.95)', 
          marginRight: '24px',
          fontSize: 14,
          fontWeight: 500,
        }}>
          欢迎，{userInfo?.username || '管理员'}
        </div>
        <Space size="middle">
          <Button
            type="primary"
            icon={<KeyOutlined />}
            onClick={() => setChangePasswordVisible(true)}
            style={{
              background: 'rgba(255, 255, 255, 0.2)',
              border: '1px solid rgba(255, 255, 255, 0.3)',
              color: 'white',
            }}
          >
            修改密码
          </Button>
          <Button
            type="primary"
            danger
            icon={<LogoutOutlined />}
            onClick={handleLogout}
            style={{
              background: 'rgba(255, 77, 79, 0.2)',
              border: '1px solid rgba(255, 77, 79, 0.3)',
            }}
          >
            退出
          </Button>
        </Space>
        
        {/* 修改密码对话框 */}
        <ChangePasswordModal
          visible={changePasswordVisible}
          onCancel={() => setChangePasswordVisible(false)}
        />
      </Header>
      <Content className="content-container" style={{ marginTop: 24 }}>
        <div className="inner-content" style={{ background: colorBgContainer }}>
          <Outlet />
        </div>
      </Content>
      <Footer style={{ 
        textAlign: 'center',
        background: 'transparent',
        padding: '20px 48px',
        borderTop: '1px solid rgba(24, 144, 255, 0.1)',
        fontSize: 14,
      }}>
        短链接管理系统 ©{new Date().getFullYear()} 版权所有
      </Footer>
    </AntLayout>
  );
};

export default Layout;
import React from 'react';
import ReactDOM from 'react-dom/client';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/lib/locale/zh_CN';
import App from './App';
import './index.css';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <React.StrictMode>
    <ConfigProvider 
      locale={zhCN}
      theme={{
        token: {
          colorPrimary: '#1890ff',
          colorInfo: '#1890ff',
          borderRadius: 6,
          wireframe: false,
          fontSize: 14,
        },
        components: {
          Layout: {
            headerBg: 'linear-gradient(135deg, #1890ff 0%, #096dd9 100%)',
            headerHeight: 64,
            headerPadding: '0 48px',
          },
          Menu: {
            itemHeight: 48,
            itemMarginInline: 12,
            itemMarginBlock: 4,
            subMenuItemBg: 'rgba(24, 144, 255, 0.06)',
            fontSize: 14,
          },
          Table: {
            headerBg: '#f5f7fa',
            headerColor: '#1f2937',
            rowHoverBg: '#e6f7ff',
            borderColor: '#e5e7eb',
            cellPaddingBlock: 12,
            cellPaddingInline: 16,
            fontSize: 14,
          },
          Card: {
            borderRadius: 8,
            boxShadow: '0 2px 8px rgba(24, 144, 255, 0.08)',
            paddingLG: 24,
          },
          Button: {
            borderRadius: 6,
            controlHeight: 36,
            paddingInline: 20,
            fontSize: 14,
          },
          Input: {
            borderRadius: 6,
            controlHeight: 36,
            paddingInline: 12,
            fontSize: 14,
          },
          Form: {
            verticalLabelPadding: '0 0 8px',
            itemMarginBottom: 20,
            labelFontSize: 14,
          },
        },
      }}
    >
      <App />
    </ConfigProvider>
  </React.StrictMode>
);
import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Statistic } from 'antd';
import { LinkOutlined, HistoryOutlined } from '@ant-design/icons';
import { getShortLinks, getHistoryLinks } from '../api';

const Dashboard: React.FC = () => {
  const [activeLinks, setActiveLinks] = useState<number>(0);
  const [expiredLinks, setExpiredLinks] = useState<number>(0);
  const [historyLinks, setHistoryLinks] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        
        // 获取活跃链接数量
        const activeResponse: any = await getShortLinks({
          page: 1,
          pageSize: 1,
          status: 'active'
        });
        setActiveLinks(activeResponse.total || 0);
        
        // 获取过期链接数量
        const expiredResponse: any = await getShortLinks({
          page: 1,
          pageSize: 1,
          status: 'expired'
        });
        setExpiredLinks(expiredResponse.total || 0);
        
        // 获取历史链接数量（当前月）
        const currentMonth = new Date().toISOString().slice(2, 4) + new Date().toISOString().slice(5, 7); // YYMM
        const historyResponse: any = await getHistoryLinks({
          month: currentMonth,
          page: 1,
          pageSize: 1
        });
        setHistoryLinks(historyResponse.total || 0);
      } catch (error) {
        console.error('获取数据失败:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  return (
    <div>
      <h1>仪表盘</h1>
      <Row gutter={16}>
        <Col span={8}>
          <Card>
            <Statistic
              title="活跃短链接"
              value={activeLinks}
              loading={loading}
              prefix={<LinkOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="过期短链接"
              value={expiredLinks}
              loading={loading}
              prefix={<LinkOutlined />}
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title="历史短链接"
              value={historyLinks}
              loading={loading}
              prefix={<HistoryOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
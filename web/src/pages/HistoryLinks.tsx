import React, { useState, useEffect, useCallback } from 'react';
import { Table, Input, Form, Button, Select, Space, message } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { getHistoryLinks } from '../api';

const { Option } = Select;

const HistoryLinks: React.FC = () => {
  const [links, setLinks] = useState<any[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [searchForm] = Form.useForm();
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
  });
  // 获取默认月份格式 YYMM
  const getDefaultMonth = () => {
    const now = new Date();
    const year = now.getFullYear().toString().slice(2); // 获取年份后两位
    const month = (now.getMonth() + 1).toString().padStart(2, '0'); // 获取月份，补0
    return year + month;
  };

  const [filters, setFilters] = useState({
    month: getDefaultMonth(), // YYMM
    shortCode: '',
    originalUrl: '',
  });

  // 获取历史短链接列表
  const fetchLinks = useCallback(async () => {
    try {
      setLoading(true);
      const response: any = await getHistoryLinks({
        month: filters.month,
        page: pagination.current,
        pageSize: pagination.pageSize,
        shortCode: filters.shortCode || undefined,
        originalUrl: filters.originalUrl || undefined,
      });
      
      setLinks(response.links || []);
      setTotal(response.total || 0);
    } catch (error) {
      message.error('获取历史短链接列表失败，请检查网络连接或联系管理员');
    } finally {
      setLoading(false);
    }
  }, [filters, pagination]);

  // 首次加载和分页/筛选条件变化时获取数据
  useEffect(() => {
    fetchLinks();
  }, [fetchLinks]);

  // 处理表格分页变化
  const handleTableChange = (pagination: any) => {
    setPagination({
      current: pagination.current,
      pageSize: pagination.pageSize,
    });
  };

  // 处理搜索
  const handleSearch = (values: any) => {
    setFilters({
      ...filters,
      shortCode: values.shortCode || '',
      originalUrl: values.originalUrl || '',
      month: values.month || filters.month,
    });
    setPagination({
      ...pagination,
      current: 1,
    });
  };

  // 生成月份选项
  const generateMonthOptions = () => {
    const options = [];
    const currentDate = new Date();
    const currentYear = currentDate.getFullYear();
    const currentMonth = currentDate.getMonth();
    
    // 生成过去12个月的选项
    for (let i = 0; i < 12; i++) {
      let year = currentYear;
      let month = currentMonth - i;
      
      if (month < 0) {
        month += 12;
        year -= 1;
      }
      
      const monthStr = String(month + 1).padStart(2, '0');
      const yearStr = String(year).slice(2);
      const value = `${yearStr}${monthStr}`;
      const label = `${year}年${monthStr}月`;
      
      options.push(
        <Option key={value} value={value}>
          {label}
        </Option>
      );
    }
    
    return options;
  };

  // 表格列定义
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '短链接代码',
      dataIndex: 'shortCode',
      key: 'shortCode',
    },
    {
      title: '原始链接',
      dataIndex: 'originalUrl',
      key: 'originalUrl',
      ellipsis: true,
      render: (text: string) => <a href={text} target="_blank" rel="noopener noreferrer">{text}</a>,
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
    },
    {
      title: '过期时间',
      dataIndex: 'expiresAt',
      key: 'expiresAt',
    },
    {
      title: '访问次数',
      dataIndex: 'accessCount',
      key: 'accessCount',
    },
    {
      title: '最后访问时间',
      dataIndex: 'lastAccess',
      key: 'lastAccess',
      render: (text: string) => text || '无访问记录',
    },
  ];

  return (
    <div>
      <h1>历史短链接</h1>
      
      {/* 搜索和操作区域 */}
      <div className="table-operations" style={{ marginBottom: 16 }}>
        <Form
          form={searchForm}
          layout="inline"
          onFinish={handleSearch}
          initialValues={{ month: filters.month }}
          style={{ marginBottom: 16 }}
        >
          <Form.Item name="month" label="月份">
            <Select style={{ width: 120 }}>
              {generateMonthOptions()}
            </Select>
          </Form.Item>
          <Form.Item name="shortCode">
            <Input
              placeholder="短链接代码"
              prefix={<SearchOutlined />}
              allowClear
            />
          </Form.Item>
          <Form.Item name="originalUrl">
            <Input
              placeholder="原始链接"
              prefix={<SearchOutlined />}
              allowClear
            />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                搜索
              </Button>
              <Button onClick={() => searchForm.resetFields()}>
                重置
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </div>
      
      {/* 历史短链接表格 */}
      <Table
        columns={columns}
        dataSource={links}
        rowKey="id"
        pagination={{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条记录`,
        }}
        loading={loading}
        onChange={handleTableChange}
      />
    </div>
  );
};

export default HistoryLinks;
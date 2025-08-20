import React, { useState, useEffect } from 'react';
import { Table, Button, Input, Space, Modal, Form, InputNumber, message, Tag } from 'antd';
import { SearchOutlined, PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import { getShortLinks, createShortLink, deleteShortLink } from '../api';

const { confirm } = Modal;

const ShortLinks: React.FC = () => {
  const [links, setLinks] = useState<any[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [searchForm] = Form.useForm();
  const [createForm] = Form.useForm();
  const [createModalVisible, setCreateModalVisible] = useState<boolean>(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
  });
  const [filters, setFilters] = useState({
    shortCode: '',
    originalUrl: '',
    status: 'active',
  });

  // 获取短链接列表
  const fetchLinks = async () => {
    try {
      setLoading(true);
      const response: any = await getShortLinks({
        page: pagination.current,
        pageSize: pagination.pageSize,
        shortCode: filters.shortCode || undefined,
        originalUrl: filters.originalUrl || undefined,
        status: filters.status as any,
      });
      
      setLinks(response.links || []);
      setTotal(response.total || 0);
    } catch (error) {
      console.error('获取短链接列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 首次加载和分页/筛选条件变化时获取数据
  useEffect(() => {
    fetchLinks();
  }, [pagination, filters]);

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
    });
    setPagination({
      ...pagination,
      current: 1,
    });
  };

  // 处理状态筛选
  const handleStatusChange = (status: string) => {
    setFilters({
      ...filters,
      status: status,
    });
    setPagination({
      ...pagination,
      current: 1,
    });
  };

  // 处理创建短链接
  const handleCreate = async (values: any) => {
    try {
      await createShortLink(values);
      message.success('创建短链接成功');
      setCreateModalVisible(false);
      createForm.resetFields();
      fetchLinks();
    } catch (error) {
      console.error('创建短链接失败:', error);
    }
  };

  // 处理删除短链接
  const handleDelete = (id: number) => {
    confirm({
      title: '确认删除',
      content: '确定要删除这个短链接吗？删除后将移动到历史表。',
      okText: '确认',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        try {
          await deleteShortLink(id);
          message.success('删除短链接成功');
          fetchLinks();
        } catch (error) {
          console.error('删除短链接失败:', error);
        }
      },
    });
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
      render: (text: string) => <a href={`/s/${text}`} target="_blank" rel="noopener noreferrer">{text}</a>,
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
      title: '状态',
      key: 'status',
      render: (record: any) => {
        const now = new Date().toISOString().replace('T', ' ').substring(0, 19);
        // 比较格式化后的日期字符串
        const isExpired = now > record.expiresAt;
        return isExpired ? (
          <Tag color="red">已过期</Tag>
        ) : (
          <Tag color="green">有效</Tag>
        );
      },
    },
    {
      title: '访问次数',
      dataIndex: 'accessCount',
      key: 'accessCount',
    },
    {
      title: '操作',
      key: 'action',
      render: (record: any) => (
        <Button
          type="primary"
          danger
          icon={<DeleteOutlined />}
          onClick={() => handleDelete(record.id)}
        >
          删除
        </Button>
      ),
    },
  ];

  return (
    <div>
      <h1>短链接管理</h1>
      
      {/* 搜索和操作区域 */}
      <div className="table-operations">
        <Form
          form={searchForm}
          layout="inline"
          onFinish={handleSearch}
          style={{ marginBottom: 16 }}
        >
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
            <Button type="primary" htmlType="submit">
              搜索
            </Button>
          </Form.Item>
          <Form.Item>
            <Button onClick={() => searchForm.resetFields()}>
              重置
            </Button>
          </Form.Item>
        </Form>
        
        <Space style={{ marginBottom: 16 }}>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setCreateModalVisible(true)}
          >
            创建短链接
          </Button>
          
          <Button
            type={filters.status === 'active' ? 'primary' : 'default'}
            onClick={() => handleStatusChange('active')}
          >
            有效链接
          </Button>
          
          <Button
            type={filters.status === 'expired' ? 'primary' : 'default'}
            onClick={() => handleStatusChange('expired')}
          >
            过期链接
          </Button>
        </Space>
      </div>
      
      {/* 短链接表格 */}
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
      
      {/* 创建短链接对话框 */}
      <Modal
        title="创建短链接"
        open={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        footer={null}
      >
        <Form
          form={createForm}
          layout="vertical"
          onFinish={handleCreate}
        >
          <Form.Item
            name="link"
            label="原始链接"
            rules={[{ required: true, message: '请输入原始链接' }]}
          >
            <Input placeholder="请输入需要转换的链接" />
          </Form.Item>
          <Form.Item
            name="expire"
            label="有效期（秒）"
            initialValue={3600}
            rules={[{ required: true, message: '请输入有效期' }]}
          >
            <InputNumber
              min={60}
              max={31536000} // 1年
              style={{ width: '100%' }}
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              创建
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default ShortLinks;
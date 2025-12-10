import React, { useState, useEffect } from 'react';
import { Card, Table, Tag, Typography, Space, Collapse, Divider } from 'antd';
import { BookOutlined, ApiOutlined, LockOutlined, InfoCircleOutlined } from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;
const { Panel } = Collapse;

interface ApiDocProps {}

const ApiDoc: React.FC<ApiDocProps> = () => {
  const [adminBaseURL, setAdminBaseURL] = useState<string>('');
  const [accessBaseURL, setAccessBaseURL] = useState<string>('');

  useEffect(() => {
    // 从当前 URL 获取基础地址
    const origin = window.location.origin;
    setAdminBaseURL(origin);
    // 访问 API 通常在不同端口，这里可以根据实际情况调整
    setAccessBaseURL(origin.replace(':8081', ':8082'));
  }, []);

  const scrollToTop = () => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const methodBadge = (method: string) => {
    const colors: { [key: string]: string } = {
      GET: 'success',
      POST: 'processing',
      DELETE: 'error',
      PUT: 'warning',
    };
    return <Tag color={colors[method] || 'default'}>{method}</Tag>;
  };

  const statusBadge = (code: number) => {
    const colors: { [key: number]: string } = {
      200: 'success',
      307: 'processing',
      400: 'warning',
      401: 'warning',
      404: 'error',
      500: 'error',
    };
    return <Tag color={colors[code] || 'default'}>{code}</Tag>;
  };

  return (
    <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
      {/* 头部 */}
      <Card
        style={{
          marginBottom: 24,
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          border: 'none',
        }}
        bodyStyle={{ padding: '40px', textAlign: 'center' }}
      >
        <Title level={1} style={{ color: 'white', marginBottom: 10 }}>
          <BookOutlined /> API 接口文档
        </Title>
        <Paragraph style={{ color: 'rgba(255, 255, 255, 0.9)', fontSize: 16, margin: 0 }}>
          短链接服务完整接口说明
        </Paragraph>
      </Card>

      {/* 基础信息 */}
      <Card
        title={
          <span>
            <InfoCircleOutlined /> 基础信息
          </span>
        }
        style={{ marginBottom: 24 }}
      >
        <Table
          dataSource={[
            { key: '1', name: '管理API地址', value: adminBaseURL },
            { key: '2', name: '访问API地址', value: accessBaseURL },
            { key: '3', name: '认证方式', value: 'JWT Token (Bearer Token)' },
            { key: '4', name: '内容类型', value: 'application/json' },
            { key: '5', name: '字符编码', value: 'UTF-8' },
          ]}
          columns={[
            {
              title: '项目',
              dataIndex: 'name',
              width: 150,
              render: (text) => <strong>{text}</strong>,
            },
            {
              title: '值',
              dataIndex: 'value',
            },
          ]}
          pagination={false}
          showHeader={false}
        />
      </Card>

      {/* 管理API接口 */}
      <Card
        title={
          <span>
            <ApiOutlined /> 管理API接口
          </span>
        }
        style={{ marginBottom: 24 }}
      >
        <Collapse defaultActiveKey={['1']} ghost>
          {/* 1. 登录 */}
          <Panel
            header={
              <Space>
                {methodBadge('POST')}
                <Text code>/api/login</Text>
                <Text type="secondary">管理员登录</Text>
              </Space>
            }
            key="1"
          >
            <Paragraph>登录获取访问令牌</Paragraph>
            <Title level={5}>请求参数</Title>
            <Table
              dataSource={[
                { param: 'username', type: 'string', required: '是', desc: '管理员用户名' },
                { param: 'password', type: 'string', required: '是', desc: '管理员密码' },
              ]}
              columns={[
                { title: '参数名', dataIndex: 'param' },
                { title: '类型', dataIndex: 'type' },
                { title: '必填', dataIndex: 'required' },
                { title: '说明', dataIndex: 'desc' },
              ]}
              pagination={false}
              size="small"
            />
            <Title level={5}>请求示例</Title>
            <Card size="small" style={{ backgroundColor: '#282c34', color: '#abb2bf' }}>
              <pre style={{ margin: 0, color: '#abb2bf' }}>
                {`{
  "username": "admin",
  "password": "123456"
}`}
              </pre>
            </Card>
            <Title level={5}>响应示例</Title>
            <Space direction="vertical" style={{ width: '100%' }}>
              {statusBadge(200)}
              <Card size="small" style={{ backgroundColor: '#282c34', color: '#abb2bf' }}>
                <pre style={{ margin: 0, color: '#abb2bf' }}>
                  {`{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "admin",
  "userId": 1
}`}
                </pre>
              </Card>
            </Space>
          </Panel>

          {/* 2. 创建短链接 */}
          <Panel
            header={
              <Space>
                {methodBadge('POST')}
                <Text code>/api/short-link/create</Text>
                <Text type="secondary">创建短链接</Text>
              </Space>
            }
            key="2"
          >
            <Paragraph>创建新的短链接</Paragraph>
            <Title level={5}>请求参数</Title>
            <Table
              dataSource={[
                { param: 'link', type: 'string', required: '是', desc: '原始URL地址' },
                { param: 'expire', type: 'int', required: '是', desc: '过期时间（秒）' },
              ]}
              columns={[
                { title: '参数名', dataIndex: 'param' },
                { title: '类型', dataIndex: 'type' },
                { title: '必填', dataIndex: 'required' },
                { title: '说明', dataIndex: 'desc' },
              ]}
              pagination={false}
              size="small"
            />
            <Title level={5}>请求示例</Title>
            <Card size="small" style={{ backgroundColor: '#282c34', color: '#abb2bf' }}>
              <pre style={{ margin: 0, color: '#abb2bf' }}>
                {`{
  "link": "https://www.example.com",
  "expire": 3600
}`}
              </pre>
            </Card>
            <Title level={5}>响应示例</Title>
            <Space direction="vertical" style={{ width: '100%' }}>
              {statusBadge(200)}
              <Card size="small" style={{ backgroundColor: '#282c34', color: '#abb2bf' }}>
                <pre style={{ margin: 0, color: '#abb2bf' }}>
                  {`{
  "shortLink": "${accessBaseURL}/s/abc123"
}`}
                </pre>
              </Card>
            </Space>
          </Panel>

          {/* 3. 获取短链接列表 */}
          <Panel
            header={
              <Space>
                {methodBadge('GET')}
                <Text code>/api/short-link/list</Text>
                <Tag icon={<LockOutlined />} color="warning">需要认证</Tag>
                <Text type="secondary">获取短链接列表</Text>
              </Space>
            }
            key="3"
          >
            <Paragraph>获取有效的短链接列表，支持分页和筛选</Paragraph>
            <Title level={5}>查询参数</Title>
            <Table
              dataSource={[
                { param: 'page', type: 'string', required: '否', default: '1', desc: '页码' },
                { param: 'pageSize', type: 'string', required: '否', default: '10', desc: '每页数量（最大100）' },
                { param: 'shortCode', type: 'string', required: '否', default: '-', desc: '短码筛选（模糊查询）' },
                { param: 'originalUrl', type: 'string', required: '否', default: '-', desc: '原始URL筛选（模糊查询）' },
                { param: 'status', type: 'string', required: '否', default: '-', desc: '状态：active(有效) 或 expired(已过期)' },
              ]}
              columns={[
                { title: '参数名', dataIndex: 'param' },
                { title: '类型', dataIndex: 'type' },
                { title: '必填', dataIndex: 'required' },
                { title: '默认值', dataIndex: 'default' },
                { title: '说明', dataIndex: 'desc' },
              ]}
              pagination={false}
              size="small"
            />
            <Title level={5}>响应示例</Title>
            <Space direction="vertical" style={{ width: '100%' }}>
              {statusBadge(200)}
              <Card size="small" style={{ backgroundColor: '#282c34', color: '#abb2bf' }}>
                <pre style={{ margin: 0, color: '#abb2bf' }}>
                  {`{
  "total": 100,
  "links": [
    {
      "id": 1,
      "shortCode": "abc123",
      "originalUrl": "https://www.example.com",
      "createdAt": "2024-01-01 10:00:00.000",
      "expiresAt": "2024-01-02 10:00:00.000",
      "accessCount": 42,
      "lastAccess": "2024-01-01 15:30:00.000"
    }
  ]
}`}
                </pre>
              </Card>
            </Space>
          </Panel>

          {/* 4. 获取历史短链接列表 */}
          <Panel
            header={
              <Space>
                {methodBadge('GET')}
                <Text code>/api/short-link/history</Text>
                <Tag icon={<LockOutlined />} color="warning">需要认证</Tag>
                <Text type="secondary">获取历史短链接列表</Text>
              </Space>
            }
            key="4"
          >
            <Paragraph>获取已归档到历史表的短链接列表</Paragraph>
            <Title level={5}>查询参数</Title>
            <Table
              dataSource={[
                { param: 'month', type: 'string', required: '否', default: '当前月份', desc: '月份（格式：YYMM，如2401）' },
                { param: 'page', type: 'string', required: '否', default: '1', desc: '页码' },
                { param: 'pageSize', type: 'string', required: '否', default: '10', desc: '每页数量（最大100）' },
                { param: 'shortCode', type: 'string', required: '否', default: '-', desc: '短码筛选（模糊查询）' },
                { param: 'originalUrl', type: 'string', required: '否', default: '-', desc: '原始URL筛选（模糊查询）' },
              ]}
              columns={[
                { title: '参数名', dataIndex: 'param' },
                { title: '类型', dataIndex: 'type' },
                { title: '必填', dataIndex: 'required' },
                { title: '默认值', dataIndex: 'default' },
                { title: '说明', dataIndex: 'desc' },
              ]}
              pagination={false}
              size="small"
            />
          </Panel>

          {/* 5. 删除短链接 */}
          <Panel
            header={
              <Space>
                {methodBadge('DELETE')}
                <Text code>/api/short-link/:id</Text>
                <Tag icon={<LockOutlined />} color="warning">需要认证</Tag>
                <Text type="secondary">删除短链接</Text>
              </Space>
            }
            key="5"
          >
            <Paragraph>删除指定的短链接（移动到历史表）</Paragraph>
            <Title level={5}>路径参数</Title>
            <Table
              dataSource={[{ param: 'id', type: 'string', required: '是', desc: '短链接ID' }]}
              columns={[
                { title: '参数名', dataIndex: 'param' },
                { title: '类型', dataIndex: 'type' },
                { title: '必填', dataIndex: 'required' },
                { title: '说明', dataIndex: 'desc' },
              ]}
              pagination={false}
              size="small"
            />
            <Title level={5}>响应示例</Title>
            <Space direction="vertical" style={{ width: '100%' }}>
              {statusBadge(200)}
              <Card size="small" style={{ backgroundColor: '#282c34', color: '#abb2bf' }}>
                <pre style={{ margin: 0, color: '#abb2bf' }}>
                  {`{
  "message": "短链接已成功删除"
}`}
                </pre>
              </Card>
            </Space>
          </Panel>

          {/* 6. 修改密码 */}
          <Panel
            header={
              <Space>
                {methodBadge('POST')}
                <Text code>/api/change-password</Text>
                <Tag icon={<LockOutlined />} color="warning">需要认证</Tag>
                <Text type="secondary">修改密码</Text>
              </Space>
            }
            key="6"
          >
            <Paragraph>修改当前登录用户的密码</Paragraph>
            <Title level={5}>请求参数</Title>
            <Table
              dataSource={[
                { param: 'currentPassword', type: 'string', required: '是', desc: '当前密码' },
                { param: 'newPassword', type: 'string', required: '是', desc: '新密码（最少6位）' },
                { param: 'confirmPassword', type: 'string', required: '是', desc: '确认新密码' },
              ]}
              columns={[
                { title: '参数名', dataIndex: 'param' },
                { title: '类型', dataIndex: 'type' },
                { title: '必填', dataIndex: 'required' },
                { title: '说明', dataIndex: 'desc' },
              ]}
              pagination={false}
              size="small"
            />
            <Title level={5}>请求示例</Title>
            <Card size="small" style={{ backgroundColor: '#282c34', color: '#abb2bf' }}>
              <pre style={{ margin: 0, color: '#abb2bf' }}>
                {`{
  "currentPassword": "old_password",
  "newPassword": "new_password",
  "confirmPassword": "new_password"
}`}
              </pre>
            </Card>
          </Panel>
        </Collapse>
      </Card>

      {/* 访问API接口 */}
      <Card
        title={
          <span>
            <ApiOutlined /> 访问API接口
          </span>
        }
        style={{ marginBottom: 24 }}
      >
        <Collapse defaultActiveKey={['1']} ghost>
          {/* 短链接重定向 */}
          <Panel
            header={
              <Space>
                {methodBadge('GET')}
                <Text code>/s/:code</Text>
                <Text type="secondary">短链接重定向</Text>
              </Space>
            }
            key="1"
          >
            <Paragraph>访问短链接时自动重定向到原始URL</Paragraph>
            <Title level={5}>路径参数</Title>
            <Table
              dataSource={[{ param: 'code', type: 'string', required: '是', desc: '短码' }]}
              columns={[
                { title: '参数名', dataIndex: 'param' },
                { title: '类型', dataIndex: 'type' },
                { title: '必填', dataIndex: 'required' },
                { title: '说明', dataIndex: 'desc' },
              ]}
              pagination={false}
              size="small"
            />
            <Title level={5}>响应状态码</Title>
            <Space direction="vertical">
              <div>
                {statusBadge(307)} - 成功重定向到原始URL
              </div>
              <div>
                {statusBadge(404)} - 短链接不存在或已过期
              </div>
            </Space>
          </Panel>
        </Collapse>
      </Card>

      {/* 错误码说明 */}
      <Card
        title={
          <span>
            <InfoCircleOutlined /> 错误码说明
          </span>
        }
        style={{ marginBottom: 24 }}
      >
        <Table
          dataSource={[
            { code: 200, desc: '请求成功' },
            { code: 307, desc: '临时重定向' },
            { code: 400, desc: '请求参数错误' },
            { code: 401, desc: '未认证或认证失败' },
            { code: 404, desc: '资源不存在' },
            { code: 500, desc: '服务器内部错误' },
          ]}
          columns={[
            {
              title: '状态码',
              dataIndex: 'code',
              render: (code) => statusBadge(code),
            },
            { title: '说明', dataIndex: 'desc' },
          ]}
          pagination={false}
        />
      </Card>

      {/* 认证说明 */}
      <Card
        title={
          <span>
            <LockOutlined /> 认证说明
          </span>
        }
        style={{ marginBottom: 24 }}
      >
        <Paragraph>
          需要认证的接口需要在请求头中携带JWT Token：
        </Paragraph>
        <Card size="small" style={{ backgroundColor: '#282c34', color: '#abb2bf', marginTop: 16 }}>
          <pre style={{ margin: 0, color: '#abb2bf' }}>
            {`Authorization: Bearer <token>`}
          </pre>
        </Card>
        <Paragraph style={{ marginTop: 16 }}>
          Token可以通过登录接口获取，默认有效期为24小时（可在配置文件中修改）。
        </Paragraph>
      </Card>

      {/* 返回顶部按钮 */}
      <div
        onClick={scrollToTop}
        style={{
          position: 'fixed',
          bottom: 30,
          right: 30,
          width: 50,
          height: 50,
          borderRadius: '50%',
          background: '#667eea',
          color: 'white',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          cursor: 'pointer',
          boxShadow: '0 4px 12px rgba(0,0,0,0.3)',
          transition: 'all 0.3s ease',
          fontSize: 20,
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.background = '#764ba2';
          e.currentTarget.style.transform = 'translateY(-3px)';
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.background = '#667eea';
          e.currentTarget.style.transform = 'translateY(0)';
        }}
      >
        ↑
      </div>
    </div>
  );
};

export default ApiDoc;

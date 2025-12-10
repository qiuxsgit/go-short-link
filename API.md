# API 接口文档

## 概述

本短链接服务提供两套API服务：

- **管理API服务**：端口 `8081`，提供短链接创建、管理和用户认证功能
- **访问API服务**：端口 `8082`，提供短链接重定向功能

所有API请求和响应均使用JSON格式，字符编码为UTF-8。

## 基础信息

- **管理API Base URL**: `http://localhost:8081`
- **访问API Base URL**: `http://localhost:8082`
- **认证方式**: JWT Token（Bearer Token）
- **Content-Type**: `application/json`

## 认证说明

需要认证的接口需要在请求头中携带JWT Token：

```
Authorization: Bearer <token>
```

Token可以通过登录接口获取，默认有效期为24小时（可在配置文件中修改）。

---

## 管理API接口

### 1. 管理员登录

登录获取访问令牌。

**接口地址**: `POST /api/login`

**认证要求**: 无需认证

**请求参数**:

```json
{
  "username": "admin",
  "password": "123456"
}
```

| 参数名   | 类型   | 必填 | 说明   |
|---------|--------|------|--------|
| username | string | 是   | 管理员用户名 |
| password | string | 是   | 管理员密码   |

**响应示例**:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "username": "admin",
  "userId": 1
}
```

**响应字段说明**:

| 字段名   | 类型   | 说明           |
|---------|--------|----------------|
| token    | string | JWT访问令牌    |
| username | string | 用户名         |
| userId   | int64  | 用户ID         |

**错误响应**:

- `400 Bad Request`: 请求参数无效
- `401 Unauthorized`: 用户名或密码错误

---

### 2. 创建短链接

创建新的短链接。

**接口地址**: `POST /api/short-link/create`

**认证要求**: 无需认证

**请求参数**:

```json
{
  "link": "https://www.example.com/very/long/url",
  "expire": 3600
}
```

| 参数名 | 类型   | 必填 | 说明                           |
|-------|--------|------|--------------------------------|
| link   | string | 是   | 原始URL地址                    |
| expire | int    | 是   | 过期时间（秒），从创建时开始计算 |

**响应示例**:

```json
{
  "shortLink": "http://localhost:8082/s/abc123"
}
```

**响应字段说明**:

| 字段名    | 类型   | 说明              |
|----------|--------|-------------------|
| shortLink | string | 生成的短链接完整URL |

**错误响应**:

- `400 Bad Request`: 请求参数无效
- `500 Internal Server Error`: 创建短链接失败

---

### 3. 获取短链接列表

获取有效的短链接列表，支持分页和筛选。

**接口地址**: `GET /api/short-link/list`

**认证要求**: 需要认证

**请求头**:

```
Authorization: Bearer <token>
```

**查询参数**:

| 参数名      | 类型   | 必填 | 默认值 | 说明                                           |
|------------|--------|------|--------|------------------------------------------------|
| page       | string | 否   | "1"    | 页码                                           |
| pageSize   | string | 否   | "10"   | 每页数量（最大100）                             |
| shortCode  | string | 否   | -      | 短码筛选（支持模糊查询）                         |
| originalUrl | string | 否   | -      | 原始URL筛选（支持模糊查询）                      |
| status     | string | 否   | -      | 状态筛选：`active`(有效) 或 `expired`(已过期)   |

**请求示例**:

```
GET /api/short-link/list?page=1&pageSize=20&status=active&shortCode=abc
```

**响应示例**:

```json
{
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
}
```

**响应字段说明**:

| 字段名       | 类型    | 说明                    |
|------------|---------|-------------------------|
| total      | int64   | 总记录数                |
| links      | array   | 短链接列表              |
| links[].id | int64   | 短链接ID                |
| links[].shortCode | string | 短码                |
| links[].originalUrl | string | 原始URL          |
| links[].createdAt | string | 创建时间（格式：YYYY-MM-DD HH:mm:ss.SSS） |
| links[].expiresAt | string | 过期时间（格式：YYYY-MM-DD HH:mm:ss.SSS） |
| links[].accessCount | int64 | 访问次数        |
| links[].lastAccess | string | 最后访问时间（格式：YYYY-MM-DD HH:mm:ss.SSS） |

**错误响应**:

- `401 Unauthorized`: 未提供认证令牌或令牌无效/过期
- `500 Internal Server Error`: 查询数据失败

---

### 4. 获取历史短链接列表

获取已删除或过期并已归档到历史表的短链接列表。

**接口地址**: `GET /api/short-link/history`

**认证要求**: 需要认证

**请求头**:

```
Authorization: Bearer <token>
```

**查询参数**:

| 参数名       | 类型   | 必填 | 默认值      | 说明                                    |
|------------|--------|------|-------------|-----------------------------------------|
| month      | string | 否   | 当前月份    | 月份（格式：YYMM，如2401表示2024年1月） |
| page       | string | 否   | "1"         | 页码                                    |
| pageSize   | string | 否   | "10"        | 每页数量（最大100）                      |
| shortCode  | string | 否   | -           | 短码筛选（支持模糊查询）                  |
| originalUrl | string | 否   | -          | 原始URL筛选（支持模糊查询）               |

**请求示例**:

```
GET /api/short-link/history?month=2401&page=1&pageSize=20
```

**响应示例**:

```json
{
  "total": 50,
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
  ],
  "debug_month": "2401",
  "debug_table": "short_links_history_2401",
  "debug_exists": true,
  "debug_count": 10
}
```

**响应字段说明**:

响应结构与获取短链接列表相同，额外包含以下调试字段：

| 字段名        | 类型   | 说明                  |
|-------------|--------|-----------------------|
| debug_month | string | 查询的月份            |
| debug_table | string | 历史表名              |
| debug_exists | bool  | 历史表是否存在        |
| debug_count | int    | 当前返回的记录数      |

**错误响应**:

- `401 Unauthorized`: 未提供认证令牌或令牌无效/过期
- `500 Internal Server Error`: 查询数据失败

---

### 5. 删除短链接

删除指定的短链接（将移动到历史表）。

**接口地址**: `DELETE /api/short-link/:id`

**认证要求**: 需要认证

**请求头**:

```
Authorization: Bearer <token>
```

**路径参数**:

| 参数名 | 类型   | 必填 | 说明      |
|-------|--------|------|-----------|
| id    | string | 是   | 短链接ID  |

**请求示例**:

```
DELETE /api/short-link/123
```

**响应示例**:

```json
{
  "message": "短链接已成功删除"
}
```

**响应字段说明**:

| 字段名  | 类型   | 说明           |
|--------|--------|----------------|
| message | string | 操作结果消息   |

**错误响应**:

- `400 Bad Request`: 无效的短链接ID
- `401 Unauthorized`: 未提供认证令牌或令牌无效/过期
- `404 Not Found`: 短链接不存在
- `500 Internal Server Error`: 删除操作失败

---

### 6. 修改密码

修改当前登录用户的密码。

**接口地址**: `POST /api/change-password`

**认证要求**: 需要认证

**请求头**:

```
Authorization: Bearer <token>
```

**请求参数**:

```json
{
  "currentPassword": "old_password",
  "newPassword": "new_password",
  "confirmPassword": "new_password"
}
```

| 参数名          | 类型   | 必填 | 说明                         |
|---------------|--------|------|------------------------------|
| currentPassword | string | 是   | 当前密码                     |
| newPassword   | string | 是   | 新密码（最少6位）             |
| confirmPassword | string | 是   | 确认新密码（必须与新密码相同） |

**响应示例**:

```json
{
  "message": "密码修改成功"
}
```

**响应字段说明**:

| 字段名  | 类型   | 说明           |
|--------|--------|----------------|
| message | string | 操作结果消息   |

**错误响应**:

- `400 Bad Request`: 请求参数无效或当前密码错误
- `401 Unauthorized`: 未提供认证令牌或令牌无效/过期
- `404 Not Found`: 用户不存在
- `500 Internal Server Error`: 密码更新失败

---

## 访问API接口

### 1. 短链接重定向

访问短链接时自动重定向到原始URL。

**接口地址**: `GET /s/:code`

**认证要求**: 无需认证

**路径参数**:

| 参数名 | 类型   | 必填 | 说明   |
|-------|--------|------|--------|
| code  | string | 是   | 短码   |

**请求示例**:

```
GET http://localhost:8082/s/abc123
```

**响应**:

- `307 Temporary Redirect`: 成功重定向到原始URL
- `404 Not Found`: 短链接不存在或已过期

**404响应示例**:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Page Not Found</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            text-align: center;
            padding-top: 100px;
            background-color: #f7f7f7;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #fff;
            border-radius: 5px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #e74c3c;
        }
        p {
            color: #7f8c8d;
            font-size: 18px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>404 - Page Not Found</h1>
        <p>抱歉，您访问的短链接不存在或已过期。</p>
    </div>
</body>
</html>
```

**说明**:

- 访问短链接时，系统会自动检查短链接是否存在且未过期
- 如果短链接有效，会返回307重定向响应，浏览器会自动跳转到原始URL
- 每次访问会自动更新访问计数和最后访问时间

---

## 错误码说明

### HTTP状态码

| 状态码 | 说明           |
|--------|----------------|
| 200    | 请求成功       |
| 307    | 临时重定向     |
| 400    | 请求参数错误   |
| 401    | 未认证或认证失败 |
| 403    | 无权限访问     |
| 404    | 资源不存在     |
| 500    | 服务器内部错误 |

### 错误响应格式

所有错误响应均为JSON格式：

```json
{
  "error": "错误描述信息"
}
```

---

## 接口调用示例

### 使用curl调用

#### 1. 登录

```bash
curl -X POST http://localhost:8081/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'
```

#### 2. 创建短链接

```bash
curl -X POST http://localhost:8081/api/short-link/create \
  -H "Content-Type: application/json" \
  -d '{
    "link": "https://www.example.com",
    "expire": 3600
  }'
```

#### 3. 获取短链接列表（需要认证）

```bash
curl -X GET "http://localhost:8081/api/short-link/list?page=1&pageSize=10" \
  -H "Authorization: Bearer <your_token>"
```

#### 4. 删除短链接（需要认证）

```bash
curl -X DELETE http://localhost:8081/api/short-link/123 \
  -H "Authorization: Bearer <your_token>"
```

#### 5. 访问短链接

```bash
curl -L http://localhost:8082/s/abc123
```

### 使用JavaScript调用

#### 登录并获取Token

```javascript
const login = async () => {
  const response = await fetch('http://localhost:8081/api/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      username: 'admin',
      password: '123456',
    }),
  });
  
  const data = await response.json();
  return data.token;
};
```

#### 创建短链接

```javascript
const createShortLink = async (link, expire) => {
  const response = await fetch('http://localhost:8081/api/short-link/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      link: link,
      expire: expire,
    }),
  });
  
  return await response.json();
};
```

#### 获取短链接列表（需要认证）

```javascript
const getShortLinks = async (token, page = 1, pageSize = 10) => {
  const response = await fetch(
    `http://localhost:8081/api/short-link/list?page=${page}&pageSize=${pageSize}`,
    {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    }
  );
  
  return await response.json();
};
```

---

## 注意事项

1. **Token过期**: JWT Token默认有效期为24小时，过期后需要重新登录获取新Token。

2. **分页限制**: 每页最大数量限制为100条记录。

3. **短链接过期**: 过期的短链接会自动失效，访问时会返回404错误。

4. **历史数据**: 已删除的短链接会移动到历史表，历史表按月份存储（格式：`short_links_history_YYMM`）。

5. **访问统计**: 每次访问短链接时，系统会异步更新访问计数和最后访问时间。

6. **缓存机制**: 系统使用内存缓存提高短链接查询性能，删除短链接时会同时清除缓存。

7. **CORS**: 如果前端应用与API服务不在同一域名，需要配置CORS支持跨域访问。

---

## 更新日志

- 2024-01-01: 初始版本API文档


package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qiuxsgit/go-short-link/conf"
)

// DocHandler 处理API文档相关的请求
type DocHandler struct {
	config *conf.Config
}

// NewDocHandler 创建一个新的文档处理器
func NewDocHandler(config *conf.Config) *DocHandler {
	return &DocHandler{
		config: config,
	}
}

// ShowAPIDoc 显示API文档页面
func (h *DocHandler) ShowAPIDoc(c *gin.Context) {
	adminBaseURL := h.config.Server.Admin.BaseURL
	accessBaseURL := h.config.Server.Access.BaseURL

	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API 接口文档 - 短链接服务</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: #fff;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            font-weight: 700;
        }
        
        .header p {
            font-size: 1.1em;
            opacity: 0.9;
        }
        
        .content {
            padding: 40px;
        }
        
        .info-section {
            background: #f8f9fa;
            padding: 25px;
            border-radius: 8px;
            margin-bottom: 30px;
            border-left: 4px solid #667eea;
        }
        
        .info-section h2 {
            color: #667eea;
            margin-bottom: 15px;
            font-size: 1.5em;
        }
        
        .info-section ul {
            list-style: none;
            padding-left: 0;
        }
        
        .info-section li {
            padding: 8px 0;
            border-bottom: 1px solid #e9ecef;
        }
        
        .info-section li:last-child {
            border-bottom: none;
        }
        
        .info-section strong {
            color: #495057;
            display: inline-block;
            width: 150px;
        }
        
        .api-section {
            margin-bottom: 50px;
        }
        
        .api-section h2 {
            color: #667eea;
            font-size: 2em;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 3px solid #667eea;
        }
        
        .api-endpoint {
            background: #fff;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            margin-bottom: 25px;
            overflow: hidden;
            transition: all 0.3s ease;
        }
        
        .api-endpoint:hover {
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
            transform: translateY(-2px);
        }
        
        .endpoint-header {
            background: #f8f9fa;
            padding: 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            border-bottom: 1px solid #e9ecef;
        }
        
        .method-badge {
            display: inline-block;
            padding: 6px 12px;
            border-radius: 4px;
            font-weight: 600;
            font-size: 0.9em;
            margin-right: 15px;
            min-width: 70px;
            text-align: center;
        }
        
        .method-get { background: #28a745; color: white; }
        .method-post { background: #007bff; color: white; }
        .method-delete { background: #dc3545; color: white; }
        .method-put { background: #ffc107; color: #333; }
        
        .endpoint-path {
            font-family: 'Courier New', monospace;
            font-size: 1.1em;
            color: #495057;
            flex: 1;
        }
        
        .auth-badge {
            background: #ffc107;
            color: #333;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.85em;
            font-weight: 600;
        }
        
        .endpoint-body {
            padding: 25px;
        }
        
        .endpoint-description {
            color: #6c757d;
            margin-bottom: 20px;
            font-size: 1.05em;
        }
        
        .params-section {
            margin-bottom: 20px;
        }
        
        .params-section h3 {
            color: #495057;
            margin-bottom: 12px;
            font-size: 1.2em;
        }
        
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 20px;
        }
        
        table th,
        table td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #e9ecef;
        }
        
        table th {
            background: #f8f9fa;
            font-weight: 600;
            color: #495057;
        }
        
        table tr:hover {
            background: #f8f9fa;
        }
        
        .code-block {
            background: #282c34;
            color: #abb2bf;
            padding: 20px;
            border-radius: 6px;
            overflow-x: auto;
            margin: 15px 0;
            font-family: 'Courier New', monospace;
            font-size: 0.9em;
            line-height: 1.5;
        }
        
        .code-block pre {
            margin: 0;
            white-space: pre-wrap;
            word-wrap: break-word;
        }
        
        .json-key { color: #e06c75; }
        .json-string { color: #98c379; }
        .json-number { color: #d19a66; }
        .json-boolean { color: #56b6c2; }
        
        .response-section {
            margin-top: 20px;
        }
        
        .status-code {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-weight: 600;
            font-size: 0.85em;
            margin-right: 10px;
        }
        
        .status-200 { background: #28a745; color: white; }
        .status-307 { background: #17a2b8; color: white; }
        .status-400 { background: #ffc107; color: #333; }
        .status-401 { background: #fd7e14; color: white; }
        .status-404 { background: #dc3545; color: white; }
        .status-500 { background: #6f42c1; color: white; }
        
        .footer {
            background: #f8f9fa;
            padding: 30px;
            text-align: center;
            color: #6c757d;
            border-top: 1px solid #e9ecef;
        }
        
        @media (max-width: 768px) {
            .header h1 {
                font-size: 1.8em;
            }
            
            .content {
                padding: 20px;
            }
            
            .endpoint-header {
                flex-direction: column;
                align-items: flex-start;
            }
            
            .method-badge {
                margin-bottom: 10px;
            }
            
            table {
                font-size: 0.9em;
            }
            
            .code-block {
                font-size: 0.8em;
            }
        }
        
        .scroll-top {
            position: fixed;
            bottom: 30px;
            right: 30px;
            background: #667eea;
            color: white;
            width: 50px;
            height: 50px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            cursor: pointer;
            box-shadow: 0 4px 12px rgba(0,0,0,0.3);
            transition: all 0.3s ease;
            text-decoration: none;
            font-size: 1.5em;
        }
        
        .scroll-top:hover {
            background: #764ba2;
            transform: translateY(-3px);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📚 API 接口文档</h1>
            <p>短链接服务完整接口说明</p>
        </div>
        
        <div class="content">
            <!-- 基础信息 -->
            <div class="info-section">
                <h2>🔧 基础信息</h2>
                <ul>
                    <li><strong>管理API地址:</strong> ` + adminBaseURL + `</li>
                    <li><strong>访问API地址:</strong> ` + accessBaseURL + `</li>
                    <li><strong>认证方式:</strong> JWT Token (Bearer Token)</li>
                    <li><strong>内容类型:</strong> application/json</li>
                    <li><strong>字符编码:</strong> UTF-8</li>
                </ul>
            </div>
            
            <!-- 管理API -->
            <div class="api-section">
                <h2>管理API接口</h2>
                
                <!-- 1. 登录 -->
                <div class="api-endpoint">
                    <div class="endpoint-header">
                        <div style="display: flex; align-items: center;">
                            <span class="method-badge method-post">POST</span>
                            <span class="endpoint-path">/api/login</span>
                        </div>
                    </div>
                    <div class="endpoint-body">
                        <div class="endpoint-description">
                            <strong>管理员登录</strong> - 登录获取访问令牌
                        </div>
                        <div class="params-section">
                            <h3>请求参数</h3>
                            <table>
                                <thead>
                                    <tr>
                                        <th>参数名</th>
                                        <th>类型</th>
                                        <th>必填</th>
                                        <th>说明</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>username</td>
                                        <td>string</td>
                                        <td>是</td>
                                        <td>管理员用户名</td>
                                    </tr>
                                    <tr>
                                        <td>password</td>
                                        <td>string</td>
                                        <td>是</td>
                                        <td>管理员密码</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                        <div class="params-section">
                            <h3>请求示例</h3>
                            <div class="code-block">
<pre>{
  <span class="json-key">"username"</span>: <span class="json-string">"admin"</span>,
  <span class="json-key">"password"</span>: <span class="json-string">"123456"</span>
}</pre>
                            </div>
                        </div>
                        <div class="response-section">
                            <h3>响应示例</h3>
                            <div>
                                <span class="status-code status-200">200 OK</span>
                            </div>
                            <div class="code-block">
<pre>{
  <span class="json-key">"token"</span>: <span class="json-string">"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."</span>,
  <span class="json-key">"username"</span>: <span class="json-string">"admin"</span>,
  <span class="json-key">"userId"</span>: <span class="json-number">1</span>
}</pre>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- 2. 创建短链接 -->
                <div class="api-endpoint">
                    <div class="endpoint-header">
                        <div style="display: flex; align-items: center;">
                            <span class="method-badge method-post">POST</span>
                            <span class="endpoint-path">/api/short-link/create</span>
                        </div>
                    </div>
                    <div class="endpoint-body">
                        <div class="endpoint-description">
                            <strong>创建短链接</strong> - 创建新的短链接
                        </div>
                        <div class="params-section">
                            <h3>请求参数</h3>
                            <table>
                                <thead>
                                    <tr>
                                        <th>参数名</th>
                                        <th>类型</th>
                                        <th>必填</th>
                                        <th>说明</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>link</td>
                                        <td>string</td>
                                        <td>是</td>
                                        <td>原始URL地址</td>
                                    </tr>
                                    <tr>
                                        <td>expire</td>
                                        <td>int</td>
                                        <td>是</td>
                                        <td>过期时间（秒）</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                        <div class="params-section">
                            <h3>请求示例</h3>
                            <div class="code-block">
<pre>{
  <span class="json-key">"link"</span>: <span class="json-string">"https://www.example.com"</span>,
  <span class="json-key">"expire"</span>: <span class="json-number">3600</span>
}</pre>
                            </div>
                        </div>
                        <div class="response-section">
                            <h3>响应示例</h3>
                            <div>
                                <span class="status-code status-200">200 OK</span>
                            </div>
                            <div class="code-block">
<pre>{
  <span class="json-key">"shortLink"</span>: <span class="json-string">"` + accessBaseURL + `s/abc123"</span>
}</pre>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- 3. 获取短链接列表 -->
                <div class="api-endpoint">
                    <div class="endpoint-header">
                        <div style="display: flex; align-items: center;">
                            <span class="method-badge method-get">GET</span>
                            <span class="endpoint-path">/api/short-link/list</span>
                            <span class="auth-badge">需要认证</span>
                        </div>
                    </div>
                    <div class="endpoint-body">
                        <div class="endpoint-description">
                            <strong>获取短链接列表</strong> - 获取有效的短链接列表，支持分页和筛选
                        </div>
                        <div class="params-section">
                            <h3>查询参数</h3>
                            <table>
                                <thead>
                                    <tr>
                                        <th>参数名</th>
                                        <th>类型</th>
                                        <th>必填</th>
                                        <th>默认值</th>
                                        <th>说明</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>page</td>
                                        <td>string</td>
                                        <td>否</td>
                                        <td>1</td>
                                        <td>页码</td>
                                    </tr>
                                    <tr>
                                        <td>pageSize</td>
                                        <td>string</td>
                                        <td>否</td>
                                        <td>10</td>
                                        <td>每页数量（最大100）</td>
                                    </tr>
                                    <tr>
                                        <td>shortCode</td>
                                        <td>string</td>
                                        <td>否</td>
                                        <td>-</td>
                                        <td>短码筛选（模糊查询）</td>
                                    </tr>
                                    <tr>
                                        <td>originalUrl</td>
                                        <td>string</td>
                                        <td>否</td>
                                        <td>-</td>
                                        <td>原始URL筛选（模糊查询）</td>
                                    </tr>
                                    <tr>
                                        <td>status</td>
                                        <td>string</td>
                                        <td>否</td>
                                        <td>-</td>
                                        <td>状态：active(有效) 或 expired(已过期)</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                        <div class="response-section">
                            <h3>响应示例</h3>
                            <div>
                                <span class="status-code status-200">200 OK</span>
                            </div>
                            <div class="code-block">
<pre>{
  <span class="json-key">"total"</span>: <span class="json-number">100</span>,
  <span class="json-key">"links"</span>: [
    {
      <span class="json-key">"id"</span>: <span class="json-number">1</span>,
      <span class="json-key">"shortCode"</span>: <span class="json-string">"abc123"</span>,
      <span class="json-key">"originalUrl"</span>: <span class="json-string">"https://www.example.com"</span>,
      <span class="json-key">"createdAt"</span>: <span class="json-string">"2024-01-01 10:00:00.000"</span>,
      <span class="json-key">"expiresAt"</span>: <span class="json-string">"2024-01-02 10:00:00.000"</span>,
      <span class="json-key">"accessCount"</span>: <span class="json-number">42</span>,
      <span class="json-key">"lastAccess"</span>: <span class="json-string">"2024-01-01 15:30:00.000"</span>
    }
  ]
}</pre>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- 4. 获取历史短链接列表 -->
                <div class="api-endpoint">
                    <div class="endpoint-header">
                        <div style="display: flex; align-items: center;">
                            <span class="method-badge method-get">GET</span>
                            <span class="endpoint-path">/api/short-link/history</span>
                            <span class="auth-badge">需要认证</span>
                        </div>
                    </div>
                    <div class="endpoint-body">
                        <div class="endpoint-description">
                            <strong>获取历史短链接列表</strong> - 获取已归档到历史表的短链接列表
                        </div>
                        <div class="params-section">
                            <h3>查询参数</h3>
                            <table>
                                <thead>
                                    <tr>
                                        <th>参数名</th>
                                        <th>类型</th>
                                        <th>必填</th>
                                        <th>默认值</th>
                                        <th>说明</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>month</td>
                                        <td>string</td>
                                        <td>否</td>
                                        <td>当前月份</td>
                                        <td>月份（格式：YYMM，如2401）</td>
                                    </tr>
                                    <tr>
                                        <td>page</td>
                                        <td>string</td>
                                        <td>否</td>
                                        <td>1</td>
                                        <td>页码</td>
                                    </tr>
                                    <tr>
                                        <td>pageSize</td>
                                        <td>string</td>
                                        <td>否</td>
                                        <td>10</td>
                                        <td>每页数量（最大100）</td>
                                    </tr>
                                    <tr>
                                        <td>shortCode</td>
                                        <td>string</td>
                                        <td>否</td>
                                        <td>-</td>
                                        <td>短码筛选（模糊查询）</td>
                                    </tr>
                                    <tr>
                                        <td>originalUrl</td>
                                        <td>string</td>
                                        <td>否</td>
                                        <td>-</td>
                                        <td>原始URL筛选（模糊查询）</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
                
                <!-- 5. 删除短链接 -->
                <div class="api-endpoint">
                    <div class="endpoint-header">
                        <div style="display: flex; align-items: center;">
                            <span class="method-badge method-delete">DELETE</span>
                            <span class="endpoint-path">/api/short-link/:id</span>
                            <span class="auth-badge">需要认证</span>
                        </div>
                    </div>
                    <div class="endpoint-body">
                        <div class="endpoint-description">
                            <strong>删除短链接</strong> - 删除指定的短链接（移动到历史表）
                        </div>
                        <div class="params-section">
                            <h3>路径参数</h3>
                            <table>
                                <thead>
                                    <tr>
                                        <th>参数名</th>
                                        <th>类型</th>
                                        <th>必填</th>
                                        <th>说明</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>id</td>
                                        <td>string</td>
                                        <td>是</td>
                                        <td>短链接ID</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                        <div class="response-section">
                            <h3>响应示例</h3>
                            <div>
                                <span class="status-code status-200">200 OK</span>
                            </div>
                            <div class="code-block">
<pre>{
  <span class="json-key">"message"</span>: <span class="json-string">"短链接已成功删除"</span>
}</pre>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- 6. 修改密码 -->
                <div class="api-endpoint">
                    <div class="endpoint-header">
                        <div style="display: flex; align-items: center;">
                            <span class="method-badge method-post">POST</span>
                            <span class="endpoint-path">/api/change-password</span>
                            <span class="auth-badge">需要认证</span>
                        </div>
                    </div>
                    <div class="endpoint-body">
                        <div class="endpoint-description">
                            <strong>修改密码</strong> - 修改当前登录用户的密码
                        </div>
                        <div class="params-section">
                            <h3>请求参数</h3>
                            <table>
                                <thead>
                                    <tr>
                                        <th>参数名</th>
                                        <th>类型</th>
                                        <th>必填</th>
                                        <th>说明</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>currentPassword</td>
                                        <td>string</td>
                                        <td>是</td>
                                        <td>当前密码</td>
                                    </tr>
                                    <tr>
                                        <td>newPassword</td>
                                        <td>string</td>
                                        <td>是</td>
                                        <td>新密码（最少6位）</td>
                                    </tr>
                                    <tr>
                                        <td>confirmPassword</td>
                                        <td>string</td>
                                        <td>是</td>
                                        <td>确认新密码</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                        <div class="params-section">
                            <h3>请求示例</h3>
                            <div class="code-block">
<pre>{
  <span class="json-key">"currentPassword"</span>: <span class="json-string">"old_password"</span>,
  <span class="json-key">"newPassword"</span>: <span class="json-string">"new_password"</span>,
  <span class="json-key">"confirmPassword"</span>: <span class="json-string">"new_password"</span>
}</pre>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            
            <!-- 访问API -->
            <div class="api-section">
                <h2>访问API接口</h2>
                
                <!-- 短链接重定向 -->
                <div class="api-endpoint">
                    <div class="endpoint-header">
                        <div style="display: flex; align-items: center;">
                            <span class="method-badge method-get">GET</span>
                            <span class="endpoint-path">/s/:code</span>
                        </div>
                    </div>
                    <div class="endpoint-body">
                        <div class="endpoint-description">
                            <strong>短链接重定向</strong> - 访问短链接时自动重定向到原始URL
                        </div>
                        <div class="params-section">
                            <h3>路径参数</h3>
                            <table>
                                <thead>
                                    <tr>
                                        <th>参数名</th>
                                        <th>类型</th>
                                        <th>必填</th>
                                        <th>说明</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>code</td>
                                        <td>string</td>
                                        <td>是</td>
                                        <td>短码</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                        <div class="response-section">
                            <h3>响应状态码</h3>
                            <div style="margin-bottom: 10px;">
                                <span class="status-code status-307">307 Temporary Redirect</span> - 成功重定向到原始URL
                            </div>
                            <div>
                                <span class="status-code status-404">404 Not Found</span> - 短链接不存在或已过期
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            
            <!-- 错误码说明 -->
            <div class="info-section">
                <h2>❌ 错误码说明</h2>
                <table>
                    <thead>
                        <tr>
                            <th>状态码</th>
                            <th>说明</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td><span class="status-code status-200">200</span></td>
                            <td>请求成功</td>
                        </tr>
                        <tr>
                            <td><span class="status-code status-307">307</span></td>
                            <td>临时重定向</td>
                        </tr>
                        <tr>
                            <td><span class="status-code status-400">400</span></td>
                            <td>请求参数错误</td>
                        </tr>
                        <tr>
                            <td><span class="status-code status-401">401</span></td>
                            <td>未认证或认证失败</td>
                        </tr>
                        <tr>
                            <td><span class="status-code status-404">404</span></td>
                            <td>资源不存在</td>
                        </tr>
                        <tr>
                            <td><span class="status-code status-500">500</span></td>
                            <td>服务器内部错误</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            
            <!-- 认证说明 -->
            <div class="info-section">
                <h2>🔐 认证说明</h2>
                <p style="margin-bottom: 15px;">
                    需要认证的接口需要在请求头中携带JWT Token：
                </p>
                <div class="code-block">
<pre>Authorization: Bearer &lt;token&gt;</pre>
                </div>
                <p style="margin-top: 15px;">
                    Token可以通过登录接口获取，默认有效期为24小时（可在配置文件中修改）。
                </p>
            </div>
        </div>
        
        <div class="footer">
            <p>© 2024 短链接服务 API 文档 | 最后更新: 2024-01-01</p>
        </div>
    </div>
    
    <a href="#" class="scroll-top" onclick="window.scrollTo({top: 0, behavior: 'smooth'}); return false;">↑</a>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}


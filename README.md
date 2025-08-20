# Go短链接服务

这是一个使用Go语言开发的短链接服务，提供链接缩短、重定向和管理功能。

## 功能特点

- 链接缩短：将长URL转换为短链接
- 链接重定向：访问短链接时自动重定向到原始URL
- 链接管理：创建、查询、更新和删除短链接
- 访问统计：记录短链接的访问次数和最后访问时间
- 过期清理：自动清理过期的短链接
- 管理后台：提供Web界面进行短链接管理

## 技术栈

### 后端
- Go语言
- Gin Web框架
- GORM数据库ORM
- JWT认证

### 前端
- React
- TypeScript
- Ant Design

## 项目结构

```
.
├── api/                # API接口定义
│   ├── access.go       # 访问相关API
│   ├── admin.go        # 管理员相关API
│   └── middleware.go   # 中间件
├── app/                # 应用程序入口
│   └── app.go
├── conf/               # 配置相关
│   ├── config.go
│   └── config.yaml
├── handlers/           # 请求处理器
│   ├── admin.go
│   └── shortlink.go
├── models/             # 数据模型
│   ├── admin.go
│   ├── gorm_store.go
│   ├── pagination.go
│   ├── response.go
│   ├── shortlink.go
│   └── store.go
├── server/             # 服务器配置
│   └── server.go
├── static/             # 静态资源
├── tasks/              # 定时任务
│   ├── clean_expired_links.go
│   └── scheduler.go
├── utils/              # 工具函数
│   ├── gorm_id_generator.go
│   ├── idgenerator.go
│   ├── jwt.go
│   └── shortcode.go
├── web/                # 前端代码
│   ├── public/
│   └── src/
├── go.mod              # Go模块依赖
├── go.sum              # Go模块校验和
├── main.go             # 主程序入口
└── build_web.sh        # 前端构建脚本
```

## 快速开始

### 前提条件

- Go 1.16+
- Node.js 14+
- MySQL 5.7+

### 安装与运行

1. 克隆仓库

```bash
git clone https://github.com/yourusername/go-short-link.git
cd go-short-link
```

2. 配置数据库

编辑 `conf/config.yaml` 文件，设置数据库连接信息。

3. 构建前端

```bash
cd web
yarn install
yarn build
cd ..
```

或者使用提供的脚本：

```bash
./build_web.sh
```

4. 运行后端服务

```bash
go run main.go
```

5. 访问服务

- 短链接服务: http://localhost:8082/s/
- 管理后台: http://localhost:8081

## API文档

### 短链接API

- `POST /api/short-link/create` - 创建新的短链接
- `GET /api/short-link/list` - 获取短链接列表
- `GET /api/short-link/history` - 获取历史短链接列表
- `DELETE /api/short-link/:id` - 删除短链接

### 管理员API

- `POST /api/login` - 管理员登录
- `POST /api/change-password` - 修改密码

## 配置说明

配置文件位于 `conf/config.yaml`，主要配置项包括：

- 服务器配置（端口、地址等）
- 数据库配置（连接信息、表前缀等）
- JWT配置（密钥、过期时间等）
- 短链接配置（默认过期时间、短码长度等）

## 许可证

本项目采用 [LICENSE](LICENSE) 许可证。
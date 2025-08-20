# Go短链接服务

一个基于Go语言开发的高性能短链接服务，支持创建短链接和重定向功能，采用MySQL持久化存储和内存LRU缓存。

## 项目概述

本项目是一个完整的短链接服务，具有以下特点：

- 分离的管理API和访问API，分别运行在不同端口
- 基于MySQL的持久化存储
- 基于LRU算法的内存缓存，支持热点数据快速访问
- IP白名单保护管理API
- 支持设置链接过期时间
- 模块化设计，易于扩展和维护

## 系统架构

### 整体架构

系统采用分层架构设计：

1. **API层**：处理HTTP请求和路由
2. **业务逻辑层**：处理短链接的创建和重定向逻辑
3. **数据访问层**：管理数据的存储和检索
4. **缓存层**：提供高性能数据访问

### 存储设计

系统采用混合存储策略：

- **MySQL数据库**：持久化存储所有短链接数据
- **内存LRU缓存**：存储最多1000条热点短链接数据，提高访问性能
- **缓存淘汰策略**：当缓存达到容量上限时，自动淘汰最不常用的短链接
- **Redis ID生成器**：使用Redis生成分布式唯一ID，支持水平扩展

### 服务分离

系统将管理API和访问API分离：

- **管理API**：处理短链接的创建，受IP白名单保护
- **访问API**：处理短链接的重定向，面向所有用户

## 项目结构

```
go-short-link/
├── api/                  # API路由定义
│   ├── access.go         # 访问API路由
│   └── admin.go          # 管理API路由和IP白名单中间件
├── app/                  # 应用程序初始化
│   └── app.go            # 应用程序配置和资源管理
├── conf/                 # 配置管理
│   ├── config.go         # 配置加载和解析
│   └── config.yaml       # 配置文件
├── handlers/             # 请求处理器
│   └── shortlink.go      # 短链接处理逻辑
├── models/               # 数据模型和存储
│   ├── shortlink.go      # 短链接数据模型
│   └── store.go          # 存储接口和实现（MySQL+缓存）
├── server/               # 服务器管理
│   └── server.go         # 服务器初始化和生命周期管理
├── utils/                # 工具函数
│   ├── shortcode.go      # 短链接生成工具
│   ├── idgenerator.go    # Redis ID生成器（旧版）
│   └── gorm_id_generator.go # GORM Redis ID生成器插件
├── main.go               # 应用程序入口
├── go.mod                # Go模块定义
├── go.sum                # 依赖校验和
└── test.sh               # 测试脚本
```

## 核心功能

### 1. 创建短链接

```
POST /short-link/create
Content-Type: application/json

{
  "link": "https://www.example.com",
  "expire": 3600
}
```

响应:

```json
{
  "shortLink": "http://localhost:8081/s/a1b2c3d4"
}
```

### 2. 访问短链接

```
GET /s/{shortCode}
```

系统会将请求重定向到原始URL。

## 设计思路

### 模块化设计

项目采用模块化设计，每个模块负责特定功能：

- **API模块**：负责HTTP路由和请求处理
- **处理器模块**：实现业务逻辑
- **模型模块**：定义数据结构和存储接口
- **配置模块**：管理应用程序配置
- **服务器模块**：管理HTTP服务器生命周期
- **工具模块**：提供通用功能

这种设计使代码更易于维护和扩展。

### 混合存储策略

系统采用MySQL+内存缓存的混合存储策略：

1. 所有短链接数据都存储在MySQL数据库中
2. 最多1000条热点短链接数据缓存在内存中
3. 使用LRU算法管理缓存，自动淘汰最不常用的数据
4. 异步更新访问计数，减少数据库压力

这种策略在保证数据持久性的同时，提供了高性能的数据访问。

### 服务分离

将管理API和访问API分离到不同端口，有以下优势：

1. 提高安全性，管理API可以受到更严格的保护
2. 便于独立扩展，可以根据不同API的负载特点进行扩展
3. 便于维护，可以独立更新和部署

### IP白名单保护

管理API受IP白名单保护，只有白名单中的IP地址才能访问管理API，提高了系统安全性。

## 部署指南

### 环境要求

- Go 1.16+
- MySQL 5.7+

### 配置文件

部署前需要修改`conf/config.yaml`文件，配置服务器和数据库信息：

```yaml
# 服务配置
server:
  # 管理API服务配置（创建短链接）
  admin:
    port: 8081
    baseURL: "http://your-domain.com:8081/"
    # IP白名单，允许访问管理API的IP列表
    ipWhitelist:
      - "127.0.0.1"
      - "::1"
      - "192.168.1.0/24"
  
  # 访问API服务配置（短链接重定向）
  access:
    port: 8082
    baseURL: "http://your-domain.com:8082/"

# 数据库配置
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "your_username"
  password: "your_password"
  dbname: "gsl"
  charset: "utf8mb4"
  parseTime: true
  loc: "Local"
  maxIdleConns: 10
  maxOpenConns: 100
  connMaxLifetime: 3600

# Redis配置
redis:
  addr: "localhost:6379"
  password: ""
  db: 0
  poolSize: 10
  idKeyPrefix: "seq:"  # ID生成器的键前缀
  idStep: 100  # 每次从Redis获取的ID数量

# 缓存配置
cache:
  type: "memory"
  capacity: 1000
```

### 构建和运行

#### 本地开发环境

1. 克隆代码库

```bash
git clone https://github.com/yourusername/go-short-link.git
cd go-short-link
```

2. 安装依赖

```bash
go mod download
```

3. 构建项目

```bash
go build -o short-link
```

4. 运行服务

```bash
./short-link
```

#### 使用Docker部署

1. 创建Dockerfile

```dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o short-link

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/short-link /app/
COPY --from=builder /app/conf /app/conf

EXPOSE 8081 8082
CMD ["./short-link"]
```

2. 构建Docker镜像

```bash
docker build -t go-short-link .
```

3. 运行Docker容器

```bash
docker run -d -p 8081:8081 -p 8082:8082 --name short-link go-short-link
```

#### 使用Docker Compose部署

1. 创建docker-compose.yml文件

```yaml
version: '3'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: gsl
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "3306:3306"
    networks:
      - short-link-network

  short-link:
    build: .
    depends_on:
      - mysql
    ports:
      - "8081:8081"
      - "8082:8082"
    networks:
      - short-link-network
    restart: always

networks:
  short-link-network:

volumes:
  mysql_data:
```

2. 启动服务

```bash
docker-compose up -d
```

### 生产环境部署建议

1. 使用Nginx作为反向代理，处理SSL终止和负载均衡
2. 配置适当的日志轮转
3. 设置监控和告警
4. 定期备份数据库
5. 使用容器编排工具（如Kubernetes）进行部署和管理

## 分布式ID生成

系统使用Redis实现了分布式ID生成器，具有以下特点：

1. **基于GORM插件**：作为GORM插件集成，自动为新记录生成ID
2. **通用设计**：使用`seq:{表名}`作为键前缀，支持多个表的ID生成
3. **批量获取**：每次从Redis获取一批ID（默认100个），减少网络请求
4. **高性能**：本地缓存ID段，大幅减少Redis访问频率
5. **容错机制**：内置重试逻辑，提高系统稳定性

### ID生成器工作原理

1. 初始化时，ID生成器注册为GORM的回调函数
2. 创建记录时，检查主键是否为零值
3. 如果是零值，则为该表生成一个新的唯一ID
4. ID生成器会批量从Redis获取ID段，并在本地缓存
5. 当本地缓存的ID用完后，再次从Redis获取新的ID段

这种设计既保证了ID的唯一性，又提高了系统性能，同时支持水平扩展。

## 性能优化

系统已经实现了一些性能优化措施：

1. 使用LRU缓存减少数据库访问
2. 异步更新访问计数
3. 使用连接池管理数据库连接
4. 批量获取Redis ID，减少网络请求

对于高负载场景，可以考虑以下优化：

1. 增加缓存容量
2. 使用Redis替代内存缓存
3. 数据库读写分离
4. 水平扩展访问API服务

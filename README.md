# 栈序管理平台

栈序管理平台是一个前后端分离的企业级中后台管理系统，提供用户、角色、菜单、部门、登录日志和操作日志等基础能力。

后端基于 Go、Gin、GORM、Casbin 和 Redis，前端基于 Vue 3、Vite、Ant Design Vue 和 Pinia。

## 功能

- JWT 登录、多端会话、刷新令牌轮换和退出登录
- Redis 验证码、分布式请求限流、动态菜单和按钮权限缓存
- 基于 Casbin 的 RBAC 接口权限控制
- 动态菜单和前端动态路由
- 用户、角色、菜单和部门管理
- 按钮级权限控制
- 登录日志和操作日志
- 操作日志请求体敏感字段脱敏
- 请求限流、访问日志、异常恢复和定时任务基础设施
- 管理员账户删除和角色变更保护

## 技术栈

| 模块 | 技术 |
|------|------|
| 后端 | Go 1.24、Gin、GORM、Casbin、JWT、Viper、Zap |
| 前端 | Vue 3、Vite、Ant Design Vue、Pinia、Vue Router、Axios |
| 数据库 | MySQL 8 |
| 缓存 | Redis 7 |

## 项目结构

```text
.
├── backend/
│   ├── cmd/server/          # 后端程序入口
│   ├── config/              # 配置和 Casbin 模型
│   ├── internal/
│   │   ├── bootstrap/       # DB、Redis、Casbin、Cron 初始化
│   │   ├── dto/             # 请求和响应结构体
│   │   ├── handler/v1/      # HTTP 处理层
│   │   ├── middleware/      # Gin 中间件
│   │   ├── model/           # GORM 模型
│   │   ├── repository/      # 数据访问层
│   │   ├── router/          # 路由注册
│   │   └── service/         # 业务逻辑层
│   ├── migrations/          # 初始化数据 SQL
│   └── pkg/                 # 公共工具包
├── frontend/
│   └── src/
│       ├── api/             # API 封装
│       ├── layout/          # 页面布局
│       ├── router/          # 静态路由和路由守卫
│       ├── stores/          # Pinia 状态管理
│       ├── utils/           # 请求和认证工具
│       └── views/           # 页面
└── CLAUDE.md                # AI 协作开发规范
```

## 环境要求

- Go `1.24.13`
- Node.js `20+`
- npm
- MySQL `8`
- Redis `7`
- Docker，可选，用于快速启动 MySQL 和 Redis

默认服务地址：

| 服务 | 地址 |
|------|------|
| 前端 | `http://localhost:3000` |
| 后端 | `http://localhost:8080` |
| MySQL | `127.0.0.1:3309` |
| Redis | `127.0.0.1:7379` |

## 快速开始

### 1. 启动基础服务

使用 Docker 启动 MySQL 和 Redis：

```bash
docker run -d \
  --name zhanxu-admin-mysql \
  -e MYSQL_ROOT_PASSWORD=your-password \
  -p 3309:3306 \
  mysql:8

docker run -d \
  --name zhanxu-admin-redis \
  -p 7379:6379 \
  redis:7
```

创建数据库：

```bash
docker exec zhanxu-admin-mysql mysql -uroot -pyour-password \
  -e "CREATE DATABASE IF NOT EXISTS zhanxu_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 2. 配置后端

```bash
cp backend/config/config.example.yaml backend/config/config.yaml
```

编辑 `backend/config/config.yaml`，至少替换以下配置：

- `server.jwt_secret`
- `database.password`
- `storage.access_key`
- `storage.secret_key`

`backend/config/config.yaml` 包含本地敏感配置，已被 Git 忽略。

### 3. 首次建表

后端启动时会通过 GORM AutoMigrate 自动创建表：

```bash
cd backend
go run ./cmd/server
```

看到服务启动成功后，使用 `Ctrl+C` 停止后端。

### 4. 导入初始化数据

初始化数据包含管理员账户、角色、菜单、按钮权限和 Casbin 策略，仅应在空数据库中执行一次：

```bash
docker exec -i zhanxu-admin-mysql mysql -uroot -pyour-password zhanxu_admin \
  < backend/migrations/init_data.sql
```

### 5. 启动后端

```bash
cd backend
go run ./cmd/server
```

后端启动后监听 `http://localhost:8080`，接口前缀为 `/api/v1`。

### 6. 启动前端

打开新的终端：

```bash
cd frontend
npm install
npm run dev
```

访问 `http://localhost:3000`。

## 默认账户

| 用户名 | 密码 | 角色 |
|--------|------|------|
| `admin` | `Admin@123456` | 超级管理员 |

首次登录后应立即修改默认密码。

## 常用命令

后端测试、静态检查和构建：

```bash
cd backend
go test ./...
go vet ./...
go build ./cmd/server
```

前端开发和构建：

```bash
cd frontend
npm run dev
npm run build
npm run preview
```

## 配置说明

后端配置文件为 `backend/config/config.yaml`，参考模板为：

```text
backend/config/config.example.yaml
```

前端环境配置：

| 文件 | 用途 |
|------|------|
| `frontend/.env` | 通用配置 |
| `frontend/.env.development` | 开发环境配置 |
| `frontend/.env.production` | 生产环境配置 |

开发环境下，Vite 会将 `/api` 请求代理到 `http://localhost:8080`。

### Redis 业务配置

Redis 配置项：

| 配置项 | 默认值 | 用途 |
|--------|--------|------|
| `key_prefix` | `zhanxu-admin:` | Redis Key 项目前缀 |
| `cache_ttl` | `900` | 菜单和权限缓存有效期，单位秒 |
| `captcha_expire` | `300` | 验证码有效期，单位秒 |
| `dial_timeout` | `5` | 连接超时，单位秒 |
| `read_timeout` | `3` | 读取超时，单位秒 |
| `write_timeout` | `3` | 写入超时，单位秒 |

`rate_limit` 中的 `login_*`、`captcha_*` 和 `sensitive_*` 分别控制登录、验证码及密码修改接口的独立 Redis 令牌桶。`rate` 表示每秒恢复的令牌数，`burst` 表示可突发请求数。

Redis 当前承担以下业务职责：

- 保存按会话隔离的 Refresh Token，并在刷新时原子轮换。
- 保存已退出 Access Token 的 JWT ID 黑名单。
- 保存验证码，支持多实例部署。
- 缓存用户动态菜单与按钮权限，角色或菜单变更后主动失效。
- 使用 Redis 令牌桶实现跨实例请求限流。

认证安全数据访问 Redis 失败时请求会失败；菜单和权限缓存访问失败时会自动回源 MySQL。

升级到使用会话 ID 和 JWT ID 的版本后，旧版本签发的 Token 将失效，用户需要重新登录。

## 权限模型

- 接口权限由 Casbin 校验。
- 用户通过角色获得权限。
- 菜单类型分为目录、菜单和按钮。
- 前端菜单和路由根据当前用户权限动态生成。
- 按钮权限通过 `v-permission` 指令控制。
- `admin` 角色拥有 `/api/v1/*` 的全部权限。

## 开发约束

开发前请阅读 [CLAUDE.md](CLAUDE.md)。主要约束包括：

- Handler 仅处理参数绑定、调用 Service 和返回响应。
- Service 负责业务逻辑，不直接操作数据库。
- Repository 统一负责 GORM 数据访问。
- Model 与外部 DTO 严格隔离。
- 前端接口统一通过 `src/api/` 和 `src/utils/request.js` 调用。
- 新增业务页面通过数据库菜单动态注册路由。

## 当前限制

- 文件上传和删除接口尚未接入对象存储。
- 项目暂未提供 Docker Compose 和生产部署配置。
- `init_data.sql` 不是幂等脚本，请勿对已有初始化数据的数据库重复执行。

## 安全提示

- 不要提交 `backend/config/config.yaml` 或本地 `.env` 文件。
- 生产环境必须替换 JWT 密钥、数据库密码和存储凭据。
- 生产环境应关闭 Gin debug 模式，并通过 HTTPS 对外提供服务。

## 许可证

本项目由 zhang zhenming 创建，基于 [Apache License 2.0](LICENSE) 开源。

Copyright 2026 zhang zhenming

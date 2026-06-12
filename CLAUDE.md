# 栈序管理平台 — AI 开发规范（CLAUDE.md）

> 本文件是项目 AI 协作开发的强制约束，开发前必须完整阅读。

---

## 一、项目概览

企业级中后台管理系统，前后端分离架构。

| 层级 | 路径 | 说明 |
|------|------|------|
| 后端 | `backend/` | Go 服务，监听 :8080 |
| 前端 | `frontend/` | Vue3 SPA，开发时监听 :3000，通过 Vite proxy 转发 /api |
| 数据库 | Docker MySQL | host: 127.0.0.1, port: 3309, db: zhanxu_admin |
| 缓存 | Docker Redis | host: 127.0.0.1, port: 7379, 无密码 |

---

## 二、后端技术栈与版本

```
Go              1.24.13
Gin             v1.10.0       HTTP 框架
GORM            v1.26.0       ORM（软删除用 gorm.DeletedAt）
MySQL Driver    v1.5.7
Casbin          v2.99.0       RBAC 权限
gorm-adapter    v3.27.0       Casbin 策略存储
JWT             golang-jwt/jwt v5.2.1
Redis           go-redis/v9 v9.7.0
Viper           v1.19.0       配置管理
Zap             v1.27.0       日志
lumberjack      v2.2.1        日志轮转（包名 gopkg.in/natefinch/lumberjack.v2）
base64Captcha   v1.3.6        验证码
robfig/cron     v3.0.1        定时任务（WithSeconds 模式）
golang.org/x/time/rate        限流（令牌桶）
```

---

## 三、后端目录结构与分层规范

```
backend/
├── cmd/server/main.go          程序入口，负责依赖注入和启动
├── config/
│   ├── config.go               Config 结构体定义
│   ├── config.yaml             配置文件
│   └── rbac_model.conf         Casbin RBAC 模型
├── internal/
│   ├── bootstrap/              基础设施初始化（DB/Redis/Casbin/Cron）
│   ├── handler/v1/             HTTP 处理层：参数绑定、调用 Service、返回响应
│   ├── service/                业务逻辑层：不直接操作 DB，只调用 Repository
│   ├── repository/             数据访问层：所有 GORM 操作收口于此
│   ├── model/                  GORM 模型，对应数据库表
│   ├── dto/                    请求/响应结构体，隔离 model 与外部
│   ├── middleware/             Gin 中间件
│   └── router/                 路由注册
├── pkg/                        可复用公共库（无业务依赖）
│   ├── cache/                  Redis 封装
│   ├── crypto/                 bcrypt 加密
│   ├── jwtutil/                JWT 生成解析
│   ├── logger/                 Zap 封装
│   ├── pagination/             分页工具
│   └── response/               统一响应封装
└── migrations/                 SQL 文件
```

### 分层约束

- **Handler** 只做：参数绑定（ShouldBindJSON/ShouldBindQuery/ShouldBindUri）、调用 Service、调用 response 包返回结果。不写任何业务逻辑。
- **Service** 只做：业务逻辑、调用 Repository、处理 BizError。不直接使用 gorm.DB。
- **Repository** 只做：GORM 操作。不写业务判断。
- **Model** 只做：GORM 结构体定义。不写任何方法逻辑。
- **DTO** 与 **Model** 严格隔离：Model 禁止直接序列化给前端，必须转换为 DTO 的 Resp 结构体。

---

## 四、后端编码规范

### 4.1 统一响应

所有接口必须使用 `pkg/response` 包返回，禁止直接 `c.JSON`。

```go
// 成功
response.Success(c, data)

// 分页成功
response.SuccessPage(c, list, total, page, pageSize)

// 业务失败（使用预定义错误码）
response.Fail(c, response.CodeUserNotFound)

// 带自定义消息的失败
response.FailWithMsg(c, response.CodeBadRequest, "xxx")

// 服务器错误
response.ServerError(c)
```

### 4.2 错误码

业务错误码定义在 `pkg/response/response.go`，新增模块按段分配：

```
200       成功
400       参数错误
401       未登录
403       无权限
404       资源不存在
429       限流
500       服务器错误
1000x     用户模块
1001x     角色模块
1002x     菜单模块
1003x     部门模块
1004x     认证模块
1005x     下一个新模块从此段开始
```

新增错误码必须同时在 `msgMap` 中注册中文描述。

### 4.3 业务错误处理

Service 层通过 `BizError` 返回业务错误，Handler 层统一用 `handleBizError` 处理：

```go
// service 层
return &BizError{Code: response.CodeUserNotFound}

// handler 层（复用已有的 handleBizError 函数，定义在 user.go 中）
if err := svc.DoSomething(); err != nil {
    handleBizError(c, err)
    return
}
```

### 4.4 路由规范

- 所有接口前缀 `/api/v1`
- 公开接口（无需登录）注册在 `public` 路由组
- 需要登录不需要权限校验的接口注册在 `authed` 路由组
- 需要权限校验的接口注册在 `permRoutes` 路由组
- 新增路由必须在 `internal/router/router.go` 中注册
- RESTful 风格：GET 查询、POST 新增、PUT 修改、DELETE 删除

### 4.5 数据库规范

- 所有业务表必须嵌入 `model.Base`（含 ID/CreatedAt/UpdatedAt/DeletedAt 软删除）
- 日志表不使用软删除（直接物理删除）
- 表名统一 `sys_` 前缀（日志表除外）
- 新增字段通过 GORM AutoMigrate 自动迁移，重大变更写 migration SQL 文件
- 禁止在 Service/Handler 中直接使用 `bootstrap.DB`，必须通过 Repository

### 4.6 Casbin 权限

- 接口权限规则：`p, 角色code, /api/v1/路径, HTTP方法`
- 用户角色绑定：`g, 用户ID字符串, 角色code`
- admin 角色使用通配规则 `p, admin, /api/v1/*, *`
- 修改角色菜单后必须调用 `enforcer.LoadPolicy()` 重新加载

### 4.7 定时任务

新增定时任务在 `internal/bootstrap/cron.go` 中注册，任务实现放在 `internal/cron/tasks/` 下。

---

## 五、前端技术栈与版本

```
Vue             3.5.34
Vite            8.x
Ant Design Vue  4.2.6
Pinia           3.0.4
Vue Router      4.6.4
Axios           1.17.0
dayjs           1.11.21
@ant-design/icons-vue  7.0.1
```

---

## 六、前端目录结构规范

```
frontend/src/
├── api/                每个模块一个文件，只封装接口调用
├── assets/styles/      global.css 定义 CSS 变量，禁止在组件内重复定义颜色/间距
├── components/         全局公共组件（跨页面复用才放这里）
├── directives/         自定义指令（permission.js 权限控制）
├── layout/             整体框架布局，非业务组件
├── router/
│   ├── index.js        静态路由（login/dashboard/profile/404）
│   └── guards.js       路由守卫，动态路由在此注册
├── stores/             Pinia store（auth/user/menu/app）
├── utils/
│   ├── auth.js         token 读写（localStorage）
│   └── request.js      Axios 实例，含 token 自动刷新
└── views/              页面组件，路径对应路由
    ├── login/
    ├── dashboard/
    ├── profile/
    ├── system/         系统管理（user/role/menu/dept）
    └── log/            日志管理（operation/login）
```

---

## 七、前端编码规范

### 7.1 API 层

```js
// 统一格式，文件路径 src/api/xxx.js
import request from '@/utils/request'

export const xxxApi = {
  list: (params) => request.get('/xxx', { params }),
  get: (id) => request.get(`/xxx/${id}`),
  create: (data) => request.post('/xxx', data),
  update: (id, data) => request.put(`/xxx/${id}`, data),
  delete: (id) => request.delete(`/xxx/${id}`),
}
```

request.js 响应拦截器已处理：code=200 直接返回 data，其他 code 自动 message.error 并 reject。业务代码中 **不需要再判断 code**，直接 try/catch 即可。

### 7.2 页面组件结构

每个管理页面遵循固定结构：

```
搜索卡片（a-card.search-card）
  └── a-form layout="inline"
        └── 查询条件 + [查询] [重置] 按钮

数据卡片（a-card style="margin-top:12px"）
  ├── 工具栏 div.toolbar（左：操作按钮，右：刷新图标）
  └── a-table（size="middle", row-key="id"）

弹窗（a-modal）
  └── a-form（label-col span:6, wrapper-col span:16）
```

### 7.3 动态路由规则

菜单类型说明（对应数据库 sys_menu.type）：
- `1` 目录：有子菜单，`component = 'Layout'`，注册为带 Layout 的顶级路由
- `2` 菜单：实际页面，`component` 为相对于 `src/views/` 的路径（如 `system/user/index`）
- `3` 按钮：不参与路由，仅用于权限标识

新增页面必须：
1. 在 `src/views/` 对应路径下创建 `index.vue`
2. 在数据库 `sys_menu` 中插入对应菜单记录（type=2，component 填写正确路径）
3. 在 `sys_role_menu` 中为相关角色分配该菜单
4. **不需要**修改 `router/index.js`（动态路由自动注册）

静态路由（不走菜单树的页面）才需要修改 `router/index.js`，目前静态路由：`/login`、`/dashboard`、`/profile`、`404`。

### 7.4 状态管理

```
useAuthStore    token 管理、登录/登出
useUserStore    用户信息、权限标识列表
useMenuStore    菜单树、动态路由注册状态
useAppStore     侧边栏折叠、多标签页
```

跨组件共享的状态放 store，组件内部临时状态用 `ref/reactive`，禁止在 store 里存放列表数据（列表数据属于页面局部状态）。

### 7.5 权限控制

按钮级权限使用 `v-permission` 指令：

```html
<a-button v-permission="'system:user:add'" type="primary">新增</a-button>
```

权限标识来源于 `sys_menu.permission` 字段，格式：`模块:子模块:操作`。

### 7.6 样式规范

- 颜色、间距必须使用 `global.css` 中的 CSS 变量，禁止硬编码颜色值
- 主色 `var(--primary-color)` = `#1677ff`
- 侧边栏背景 `var(--sider-bg)` = `#001529`（深色，禁止改为浅色）
- 页面内容背景 `var(--bg-color)` = `#f5f6fa`
- 组件圆角统一 `var(--border-radius)` = `6px`
- 禁止使用内联 style 设置颜色，间距除对齐用途外也避免内联

---

## 八、新增业务模块标准流程

### 后端

1. `internal/model/` 新增 Model 结构体，嵌入 `model.Base`
2. `internal/dto/` 新增 ListReq、CreateReq、UpdateReq、Resp 结构体
3. `internal/repository/` 新增 Repo，实现 CRUD 方法
4. `internal/service/` 新增 Service，调用 Repo，处理业务逻辑
5. `internal/handler/v1/` 新增 Handler，绑定参数调用 Service
6. `internal/handler/v1/handler.go` 注入新 Handler
7. `internal/router/router.go` 注册新路由
8. `cmd/server/main.go` 初始化新 Repo 和 Service
9. `pkg/response/response.go` 新增对应模块错误码段

### 前端

1. `src/api/` 新增接口文件
2. `src/views/` 新增页面组件
3. 数据库插入菜单数据（`sys_menu` + `sys_role_menu`）

---

## 九、禁止事项

### 后端禁止

- 禁止在 Handler/Service 直接 `import bootstrap` 使用 `bootstrap.DB`
- 禁止 Model 直接序列化返回给前端，必须通过 DTO Resp 转换
- 禁止跳过 `response` 包直接 `c.JSON`
- 禁止新增全局变量（配置、DB、Redis 等已通过 bootstrap 初始化并注入）
- 禁止修改 `pkg/` 下已有公共接口签名（会影响全局调用）
- 新增依赖必须通过 `go get` 指定精确版本，禁止使用 `latest`

### 前端禁止

- 禁止在页面组件中直接使用 `axios`，必须通过 `src/utils/request.js`
- 禁止在组件中手动处理 token（读写统一走 `src/utils/auth.js`）
- 禁止修改 `layout/` 下的布局组件结构（Header/Sidebar/TabsView）
- 禁止在 `stores/` 外直接修改 store state（必须通过 action）
- 禁止硬编码接口路径字符串，统一走 `src/api/` 层
- 禁止引入 Element Plus、Vuetify 等其他 UI 库，统一使用 Ant Design Vue

---

## 十、关键文件速查

| 文件 | 用途 |
|------|------|
| `backend/config/config.yaml` | 数据库/Redis/JWT 等配置 |
| `backend/config/rbac_model.conf` | Casbin RBAC 模型定义 |
| `backend/pkg/response/response.go` | 统一响应 + 全部错误码 |
| `backend/internal/router/router.go` | 全部路由注册 |
| `backend/internal/handler/v1/handler.go` | Handler 依赖注入聚合 |
| `backend/cmd/server/main.go` | 服务启动 + 依赖注入 |
| `frontend/src/utils/request.js` | Axios 实例 + token 刷新逻辑 |
| `frontend/src/stores/menu.js` | 动态路由生成逻辑 |
| `frontend/src/assets/styles/global.css` | 全局 CSS 变量 |
| `frontend/src/router/index.js` | 静态路由（仅4个） |
| `backend/migrations/init_data.sql` | 初始化数据（admin账号等） |

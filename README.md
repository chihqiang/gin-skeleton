# gin-skeleton

🛠 A skeleton of Golang gin framework with admin panel, JWT authentication, and queue processing.

## 特性

- 基于 Gin 框架构建
- 支持 JWT 身份验证
- 集成 Casbin 权限控制
- 支持 SQLite 和 MySQL 数据库
- 内置后台管理模块
- 支持 Redis 队列任务处理
- 结构化日志记录

## 项目结构

```
├── app/                    # 应用核心代码
│   ├── admin/              # 后台管理模块
│   ├── apis/               # API 接口
│   ├── middlewares/        # 中间件
│   ├── models/             # 数据模型
│   └── tasks/              # 队列任务
├── bootstrap/              # 应用启动引导
├── cmd/                    # 命令行工具
├── config/                 # 配置加载
├── docs/                   # 文档
├── pkg/                    # 公共包
│   ├── casbinx/            # Casbin 权限控制
│   ├── cryptox/            # 加密工具
│   ├── database/           # 数据库连接
│   ├── helper/             # 辅助函数
│   ├── httpx/              # HTTP 工具
│   ├── jwts/               # JWT 处理
│   ├── logger/             # 日志处理
│   ├── queue/              # 队列处理
│   └── redisx/             # Redis 连接
├── routes/                 # 路由定义
└── runtime/                # 运行时文件（日志、数据库等）
```

## 安装

1. 下载项目并编译:
   ```sh
   git clone https://github.com/chihqiang/gin-skeleton.git
   cd gin-skeleton
   go mod tidy
   make build
   ```
 2. 运行

```bash
# 运行 HTTP 服务
./skeleton http
# 运行队列任务
./skeleton queue:start
```

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

FootballPaul 是一个基于 Go 的足球比赛预测与积分系统，用户预测比赛比分并根据准确程度获得积分。

**技术栈**: Go 1.21+, Gin Web 框架, GORM ORM, SQLite, JWT 认证, React + Vite + TailwindCSS。

## 常用命令

```bash
# 安装依赖
go mod download
go mod tidy
cd frontend && npm install

# 开发运行（后端使用 air 热重载）
air
cd frontend && npm run dev

# 生产构建
go build -o bin/footballpaul main.go
cd frontend && npm run build

# 运行测试
go test ./...

# 运行测试（带覆盖率）
go test -cover ./...

# 运行特定包的测试
go test ./utils
go test ./models
go test ./middleware
go test ./handlers
go test -v ./...

# 运行测试（覆盖单个文件）
go test -v ./utils/scoring_test.go
```

## 项目结构

```
footballPaul/
├── main.go                    # 程序入口
├── router/router.go           # 所有路由定义
├── handlers/                  # HTTP 请求处理层
│   ├── user_handler.go
│   ├── match_handler.go
│   ├── prediction_handler.go
│   ├── leaderboard_handler.go
│   └── competition_handler.go
├── services/                  # 业务逻辑层
│   ├── user_service.go
│   ├── match_service.go
│   ├── prediction_service.go
│   ├── leaderboard_service.go
│   └── competition_service.go
├── models/                    # 数据模型层
│   ├── user.go
│   ├── match.go
│   ├── prediction.go
│   ├── competition.go
│   └── scoring_rule.go
├── middleware/                 # 中间件
│   ├── auth.go                # JWT 认证
│   └── cors.go
├── config/                     # 配置管理
│   └── config.go
├── database/                   # 数据库连接和迁移
│   └── database.go
├── utils/                      # 工具函数
│   └── scoring.go             # 积分计算逻辑
└── frontend/                   # React 前端
    ├── src/
    │   ├── pages/             # 页面组件
    │   ├── components/        # 可复用组件
    │   ├── context/           # React Context
    │   └── api/               # API 客户端
    └── package.json
```

## 架构

**分层结构**: Handlers → Services → Models。Handler 接收 HTTP 请求，委托给 Services 处理业务逻辑，通过 Models 与数据交互。

**创建预测的核心流程**:
1. `handlers/prediction_handler.go` 接收 POST `/api/predictions`
2. `services/prediction_service.go` 验证截止时间并检查是否已预测
3. `models/prediction.go` 存储预测记录
4. 当录入比赛结果时（PUT `/api/matches/:id/result`），`scoring_service.go` 遍历所有预测，调用 `utils/scoring.go` 计算积分

**数据库**: SQLite（通过 `config/config.go` 配置 `DB_DSN`）。通过 `database/database.go` 在启动时自动迁移。自动插入 4 条默认积分规则和 8 个预置赛事。

**认证**: 基于 JWT（`middleware/auth.go`）。受保护的路由使用 `middleware.AuthMiddleware(cfg)`。

**API 认证路由**:
- `POST /api/auth/register` - 注册用户（需提供 username, email, password）
- `POST /api/auth/login` - 登录获取 JWT token（需提供 email, password）

**配置**: 通过 `config/config.go` 读取环境变量 - DB_DSN（数据库路径）, SERVER_PORT, JWT_SECRET。

## 测试

本项目使用 Go 标准 `testing` 包，配合 `testify/assert` 做断言。

**测试原则**:
- `utils/` - 纯函数，优先写单元测试（如 `scoring.go`）
- `models/` - 测试模型方法和序列化（如 `ToResponse()`）
- `middleware/` - 使用 `httptest` 测试 gin 中间件行为
- `handlers/` - 测试路由结构和请求/响应处理
- `services/` - 数据库依赖较重，可通过 mock 或 integration test 覆盖
- `config/` - 测试配置加载和环境变量覆盖

**测试文件命名**: `*_test.go`，与被测文件同包同目录。

## 积分规则

| 预测准确度 | 积分 |
|-----------|------|
| 完全正确（比分一致） | 10 |
| 猜中胜负+净胜球 | 7 |
| 只猜中胜负 | 5 |
| 猜中一方得分 | 3 |
| 预测错误 | 0 |

实现位于 `utils/scoring.go` 中的 `CalculatePoints(predictedHome, predictedAway, actualHome, actualAway int) int`。

## 关键文件

- `main.go` - 程序入口，初始化数据库并启动服务
- `router/router.go` - 所有路由定义
- `database/database.go` - 数据库连接、迁移、数据初始化
- `utils/scoring.go` - 积分计算逻辑
- `.env.example` - 所需环境变量示例

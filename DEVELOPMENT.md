# FootballPaul 开发指南

## 项目结构

```
footballPaul/
├── config/              # 配置管理
│   └── config.go        # 配置加载和环境变量
├── database/            # 数据库层
│   └── database.go      # 数据库连接和迁移
├── handlers/            # HTTP处理器（控制器）
│   ├── match_handler.go
│   ├── prediction_handler.go
│   ├── leaderboard_handler.go
│   └── user_handler.go
├── middleware/          # 中间件
│   ├── auth.go          # JWT认证中间件
│   └── cors.go          # CORS中间件
├── models/              # 数据模型
│   ├── user.go
│   ├── match.go
│   ├── prediction.go
│   └── scoring_rule.go
├── router/              # 路由配置
│   └── router.go
├── services/            # 业务逻辑层
│   ├── user_service.go
│   ├── match_service.go
│   ├── prediction_service.go
│   ├── scoring_service.go
│   └── leaderboard_service.go
├── utils/               # 工具函数
│   └── scoring.go       # 积分计算逻辑
├── main.go              # 应用入口
├── go.mod               # Go模块依赖
├── .env.example         # 环境变量示例
├── .gitignore           # Git忽略文件
├── README.md            # 项目说明
├── API.md               # API文档
└── DEVELOPMENT.md       # 开发指南（本文件）
```

## 环境准备

### 1. 安装 Go

需要 Go 1.21 或更高版本：

```bash
# 检查 Go 版本
go version

# 如果未安装，请访问：https://golang.org/dl/
```

### 2. 安装 PostgreSQL

```bash
# macOS
brew install postgresql@14
brew services start postgresql@14

# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql

# 创建数据库
createdb footballpaul
```

### 3. 克隆并配置项目

```bash
# 克隆项目
git clone https://github.com/liyufei916/footballPaul.git
cd footballPaul

# 复制环境变量文件
cp .env.example .env

# 编辑 .env 文件，配置数据库连接
nano .env
```

### 4. 安装依赖

```bash
go mod download
go mod tidy
```

## 运行项目

### 开发模式

```bash
# 直接运行
go run main.go

# 使用 air 进行热重载（推荐）
# 首先安装 air
go install github.com/cosmtrek/air@latest

# 运行
air
```

### 生产模式

```bash
# 编译
go build -o bin/footballpaul main.go

# 运行
./bin/footballpaul
```

## 数据库操作

### 自动迁移

应用启动时会自动执行数据库迁移，创建所有必要的表。

### 手动迁移

如果需要手动控制迁移过程：

```go
// 在 main.go 中
database.AutoMigrate()
```

### 初始化数据

应用启动时会自动插入默认的积分规则。

## 开发工作流

### 1. 添加新功能

#### 步骤1：定义数据模型（如需要）

在 `models/` 目录下创建或修改模型：

```go
// models/new_model.go
package models

type NewModel struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"not null" json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### 步骤2：实现业务逻辑

在 `services/` 目录下创建服务：

```go
// services/new_service.go
package services

type NewService struct{}

func NewNewService() *NewService {
    return &NewService{}
}

func (s *NewService) DoSomething() error {
    // 业务逻辑
    return nil
}
```

#### 步骤3：创建HTTP处理器

在 `handlers/` 目录下创建处理器：

```go
// handlers/new_handler.go
package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/liyufei916/footballPaul/services"
)

type NewHandler struct {
    service *services.NewService
}

func NewNewHandler() *NewHandler {
    return &NewHandler{
        service: services.NewNewService(),
    }
}

func (h *NewHandler) HandleRequest(c *gin.Context) {
    // 处理请求
    c.JSON(200, gin.H{"status": "ok"})
}
```

#### 步骤4：注册路由

在 `router/router.go` 中添加路由：

```go
newHandler := handlers.NewNewHandler()
api.GET("/new-endpoint", newHandler.HandleRequest)
```

### 2. 代码规范

#### 命名规范

- **文件名**: 使用小写和下划线，如 `user_service.go`
- **包名**: 使用小写，如 `package services`
- **类型名**: 使用大驼峰，如 `UserService`
- **函数名**: 导出函数使用大驼峰，私有函数使用小驼峰
- **变量名**: 使用小驼峰，如 `userID`

#### 错误处理

始终检查并处理错误：

```go
result, err := someFunction()
if err != nil {
    return err
}
```

#### 日志记录

使用标准库的 `log` 包：

```go
log.Println("Info message")
log.Printf("Error: %v", err)
```

### 3. 测试

#### 创建测试文件

```go
// services/user_service_test.go
package services

import "testing"

func TestCreateUser(t *testing.T) {
    // 测试逻辑
}
```

#### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./services

# 显示详细输出
go test -v ./...

# 测试覆盖率
go test -cover ./...
```

## API 测试

### 使用 curl

```bash
# 注册用户
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"123456"}'

# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"123456"}'
```

### 使用 Postman

1. 导入 API 文档创建集合
2. 设置环境变量 `baseUrl` 为 `http://localhost:8080/api`
3. 设置认证 token 为环境变量 `token`

## 常见问题

### 1. 数据库连接失败

**问题**: `failed to connect to database`

**解决**:
- 检查 PostgreSQL 是否运行
- 验证 `.env` 文件中的数据库配置
- 确保数据库已创建

### 2. 端口被占用

**问题**: `bind: address already in use`

**解决**:
```bash
# 查找占用端口的进程
lsof -i :8080

# 杀死进程
kill -9 <PID>
```

### 3. 依赖问题

**问题**: 包导入错误

**解决**:
```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download
go mod tidy
```

## 性能优化

### 1. 数据库查询优化

- 使用索引
- 避免 N+1 查询问题
- 使用预加载（Preload）

```go
// 不好的做法
for _, prediction := range predictions {
    database.DB.Model(&prediction).Association("Match").Find(&prediction.Match)
}

// 好的做法
database.DB.Preload("Match").Find(&predictions)
```

### 2. 缓存

对于频繁访问的数据（如排行榜），可以使用 Redis 缓存：

```go
// 伪代码
if cached := redis.Get("leaderboard"); cached != nil {
    return cached
}
data := database.GetLeaderboard()
redis.Set("leaderboard", data, 5*time.Minute)
return data
```

### 3. 并发处理

使用 goroutine 处理耗时操作：

```go
go func() {
    // 异步任务
    sendNotification(user)
}()
```

## 部署

### Docker 部署

创建 `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o footballpaul main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/footballpaul .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./footballpaul"]
```

构建和运行：

```bash
docker build -t footballpaul .
docker run -p 8080:8080 --env-file .env footballpaul
```

### 使用 Docker Compose

创建 `docker-compose.yml`:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=footballpaul
    depends_on:
      - postgres

  postgres:
    image: postgres:14-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=footballpaul
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  postgres_data:
```

运行：

```bash
docker-compose up -d
```

## 贡献指南

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/amazing-feature`
3. 提交更改：`git commit -m 'Add amazing feature'`
4. 推送到分支：`git push origin feature/amazing-feature`
5. 创建 Pull Request

## 许可证

MIT License

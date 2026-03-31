# FootballPaul API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api`
- **认证方式**: JWT Bearer Token
- **内容类型**: `application/json`

---

## 认证相关

### 用户注册

**POST** `/auth/register`

注册新用户。

**请求体**:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "123456"
}
```

**响应** (201):
```json
{
  "success": true,
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "total_points": 0,
    "is_admin": false,
    "created_at": "2024-10-24T00:00:00Z"
  }
}
```

**错误响应**:
- `400` - 参数错误 / 邮箱已被注册 / 用户名已被使用

---

### 用户登录

**POST** `/auth/login`

**请求体**:
```json
{
  "email": "test@example.com",
  "password": "123456"
}
```

**响应** (200):
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "total_points": 0,
    "is_admin": false,
    "created_at": "2024-10-24T00:00:00Z"
  }
}
```

**错误响应**:
- `401` - 邮箱或密码错误

---

## 用户

### 获取个人资料 🔒

**GET** `/users/profile`

获取当前登录用户的资料。

**请求头**:
```
Authorization: Bearer <token>
```

**响应** (200):
```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "total_points": 150,
  "is_admin": false,
  "created_at": "2024-10-24T00:00:00Z"
}
```

---

## 赛事管理

### 获取赛事列表

**GET** `/competitions`

**响应** (200):
```json
{
  "competitions": [
    {
      "id": 1,
      "name": "英超联赛",
      "code": "EPL",
      "logo": "",
      "match_count": 0,
      "created_at": "2024-10-24T00:00:00Z"
    },
    {
      "id": 2,
      "name": "欧冠联赛",
      "code": "CHAMPIONS_LEAGUE",
      "logo": "",
      "match_count": 0,
      "created_at": "2024-10-24T00:00:00Z"
    }
  ],
  "count": 8
}
```

---

### 获取单个赛事

**GET** `/competitions/:id`

**响应** (200):
```json
{
  "id": 1,
  "name": "英超联赛",
  "code": "EPL",
  "logo": "",
  "match_count": 0,
  "created_at": "2024-10-24T00:00:00Z"
}
```

**错误响应**:
- `400` - 无效的赛事 ID
- `404` - 赛事不存在

---

## 比赛管理

### 获取比赛列表

**GET** `/matches`

**查询参数**:
- `status` (可选): `pending` | `ongoing` | `finished`
- `competition_id` (可选): 按赛事筛选
- `limit` (可选): 限制返回数量，默认 10

**响应** (200):
```json
{
  "matches": [
    {
      "id": 1,
      "competition_id": 1,
      "home_team": "曼联",
      "away_team": "利物浦",
      "match_date": "2024-10-25T15:00:00Z",
      "home_score": null,
      "away_score": null,
      "status": "pending",
      "deadline": "2024-10-25T14:45:00Z",
      "created_at": "2024-10-24T00:00:00Z",
      "competition": {
        "id": 1,
        "name": "英超联赛",
        "code": "EPL"
      }
    }
  ],
  "count": 1
}
```

---

### 获取单个比赛

**GET** `/matches/:id`

**响应** (200):
```json
{
  "id": 1,
  "competition_id": 1,
  "home_team": "曼联",
  "away_team": "利物浦",
  "match_date": "2024-10-25T15:00:00Z",
  "home_score": null,
  "away_score": null,
  "status": "pending",
  "deadline": "2024-10-25T14:45:00Z",
  "created_at": "2024-10-24T00:00:00Z",
  "competition": {
    "id": 1,
    "name": "英超联赛",
    "code": "EPL"
  }
}
```

**错误响应**:
- `400` - 无效的比赛 ID
- `404` - 比赛不存在

---

### 创建比赛 🔒

**POST** `/matches`

创建新比赛（管理员功能）。

**请求头**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "competition_id": 1,
  "home_team": "曼联",
  "away_team": "利物浦",
  "match_date": "2024-10-25T15:00:00Z",
  "deadline": "2024-10-25T14:45:00Z"
}
```

**响应** (201):
```json
{
  "success": true,
  "match": {
    "id": 1,
    "competition_id": 1,
    "home_team": "曼联",
    "away_team": "利物浦",
    "match_date": "2024-10-25T15:00:00Z",
    "status": "pending",
    "deadline": "2024-10-25T14:45:00Z",
    "created_at": "2024-10-24T00:00:00Z"
  }
}
```

---

### 录入比赛结果 🔒

**PUT** `/matches/:id/result`

录入比赛结果并自动评分。

**请求头**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "home_score": 3,
  "away_score": 1
}
```

**响应** (200):
```json
{
  "success": true,
  "message": "比分已录入，评分完成"
}
```

**说明**: 此操作会自动触发该比赛所有预测的评分，并更新用户积分。

**错误响应**:
- `400` - 比分不能为负数
- `400` - 比赛已结束

---

### 获取比赛的所有预测 🔒

**GET** `/matches/:matchId/predictions`

获取指定比赛的所有用户预测。

**请求头**:
```
Authorization: Bearer <token>
```

**响应** (200):
```json
{
  "predictions": [
    {
      "id": 1,
      "user_id": 1,
      "match_id": 1,
      "predicted_home_score": 2,
      "predicted_away_score": 1,
      "points_earned": 10,
      "is_scored": true,
      "predicted_at": "2024-10-24T10:00:00Z"
    }
  ],
  "count": 1
}
```

---

## 预测管理

### 提交预测 🔒

**POST** `/predictions`

**请求头**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "match_id": 1,
  "predicted_home_score": 2,
  "predicted_away_score": 1
}
```

**响应** (201):
```json
{
  "success": true,
  "prediction_id": 1,
  "message": "预测提交成功"
}
```

**错误响应**:
- `400` - 预测截止时间已过
- `400` - 已经为该比赛提交过预测

---

### 更新预测 🔒

**PUT** `/predictions/:id`

更新已提交的预测（仅在截止时间前且未评分时可更新）。

**请求头**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "match_id": 1,
  "predicted_home_score": 3,
  "predicted_away_score": 2
}
```

**响应** (200):
```json
{
  "success": true,
  "prediction": {
    "id": 1,
    "user_id": 1,
    "match_id": 1,
    "predicted_home_score": 3,
    "predicted_away_score": 2,
    "points_earned": 0,
    "is_scored": false,
    "predicted_at": "2024-10-24T10:00:00Z"
  },
  "message": "预测更新成功"
}
```

**错误响应**:
- `400` - 预测截止时间已过
- `400` - 无法更新已评分的预测

---

### 获取我的预测 🔒

**GET** `/predictions/my`

获取当前用户的预测历史。

**请求头**:
```
Authorization: Bearer <token>
```

**响应** (200):
```json
{
  "predictions": [
    {
      "id": 1,
      "user_id": 1,
      "match_id": 1,
      "predicted_home_score": 2,
      "predicted_away_score": 1,
      "points_earned": 10,
      "is_scored": true,
      "predicted_at": "2024-10-24T10:00:00Z",
      "match": {
        "id": 1,
        "home_team": "曼联",
        "away_team": "利物浦",
        "match_date": "2024-10-25T15:00:00Z",
        "home_score": 2,
        "away_score": 1,
        "status": "finished"
      }
    }
  ],
  "count": 1
}
```

---

## 排行榜

### 获取排行榜

**GET** `/leaderboard`

**查询参数**:
- `limit` (可选): 限制返回数量，默认 10

**响应** (200):
```json
{
  "rankings": [
    {
      "rank": 1,
      "user_id": 1,
      "username": "player1",
      "total_points": 850,
      "predictions_count": 50
    },
    {
      "rank": 2,
      "user_id": 2,
      "username": "player2",
      "total_points": 720,
      "predictions_count": 48
    }
  ]
}
```

---

### 获取我的排名 🔒

**GET** `/leaderboard/my-rank`

获取当前登录用户的排名。

**请求头**:
```
Authorization: Bearer <token>
```

**响应** (200):
```json
{
  "rank": 5
}
```

---

## 健康检查

### 服务健康状态

**GET** `/health`

**响应** (200):
```json
{
  "status": "ok"
}
```

---

## 错误响应格式

所有错误响应都遵循以下格式：

```json
{
  "error": "错误信息描述"
}
```

常见 HTTP 状态码：
- `200` - 成功
- `201` - 创建成功
- `400` - 请求参数错误
- `401` - 未认证或认证失败
- `404` - 资源不存在
- `500` - 服务器内部错误

---

## 积分规则

| 预测准确度 | 积分 | 说明 |
|---------|------|------|
| 完全正确 | 10分 | 比分完全一致 |
| 猜中胜负+净胜球 | 7分 | 结果和净胜球都正确 |
| 猜中胜负 | 5分 | 只猜中胜/平/负 |
| 猜中一方得分 | 3分 | 猜中任一队伍得分 |
| 其他 | 0分 | 预测不准确 |

---

## 预置赛事

系统初始化时会创建以下赛事：

| 名称 | 代码 |
|------|------|
| 2026世界杯 | 2026WORLD_CUP |

---

🔒 = 需要认证（需在请求头中携带 `Authorization: Bearer <token>`）

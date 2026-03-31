# 组队功能设计文档

## 1. 概述

### 1.1 功能目标
用户可以创建竞猜组、加入已有组，一个用户可属于多个组。组内成员可查看组成员列表，以及各赛事在组内的独立排行榜。

### 1.2 核心场景
- 用户 A 创建「欧冠竞猜群」，获得邀请码 `ABC123`
- 用户 B、C 收到邀请码，加入「欧冠竞猜群」
- 群内开启「2024-25 欧冠」赛事追踪
- 群成员竞猜欧冠比赛，按积分生成组内排行榜
- 各成员可以看到自己在组内的排名，以及与其他群友的差距

---

## 2. 数据模型

### 2.1 新增表

#### `groups`
| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | pk, auto | |
| name | string | not null | 组名，最长50字符 |
| invite_code | string | unique, not null | 6位邀请码（大写字母+数字） |
| owner_id | uint | fk → users, not null | 创建者 |
| created_at | timestamp | | |
| updated_at | timestamp | | |

#### `group_members`
| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | pk, auto | |
| group_id | uint | fk → groups, index | |
| user_id | uint | fk → users, index | |
| role | string | default:'member' | `admin` 或 `member`，owner 在此表中也是 admin |
| joined_at | timestamp | | |

**唯一约束**：`UNIQUE(group_id, user_id)` — 同一用户不能重复加入同一组。

#### `group_competitions`
| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | pk, auto | |
| group_id | uint | fk → groups, index | |
| competition_id | uint | fk → competitions, index | |
| created_at | timestamp | | |

**唯一约束**：`UNIQUE(group_id, competition_id)` — 同一赛事在同一组内不能重复添加。

### 2.2 现有表变更

**`users`** — 无变更。

**`competitions`** — 无变更（已经存在联赛数据）。

**`predictions`** — 无变更（已有 user_id, match_id）。

---

## 3. API 设计

### 3.1 认证要求
- 带 `AuthMiddleware` 的接口：需要 Bearer Token
- 公开接口：无需认证

### 3.2 接口列表

#### 创建组
```
POST /api/groups
Auth: required

Request Body:
{
  "name": "欧冠竞猜群"      // required, 1-50 chars
}

Response 201:
{
  "group": {
    "id": 1,
    "name": "欧冠竞猜群",
    "invite_code": "ABC123",
    "owner_id": 42,
    "created_at": "2026-03-31T..."
  },
  "message": "组队创建成功"
}
```

#### 获取我的所有组
```
GET /api/groups
Auth: required

Response 200:
{
  "groups": [
    {
      "id": 1,
      "name": "欧冠竞猜群",
      "role": "admin",
      "member_count": 3,
      "competition_count": 1,
      "joined_at": "2026-03-31T..."
    }
  ]
}
```

#### 获取组详情
```
GET /api/groups/:id
Auth: required (必须是组成员)

Response 200:
{
  "group": {
    "id": 1,
    "name": "欧冠竞猜群",
    "invite_code": "ABC123",
    "owner_id": 42,
    "created_at": "2026-03-31T..."
  }
}
```

#### 加入组（凭邀请码）
```
POST /api/groups/join
Auth: required

Request Body:
{
  "invite_code": "ABC123"   // required, 6 chars
}

Response 200:
{
  "message": "加入成功",
  "group": { ... }
}

Response 400:
{
  "error": "邀请码无效或已失效"
}

Response 409:
{
  "error": "你已经在该组中"
}
```

#### 离开组
```
DELETE /api/groups/:id/leave
Auth: required (必须是组成员，owner 不可退组)

Response 200:
{
  "message": "已离开该组"
}

Response 400:
{
  "error": "组长无法离开，请先转让组长或解散该组"
}
```

#### 解散组（仅 owner）
```
DELETE /api/groups/:id
Auth: required (必须是 owner)

Response 200:
{
  "message": "组队已解散"
}
```

#### 获取组成员列表
```
GET /api/groups/:id/members
Auth: required (必须是组成员)

Response 200:
{
  "members": [
    {
      "user_id": 42,
      "username": "liyufei",
      "role": "admin",
      "joined_at": "2026-03-31T..."
    },
    {
      "user_id": 43,
      "username": "alice",
      "role": "member",
      "joined_at": "2026-03-31T..."
    }
  ]
}
```

#### 获取组内追踪的赛事列表
```
GET /api/groups/:id/competitions
Auth: required (必须是组成员)

Response 200:
{
  "competitions": [
    {
      "id": 1,
      "name": "UEFA Champions League 2024/25",
      "logo_url": "https://...",
      "added_at": "2026-03-31T..."
    }
  ]
}
```

#### 在组内新增追踪赛事
```
POST /api/groups/:id/competitions
Auth: required (必须是 admin)

Request Body:
{
  "competition_id": 1
}

Response 201:
{
  "message": "赛事已添加到本组"
}
```

#### 从组内移除追踪赛事
```
DELETE /api/groups/:id/competitions/:competitionId
Auth: required (必须是 admin)

Response 200:
{
  "message": "赛事已从本组移除"
}
```

#### 获取赛事在组内的排行榜
```
GET /api/groups/:id/leaderboard/:competitionId
Auth: required (必须是组成员)

Query Params:
  ?limit=50   // 默认50，最多100

Response 200:
{
  "competition": {
    "id": 1,
    "name": "UEFA Champions League 2024/25"
  },
  "leaderboard": [
    {
      "rank": 1,
      "user_id": 42,
      "username": "liyufei",
      "total_points": 85,
      "predictions_count": 12,
      "exact_scores": 3,
      "correct_winners": 5
    },
    {
      "rank": 2,
      "user_id": 43,
      "username": "alice",
      "total_points": 72,
      "predictions_count": 12,
      "exact_scores": 2,
      "correct_winners": 4
    }
  ]
}
```

#### 转让组长
```
PUT /api/groups/:id/transfer-owner
Auth: required (必须是 owner)

Request Body:
{
  "new_owner_id": 43
}

Response 200:
{
  "message": "组长已转让"
}
```

---

## 4. 核心业务逻辑

### 4.1 邀请码生成规则
- 6位，由大写字母 A-Z 和数字 0-9 组成
- 不含易混淆字符（0/O、1/I/L）
- 全局唯一，不重复
- 邀请码永不过期

### 4.2 组内排行榜计算逻辑

```
SQL 伪代码：

SELECT 
  u.id AS user_id,
  u.username,
  COALESCE(SUM(p.points_earned), 0) AS total_points,
  COUNT(p.id) AS predictions_count,
  COUNT(CASE WHEN p.points_earned = 10 THEN 1 END) AS exact_scores,
  COUNT(CASE WHEN p.points_earned IN (5, 7) THEN 1 END) AS correct_winners
FROM users u
INNER JOIN group_members gm ON gm.user_id = u.id AND gm.group_id = ?
INNER JOIN predictions p ON p.user_id = u.id
INNER JOIN matches m ON m.id = p.match_id AND m.competition_id = ?
WHERE p.is_scored = true
GROUP BY u.id
ORDER BY total_points DESC, exact_scores DESC
LIMIT ?
```

### 4.3 权限矩阵

| 操作 | admin | member | owner |
|------|:-----:|:------:|:-----:|
| 查看组成员 | ✅ | ✅ | ✅ |
| 查看组排行榜 | ✅ | ✅ | ✅ |
| 查看/追踪赛事 | ✅ | ✅ | ✅ |
| 加入组（凭邀请码） | - | - | - |
| 离开组 | ✅ | ✅ | ❌ |
| 添加/移除追踪赛事 | ✅ | ❌ | ✅ |
| 解散组 | ❌ | ❌ | ✅ |
| 转让组长 | ❌ | ❌ | ✅ |

---

## 5. 前端页面规划

| 页面 | 路由 | 说明 |
|------|------|------|
| 我的组列表 | `/groups` | 展示我加入的所有组，卡片形式 |
| 创建组 | `/groups/new` | 表单：输入组名 |
| 组详情 | `/groups/:id` | 组信息 + 成员数 + 追踪赛事入口 |
| 组成员 | `/groups/:id/members` | 成员列表，含角色标识 |
| 组排行榜 | `/groups/:id/leaderboard/:competitionId` | 选择赛事后展示排行 |
| 加入组 | `/groups/join` | 输入邀请码加入 |

---

## 6. 实现顺序

### Phase 1：数据层
1. 新增 `models/group.go` — Group, GroupMember, GroupCompetition
2. 新增 `database/group.go` — 表结构迁移
3. 更新 `database/database.go` — 注册迁移

### Phase 2：Service 层
4. 新增 `services/group_service.go` — 组创建/加入/退出的核心逻辑
5. 新增 `services/group_leaderboard_service.go` — 组内排行榜计算

### Phase 3：Handler 层
6. 新增 `handlers/group_handler.go` — 所有 `/api/groups` 路由处理
7. 注册路由到 `router/router.go`

### Phase 4：前端
8. React 页面组件（参考第5节）
9. API 客户端方法

### Phase 5：完善
10. 更新 CLAUDE.md
11. 补充单元测试

---

## 7. 错误处理

| HTTP Code | 场景 |
|-----------|------|
| 400 | 请求参数缺失或非法 |
| 401 | 未认证或 Token 无效 |
| 403 | 无权限操作（如非管理员添加赛事） |
| 404 | 组/赛事不存在 |
| 409 | 重复加入、重复添加赛事等冲突 |

---

## 8. 扩展考虑（Future）

- 组内聊天功能
- 组内预测讨论帖
- 公开组（无需邀请码可加入）
- 组间排行榜（跨组排名）

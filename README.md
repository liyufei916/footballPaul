# ⚽ FootballPaul - 足球比赛竞猜积分系统

一个功能完善的足球比赛结果竞猜与积分管理系统，让用户可以预测比赛结果并通过准确的预测获得积分。

## 📋 项目概述

FootballPaul 是一个交互式的足球比赛预测平台，用户可以：
- 在比赛开始前提交对比赛结果的预测
- 根据预测的准确度获得不同等级的积分
- 在排行榜上与其他用户竞争
- 查看个人的预测历史和统计数据

## ✨ 核心功能

### 1. 用户竞猜
- 在比赛开始前提交预测结果（比分）
- 支持修改预测（截止时间前）
- 查看个人预测历史
- 实时显示竞猜截止倒计时

### 2. 比赛管理
- 管理员创建和管理比赛信息
- 录入实际比赛结果
- 自动触发评分流程
- 比赛状态追踪（待开始/进行中/已结束）

### 3. 智能评分系统
- 多层次积分规则
- 自动计算和分配积分
- 实时更新用户总积分
- 积分历史记录

### 4. 排行榜系统
- 全局积分排行
- 按时间段筛选（周/月/赛季）
- 个人统计数据展示
- 实时排名更新

## 🎯 积分规则

系统采用分层积分制度，根据预测的准确程度给予不同积分：

| 预测准确度 | 积分 | 说明 | 示例 |
|---------|------|------|------|
| 完全正确 | 10分 | 比分完全一致 | 预测2:1，实际2:1 |
| 猜中胜负+净胜球 | 7分 | 结果和净胜球都正确 | 预测2:0，实际3:1（净胜2球） |
| 猜中胜负 | 5分 | 只猜中胜/平/负 | 预测2:1，实际3:2（都是主胜） |
| 猜中一方得分 | 3分 | 猜中任一队伍得分 | 预测2:1，实际2:3（主队得2分） |
| 其他 | 0分 | 预测不准确 | - |

### 评分逻辑示例

```python
def calculate_points(predicted_home, predicted_away, actual_home, actual_away):
    # 完全正确：10分
    if predicted_home == actual_home and predicted_away == actual_away:
        return 10
    
    # 计算胜负关系
    predicted_result = get_match_result(predicted_home, predicted_away)
    actual_result = get_match_result(actual_home, actual_away)
    
    # 预测结果错误
    if predicted_result != actual_result:
        # 检查是否猜中某一方得分
        if predicted_home == actual_home or predicted_away == actual_away:
            return 3
        return 0
    
    # 预测结果正确，检查净胜球
    predicted_diff = abs(predicted_home - predicted_away)
    actual_diff = abs(actual_home - actual_away)
    
    # 净胜球也正确：7分
    if predicted_diff == actual_diff:
        return 7
    
    # 只猜中胜负：5分
    return 5
```

## 🗄️ 数据模型

### 用户表 (Users)
```
id              唯一标识
username        用户名
email           邮箱
total_points    总积分
created_at      注册时间
```

### 比赛表 (Matches)
```
id              唯一标识
home_team       主队名称
away_team       客队名称
match_date      比赛时间
home_score      主队实际得分（赛后填写）
away_score      客队实际得分（赛后填写）
status          状态（pending/ongoing/finished）
deadline        竞猜截止时间
created_at      创建时间
```

### 预测表 (Predictions)
```
id                      唯一标识
user_id                 用户ID（外键）
match_id                比赛ID（外键）
predicted_home_score    预测主队得分
predicted_away_score    预测客队得分
points_earned           获得积分（评分后填写）
predicted_at            预测时间
is_scored               是否已评分
```

### 积分规则表 (ScoringRules)
```
id              唯一标识
rule_name       规则名称
rule_type       规则类型
points          对应积分
description     规则描述
```

## 🏗️ 系统架构

### 推荐技术栈

#### 后端选项

**Option 1: Python + Flask/Django**
```
footballPaul/
├── app/
│   ├── models/              # 数据模型
│   │   ├── user.py
│   │   ├── match.py
│   │   └── prediction.py
│   ├── services/            # 业务逻辑
│   │   ├── match_service.py
│   │   ├── prediction_service.py
│   │   └── scoring_service.py
│   ├── api/                 # API接口
│   │   ├── user_api.py
│   │   ├── match_api.py
│   │   └── prediction_api.py
│   └── utils/               # 工具函数
│       └── scoring_calculator.py
├── tests/                   # 测试
├── config.py                # 配置
└── requirements.txt         # 依赖
```

**Option 2: Node.js + Express**
```
footballPaul/
├── src/
│   ├── models/              # 数据模型
│   ├── controllers/         # 控制器
│   ├── services/            # 业务逻辑
│   ├── routes/              # 路由
│   └── utils/               # 工具
├── tests/
└── package.json
```

### 核心模块

#### 1. 比赛管理模块
- 创建和更新比赛信息
- 管理比赛状态
- 录入实际比分
- 触发自动评分

#### 2. 预测管理模块
- 提交预测（验证截止时间）
- 修改预测（截止前允许）
- 查询预测历史

#### 3. 评分引擎模块
- 自动计算积分
- 更新用户总积分
- 记录积分历史

#### 4. 排行榜模块
- 全局排行榜
- 按时间段统计
- 个人数据分析

## 📡 API 接口

### 用户预测接口

**提交预测**
```http
POST /api/predictions
Content-Type: application/json

{
  "match_id": 123,
  "predicted_home_score": 2,
  "predicted_away_score": 1
}

Response:
{
  "success": true,
  "prediction_id": 456,
  "message": "预测提交成功"
}
```

### 比赛管理接口

**获取比赛列表**
```http
GET /api/matches?status=pending&limit=10

Response:
{
  "matches": [
    {
      "id": 123,
      "home_team": "曼联",
      "away_team": "利物浦",
      "match_date": "2024-10-25T15:00:00Z",
      "deadline": "2024-10-25T14:45:00Z",
      "status": "pending"
    }
  ]
}
```

**录入比赛结果（管理员）**
```http
PUT /api/matches/{match_id}/result
Content-Type: application/json

{
  "home_score": 3,
  "away_score": 1
}

Response:
{
  "success": true,
  "message": "比分已录入，评分完成",
  "scored_predictions": 25
}
```

### 排行榜接口

**获取排行榜**
```http
GET /api/leaderboard?limit=10&period=monthly

Response:
{
  "rankings": [
    {
      "rank": 1,
      "user_id": 1,
      "username": "player1",
      "total_points": 850,
      "predictions_count": 50,
      "accuracy_rate": 68.5
    }
  ]
}
```

## 🚀 快速开始

### 环境要求

- Python 3.8+ / Node.js 14+ （根据选择的技术栈）
- PostgreSQL 12+ / MySQL 8+ / MongoDB 4+
- Redis（可选，用于缓存）

### 安装步骤

#### Python + Flask 示例

1. 克隆仓库
```bash
git clone https://github.com/liyufei916/footballPaul.git
cd footballPaul
```

2. 创建虚拟环境
```bash
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate
```

3. 安装依赖
```bash
pip install -r requirements.txt
```

4. 配置数据库
```bash
# 复制配置文件
cp config.example.py config.py

# 编辑 config.py 设置数据库连接信息
# 初始化数据库
flask db init
flask db migrate
flask db upgrade
```

5. 运行应用
```bash
flask run
```

访问 `http://localhost:5000` 查看应用。

#### Node.js + Express 示例

1. 克隆仓库
```bash
git clone https://github.com/liyufei916/footballPaul.git
cd footballPaul
```

2. 安装依赖
```bash
npm install
```

3. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件设置数据库连接等配置
```

4. 初始化数据库
```bash
npm run db:migrate
```

5. 运行应用
```bash
npm run dev
```

访问 `http://localhost:3000` 查看应用。

## 🔐 关键技术要点

### 数据完整性
- 预测提交后不可删除（截止前可修改）
- 比赛结果确认后不可修改（需要特殊权限）
- 使用数据库事务确保积分计算准确性

### 性能优化
- 使用 Redis 缓存排行榜数据
- 批量评分处理
- 数据库索引优化（user_id, match_id, match_date）

### 安全考虑
- 预测截止时间严格校验
- 管理员权限控制
- 防止重复提交
- API 请求限流

### 扩展性
- 积分规则可配置化
- 支持多种比赛类型
- 预留奖励机制（连胜加成等）

## 📊 业务流程

### 用户提交预测流程
```
1. 用户选择比赛
2. 系统验证：
   - 比赛是否存在
   - 是否已过截止时间
   - 用户是否已预测
3. 保存预测记录
4. 返回成功确认
```

### 比赛结束评分流程
```
1. 管理员录入实际比分
2. 比赛状态更新为 'finished'
3. 触发评分引擎：
   a. 查询该比赛所有预测
   b. 逐条计算积分
   c. 更新预测记录的积分
   d. 累加到用户总积分
4. 发送通知给用户（可选）
```

## 🎨 前端设计建议

### 页面结构
1. **首页** - 显示即将开始的比赛列表
2. **比赛详情页** - 提交/查看预测
3. **我的预测** - 历史预测记录和得分
4. **排行榜** - 积分排名
5. **个人中心** - 统计数据和成就

### 关键交互
- 实时倒计时显示截止时间
- 预测提交后的确认动画
- 比赛结束后的得分展示
- 排行榜的实时更新

## 💡 未来功能规划

### 第一阶段（MVP）
- [x] 基础数据模型
- [x] 用户预测功能
- [x] 简单积分计算
- [x] 基础排行榜

### 第二阶段（功能增强）
- [ ] 复杂积分规则
- [ ] 预测统计分析
- [ ] 通知系统
- [ ] 社交功能（评论、分享）

### 第三阶段（高级功能）
- [ ] 成就系统
- [ ] 虚拟货币/奖品
- [ ] 预测推荐算法
- [ ] 数据可视化

### 额外功能建议
1. **连胜奖励** - 连续猜对N场获得额外积分
2. **难度系数** - 冷门比赛预测准确获得更高积分
3. **预测分析** - 展示所有用户的预测分布
4. **好友对战** - 与好友比较预测能力
5. **赛季总结** - 生成个人预测报告

## 🧪 测试

### 运行测试
```bash
# Python
pytest

# Node.js
npm test
```

### 测试覆盖率
```bash
# Python
pytest --cov=app tests/

# Node.js
npm run test:coverage
```

## 📝 开发指南

### 代码规范
- 遵循 PEP 8 (Python) / ESLint (JavaScript) 规范
- 使用有意义的变量和函数名
- 添加必要的注释和文档字符串
- 编写单元测试

### Git 工作流
1. 从 `develop` 分支创建功能分支
2. 完成开发并通过测试
3. 提交 Pull Request
4. 代码审查通过后合并

### 提交信息格式
```
<type>(<scope>): <subject>

类型：
- feat: 新功能
- fix: 修复
- docs: 文档
- style: 格式
- refactor: 重构
- test: 测试
- chore: 构建/工具
```

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 👥 作者

- [@liyufei916](https://github.com/liyufei916)

## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

## 📞 联系方式

如有问题或建议，请：
- 提交 [Issue](https://github.com/liyufei916/footballPaul/issues)
- 发送邮件至：[您的邮箱]

---

**享受预测足球比赛的乐趣！⚽🎉**

# Release v0.1.0 - Initial Framework Release

**发布日期**: 2024-08-13
**版本类型**: 初始框架版本

---

## 🎉 首个框架版本发布

这是License Manager项目的第一个发布版本，包含了基础的系统框架和核心功能。该版本主要展示了完整的技术架构和基础功能实现，为后续功能开发奠定了坚实的基础。

## ✨ 新功能

### 🔐 认证系统
- JWT身份认证机制
- 安全的登录/登出功能
- 会话管理和状态保持
- 认证中间件保护

### 🌐 国际化支持
- 多语言框架实现
- 支持中文、英文、日文
- 前后端统一的国际化方案
- 动态语言切换

### 📊 管理界面
- 现代化的响应式设计
- 基于Element Plus的UI组件
- 仪表盘数据展示
- 侧边栏导航系统
- 深色模式支持

### 🛠️ 系统架构
- 完整的前后端分离架构
- RESTful API设计规范
- 结构化错误处理
- 配置管理系统
- 日志记录机制

### 🐳 部署支持
- Docker容器化部署
- Docker Compose编排
- 开发/生产环境配置分离
- 健康检查机制

## 🏗️ 技术栈

### 前端技术
- **框架**: Vue 3.4+ with Composition API
- **语言**: TypeScript 5.0+
- **UI组件**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **国际化**: Vue I18n 9
- **构建工具**: Vite 5.0+
- **样式**: SCSS

### 后端技术
- **语言**: Go 1.23+
- **框架**: Gin Web Framework
- **ORM**: GORM
- **数据库**: MySQL 8.0+
- **配置**: Viper
- **日志**: Logrus
- **认证**: JWT

### 基础设施
- **容器化**: Docker & Docker Compose
- **反向代理**: Nginx
- **API文档**: Swagger
- **数据库迁移**: 自定义迁移系统

## 📖 快速开始

### 使用Docker部署（推荐）

```bash
# 克隆项目
git clone https://github.com/cedar-v/license-manager.git
cd license-manager

# 开发环境部署
docker compose up -d --build

# 访问系统
echo "Frontend: http://localhost:18080"
echo "Backend: http://localhost:18888"
echo "Health Check: http://localhost:18888/health"
```

### 本地开发

```bash
# 后端启动
cd backend/cmd
cp ../configs/config.example.yaml ./config.yaml
# 编辑config.yaml配置数据库连接
go build -o license-manager
./license-manager

# 前端启动（新终端）
cd frontend
npm install
npm run dev
```

## 🗂️ 项目结构

```
license-manager/
├── backend/                 # Go后端服务
│   ├── cmd/                # 应用入口
│   ├── internal/           # 内部模块
│   ├── pkg/                # 公共包
│   ├── configs/            # 配置文件
│   └── migrations/         # 数据库迁移
├── frontend/               # Vue前端应用
│   ├── src/                # 源代码
│   ├── public/             # 静态资源
│   └── docs/               # 前端文档
├── docs/                   # 项目文档
└── scripts/                # 构建脚本
```

## 🔧 配置说明

### 数据库配置
需要MySQL 8.0+数据库，并在`backend/configs/config.yaml`中配置连接信息。

### 环境变量
- `ENV`: 运行环境（development/production）
- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户名
- `DB_PASS`: 数据库密码
- `DB_NAME`: 数据库名称

## 🔮 下一步计划

### v0.2.0 - 客户管理模块
- [ ] 客户信息增删改查
- [ ] 客户状态管理
- [ ] 批量操作功能

### v0.3.0 - 授权管理核心
- [ ] 授权码生成算法
- [ ] 硬件指纹绑定
- [ ] 在线/离线授权模式

### v0.4.0 - API完善
- [ ] 授权验证API
- [ ] 激活与心跳API
- [ ] API限流和安全增强

## ⚠️ 重要说明

**这是早期开发版本，主要用于：**
- 展示系统架构和技术选型
- 收集用户反馈和建议
- 建立开发基线和版本管理

**当前版本限制：**
- 仅包含基础框架和登录功能
- 数据为模拟数据，非完整业务逻辑
- API接口尚未完全实现
- 未进行完整的安全测试

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

本项目采用GPL-3.0开源许可证。

---

**项目地址**: https://github.com/cedar-v/license-manager
**问题反馈**: https://github.com/cedar-v/license-manager/issues
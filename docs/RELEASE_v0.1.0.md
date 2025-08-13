# Release v0.1.0 - Initial Framework Release

**发布日期**: 2024-08-13  
**版本类型**: 初始框架版本

---

## 🎉 首个框架版本发布

这是License Manager项目的第一个发布版本，包含了基础的系统框架和核心功能。

## ✨ 本版本新增

- 🔐 JWT用户认证系统
- 🌐 多语言支持（中/英/日）  
- 📊 响应式管理界面
- 🛠️ 完整的前后端分离架构
- 🐳 Docker容器化部署支持
- 📝 API基础框架

## 📦 快速体验

```bash
git clone https://github.com/cedar-v/license-manager.git
cd license-manager
docker compose up -d --build
```

访问: http://localhost:18080

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
# 软件授权管理平台

[English](README.md) | [中文](README_ZH.md)

---

## 项目概述

软件授权管理平台是一个独立的软件授权管理系统，为IoT平台等软件系统提供授权码生成、分发、验证和管理服务。支持在线和离线两种授权模式，基于硬件绑定的授权机制确保安全性。

## 核心功能

- 🔧 **客户管理**：完整的客户信息管理和状态控制
- 🔐 **授权生成**：在线/离线授权模式，支持硬件指纹绑定
- 📊 **授权管理**：实时状态监控和授权生命周期管理
- 📦 **部署包生成**：自动生成包含配置的部署包
- 🌐 **API服务**：提供验证、激活、心跳监控等RESTful API
- ⚙️ **系统管理**：管理员认证和监控仪表盘
- 🛠️ **跨平台工具**：多平台硬件信息获取工具

## 技术栈

- **前端**：Vue.js 3+ 配合现代化UI组件
- **后端**：Go 1.23+ 使用Gorilla Mux路由和Logrus日志
- **数据库**：PostgreSQL 12+ / MySQL 12+
- **配置**：YAML格式配置文件
- **部署**：Docker、单机部署或系统服务

## API接口

```
POST /api/validate      - 授权验证
POST /api/activate      - 授权激活
POST /api/heartbeat     - 心跳上报
GET  /api/license/{code} - 授权信息查询
GET  /api/customers     - 客户列表API
GET  /tools/{tool}      - 工具下载
```

## 安全与性能

- **安全特性**：HMAC-SHA256签名，硬件指纹绑定，HTTPS传输加密
- **性能表现**：支持100+并发用户，API响应时间<2秒
- **可靠性**：完善的错误处理和日志记录

## 安装部署

```bash
# 克隆项目
git clone <repository-url>
cd license-manager

# 构建应用
go build -o license-manager

# 配置应用
cp config.example.yaml config.yaml
# 编辑 config.yaml 配置文件

# 运行应用
./license-manager
```

## Docker部署

```bash
# 构建Docker镜像
docker build -t license-manager .

# 使用Docker运行
docker run -p 8080:8080 -v $(pwd)/config.yaml:/app/config.yaml license-manager
```

## 开源许可证

本项目采用GNU通用公共许可证v3.0 - 详见[LICENSE](LICENSE)文件。

---

## 贡献

欢迎贡献代码！请随时提交Pull Request。

## 支持

如有任何问题或需要支持，请提交issue。 
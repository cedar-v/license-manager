# 软件授权管理平台

[中文](README.md) | [English](README_EN.md) 

---

## 项目概述

软件授权管理平台是一个独立的软件授权管理系统，为软件系统提供授权码生成、分发、验证和管理服务。支持在线和离线两种授权模式，基于硬件绑定的授权机制确保安全性。

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
- **后端**：Go 1.23+ 使用Gin框架、GORM、Viper配置和Logrus日志
- **数据库**：MySQL 8+
- **配置**：YAML格式配置文件
- **部署**：Docker、单机部署或系统服务

## 安全与性能

- **安全特性**：
  - JWT身份认证和授权机制
  - HMAC-SHA256签名验证  
  - 硬件指纹绑定防盗版
  - HTTPS传输加密
  - 多语言错误信息支持（中/英/日）
  
- **性能表现**：
  - **高并发**：Go原生协程支持，理论支持10,000+并发连接
  - **低延迟**：API平均响应时间 < 50ms
  - **内存优化**：Go GC优化，内存占用 < 100MB
  - **数据库连接池**：支持连接复用，最大化数据库性能
  
- **可靠性**：
  - 完善的错误处理和多语言错误信息
  - 结构化日志记录和监控
  - 自动数据库迁移
  - 优雅关闭和资源清理

## 安装部署

```bash
# 克隆项目
git clone https://github.com/cedar-v/license-manager.git
cd license-manager

# 本地后端构建与运行（可选）
cd backend/cmd
go build -o license-manager

# 配置（二选一）
# 1) 在当前目录创建配置文件（程序优先读取工作目录下的 config.yaml）
# cp ../configs/config.example.yaml ./config.yaml && 编辑
# 2) 直接编辑 ../configs/config.yaml（程序会自动回退读取该文件）

./license-manager
```

## Docker 部署

推荐使用 Docker Compose，已提供开发与生产编排文件：

```bash
# 开发环境（首次建议带 --build）
docker compose up -d --build

# 生产环境
docker compose -f docker-compose.prod.yml up -d

# 验证后端健康
curl http://localhost:18888/health
```

完整部署说明（反向代理、健康检查、配置挂载等）请参见 `README-Docker.md`。

## 开源许可证

本项目采用 **GNU General Public License v3.0 (GPL-3.0)** 开源许可证。

### 许可证说明

- **自由使用**：您可以自由地使用、学习、修改和分发本软件
- **开源要求**：如果您分发修改版本，必须同样开源并采用GPL-3.0许可证
- **版权保护**：使用本软件的衍生作品必须保留原始版权声明
- **无担保**：软件按"现状"提供，不提供任何明示或暗示的担保

### 商业使用

- 允许商业使用，但衍生作品必须同样开源
- 如需专有许可证或商业支持，请联系项目维护者

### 许可证全文

详细许可证条款请查看项目根目录的 [LICENSE](LICENSE) 文件，或访问：
https://www.gnu.org/licenses/gpl-3.0.html


---

## 贡献

欢迎贡献代码！请随时提交Pull Request。

## 支持

如有任何问题或需要支持，请提交issue。 
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
- **数据库**：PostgreSQL 12+ / MySQL 8+
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
cd license-manager/backend/cmd

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
# 构建Docker镜像（多阶段构建，包含前后端）
docker build -t license-manager .

# 创建配置文件
cp backend/configs/config.example.yaml config.yaml
# 编辑 config.yaml 配置数据库等信息

# 使用Docker运行
docker run -d \
  --name license-manager \
  -p 18888:18888 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  license-manager

# 查看运行状态
docker logs license-manager

# 健康检查
curl http://localhost:18888/health
```

### Docker Compose部署（推荐）

```yaml
# docker-compose.yml
version: '3.8'
services:
  license-manager:
    build: .
    ports:
      - "18888:18888"
    volumes:
      - ./config.yaml:/app/config.yaml
    environment:
      - GIN_MODE=release
    restart: unless-stopped
    
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root@123
      MYSQL_DATABASE: license_manager
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql_data:
```

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
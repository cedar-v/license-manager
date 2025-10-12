# GitHub镜像部署指南

使用预构建的GitHub Container Registry镜像快速部署许可证管理系统。

## 镜像信息

- **后端镜像**: `ghcr.io/cedar-v/license-manager-backend:v0.1.0`
- **前端镜像**: `ghcr.io/cedar-v/license-manager-frontend:v0.1.0`

## 快速部署

### 1. 初始化配置文件

```bash
# 创建配置目录
mkdir -p backend-config

# 从镜像中提取配置文件
docker run --rm -v $(pwd)/backend-config:/tmp/config ghcr.io/cedar-v/license-manager-backend:v0.1.0 sh -c "cp -r /app/backend/configs/* /tmp/config/"
```

### 2. 启动服务

```bash
docker-compose -f docker-compose.github.image.yml up -d
```

### 3. 检查服务状态

```bash
docker-compose -f docker-compose.github.image.yml ps
```

## 访问地址

- **前端**: http://localhost
- **后端API**: http://localhost:18888
- **MySQL**: localhost:3306

## 默认账号

- **用户名**: admin
- **密码**: changeMe!

## 停止服务

```bash
docker-compose -f docker-compose.github.image.yml down
```

## 清理数据

```bash
docker-compose -f docker-compose.github.image.yml down -v
```

## 故障排查

### 查看日志

```bash
# 查看所有服务日志
docker-compose -f docker-compose.github.image.yml logs

# 查看特定服务日志
docker-compose -f docker-compose.github.image.yml logs backend
docker-compose -f docker-compose.github.image.yml logs frontend
docker-compose -f docker-compose.github.image.yml logs mysql
```

### 重启服务

```bash
docker-compose -f docker-compose.github.image.yml restart
```

## 修改配置

### 提取配置文件

```bash
# 如果需要修改配置，先提取配置文件
mkdir -p backend-config
docker run --rm -v $(pwd)/backend-config:/tmp/config ghcr.io/cedar-v/license-manager-backend:v0.1.0 sh -c "cp -r /app/backend/configs/* /tmp/config/"
```

### 修改配置

编辑 `backend-config/config.prod.yaml` 文件修改数据库连接、密钥等配置。

### 重启应用配置

```bash
# 重启后端服务使配置生效
docker-compose -f docker-compose.github.image.yml restart backend
```
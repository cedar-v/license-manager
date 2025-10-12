# GitHub镜像部署指南

使用预构建的GitHub Container Registry镜像快速部署许可证管理系统。

## 镜像信息

- **后端镜像**: `ghcr.io/cedar-v/license-manager-backend:v0.1.0`
- **前端镜像**: `ghcr.io/cedar-v/license-manager-frontend:v0.1.0`

## 快速部署

### 1. 拉取镜像（可选）

```bash
docker pull ghcr.io/cedar-v/license-manager-backend:v0.1.0
docker pull ghcr.io/cedar-v/license-manager-frontend:v0.1.0
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
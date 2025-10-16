# GitHub 镜像快速部署指南

使用预构建的 GitHub Container Registry 镜像快速部署许可证管理系统。

## 快速部署

### 1. 获取部署文件

复制项目根目录的 `docker-compose.github.image.yml` 文件到你的部署目录：

```bash
# 方式一：直接下载
curl -O https://raw.githubusercontent.com/cedar-v/license-manager/main/docker-compose.github.image.yml

# 方式二：从项目中复制
cp docker-compose.github.image.yml /your/deploy/path/
```

### 2. 提取配置文件

```bash
# 提取后端配置文件
mkdir -p backend-config
docker run --rm -v $(pwd)/backend-config:/tmp/config ghcr.io/cedar-v/license-manager-backend:v0.1.0 sh -c "cp -r /app/backend/configs/* /tmp/config/"

# 提取前端 nginx 配置文件
docker run --rm -v $(pwd):/tmp/extract ghcr.io/cedar-v/license-manager-frontend:v0.1.0 sh -c "cp /etc/nginx/conf.d/default.conf /tmp/extract/nginx.conf"
```

### 3. 启动服务

```bash
docker-compose -f docker-compose.github.image.yml up
```

## 访问信息

- **前端**: http://localhost:18080
- **后端 API**: http://localhost:18888
- **默认账号**: admin / admin@123

## 常用命令

```bash
# 查看服务状态
docker-compose -f docker-compose.github.image.yml ps

# 查看日志
docker-compose -f docker-compose.github.image.yml logs -f

# 停止服务
docker-compose -f docker-compose.github.image.yml down

# 重启服务
docker-compose -f docker-compose.github.image.yml restart

# 完全清理（包括数据）
docker-compose -f docker-compose.github.image.yml down -v
```

## 配置修改

如需修改后端配置，编辑 `backend-config/config.prod.yaml` 文件，然后重启后端服务：

```bash
docker-compose -f docker-compose.github.image.yml restart backend
```
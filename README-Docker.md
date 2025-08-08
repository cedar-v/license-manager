# Docker 部署指南

本项目提供两种部署：开发环境与生产环境。以下说明已与当前 `docker-compose` 配置、镜像与反向代理实现对齐。

## 开发环境

使用 `docker compose`（Docker Desktop 或 Docker CE 20.10+）：

```bash
# 启动（首次建议带 --build）
docker compose up -d --build

# 查看状态/日志
docker compose ps
docker compose logs -f backend

# 停止
docker compose down
```

### 访问地址（开发）
- 前端: http://localhost:18080
- 后端 API: http://localhost:18888
- MySQL: localhost:13306（数据库、用户、密码见 compose）

## 生产环境

使用 `docker-compose.prod.yml`：

```bash
# 启动生产环境
docker compose -f docker-compose.prod.yml up -d

# 查看服务状态
docker compose -f docker-compose.prod.yml ps
```

### 生产配置说明
- 后端配置文件已通过卷挂载至容器内：
  - 开发：`backend/configs/config.dev.yaml` → 容器内 `/app/backend/cmd/config.yaml`
  - 生产：`backend/configs/config.prod.yaml` → 容器内 `/app/backend/cmd/config.yaml`
- 前端 Nginx 反向代理已指向容器网络内的 `backend:18888`，无需依赖 `host.docker.internal`。
- 健康检查统一使用 `curl`，后端镜像已内置 `curl`。
- 如需隐藏后端对外端口，可在生产 compose 中移除后端端口映射，仅通过前端或网关暴露。

## 服务说明

### MySQL 8.0
- 端口：3306（生产）/ 13306（开发宿主机映射）
- 库：`license_manager`，用户：`license_user`
- 字符集/排序：`utf8mb4`/`utf8mb4_unicode_ci`
- 开启慢查询日志（生产）

### 后端服务
- 端口：18888
- 健康检查：`/health`
- 配置：通过卷挂载 `config.dev.yaml` / `config.prod.yaml` 为容器内 `config.yaml`
- 日志：标准输出 + `./logs` 挂载目录

### 前端服务
- 端口：80（生产）/ 18080（开发映射）
- 反向代理：`/api/` → `backend:18888`
- 静态资源：启用 Gzip 与长缓存

### Redis（可选，仅生产）
- 端口：6379，AOF 持久化

## 常用命令

```bash
# 重新构建镜像
docker compose build

# 启动指定服务
docker compose up -d mysql backend

# 进入容器
docker compose exec backend sh
docker compose exec mysql mysql -uroot -p

# 备份/恢复数据库
docker compose exec mysql mysqldump -uroot -p license_manager > backup.sql
docker compose exec -T mysql mysql -uroot -p license_manager < backup.sql

# 清理未使用资源
docker system prune -f
```

## 故障排查

1) MySQL 启动失败
```bash
docker compose logs mysql
docker compose down -v && docker compose up -d
```

2) 后端连接数据库失败
```bash
# 确认 MySQL 健康
docker compose ps

# 检查容器内配置（已挂载到 /app/backend/cmd/config.yaml）
docker compose exec backend cat /app/backend/cmd/config.yaml
```

3) 前端无法访问后端 API
```bash
# 容器内连通性
docker compose exec frontend sh -c "curl -sS http://backend:18888/health"

# 检查 nginx 配置
docker compose exec frontend cat /etc/nginx/conf.d/default.conf
```

## 注意事项

- 生产环境请修改默认账号、数据库口令与 JWT 秘钥（见 `backend/configs/config.prod.yaml`）。
- 如使用较新 Docker，YAML 顶部 `version` 字段会被忽略，可按需移除提示。
- `deploy.resources` 仅在 Swarm/K8s 生效，单机 Compose 可忽略。
- 如置于公网，建议仅暴露前端端口，通过前端或网关转发后端 API。
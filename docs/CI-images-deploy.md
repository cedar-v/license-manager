## 使用 CI 产物镜像替换本地构建的部署方案

本文说明如何使用 GitHub Actions 构建并推送到 GHCR 的镜像，直接替换本地 `docker compose up -d --build` 的构建流程，以实现更一致、更快的部署。

### 前置条件
- 已通过仓库内工作流发布镜像：
  - 后端：`.github/workflows/release_backend.yml`
  - 前端：`.github/workflows/release_frontend.yml`
- 镜像仓库（GHCR）可访问：`ghcr.io/<owner>/<repo>-backend:<tag>`、`ghcr.io/<owner>/<repo>-frontend:<tag>`
- 已在部署机登录 GHCR：
  ```bash
  docker login ghcr.io -u <YOUR_GITHUB_USER> -p <YOUR_GITHUB_TOKEN>
  ```

### 方式 A（推荐）：直接使用生产编排（支持镜像变量）
生产编排文件 `docker-compose.prod.yml` 已支持通过环境变量替换镜像。

1) 在项目根目录创建 `.env`：
```env
# 指向 GHCR 镜像（可用具体版本号，如 v1.0.0）
BACKEND_IMAGE=ghcr.io/<owner>/<repo>-backend:latest
FRONTEND_IMAGE=ghcr.io/<owner>/<repo>-frontend:latest

# 可选：本地端口映射（如用于开发环境）
MYSQL_PORT=13306
BACKEND_PORT=18888
FRONTEND_PORT=18080
```

2) 启动：
```bash
docker compose -f docker-compose.prod.yml up -d
```

3) 验证：
```bash
curl http://localhost:18888/health
```

说明：
- 该方式同时适用于生产/准生产/本地测试场景，只需通过 `.env` 控制端口与镜像标签即可。
- 若镜像私有，请确保已成功 `docker login ghcr.io`。

### 方式 B：临时覆盖镜像（不创建 .env）
直接在一次性命令中覆盖镜像变量：
```bash
BACKEND_IMAGE=ghcr.io/<owner>/<repo>-backend:v1.0.0 \
FRONTEND_IMAGE=ghcr.io/<owner>/<repo>-frontend:v1.0.0 \
docker compose -f docker-compose.prod.yml up -d
```

### 常用操作
- 拉取新版本并滚动更新：
```bash
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

- 回滚到上一版本：
```bash
# 修改 .env 中的镜像标签为上一个版本
docker compose -f docker-compose.prod.yml up -d
```

### 与开发编排的关系
- 当前开发编排 `docker-compose.yml` 使用 `build` 本地构建。如果希望开发环境也使用远端镜像，推荐直接沿用方式 A（使用 `docker-compose.prod.yml`，并通过 `.env` 调整端口为开发习惯）。
- 也可以手动修改 `docker-compose.yml`：将 `backend`、`frontend` 的 `build` 段替换为 `image: ${BACKEND_IMAGE}` / `image: ${FRONTEND_IMAGE}`，并在 `.env` 中声明对应镜像。

### 标签与架构
- CI 工作流默认推送 `latest` 与手动输入的版本号（如 `v1.0.0`），并构建多架构镜像：`linux/amd64, linux/arm64`。
- 生产环境建议使用显式版本标签部署，避免 `latest` 带来的不可预期变更。

### 故障排查
- 无法拉取镜像：检查 GHCR 登录权限、镜像是否设为公开或部署机是否具备读取私有包权限。
- 拉取缓慢或超时：配置镜像加速器或重试。
- 启动后访问异常：参考 `README-Docker.md` 的故障排查章节（健康检查、反向代理、容器间网络连通性等）。



# 前端自动化部署

基于Docker的前端自动化部署方案，使用GitHub Actions实现CI/CD。

## 部署特点

- 🐳 Docker容器化部署
- 🚀 自动构建和部署
- 🔄 零停机更新
- 🌐 内置Nginx服务器
- 📱 支持Vue Router history模式
- 🔗 API代理到后端服务

## 配置要求

### GitHub Secrets
在GitHub仓库设置以下Secrets：
- `HOST_1`: 服务器IP地址
- `PASS_1`: 服务器root密码

### 服务器要求
- 安装Docker

## 部署流程

1. 推送代码到main分支（前端相关文件有变更时）
2. GitHub Actions自动触发
3. 构建Docker镜像（Vue3 + Nginx）
4. 上传镜像到服务器
5. 停止旧容器并启动新容器
6. 健康检查

## 访问地址

部署完成后，前端服务将在以下地址提供：
- **前端地址**: http://服务器IP:18080
- **API代理**: 前端的 `/api/*` 请求会代理到后端 `http://localhost:18888/`

## 手动触发

可以在GitHub Actions页面手动触发部署，无需推送代码。

## 监控和维护

### 查看容器状态
```bash
docker ps | grep license-manager-frontend
```

### 查看容器日志
```bash
docker logs license-manager-frontend
```

### 重启容器
```bash
docker restart license-manager-frontend
```

## 故障排除

1. **部署失败**: 检查GitHub Actions日志
2. **容器无法启动**: 检查Docker镜像和容器配置
3. **页面无法访问**: 检查服务器防火墙18080端口
4. **API请求失败**: 确认后端服务运行在18888端口 
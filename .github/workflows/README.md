# 部署配置说明

## 环境变量配置

在 GitHub 仓库的 Settings -> Secrets and variables -> Actions 中添加以下密钥：

### 必需的 Secrets:
- `HOST_1`: 服务器IP地址或域名
- `PASS_1`: 服务器root用户密码

## 部署流程

1. **自动触发**: 当推送到 main 分支且修改了 backend/ 目录下的文件时自动部署
2. **手动触发**: 在 Actions 页面可以手动触发部署

## 部署步骤

1. 编译并测试 Go 后端代码
2. 构建 Docker 镜像
3. 通过 SCP 上传镜像到服务器
4. SSH 连接服务器执行部署脚本：
   - 停止旧容器
   - 加载新镜像
   - 启动新容器
   - 健康检查

## 服务器要求

- 安装 Docker
- 开放 8080 端口
- 创建部署目录: `/opt/license-manager/`

## 配置文件

首次部署会自动创建 `config.yaml`，请根据实际需要修改配置。
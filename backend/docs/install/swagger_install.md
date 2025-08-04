# Swagger安装和使用说明

## 安装swaggo工具

### 方法1: 使用go install（推荐）
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 方法2: 使用go get（旧版本）
```bash
go get -u github.com/swaggo/swag/cmd/swag
```

## 生成swagger文档

在backend目录下运行：

```bash
cd backend
swag init -g cmd/main.go -o ./docs/swagger
```

## 访问Swagger UI

启动服务后，访问：
- Swagger UI: http://localhost:18888/swagger/index.html
- Swagger JSON: http://localhost:18888/swagger/doc.json

## 常用命令

```bash
# 生成swagger文档
swag init -g cmd/main.go -o ./docs/swagger

# 格式化swagger注释
swag fmt

# 验证swagger注释
swag init --parseVendor -g cmd/main.go
```
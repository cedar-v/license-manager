# 许可证加密设计文档

## 加密方案

使用 **RSA 数字签名**确保许可证文件的完整性和来源认证。

### 签名算法
- **算法**: RSA-PSS
- **哈希**: SHA-256
- **密钥长度**: 2048 位（推荐）或 4096 位

## 服务器端实现

### 密钥管理
- 私钥存储在配置文件中：`configs/rsa_private_key.pem`
- 配置文件路径：`license.rsa.private_key_path`

### 签名流程
1. 构建许可证数据（JSON）
2. 使用 RSA 私钥对数据进行数字签名
3. 生成许可证文件结构：
   ```json
   {
     "data": "原始JSON数据（字符串）",
     "signature": "base64编码的签名",
     "algorithm": "RSA-PSS-SHA256"
   }
   ```
4. Base64 编码整个结构

### 代码位置
- 签名实现：`backend/internal/service/license_service.go::signLicenseFile()`
- 密钥工具：`backend/pkg/utils/crypto.go`

## 客户端实现

### 公钥加载
支持三种方式（按优先级）：
1. 当前目录的 `rsa_public_key.pem` 文件
2. 环境变量 `LICENSE_RSA_PUBLIC_KEY`（PEM 格式字符串）
3. 环境变量 `LICENSE_RSA_PUBLIC_KEY_PATH`（文件路径）

### 验证流程
1. Base64 解码许可证文件
2. 解析 JSON 结构
3. 使用 RSA 公钥验证签名
4. 验证通过后解析许可证数据

### 代码位置
- 验证实现：`backend/cmd/client-demo/rsa.go::verifyRSASignature()`
- 许可证解析：`backend/cmd/client-demo/main.go::decryptLicenseFile()`

## 密钥生成

使用工具生成密钥对：
```bash
go run backend/cmd/gen-rsa-keys/main.go
```

生成文件：
- `rsa_private_key.pem` - 服务器使用（必须保密）
- `rsa_public_key.pem` - 客户端使用（可公开）

## 安全性说明

1. **防伪造**：只有服务器拥有私钥，无法伪造有效签名
2. **防篡改**：任何数据修改都会导致签名验证失败
3. **密钥分离**：公钥可公开，私钥仅服务器持有

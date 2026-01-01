package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// LicenseConfig 自包含授权码配置结构（极简字段名）
type LicenseConfig struct {
	EndDate          time.Time              `json:"e"`           // 到期时间
	FeatureConfig    map[string]interface{} `json:"f,omitempty"` // 功能配置
	UsageLimits      map[string]interface{} `json:"l,omitempty"` // 使用限制
	CustomParameters map[string]interface{} `json:"p,omitempty"` // 自定义参数
}

// EncodeLicenseData 将配置数据编码为自包含授权码（极简模式）
func EncodeLicenseData(config *LicenseConfig, privateKeyPath string) (string, error) {
	if config == nil {
		return "", errors.New("config cannot be nil")
	}

	// 使用极简紧凑序列化
	configBytes, err := ultraCompactMarshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	// 使用HMAC-SHA256替代RSA签名（短签名模式）
	// 从私钥文件中提取密钥作为HMAC密钥（简化处理）
	hmacKey, err := getHMACKeyFromRSAKey(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to get HMAC key: %w", err)
	}

	signature, err := signWithHMAC(configBytes, hmacKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	// 极简格式：[版本(4位)][时间戳(32位)][配置数据][HMAC签名(32字节)]
	version := uint8(3) // 版本3（极简模式）
	timestamp := uint32(config.EndDate.Unix())

	// 打包头部：版本(4位) + 时间戳(28位) 到5字节
	header := make([]byte, 5)
	header[0] = (version << 4) | byte(timestamp>>28)
	header[1] = byte(timestamp >> 20)
	header[2] = byte(timestamp >> 12)
	header[3] = byte(timestamp >> 4)
	header[4] = byte(timestamp << 4) // 最后4位保留给标志位

	// 设置字段存在标志
	flagBits := uint8(0)
	if config.FeatureConfig != nil {
		flagBits |= 1 << 0
	}
	if config.UsageLimits != nil {
		flagBits |= 1 << 1
	}
	if config.CustomParameters != nil {
		flagBits |= 1 << 2
	}
	header[4] |= flagBits

	// 构建数据：[头部(5字节)][配置数据][签名(32字节)]
	data := append(header, configBytes...)
	data = append(data, signature...)

	// 使用URL安全的Base64编码
	encodedData := base64.RawURLEncoding.EncodeToString(data)

	// 生成随机码（2位）
	randomStr, err := generateRandomString(2)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	// 最终格式：LIC-{2位随机}-{数据}
	return fmt.Sprintf("LIC-%s-%s", randomStr, encodedData), nil
}

// ultraCompactMarshal 极简紧凑序列化（无长度前缀）
func ultraCompactMarshal(config *LicenseConfig) ([]byte, error) {
	var buf []byte

	// 直接序列化存在的字段，不使用长度前缀
	// 依赖解码时知道字段顺序和格式

	if config.FeatureConfig != nil {
		mapBytes, err := json.Marshal(config.FeatureConfig)
		if err != nil {
			return nil, err
		}
		buf = append(buf, mapBytes...)
		buf = append(buf, 0) // 字段分隔符
	}

	if config.UsageLimits != nil {
		mapBytes, err := json.Marshal(config.UsageLimits)
		if err != nil {
			return nil, err
		}
		buf = append(buf, mapBytes...)
		buf = append(buf, 0) // 字段分隔符
	}

	if config.CustomParameters != nil {
		mapBytes, err := json.Marshal(config.CustomParameters)
		if err != nil {
			return nil, err
		}
		buf = append(buf, mapBytes...)
		buf = append(buf, 0) // 字段分隔符
	}

	return buf, nil
}

// getHMACKeyFromRSAKey 从RSA密钥文件提取HMAC密钥
func getHMACKeyFromRSAKey(keyPath string) ([]byte, error) {
	// 首先尝试作为私钥加载
	if privateKey, err := getLicensePrivateKey(keyPath); err == nil {
		// 使用私钥的模数作为HMAC密钥
		key := privateKey.N.Bytes()
		// 截取前32字节作为HMAC-SHA256密钥
		if len(key) > 32 {
			key = key[:32]
		} else if len(key) < 32 {
			// 补齐到32字节
			padding := make([]byte, 32-len(key))
			key = append(key, padding...)
		}
		return key, nil
	}

	// 如果私钥加载失败，尝试作为公钥加载
	if publicKey, err := getLicensePublicKey(keyPath); err == nil {
		// 使用公钥的模数作为HMAC密钥
		key := publicKey.N.Bytes()
		// 截取前32字节作为HMAC-SHA256密钥
		if len(key) > 32 {
			key = key[:32]
		} else if len(key) < 32 {
			// 补齐到32字节
			padding := make([]byte, 32-len(key))
			key = append(key, padding...)
		}
		return key, nil
	}

	return nil, errors.New("failed to load RSA key for HMAC")
}

// getLicensePublicKey 获取RSA公钥（复用crypto.go的函数）
func getLicensePublicKey(publicKeyPath string) (*RSAPublicKey, error) {
	// 检查文件是否存在
	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("RSA公钥文件不存在: %s", publicKeyPath)
	}

	// 加载公钥
	publicKey, err := LoadRSAPublicKeyFromFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("RSA公钥文件损坏或格式错误: %w", err)
	}

	return publicKey, nil
}

// signWithHMAC 使用HMAC-SHA256签名
func signWithHMAC(data []byte, key []byte) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil), nil
}

// ultraCompactUnmarshal 从极简格式反序列化配置
func ultraCompactUnmarshal(data []byte, flags uint8) (*LicenseConfig, error) {
	config := &LicenseConfig{}
	offset := 0

	// 解析f字段（如果存在）
	if flags&1 != 0 {
		end := findFieldEnd(data, offset)
		if end == -1 {
			return nil, errors.New("invalid feature config format")
		}
		var m map[string]interface{}
		if err := json.Unmarshal(data[offset:end], &m); err != nil {
			return nil, fmt.Errorf("invalid feature config: %w", err)
		}
		config.FeatureConfig = m
		offset = end + 1
	}

	// 解析l字段（如果存在）
	if flags&(1<<1) != 0 {
		end := findFieldEnd(data, offset)
		if end == -1 {
			return nil, errors.New("invalid usage limits format")
		}
		var m map[string]interface{}
		if err := json.Unmarshal(data[offset:end], &m); err != nil {
			return nil, fmt.Errorf("invalid usage limits: %w", err)
		}
		config.UsageLimits = m
		offset = end + 1
	}

	// 解析p字段（如果存在）
	if flags&(1<<2) != 0 {
		if offset >= len(data) {
			return nil, errors.New("missing custom parameters")
		}
		// 最后一个字段到数据结尾
		var m map[string]interface{}
		if err := json.Unmarshal(data[offset:], &m); err != nil {
			return nil, fmt.Errorf("invalid custom parameters: %w", err)
		}
		config.CustomParameters = m
	}

	return config, nil
}

// findFieldEnd 查找字段结束位置（以null字节分隔）
func findFieldEnd(data []byte, start int) int {
	for i := start; i < len(data); i++ {
		if data[i] == 0 {
			return i
		}
	}
	return -1
}

// DecodeLicenseData 从授权码解析配置数据
func DecodeLicenseData(code string, publicKeyPath string) (*LicenseConfig, error) {
	if !strings.HasPrefix(code, "LIC-") {
		return nil, errors.New("invalid license format")
	}

	// 解析授权码格式：LIC-{2位随机码}-{数据}
	parts := strings.Split(code, "-")
	if len(parts) != 3 || parts[0] != "LIC" {
		return nil, errors.New("invalid license format")
	}

	randomStr := parts[1]
	encodedData := parts[2]

	// 验证随机码长度
	if len(randomStr) != 2 {
		return nil, errors.New("invalid random string length")
	}

	// 使用URL安全的Base64解码
	dataBytes, err := base64.RawURLEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// 解析极简格式：[头部(5字节)][配置数据][签名(32字节)]
	if len(dataBytes) < 5+32 { // 头部 + 最少签名
		return nil, errors.New("invalid data format")
	}

	// 解析头部
	header := dataBytes[:5]
	version := header[0] >> 4
	timestamp := uint32(header[0]&0x0F)<<28 | uint32(header[1])<<20 | uint32(header[2])<<12 | uint32(header[3])<<4 | uint32(header[4]>>4)
	flags := header[4] & 0x0F

	if version != 3 {
		return nil, errors.New("unsupported version")
	}

	config := &LicenseConfig{
		EndDate: time.Unix(int64(timestamp), 0),
	}

	// 分离配置数据和签名
	configDataEnd := len(dataBytes) - 32
	configBytes := dataBytes[5:configDataEnd]
	signature := dataBytes[configDataEnd:]

	// 获取HMAC密钥用于验证
	hmacKey, err := getHMACKeyFromRSAKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get HMAC key: %w", err)
	}

	// 验证HMAC签名
	if !verifyHMAC(configBytes, signature, hmacKey) {
		return nil, errors.New("signature verification failed")
	}

	// 从极简格式反序列化配置
	fullConfig, err := ultraCompactUnmarshal(configBytes, flags)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 合并时间戳
	fullConfig.EndDate = config.EndDate
	return fullConfig, nil
}

// verifyHMAC 验证HMAC签名
func verifyHMAC(data []byte, signature []byte, key []byte) bool {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	expectedMAC := h.Sum(nil)
	return hmac.Equal(signature, expectedMAC)
}

// SignLicenseData 使用RSA私钥对数据进行签名
func SignLicenseData(data []byte, privateKeyPath string) ([]byte, error) {
	// 获取RSA私钥
	privateKey, err := getLicensePrivateKey(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	// 使用现有的RSA签名方法（PSS模式，与许可证保持一致）
	signature, err := privateKey.SignData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}

	// 解码Base64签名
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}

	return signatureBytes, nil
}

// VerifyLicenseSignature 使用RSA公钥验证签名
func VerifyLicenseSignature(data []byte, signature []byte, publicKeyPath string) error {
	// 获取RSA公钥
	publicKey, err := getLicensePublicKey(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	// 编码签名为Base64（与SignData方法匹配）
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	// 使用现有的RSA验证方法
	err = publicKey.VerifySignature(data, signatureBase64)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

// getLicensePrivateKey 获取许可证RSA私钥
func getLicensePrivateKey(privateKeyPath string) (*RSAPrivateKey, error) {
	// 检查文件是否存在
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("RSA私钥文件不存在，请先生成RSA密钥对: %s", privateKeyPath)
	}

	// 加载私钥
	privateKey, err := LoadRSAPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("RSA私钥文件损坏或格式错误: %w", err)
	}

	return privateKey, nil
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}
	return string(bytes), nil
}

// generateChecksum 生成校验码
func generateChecksum(input string) string {
	sum := 0
	for _, char := range input {
		sum += int(char)
	}
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	checksum := ""
	for i := 0; i < 4; i++ {
		checksum += string(chars[sum%len(chars)])
		sum = sum / len(chars)
	}
	return checksum
}

package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用bcrypt加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// RSAPrivateKey RSA私钥结构
type RSAPrivateKey struct {
	*rsa.PrivateKey
}

// RSAPublicKey RSA公钥结构
type RSAPublicKey struct {
	*rsa.PublicKey
}

// LoadRSAPrivateKeyFromFile 从文件加载RSA私钥
func LoadRSAPrivateKeyFromFile(filePath string) (*RSAPrivateKey, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("无法解析PEM格式的私钥")
	}

	var key *rsa.PrivateKey
	var errParse error

	// 尝试PKCS1格式
	key, errParse = x509.ParsePKCS1PrivateKey(block.Bytes)
	if errParse != nil {
		// 尝试PKCS8格式
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("解析私钥失败: %w", err)
		}
		var ok bool
		key, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("不是有效的RSA私钥")
		}
	}

	return &RSAPrivateKey{PrivateKey: key}, nil
}

// LoadRSAPublicKeyFromFile 从文件加载RSA公钥
func LoadRSAPublicKeyFromFile(filePath string) (*RSAPublicKey, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取公钥文件失败: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("无法解析PEM格式的公钥")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("不是有效的RSA公钥")
	}

	return &RSAPublicKey{PublicKey: rsaPub}, nil
}

// LoadRSAPublicKeyFromString 从字符串加载RSA公钥（支持PEM格式字符串）
func LoadRSAPublicKeyFromString(pubKeyStr string) (*RSAPublicKey, error) {
	block, _ := pem.Decode([]byte(pubKeyStr))
	if block == nil {
		return nil, errors.New("无法解析PEM格式的公钥")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("不是有效的RSA公钥")
	}

	return &RSAPublicKey{PublicKey: rsaPub}, nil
}

// SignData 使用RSA私钥对数据进行数字签名
func (key *RSAPrivateKey) SignData(data []byte) (string, error) {
	// 计算数据的SHA256哈希
	hashed := sha256.Sum256(data)

	// 使用PSS模式进行签名（更安全）
	signature, err := rsa.SignPSS(rand.Reader, key.PrivateKey, crypto.SHA256, hashed[:], &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	})
	if err != nil {
		return "", fmt.Errorf("签名失败: %w", err)
	}

	// 返回base64编码的签名
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySignature 使用RSA公钥验证数字签名
func (key *RSAPublicKey) VerifySignature(data []byte, signatureBase64 string) error {
	// 解码base64签名
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return fmt.Errorf("解码签名失败: %w", err)
	}

	// 计算数据的SHA256哈希
	hashed := sha256.Sum256(data)

	// 验证签名
	err = rsa.VerifyPSS(key.PublicKey, crypto.SHA256, hashed[:], signature, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	})
	if err != nil {
		return fmt.Errorf("签名验证失败: %w", err)
	}

	return nil
}

// GenerateRSAKeyPair 生成RSA密钥对并保存到文件
func GenerateRSAKeyPair(privateKeyPath, publicKeyPath string, keySize int) error {
	if keySize < 2048 {
		keySize = 2048 // 最小2048位
	}

	// 生成密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("生成RSA密钥对失败: %w", err)
	}

	// 保存私钥（PKCS1格式）
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("保存私钥失败: %w", err)
	}

	// 保存公钥
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("序列化公钥失败: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		return fmt.Errorf("保存公钥失败: %w", err)
	}

	return nil
}

package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// 默认RSA公钥（应该从配置文件或环境变量加载）
// 在生产环境中，这个公钥应该编译到程序中，或者从配置文件读取
var defaultRSAPublicKeyPEM = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyour-public-key-here
-----END PUBLIC KEY-----
`

var rsaPublicKey *rsa.PublicKey

// initRSAPublicKey 初始化RSA公钥（从文件或环境变量加载）
func initRSAPublicKey() error {
	if rsaPublicKey != nil {
		return nil
	}

	var pubKeyPEM string
	var err error

	// 1. 首先尝试从当前目录下的 rsa_public_key.pem 文件加载
	defaultPubKeyPath := "rsa_public_key.pem"
	if _, err := os.Stat(defaultPubKeyPath); err == nil {
		data, err := os.ReadFile(defaultPubKeyPath)
		if err == nil {
			pubKeyPEM = string(data)
			// 获取绝对路径用于日志
			if absPath, err := filepath.Abs(defaultPubKeyPath); err == nil {
				log.Printf("从当前目录加载RSA公钥: %s", absPath)
			} else {
				log.Printf("从当前目录加载RSA公钥: %s", defaultPubKeyPath)
			}
		}
	}

	// 2. 如果当前目录没有，尝试从环境变量获取公钥内容
	if pubKeyPEM == "" {
		pubKeyPEM = os.Getenv("LICENSE_RSA_PUBLIC_KEY")
		if pubKeyPEM != "" {
			log.Println("从环境变量 LICENSE_RSA_PUBLIC_KEY 加载RSA公钥")
		}
	}

	// 3. 如果还是没有，尝试从环境变量指定的文件路径加载
	if pubKeyPEM == "" {
		pubKeyPath := os.Getenv("LICENSE_RSA_PUBLIC_KEY_PATH")
		if pubKeyPath != "" {
			data, err := os.ReadFile(pubKeyPath)
			if err == nil {
				pubKeyPEM = string(data)
				log.Printf("从环境变量指定的路径加载RSA公钥: %s", pubKeyPath)
			}
		}
	}

	// 如果还是没有找到公钥，返回错误
	if pubKeyPEM == "" {
		// 获取当前工作目录用于错误提示
		cwd, _ := os.Getwd()
		return fmt.Errorf("RSA公钥未配置，请将公钥文件 rsa_public_key.pem 放在当前目录 (%s) 下，或设置 LICENSE_RSA_PUBLIC_KEY 环境变量或 LICENSE_RSA_PUBLIC_KEY_PATH 环境变量", cwd)
	}

	key, err := loadRSAPublicKeyFromString(pubKeyPEM)
	if err != nil {
		return fmt.Errorf("加载RSA公钥失败: %w", err)
	}

	rsaPublicKey = key
	return nil
}

// loadRSAPublicKeyFromString 从PEM格式字符串加载RSA公钥
func loadRSAPublicKeyFromString(pubKeyStr string) (*rsa.PublicKey, error) {
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

	return rsaPub, nil
}

// verifyRSASignature 使用RSA公钥验证数字签名
func verifyRSASignature(data []byte, signatureBase64 string) error {
	if rsaPublicKey == nil {
		if err := initRSAPublicKey(); err != nil {
			return fmt.Errorf("RSA公钥未初始化: %w", err)
		}
	}

	// 解码base64签名
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return fmt.Errorf("解码签名失败: %w", err)
	}

	// 计算数据的SHA256哈希
	hashed := sha256.Sum256(data)

	// 验证签名
	err = rsa.VerifyPSS(rsaPublicKey, crypto.SHA256, hashed[:], signature, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
		Hash:       crypto.SHA256,
	})
	if err != nil {
		return fmt.Errorf("签名验证失败: %w", err)
	}

	return nil
}

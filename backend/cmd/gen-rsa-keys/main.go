package main

import (
	"flag"
	"fmt"
	"license-manager/pkg/utils"
)

func main() {
	var (
		privateKeyPath = flag.String("private", "../../configs/rsa_private_key.pem", "私钥文件路径")
		publicKeyPath  = flag.String("public", "../../configs/rsa_public_key.pem", "公钥文件路径")
		keySize        = flag.Int("size", 2048, "密钥大小（2048或4096）")
	)
	flag.Parse()

	fmt.Printf("正在生成 RSA 密钥对...\n")
	fmt.Printf("密钥大小: %d 位\n", *keySize)
	fmt.Printf("私钥路径: %s\n", *privateKeyPath)
	fmt.Printf("公钥路径: %s\n", *publicKeyPath)

	if err := utils.GenerateRSAKeyPair(*privateKeyPath, *publicKeyPath, *keySize); err != nil {
		fmt.Printf("生成失败: %v\n", err)
		return
	}

	fmt.Printf("✓ 密钥对生成成功！\n")
	fmt.Printf("\n请将私钥文件保存在安全的位置（服务器端）\n")
	fmt.Printf("请将公钥提供给客户端开发人员（可以编译到客户端程序中）\n")
}

package utils

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// GenerateLegacyAuthorizationCode 生成旧规则授权码：
// LIC-{客户ID前4位}-{12位随机}-{4位校验码}
func GenerateLegacyAuthorizationCode(customerID string) (string, error) {
	customerCode := "COMP"
	if len(customerID) >= 4 {
		customerCode = strings.ToUpper(customerID[:4])
	}

	randomStr, err := generateSecureRandomString(12)
	if err != nil {
		return "", err
	}

	checksum := generateChecksum62(customerCode + randomStr)
	return fmt.Sprintf("LIC-%s-%s-%s", customerCode, randomStr, checksum), nil
}

func generateSecureRandomString(length int) (string, error) {
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

func generateChecksum62(input string) string {
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


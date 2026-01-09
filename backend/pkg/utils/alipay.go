package utils

import (
	"fmt"
	"net/http"
	"time"
)

// GeneratePaymentNo 生成支付单号
func GeneratePaymentNo() string {
	timestamp := time.Now().Format("20060102150405")
	random := GenerateRandomString(6)
	return fmt.Sprintf("PAY%s%s", timestamp, random)
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // 简单的随机性保证
	}
	return string(result)
}

// SendNotificationResponse 发送通知响应
func SendNotificationResponse(w http.ResponseWriter, success bool) {
	response := "success"
	if !success {
		response = "failure"
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(response))
}

package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 允许的域名列表
		allowedOrigins := map[string]bool{
			"http://localhost:3000":  true, // 前端开发服务器
			"http://localhost:8080":  true, // 前端开发服务器 (Vite)
			"http://localhost:18888": true, // Swagger UI
			"http://127.0.0.1:18888": true, // Swagger UI (IPv4)
			"":                       true, // 同源请求
		}

		// 如果是允许的域名或者没有Origin头（同源请求），则允许
		if allowedOrigins[origin] || origin == "" {
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				c.Header("Access-Control-Allow-Origin", "*")
			}
		} else {
			// 对于其他域名，也允许（开发环境）
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Header("Access-Control-Max-Age", "86400") // 24小时

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

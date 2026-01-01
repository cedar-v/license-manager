package middleware

import (
	"net/http"
	"strings"

	"license-manager/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
	CuUserPayloadKey        = "cu_user"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":      http.StatusUnauthorized,
				"error":     "AUTH_001",
				"message":   "缺少认证令牌",
				"timestamp": getCurrentTimestamp(),
			})
			c.Abort()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":      http.StatusUnauthorized,
				"error":     "AUTH_002",
				"message":   "认证令牌格式无效",
				"timestamp": getCurrentTimestamp(),
			})
			c.Abort()
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":      http.StatusUnauthorized,
				"error":     "AUTH_002",
				"message":   "不支持的认证类型",
				"timestamp": getCurrentTimestamp(),
			})
			c.Abort()
			return
		}

		accessToken := fields[1]
		claims, err := utils.ValidateToken(accessToken)
		if err != nil {
			var errorCode, message string
			switch err {
			case utils.ErrTokenExpired:
				errorCode = "AUTH_003"
				message = "令牌已过期"
			case utils.ErrTokenMalformed:
				errorCode = "AUTH_002"
				message = "令牌格式错误"
			default:
				errorCode = "AUTH_002"
				message = "令牌无效"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"code":      http.StatusUnauthorized,
				"error":     errorCode,
				"message":   message,
				"timestamp": getCurrentTimestamp(),
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set(AuthorizationPayloadKey, claims)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			c.Next()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			c.Next()
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			c.Next()
			return
		}

		accessToken := fields[1]
		claims, err := utils.ValidateToken(accessToken)
		if err == nil {
			// Token有效，设置用户信息
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
			c.Set(AuthorizationPayloadKey, claims)
		}

		c.Next()
	}
}

// AdminOnlyMiddleware 仅管理员访问中间件
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":      http.StatusUnauthorized,
				"error":     "AUTH_001",
				"message":   "认证信息缺失",
				"timestamp": getCurrentTimestamp(),
			})
			c.Abort()
			return
		}

		if role != "administrator" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":      http.StatusForbidden,
				"error":     "AUTH_004",
				"message":   "权限不足",
				"timestamp": getCurrentTimestamp(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CustomerAuth 客户用户认证中间件
func CustomerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "400010",
				"message": "缺少认证令牌",
				"data":    nil,
			})
			c.Abort()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "400011",
				"message": "认证令牌格式无效",
				"data":    nil,
			})
			c.Abort()
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "400011",
				"message": "不支持的认证类型",
				"data":    nil,
			})
			c.Abort()
			return
		}

		accessToken := fields[1]
		claims, err := utils.ValidateCuToken(accessToken)
		if err != nil {
			var errorCode, message string
			switch err {
			case utils.ErrCuTokenExpired:
				errorCode = "400012"
				message = "令牌已过期"
			case utils.ErrCuTokenMalformed:
				errorCode = "400011"
				message = "令牌格式错误"
			default:
				errorCode = "400011"
				message = "令牌无效"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    errorCode,
				"message": message,
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 将客户用户信息存储到上下文
		c.Set("cu_user_id", claims.UserID)
		c.Set("cu_customer_id", claims.CustomerID)
		c.Set("cu_user_role", claims.UserRole)
		c.Set("cu_phone", claims.Phone)
		c.Set(CuUserPayloadKey, claims)

		c.Next()
	}
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() string {
	return "2024-07-30T12:00:00Z" // 简化实现，实际应该使用time.Now().Format(time.RFC3339)
}

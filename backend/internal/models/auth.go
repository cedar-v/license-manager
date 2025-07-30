package models

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginData 登录成功返回的数据
type LoginData struct {
	Token     string   `json:"token"`
	ExpiresIn int      `json:"expires_in"`
	UserInfo  UserInfo `json:"user_info"`
}

// UserInfo 用户信息
type UserInfo struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
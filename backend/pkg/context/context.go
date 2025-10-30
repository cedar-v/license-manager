package context

import (
	"context"

	"github.com/gin-gonic/gin"
)

// 定义Context键的类型，避免键冲突
type contextKey string

const (
	// LanguageKey 存储语言信息的上下文键
	LanguageKey contextKey = "language"
	// RequestIDKey 存储请求ID的上下文键
	RequestIDKey contextKey = "request_id"
	// UserIDKey 存储用户ID的上下文键
	UserIDKey contextKey = "user_id"
	// TraceIDKey 存储链路追踪ID的上下文键
	TraceIDKey contextKey = "trace_id"
)

// DefaultLanguage 默认语言
const DefaultLanguage = "zh-CN"

// WithLanguage 将语言信息添加到Context中
func WithLanguage(ctx context.Context, lang string) context.Context {
	if lang == "" {
		lang = DefaultLanguage
	}
	return context.WithValue(ctx, LanguageKey, lang)
}

// GetLanguageFromContext 从Context中获取语言信息
func GetLanguageFromContext(ctx context.Context) string {
	if lang, ok := ctx.Value(LanguageKey).(string); ok && lang != "" {
		return lang
	}
	return DefaultLanguage
}

// WithRequestID 将请求ID添加到Context中
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestIDFromContext 从Context中获取请求ID
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// WithUserID 将用户ID添加到Context中
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserIDFromContext 从context中获取当前用户ID
func GetUserIDFromContext(ctx context.Context) string {
	// 尝试从gin.Context获取
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if userID, exists := ginCtx.Get("user_id"); exists {
			if id, ok := userID.(string); ok {
				return id
			}
		}
	}

	// 尝试从标准context.Context获取
	if userID := ctx.Value("user_id"); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}

	return ""
}

// GetUsernameFromContext 从context中获取当前用户名
func GetUsernameFromContext(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if username, exists := ginCtx.Get("username"); exists {
			if name, ok := username.(string); ok {
				return name
			}
		}
	}

	if username := ctx.Value("username"); username != nil {
		if name, ok := username.(string); ok {
			return name
		}
	}

	return ""
}

// GetUserRoleFromContext 从context中获取当前用户角色
func GetUserRoleFromContext(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if role, exists := ginCtx.Get("role"); exists {
			if r, ok := role.(string); ok {
				return r
			}
		}
	}

	if role := ctx.Value("role"); role != nil {
		if r, ok := role.(string); ok {
			return r
		}
	}

	return ""
}

// WithTraceID 将链路追踪ID添加到Context中
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceIDFromContext 从Context中获取链路追踪ID
func GetTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// ContextInfo 上下文信息结构体
type ContextInfo struct {
	Language  string
	RequestID string
	UserID    string
	TraceID   string
}

// GetContextInfo 获取Context中的所有信息
func GetContextInfo(ctx context.Context) *ContextInfo {
	return &ContextInfo{
		Language:  GetLanguageFromContext(ctx),
		RequestID: GetRequestIDFromContext(ctx),
		UserID:    GetUserIDFromContext(ctx),
		TraceID:   GetTraceIDFromContext(ctx),
	}
}

// WithContextInfo 将ContextInfo添加到Context中
func WithContextInfo(ctx context.Context, info *ContextInfo) context.Context {
	if info == nil {
		return ctx
	}

	ctx = WithLanguage(ctx, info.Language)
	ctx = WithRequestID(ctx, info.RequestID)
	ctx = WithUserID(ctx, info.UserID)
	ctx = WithTraceID(ctx, info.TraceID)

	return ctx
}

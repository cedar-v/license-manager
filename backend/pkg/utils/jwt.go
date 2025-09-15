package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"license-manager/internal/config"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenNotValid    = errors.New("token is not valid")
	ErrTokenMalformed   = errors.New("token is malformed")
	ErrTokenNotFound    = errors.New("token not found")
)

// GenerateToken 生成JWT Token
func GenerateToken(userID string, username, role string) (string, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return "", errors.New("config not initialized")
	}

	expireTime := time.Now().Add(time.Duration(cfg.Auth.JWT.ExpireHours) * time.Hour)
	
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "license-manager",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Auth.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, errors.New("config not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.Auth.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrTokenNotValid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenNotValid
}

// ValidateToken 验证Token是否有效
func ValidateToken(tokenString string) (*Claims, error) {
	return ParseToken(tokenString)
}

// RefreshTokenIfNeeded 检查Token是否需要刷新，如果需要则返回新Token
func RefreshTokenIfNeeded(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	cfg := config.GetConfig()
	if cfg == nil {
		return "", errors.New("config not initialized")
	}

	// 检查是否需要刷新（剩余时间少于阈值）
	thresholdDuration := time.Duration(cfg.Auth.JWT.RefreshThresholdMinutes) * time.Minute
	if time.Until(claims.ExpiresAt.Time) < thresholdDuration {
		// 生成新Token
		return GenerateToken(claims.UserID, claims.Username, claims.Role)
	}

	// 不需要刷新，返回原Token
	return tokenString, nil
}
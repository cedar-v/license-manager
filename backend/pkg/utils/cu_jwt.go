package utils

import (
	"errors"
	"time"

	"license-manager/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type CuClaims struct {
	UserID     string `json:"user_id"`
	CustomerID string `json:"customer_id"`
	UserRole   string `json:"user_role"`
	Phone      string `json:"phone"`
	jwt.RegisteredClaims
}

var (
	ErrCuTokenExpired   = errors.New("cu token has expired")
	ErrCuTokenNotValid  = errors.New("cu token is not valid")
	ErrCuTokenMalformed = errors.New("cu token is malformed")
	ErrCuTokenNotFound  = errors.New("cu token not found")
)

// GenerateCuToken 生成客户用户JWT Token
func GenerateCuToken(userID, customerID, userRole, phone string) (string, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return "", errors.New("config not initialized")
	}

	expireTime := time.Now().Add(time.Duration(cfg.Auth.CuJWT.ExpireHours) * time.Hour)

	claims := &CuClaims{
		UserID:     userID,
		CustomerID: customerID,
		UserRole:   userRole,
		Phone:      phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "license-manager-cu",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Auth.CuJWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseCuToken 解析客户用户JWT Token
func ParseCuToken(tokenString string) (*CuClaims, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, errors.New("config not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &CuClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.Auth.CuJWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrCuTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrCuTokenMalformed
		}
		return nil, ErrCuTokenNotValid
	}

	if claims, ok := token.Claims.(*CuClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrCuTokenNotValid
}

// ValidateCuToken 验证客户用户Token是否有效
func ValidateCuToken(tokenString string) (*CuClaims, error) {
	return ParseCuToken(tokenString)
}

// RefreshCuTokenIfNeeded 检查客户用户Token是否需要刷新，如果需要则返回新Token
func RefreshCuTokenIfNeeded(tokenString string) (string, error) {
	claims, err := ParseCuToken(tokenString)
	if err != nil {
		return "", err
	}

	cfg := config.GetConfig()
	if cfg == nil {
		return "", errors.New("config not initialized")
	}

	// 检查是否需要刷新（剩余时间少于阈值）
	thresholdDuration := time.Duration(cfg.Auth.CuJWT.RefreshThresholdMinutes) * time.Minute
	if time.Until(claims.ExpiresAt.Time) < thresholdDuration {
		// 生成新Token
		return GenerateCuToken(claims.UserID, claims.CustomerID, claims.UserRole, claims.Phone)
	}

	// 不需要刷新，返回原Token
	return tokenString, nil
}

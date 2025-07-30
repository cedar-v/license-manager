package service

import (
	"runtime"
	"time"

	"license-manager/internal/models"
)

type systemService struct {
	startTime time.Time
}

// NewSystemService 创建系统服务实例
func NewSystemService() SystemService {
	return &systemService{
		startTime: time.Now(),
	}
}

// GetHealthStatus 获取健康状态
func (s *systemService) GetHealthStatus() *models.HealthResponse {
	uptime := time.Since(s.startTime)

	return &models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
		System: models.SystemInfo{
			OS:           runtime.GOOS,
			Arch:         runtime.GOARCH,
			Version:      runtime.Version(),
			NumCPU:       runtime.NumCPU(),
			NumGoroutine: runtime.NumGoroutine(),
		},
		Services: map[string]string{
			"database": "healthy", // 这里可以检查数据库连接状态
			"cache":    "healthy", // 这里可以检查缓存状态
		},
	}
}

// GetSystemInfo 获取系统详细信息
func (s *systemService) GetSystemInfo() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"version": "1.0.0",
		"uptime":  time.Since(s.startTime).String(),
		"system": map[string]interface{}{
			"os":            runtime.GOOS,
			"arch":          runtime.GOARCH,
			"version":       runtime.Version(),
			"num_cpu":       runtime.NumCPU(),
			"num_goroutine": runtime.NumGoroutine(),
		},
		"memory": map[string]interface{}{
			"alloc":       bToMb(m.Alloc),
			"total_alloc": bToMb(m.TotalAlloc),
			"sys":         bToMb(m.Sys),
			"num_gc":      m.NumGC,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}
}

// bToMb converts bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
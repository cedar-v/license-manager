package models

// HealthResponse 健康检测响应结构
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	System    SystemInfo        `json:"system"`
	Services  map[string]string `json:"services"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	Version      string `json:"version"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
}
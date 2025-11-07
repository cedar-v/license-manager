package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// ClientConfig 客户端配置
type ClientConfig struct {
	ServerURL         string `json:"server_url"`         // 服务器地址
	AuthorizationCode string `json:"authorization_code"` // 授权码
	SoftwareVersion   string `json:"software_version"`   // 软件版本
	HeartbeatInterval int    `json:"heartbeat_interval"` // 心跳间隔(秒)，默认300
	LicenseKey        string `json:"license_key"`        // 许可证密钥（激活后保存）
	ConfigUpdatedAt   string `json:"config_updated_at"`  // 配置更新时间
	LicenseFile       string `json:"license_file"`       // 许可证文件内容
}

// ActivateRequest 激活请求
type ActivateRequest struct {
	AuthorizationCode   string                 `json:"authorization_code"`
	HardwareFingerprint string                 `json:"hardware_fingerprint"`
	DeviceInfo          map[string]interface{} `json:"device_info,omitempty"`
	SoftwareVersion     *string                `json:"software_version,omitempty"`
}

// ActivateResponse 激活响应
type ActivateResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		LicenseKey        string `json:"license_key"`
		LicenseFile       string `json:"license_file"`
		HeartbeatInterval int    `json:"heartbeat_interval"`
	} `json:"data"`
}

// HeartbeatRequest 心跳请求
type HeartbeatRequest struct {
	LicenseKey          string                 `json:"license_key"`
	HardwareFingerprint string                 `json:"hardware_fingerprint"`
	ConfigUpdatedAt     *string                `json:"config_updated_at,omitempty"`
	UsageData           map[string]interface{} `json:"usage_data,omitempty"`
	SoftwareVersion     *string                `json:"software_version,omitempty"`
}

// HeartbeatResponse 心跳响应
type HeartbeatResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Status            string  `json:"status"`
		ConfigUpdated     bool    `json:"config_updated"`
		LicenseFile       *string `json:"license_file,omitempty"`
		HeartbeatInterval int     `json:"heartbeat_interval"`
	} `json:"data"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// LicenseFileData 许可证文件数据结构
type LicenseFileData struct {
	LicenseKey          string                 `json:"license_key"`
	AuthorizationCodeID string                 `json:"authorization_code_id"`
	HardwareFingerprint string                 `json:"hardware_fingerprint"`
	Status              string                 `json:"status"`
	ActivatedAt         *string                `json:"activated_at"`
	ConfigUpdatedAt     *string                `json:"config_updated_at"`
	GeneratedAt         string                 `json:"generated_at"`
	AuthorizationCode   *string                `json:"authorization_code"`
	StartDate           *string                `json:"start_date"`
	EndDate             *string                `json:"end_date"`
	DeploymentType      *string                `json:"deployment_type"`
	MaxActivations      *int                   `json:"max_activations"`
	FeatureConfig       map[string]interface{} `json:"feature_config,omitempty"`
	UsageLimits         map[string]interface{} `json:"usage_limits,omitempty"`
	CustomParameters    map[string]interface{} `json:"custom_parameters,omitempty"`
}

const (
	configFile      = "client_config.json"
	licenseDir      = "license_code"
	licenseFileName = "LICENSE"
)

var (
	config ClientConfig
)

func main() {
	var (
		serverURL       = flag.String("server", "http://localhost:18888", "服务器地址 (例如: http://localhost:18888)")
		softwareVersion = flag.String("version", "1.0.0", "软件版本")
		interval        = flag.Int("interval", 300, "心跳间隔(秒)")
		activateOnly    = flag.Bool("activate-only", false, "仅执行激活，不启动心跳")
	)
	flag.Parse()

	// 加载配置
	if err := loadConfig(); err != nil {
		log.Printf("加载配置失败，使用命令行参数: %v", err)
	}

	// 命令行参数优先
	if *serverURL != "" {
		config.ServerURL = *serverURL
	}
	// 从 license_code/AUTH_CODE 文件读取授权码
	authCodeFromFile, err := readAuthCodeFromFile()
	if err == nil && authCodeFromFile != "" {
		config.AuthorizationCode = strings.TrimSpace(authCodeFromFile)
	}
	if *softwareVersion != "" {
		config.SoftwareVersion = *softwareVersion
	}
	if *interval > 0 {
		config.HeartbeatInterval = *interval
	}
	if config.HeartbeatInterval <= 0 {
		config.HeartbeatInterval = 300 // 默认5分钟
	}

	// 验证必要参数
	if config.ServerURL == "" {
		log.Fatal("请提供服务器地址 (-server 或配置文件)")
	}

	log.Println("=== 软件授权客户端示例程序 ===")
	log.Printf("服务器地址: %s", config.ServerURL)
	log.Printf("软件版本: %s", config.SoftwareVersion)

	// 初始化RSA公钥（从当前目录的 rsa_public_key.pem 加载）
	if err := initRSAPublicKey(); err != nil {
		log.Fatalf("初始化RSA公钥失败: %v", err)
	}
	log.Println("✓ RSA公钥加载成功")

	// 采集硬件指纹
	hardwareFingerprint, deviceInfo := collectHardwareInfo()
	log.Printf("硬件指纹: %s", hardwareFingerprint)
	log.Printf("设备信息: CPU=%s, Memory=%s, OS=%s",
		deviceInfo["cpu"], deviceInfo["memory"], deviceInfo["os"])

	// ===== 新的验证流程 =====
	// 1. 检查 license_code/LICENSE 文件是否存在
	licenseFilePath := filepath.Join(licenseDir, licenseFileName)
	licenseData, err := checkAndParseLicenseFile(licenseFilePath)
	if err != nil {
		log.Printf("许可证文件检查失败: %v", err)
		log.Println("尝试重新激活...")
		licenseData = nil // 标记需要重新激活
	}

	// 2. 如果许可证文件存在，验证时间
	needActivation := false
	if licenseData != nil {
		isValid, err := validateLicenseTime(licenseData)
		if err != nil {
			log.Printf("许可证时间验证失败: %v", err)
			needActivation = true
		} else if !isValid {
			log.Println("许可证已过期或未生效，需要重新激活")
			needActivation = true
		} else {
			log.Println("✓ 许可证验证通过")
			// 从许可证文件加载许可证密钥
			config.LicenseKey = licenseData.LicenseKey
			config.LicenseFile = "" // 可以从文件读取
		}
	} else {
		needActivation = true
	}

	// 3. 如果需要激活，使用AUTH_CODE重新获取许可证
	if needActivation {
		if config.AuthorizationCode == "" {
			log.Fatal("需要授权码进行激活，请将授权码保存到 license_code/AUTH_CODE 文件")
		}

		log.Println("\n[激活] 开始激活授权...")
		if err := activateLicense(hardwareFingerprint, deviceInfo); err != nil {
			log.Fatalf("激活失败: %v", err)
		}
		log.Println("✓ 激活成功！")
		log.Printf("许可证密钥: %s", config.LicenseKey)
		log.Printf("心跳间隔: %d 秒", config.HeartbeatInterval)

		// 保存许可证文件到 license_code/LICENSE
		if err := saveLicenseFile(licenseFilePath); err != nil {
			log.Printf("保存许可证文件失败: %v", err)
		} else {
			log.Printf("许可证文件已保存到 %s", licenseFilePath)
		}

		// 保存配置
		if err := saveConfig(); err != nil {
			log.Printf("保存配置失败: %v", err)
		} else {
			log.Println("配置已保存到 " + configFile)
		}

		// 再次验证许可证文件
		licenseData, err = checkAndParseLicenseFile(licenseFilePath)
		if err != nil {
			log.Fatalf("激活后验证许可证文件失败: %v", err)
		}
		isValid, err := validateLicenseTime(licenseData)
		if err != nil || !isValid {
			log.Fatalf("激活后验证许可证时间失败: %v", err)
		}
		log.Println("✓ 激活后验证通过")
	}

	if *activateOnly {
		log.Println("激活完成，退出程序")
		return
	}

	// 4. 验证通过，启动心跳服务
	log.Println("\n[心跳] 启动心跳服务...")
	startHeartbeat(hardwareFingerprint, deviceInfo)
}

// collectHardwareInfo 采集硬件信息
func collectHardwareInfo() (string, map[string]interface{}) {
	var parts []string

	// 1. 获取 MAC 地址
	if macs, err := getMACAddresses(); err == nil && len(macs) > 0 {
		parts = append(parts, fmt.Sprintf("MAC:%s", macs[0]))
	}

	// 2. 获取 CPU 信息
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		cpuModel := strings.TrimSpace(cpuInfo[0].ModelName)
		if cpuModel != "" {
			parts = append(parts, fmt.Sprintf("CPU:%s", cpuModel))
		}
	}

	// 3. 获取主板序列号（Windows）或机器ID（Linux）
	if hostInfo, err := host.Info(); err == nil {
		if hostInfo.HostID != "" {
			parts = append(parts, fmt.Sprintf("HOSTID:%s", hostInfo.HostID))
		}
	}

	// 4. 如果没有足够信息，使用操作系统信息
	if len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("OS:%s-%s", runtime.GOOS, runtime.GOARCH))
	}

	hardwareFingerprint := strings.Join(parts, ",")

	// 收集详细的设备信息
	deviceInfo := make(map[string]interface{})

	// CPU 信息
	if cpuInfo, err := cpu.Info(); err == nil && len(cpuInfo) > 0 {
		deviceInfo["cpu"] = cpuInfo[0].ModelName
		deviceInfo["cpu_cores"] = runtime.NumCPU()
	}

	// 内存信息
	if memInfo, err := mem.VirtualMemory(); err == nil {
		deviceInfo["memory"] = fmt.Sprintf("%d GB", memInfo.Total/(1024*1024*1024))
	}

	// 操作系统信息
	if hostInfo, err := host.Info(); err == nil {
		deviceInfo["os"] = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)
	}

	deviceInfo["arch"] = runtime.GOARCH

	return hardwareFingerprint, deviceInfo
}

// getMACAddresses 获取 MAC 地址
func getMACAddresses() ([]string, error) {
	var macs []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// 跳过虚拟接口和回环接口
		if iface.Flags&net.FlagLoopback != 0 ||
			iface.Flags&net.FlagUp == 0 ||
			strings.HasPrefix(iface.Name, "veth") ||
			strings.HasPrefix(iface.Name, "docker") {
			continue
		}

		mac := iface.HardwareAddr.String()
		if mac != "" {
			macs = append(macs, mac)
		}
	}

	return macs, nil
}

// activateLicense 激活授权
func activateLicense(hardwareFingerprint string, deviceInfo map[string]interface{}) error {
	req := ActivateRequest{
		AuthorizationCode:   config.AuthorizationCode,
		HardwareFingerprint: hardwareFingerprint,
		DeviceInfo:          deviceInfo,
	}

	if config.SoftwareVersion != "" {
		req.SoftwareVersion = &config.SoftwareVersion
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/activate", config.ServerURL)
	log.Printf("POST %s", url)
	log.Printf("请求体: %s", string(jsonData))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return fmt.Errorf("激活失败 [%s]: %s", errResp.Code, errResp.Message)
		}
		return fmt.Errorf("激活失败 [HTTP %d]: %s", resp.StatusCode, string(body))
	}

	var activateResp ActivateResponse
	if err := json.Unmarshal(body, &activateResp); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if activateResp.Code != "000000" {
		return fmt.Errorf("激活失败 [%s]: %s", activateResp.Code, activateResp.Message)
	}

	// 更新配置
	config.LicenseKey = activateResp.Data.LicenseKey
	config.LicenseFile = activateResp.Data.LicenseFile
	if activateResp.Data.HeartbeatInterval > 0 {
		config.HeartbeatInterval = activateResp.Data.HeartbeatInterval
	}
	config.ConfigUpdatedAt = time.Now().Format(time.RFC3339)

	return nil
}

// startHeartbeat 启动心跳服务
func startHeartbeat(hardwareFingerprint string, deviceInfo map[string]interface{}) {
	ticker := time.NewTicker(time.Duration(config.HeartbeatInterval) * time.Second)
	defer ticker.Stop()

	// 立即发送一次心跳
	log.Println("发送初始心跳...")
	sendHeartbeat(hardwareFingerprint, deviceInfo)

	// 定期发送心跳
	for range ticker.C {
		sendHeartbeat(hardwareFingerprint, deviceInfo)
	}
}

// sendHeartbeat 发送心跳
func sendHeartbeat(hardwareFingerprint string, deviceInfo map[string]interface{}) {
	req := HeartbeatRequest{
		LicenseKey:          config.LicenseKey,
		HardwareFingerprint: hardwareFingerprint,
	}

	if config.ConfigUpdatedAt != "" {
		req.ConfigUpdatedAt = &config.ConfigUpdatedAt
	}

	// 模拟使用数据
	req.UsageData = map[string]interface{}{
		"active_users":    50,
		"api_calls_today": 5000,
		"cpu_usage":       25.5,
		"memory_usage":    45.2,
	}

	if config.SoftwareVersion != "" {
		req.SoftwareVersion = &config.SoftwareVersion
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Printf("✗ 心跳失败: 序列化请求失败: %v", err)
		return
	}

	// 打印完整的请求信息
	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("心跳请求信息:\n%s", string(reqJSON))

	url := fmt.Sprintf("%s/api/v1/heartbeat", config.ServerURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("✗ 心跳失败: 请求失败: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("✗ 心跳失败: 读取响应失败: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			log.Printf("✗ 心跳失败 [%s]: %s", errResp.Code, errResp.Message)
		} else {
			log.Printf("✗ 心跳失败 [HTTP %d]: %s", resp.StatusCode, string(body))
		}
		return
	}

	var heartbeatResp HeartbeatResponse
	if err := json.Unmarshal(body, &heartbeatResp); err != nil {
		log.Printf("✗ 心跳失败: 解析响应失败: %v", err)
		return
	}

	if heartbeatResp.Code != "000000" {
		log.Printf("✗ 心跳失败 [%s]: %s", heartbeatResp.Code, heartbeatResp.Message)
		return
	}

	// 打印完整的响应信息
	respJSON, _ := json.MarshalIndent(heartbeatResp, "", "  ")
	log.Printf("心跳返回信息:\n%s", string(respJSON))

	// 更新配置更新时间
	if heartbeatResp.Data.ConfigUpdated && heartbeatResp.Data.LicenseFile != nil {
		log.Println("✓ 收到配置更新，正在更新许可证文件...")
		config.LicenseFile = *heartbeatResp.Data.LicenseFile
		config.ConfigUpdatedAt = time.Now().Format(time.RFC3339)

		// 保存许可证文件到 license_code/LICENSE
		licenseFilePath := filepath.Join(licenseDir, licenseFileName)
		if err := saveLicenseFile(licenseFilePath); err != nil {
			log.Printf("保存许可证文件失败: %v", err)
		} else {
			log.Printf("许可证文件已更新到 %s", licenseFilePath)
		}

		// 保存配置
		if err := saveConfig(); err != nil {
			log.Printf("保存配置失败: %v", err)
		}
	}

	log.Printf("✓ 心跳成功 [%s] - 状态: %s, 配置更新: %v, 间隔: %d秒",
		time.Now().Format("15:04:05"),
		heartbeatResp.Data.Status,
		heartbeatResp.Data.ConfigUpdated,
		heartbeatResp.Data.HeartbeatInterval)

	// 如果服务器返回了不同的心跳间隔，更新配置
	if heartbeatResp.Data.HeartbeatInterval > 0 &&
		heartbeatResp.Data.HeartbeatInterval != config.HeartbeatInterval {
		log.Printf("更新心跳间隔: %d 秒 -> %d 秒",
			config.HeartbeatInterval, heartbeatResp.Data.HeartbeatInterval)
		config.HeartbeatInterval = heartbeatResp.Data.HeartbeatInterval
	}
}

// loadConfig 加载配置
func loadConfig() error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &config)
}

// saveConfig 保存配置
func saveConfig() error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}

// checkAndParseLicenseFile 检查并解析许可证文件
func checkAndParseLicenseFile(filePath string) (*LicenseFileData, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("许可证文件不存在: %s", filePath)
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取许可证文件失败: %w", err)
	}

	// 尝试解密许可证文件
	decryptedData, err := decryptLicenseFile(data)
	if err != nil {
		return nil, fmt.Errorf("解密许可证文件失败: %w", err)
	}

	// 解析JSON
	var licenseData LicenseFileData
	if err := json.Unmarshal(decryptedData, &licenseData); err != nil {
		return nil, fmt.Errorf("解析许可证文件JSON失败: %w", err)
	}

	// 打印解析后的许可证信息
	log.Println("========== 解析后的许可证信息 ==========")
	licenseJSON, _ := json.MarshalIndent(licenseData, "", "  ")
	log.Printf("%s\n", string(licenseJSON))
	log.Println("========================================")

	return &licenseData, nil
}

// validateLicenseTime 验证许可证时间
func validateLicenseTime(licenseData *LicenseFileData) (bool, error) {
	now := time.Now()

	// 检查开始日期
	if licenseData.StartDate != nil && *licenseData.StartDate != "" {
		startDate, err := time.Parse(time.RFC3339, *licenseData.StartDate)
		if err != nil {
			// 尝试其他时间格式
			startDate, err = time.Parse("2006-01-02T15:04:05Z07:00", *licenseData.StartDate)
			if err != nil {
				startDate, err = time.Parse("2006-01-02 15:04:05", *licenseData.StartDate)
				if err != nil {
					return false, fmt.Errorf("解析开始日期失败: %w", err)
				}
			}
		}
		if now.Before(startDate) {
			return false, fmt.Errorf("许可证尚未生效，生效日期: %s", startDate.Format("2006-01-02 15:04:05"))
		}
	}

	// 检查结束日期
	if licenseData.EndDate != nil && *licenseData.EndDate != "" {
		endDate, err := time.Parse(time.RFC3339, *licenseData.EndDate)
		if err != nil {
			// 尝试其他时间格式
			endDate, err = time.Parse("2006-01-02T15:04:05Z07:00", *licenseData.EndDate)
			if err != nil {
				endDate, err = time.Parse("2006-01-02 15:04:05", *licenseData.EndDate)
				if err != nil {
					return false, fmt.Errorf("解析结束日期失败: %w", err)
				}
			}
		}
		if now.After(endDate) {
			return false, fmt.Errorf("许可证已过期，过期日期: %s", endDate.Format("2006-01-02 15:04:05"))
		}
		startDateStr := ""
		if licenseData.StartDate != nil {
			startDateStr = *licenseData.StartDate
		}
		endDateStr := ""
		if licenseData.EndDate != nil {
			endDateStr = *licenseData.EndDate
		}
		log.Printf("许可证有效期: %s - %s", startDateStr, endDateStr)
	} else {
		return false, fmt.Errorf("许可证文件缺少结束日期")
	}

	// 检查状态
	if licenseData.Status == "revoked" {
		return false, fmt.Errorf("许可证已被撤销")
	}

	return true, nil
}

// LicenseFileWithSignature RSA签名后的许可证文件结构
type LicenseFileWithSignature struct {
	Data      string `json:"data"`      // 原始数据（JSON字符串）
	Signature string `json:"signature"` // 数字签名
	Algorithm string `json:"algorithm"` // 签名算法
}

// decryptLicenseFile 验证RSA签名的许可证文件
func decryptLicenseFile(encryptedData []byte) ([]byte, error) {
	// Base64解码
	decoded, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %w", err)
	}

	// 解析签名后的许可证文件
	var signedLicense LicenseFileWithSignature
	if err := json.Unmarshal(decoded, &signedLicense); err != nil {
		return nil, fmt.Errorf("解析许可证文件失败: %w", err)
	}

	// 检查签名算法
	if signedLicense.Algorithm != "RSA-PSS-SHA256" {
		return nil, fmt.Errorf("不支持的签名算法: %s", signedLicense.Algorithm)
	}

	// 验证签名
	dataBytes := []byte(signedLicense.Data)
	if err := verifyRSASignature(dataBytes, signedLicense.Signature); err != nil {
		return nil, fmt.Errorf("签名验证失败: %w", err)
	}

	// 返回原始数据
	return dataBytes, nil
}

// saveLicenseFile 保存许可证文件到 license_code/LICENSE
func saveLicenseFile(filePath string) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 如果config.LicenseFile不为空，直接保存
	if config.LicenseFile != "" {
		// 写入文件（已经是base64编码的加密数据）
		return os.WriteFile(filePath, []byte(config.LicenseFile), 0644)
	}

	// 如果LicenseFile为空，但LicenseKey存在，可能需要重新生成
	// 这里我们假设LicenseFile已经在激活时设置
	return fmt.Errorf("许可证文件内容为空")
}

// readAuthCodeFromFile 从 license_code/AUTH_CODE 文件读取授权码
func readAuthCodeFromFile() (string, error) {
	authCodePath := filepath.Join(licenseDir, "AUTH_CODE")
	data, err := os.ReadFile(authCodePath)
	if err != nil {
		return "", fmt.Errorf("读取授权码文件失败: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

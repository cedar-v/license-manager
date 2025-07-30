package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"license-manager/internal/config"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// CustomFormatter 自定义格式器，在消息末尾显示文件路径和行号
type CustomFormatter struct {
	*logrus.TextFormatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 先用标准格式器格式化
	formatted, err := f.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	// 如果有调用者信息，添加到末尾
	if entry.Caller != nil {
		filename := entry.Caller.File
		if idx := strings.Index(filename, "backend/"); idx != -1 {
			filename = "./backend/" + filename[idx+len("backend/"):]
		} else {
			filename = filepath.Base(filename)
		}

		// 移除换行符，添加文件信息，然后重新添加换行符
		formattedStr := strings.TrimRight(string(formatted), "\n")
		formattedStr = fmt.Sprintf("%s %s:%d\n", formattedStr, filename, entry.Caller.Line)
		return []byte(formattedStr), nil
	}

	return formatted, nil
}

func Init() {
	Logger = logrus.New()
	Logger.SetReportCaller(true) // 启用调用者信息

	cfg := config.GetConfig()
	if cfg != nil {
		// 设置日志级别
		switch cfg.Log.Level {
		case "debug":
			Logger.SetLevel(logrus.DebugLevel)
		case "info":
			Logger.SetLevel(logrus.InfoLevel)
		case "warn":
			Logger.SetLevel(logrus.WarnLevel)
		case "error":
			Logger.SetLevel(logrus.ErrorLevel)
		default:
			Logger.SetLevel(logrus.InfoLevel)
		}

		// 设置自定义格式器
		Logger.SetFormatter(&CustomFormatter{
			TextFormatter: &logrus.TextFormatter{
				FullTimestamp:    true,
				TimestampFormat:  "2006-01-02 15:04:05",
				ForceColors:      true,
				DisableColors:    false,
				DisableLevelTruncation: true,
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					return "", "" // 禁用默认的调用者信息显示
				},
			},
		})
	} else {
		// 默认配置
		Logger.SetLevel(logrus.InfoLevel)
		Logger.SetFormatter(&CustomFormatter{
			TextFormatter: &logrus.TextFormatter{
				FullTimestamp:    true,
				TimestampFormat:  "2006-01-02 15:04:05",
				ForceColors:      true,
				DisableLevelTruncation: true,
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					return "", "" // 禁用默认的调用者信息显示
				},
			},
		})
	}

	Logger.SetOutput(os.Stdout)
}

func GetLogger() *logrus.Logger {
	if Logger == nil {
		Init()
	}
	return Logger
}

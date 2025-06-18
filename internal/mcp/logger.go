package mcp

import (
	"log/slog"
	"os"
	"path/filepath"
)

var logger *slog.Logger

func InitLogger() {
	if os.Getenv("LOGGER") != "true" {
		logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
		return
	}

	// 确保home目录存在
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	logPath := filepath.Join(homeDir, "kaf-mcp.log")

	// 创建/追加日志文件
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	logger = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

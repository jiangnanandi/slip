package utils

import (
    "fmt"
    "os"
    "path/filepath"
    "time"

    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
    "github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger(logDir string) error {
    Log = logrus.New()
    
    // 设置日志格式
    Log.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
    })

    // 确保日志目录存在
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return fmt.Errorf("创建日志目录失败: %v", err)
    }

    // 配置 rotatelogs
    path := filepath.Join(logDir, "slip.log")
    writer, err := rotatelogs.New(
        filepath.Join(logDir, "slip.%Y%m%d.log"),
        rotatelogs.WithLinkName(path),           // 生成软链，指向最新日志文件
        rotatelogs.WithMaxAge(30*24*time.Hour),  // 保存30天
        rotatelogs.WithRotationTime(24*time.Hour), // 每天切割
    )
    if err != nil {
        return fmt.Errorf("配置日志轮转失败: %v", err)
    }

    Log.SetOutput(writer)
    return nil
}
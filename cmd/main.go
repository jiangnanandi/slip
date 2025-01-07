package main

import (
	"fmt"
	"slip/internal/config"
	slip "slip/internal/handler"
	initialize "slip/internal/pkg"
	"slip/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"context"
	"slip/internal/pkg/utils"
	"path/filepath"
)

func main() {
	// 加载持久化队列
	ctx := context.Background()
	go utils.ProcessRetryQueue(ctx)

	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	err = initialize.InitializeDirectories()
	if err != nil {
		panic(err)
	}

	logDir := filepath.Join(config.AppConfig.LogDir)
    if err := utils.InitLogger(logDir); err != nil {
        panic(fmt.Sprintf("初始化日志失败: %v", err))
    }

	r := gin.Default()

	router.InitRouter(r)

	c := cron.New()

	c.AddFunc("*/1 * * * *", func() {
		slip.BuildIndex()
	})

	c.Start()

	r.Run(":8084")
}

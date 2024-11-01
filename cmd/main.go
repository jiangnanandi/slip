package main

import (
	"slip/internal/config"
	slip "slip/internal/handler"
	initialize "slip/internal/pkg"
	"slip/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	err = initialize.InitializeDirectories()
	if err != nil {
		panic(err)
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

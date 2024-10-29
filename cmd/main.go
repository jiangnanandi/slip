package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	slip "slip/internal/handler"
	"slip/internal/router"
	"slip/internal/config"
)

func main() {
	err := config.LoadConfig()
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

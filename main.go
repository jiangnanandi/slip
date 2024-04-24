package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"net/http"
	slip "slip/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World 22",
		})
	})
	r.POST("/send-notes", slip.CreateNote)
	r.GET("/index", slip.Index)
	r.GET("/notes/:title", slip.GetNote)

	// Create a new cron instance
	c := cron.New()

	// Define the task to be executed and its schedule
	c.AddFunc("1 * * * *", func() {
		// Code to be executed every minute
		// Add your desired functionality here
		slip.BuildIndex()
		println("Running timed task...")
	})

	// Start the cron scheduler
	c.Start()

	r.Run(":8084")
}

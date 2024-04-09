package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	slip "slip/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	r.POST("/send-notes", slip.CreateNote)
	r.Run(":8084")
}

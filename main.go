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
			"message": "Hello World 22",
		})
	})
	r.POST("/send-notes", slip.CreateNote)
	r.GET("/index", slip.Index)
	r.GET("/notes/:title", slip.GetNote)
	r.Run(":8084")
}
package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"slip/internal/handler"
	"slip/internal/middleware"
)

func InitRouter(r *gin.Engine) {
	// 公开路由
	r.GET("/login", controller.Login)
	r.GET("/index", controller.Index)
	r.GET("/notes/:title", controller.GetNote)

	// 受保护的路由
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello World!",
			})
		})
		protected.POST("/send-notes", controller.CreateNote)
	}
}

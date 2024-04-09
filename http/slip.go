package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"slip/service"
	"slip/utils"
)

/* 处理请求，将传入的「笔记内容」存储到指定目录，并保存成 title.md 格式的笔记文件 */
func CreateNote(c *gin.Context) {
	var note utils.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 保存笔记内容到指定目录
	if err := service.SaveNote(note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Note created successfully",
	})
}

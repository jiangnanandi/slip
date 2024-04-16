package http

import (
	"net/http"
	"slip/service"
	"slip/utils"

	"github.com/gin-gonic/gin"
)

/* 处理请求，将传入的「笔记内容」存储到指定目录，并保存成 title.md 格式的笔记文件 */
func CreateNote(ctx *gin.Context) {
	var note utils.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 保存笔记内容到指定目录
	if err := service.SaveNote(ctx, note); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Note created successfully",
	})
}

func Index(ctx *gin.Context) {
	fileContent, err := service.Index(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Data(http.StatusOK, "text/html", fileContent)
}

func GetNote(ctx *gin.Context) {

	finalHTML, err := service.GetNote(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Data(http.StatusOK, "text/html", finalHTML)
}

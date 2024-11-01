package controller

import (
	"net/http"
	"slip/api/defines"
	"slip/internal/service"

	"github.com/gin-gonic/gin"
)

/* 处理请求，将传入的「笔记内容」存储到指定目录，并保存成 title.md 格式的笔记文件 */
func CreateNote(ctx *gin.Context) {
	var note defines.Notes
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := note.Build()
	if err != nil {
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
	title := ctx.Param("title")
	finalHTML, err := service.GetNote(ctx, title)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Data(http.StatusOK, "text/html", finalHTML)
}

func BuildIndex() error {
	err := service.BuildIndex()
	if err != nil {
		return err
	}
	return nil
}

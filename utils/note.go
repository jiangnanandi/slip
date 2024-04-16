package utils

import (
	"os"
)

func WriteNote(noteDir string, note Note) error {
	file, err := os.Create(noteDir + "/" + note.Title + ".md")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(note.Body)
	if err != nil {
		return err
	}

	// 生成 index.html 文件，查找 noteDir 目录下的所有文件，列出来（在 HTML 中展示笔记链接）
	indexFile, err := os.Create(noteDir + "/index.html")
	if err != nil {
		return err
	}
	defer indexFile.Close()

	_, err = indexFile.WriteString("<html><head><title>王掌柜的小纸条</title></head><body><h1>小纸条</h1><ul>")
	if err != nil {
		return err
	}

	files, err := os.ReadDir(noteDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			// 获取文件名，不需要扩展名
			fileName := file.Name()
			// 去掉扩展名
			fileName = fileName[:len(fileName)-3]
			_, err = indexFile.WriteString("<li><a href=\"notes\\" + fileName + "\">" + fileName + "</a></li>")
			if err != nil {
				return err
			}
		}
	}

	_, err = indexFile.WriteString("</ul></body></html>")
	if err != nil {
		return err
	}

	return nil
}

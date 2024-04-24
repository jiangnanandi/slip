package service

import (
	"bytes"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"html/template"
	"os"
	"slip/utils"

	"github.com/gin-gonic/gin"
)

func SaveNote(ctx *gin.Context, note utils.Note) error {
	// 将 Note 内容写入到指定目录
	if err := utils.WriteNote(utils.DefaultNoteDir, note); err != nil {
		return err
	}
	return nil
}

func Index(ctx *gin.Context) ([]byte, error) {
	// 读取 utils.DefaultNoteDir 目录下的 index.html 文件，并按照 html 标准输出，使请求方能够看到 index.html 内容
	indexFile := utils.DefaultNoteDir + "/index.html"

	files, err := os.ReadFile(indexFile)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func GetNote(ctx *gin.Context) ([]byte, error) {
	title := ctx.Param("title")
	noteFile := utils.DefaultNoteDir + "/" + title + ".md"
	fileContent, err := os.ReadFile(noteFile)
	if err != nil {
		return nil, err
	}

	// Create a new Markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// Parse the Markdown content to an AST
	doc := p.Parse(fileContent)

	// Create a new HTML renderer with options
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Render the AST to HTML
	htmlContent := markdown.Render(doc, renderer)

	// Embed the HTML content into the HTML template
	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Markdown Note Display</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            margin: 20px;
        }
        h1, h2, h3, h4, h5, h6 {
            margin-top: 20px;
            margin-bottom: 10px;
        }
        p {
            margin: 10px 0;
        }
        img {
            max-width: 100%%;
            height: auto;
        }
        pre {
            background-color: #f4f4f4;
            border: 1px solid #ddd;
            padding: 10px;
            overflow: auto;
        }
        code {
            background-color: #f4f4f4;
            padding: 2px 5px;
            border-radius: 3px;
        }
        blockquote {
            border-left: 3px solid #ccc;
            margin: 10px 0;
            padding-left: 20px;
            color: #666;
            font-style: italic;
        }
        ul, ol {
            margin: 10px 0;
            padding-left: 20px;
        }
    </style>
</head>
<body>
    {{.Content}}
</body>
</html>`

	tmpl, err := template.New("webpage").Funcs(
		template.FuncMap{
			"unescaped": func(htmlContent string) template.HTML {
				return template.HTML(htmlContent)
			},
		},
	).Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	data := struct {
		Content template.HTML
	}{
		Content: template.HTML(htmlContent),
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return nil, err
	}

	finalHTML := tpl.String()
	return []byte(finalHTML), nil
}

// Build index.html
func BuildIndex(noteDir string) error {
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

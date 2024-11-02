package service

import (
	"bytes"
	"html/template"
	"os"
	"slip/api/defines"
	"slip/internal/config"
	"slip/internal/pkg/utils"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/gin-gonic/gin"
)

func SaveNote(ctx *gin.Context, note defines.Notes) error {
	return utils.SaveNote(note);
}

func Index(ctx *gin.Context) ([]byte, error) {
	indexFile := config.AppConfig.Notes.PublishedDir + "/index.html"

	files, err := os.ReadFile(indexFile)
	if err != nil {
		return nil, err
	}
	// 设置响应头
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	return files, nil
}

func GetNote(ctx *gin.Context, title string) ([]byte, error) {
	noteFile := config.AppConfig.Notes.PublishedDir + "/" + title
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
            margin: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 80vh;
        }
		.container {
            max-width: 800px;
            padding: 20px;
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
    <div class="container">
        {{.Content}}
    </div>
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

// Build Published Notes Index
func BuildIndex() error {
	dir := config.AppConfig.Notes.PublishedDir
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	indexFile, err := os.Create(dir + "/index.html")
	if err != nil {
		return err
	}
	defer indexFile.Close()

	templateData := defines.TemplateData{
		Title: config.AppConfig.Title,
		Notes: []defines.Notes{},
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			title := strings.TrimSuffix(file.Name(), ".md")
			body, err := os.ReadFile(dir + "/" + file.Name())
			if err != nil {
				continue
			}

			note := defines.Notes{
				Title: title,
				Body:  string(body),
			}

			err = note.Build()
			if err != nil {
				continue
			}

			templateData.Notes = append(templateData.Notes, note)
		}
	}
	if len(templateData.Notes) == 0 {
		return nil
	}

	templateData.Sort()

	// 检查模板文件路径
	templatePath := "templates/index.html.tmpl"
	if _, err := os.Stat(templatePath); err != nil {
		return err
	}

	// 读取模板文件内容进行确认
	_, err = os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	// 解析模板
	tmpl, err := template.New("index.html.tmpl").ParseFiles(templatePath)
	if err != nil {
		return err
	}

	// 添加模板内容检查
	if tmpl.Tree == nil {
		return err
	}

	// 先尝试渲染到缓冲区
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "index.html.tmpl", templateData); err != nil {
		return err
	}
	
	// 写入文件
	if _, err := indexFile.Write(buf.Bytes()); err != nil {
		return err
	}

	// 确保文件写入完成
	if err := indexFile.Sync(); err != nil {
		return err
	}

	// 获取文件信息
	_, err = indexFile.Stat()
	if err != nil {
		return err
	}
	
	return nil
}
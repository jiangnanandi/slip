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

// Build index.html
func BuildIndex(noteDir string) error {
	// 生成 index.html 文件
	indexFile, err := os.Create(noteDir + "/index.html")
	if err != nil {
		return err
	}
	defer indexFile.Close()

	// 写入 HTML 头部和样式
	_, err = indexFile.WriteString(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
    <title>王掌柜的小纸条</title>
    <style>
        body {
            font-family: Arial, sans-serif;
			margin: 0 auto;
    		max-width: 800px;
        }
        h1 {
            color: #333;
        }
        .blog-list {
            list-style-type: none;
            padding: 0;
        }
        .blog-item {
            margin-bottom: 10px;
        }
        .blog-item a {
            color: #337ab7;
            text-decoration: none;
			border-bottom: 2px solid yellow; /* 添加标题下划线 */
            padding-bottom: 5px; /* 增加标题与下划线的间距 */
        }
        .blog-item a:hover {
            text-decoration: underline;
        }
		.blog-section-title { /* 添加列表项标题样式 */
            border-bottom: 1px solid yellow;
            color: yellow;
            font-weight: bold;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <h1>小纸条</h1>
    <ul class="blog-list">
`)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(noteDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			fileName = fileName[:len(fileName)-3]
			_, err = indexFile.WriteString(`<li class="blog-item"><a href="notes/` + fileName + `">` + fileName + `</a></li>`)
			if err != nil {
				return err
			}
		}
	}

	_, err = indexFile.WriteString(`
    </ul>
</body>
</html>
`)
	if err != nil {
		return err
	}
	return nil
}

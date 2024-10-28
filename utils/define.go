package utils

type Note struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// TODO: 以下定义应被单独存放，放到这里不合理
type Status string

const (
	StatusDraft Status = "draft"
	StatusPublished Status = "published"
)

// 笔记元数据定义
type NoteMeta struct {
	Tags        []string `json:"tags"`
	Date        string   `json:"date"`
	Author      string   `json:"author"`
	Status      Status   `json:"status"`
	LastModified string   `json:"last_modified"`
}

const (
	DefaultNoteDir = "/var/www/slip/notes"
)

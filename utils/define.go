package utils

type Note struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

const (
	// 默认笔记写入的目录 /var/www/slip/notes
	DefaultNoteDir = "/var/www/slip/notes"
)

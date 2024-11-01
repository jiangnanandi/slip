package defines

import (
	"regexp"
	"strings"
	"slip/internal/config"

	"gopkg.in/yaml.v3"
)

type Notes struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Dir string
	Meta  NoteMeta
}

type Status string

// draft -> published -> archived -> deleted -> private
const (
	StatusDraft     Status = "draft"      // 草稿
	StatusPublished Status = "published"  // 已发布
	StatusArchived  Status = "archived"   // 已归档
	StatusDeleted   Status = "deleted"    // 已删除
	StatusPrivate   Status = "private"    // 私密
)

// 笔记元数据定义
type NoteMeta struct {
	Tags         []string `json:"tags"`
	Date         string   `json:"date"`
	Author       string   `json:"author"`
	Status       Status   `json:"status"`
	LastModified string   `json:"last_modified"`
}

func (n *Notes) Build() error {
	err := n.decodeMeta()
	if err != nil {
		return err
	}
	switch strings.ToLower(string(n.Meta.Status)) {
	case string(StatusPublished):
		n.Dir = config.AppConfig.Notes.PublishedDir
	case string(StatusDraft):
		n.Dir = config.AppConfig.Notes.DraftDir
	case string(StatusArchived):
		n.Dir = config.AppConfig.Notes.ArchivedDir
	case string(StatusDeleted):
		n.Dir = config.AppConfig.Notes.DeletedDir
	case string(StatusPrivate):
		n.Dir = config.AppConfig.Notes.PrivateDir
	default:
		// 如果状态无效，返回默认目录
		n.Dir = config.AppConfig.Notes.DraftDir
	}
	return nil
}

func (n *Notes) decodeMeta() error {
	re := regexp.MustCompile(`(?s)---\s*(.*?)\s*---`)
	matches := re.FindStringSubmatch(n.Body)
	if len(matches) == 0 {
		return nil
	}
	meta := NoteMeta{}
	err := yaml.Unmarshal([]byte(matches[1]), &meta)
	if err != nil {
		return err
	}
	// 将 note 中的 meta 部分删除
	n.Body = strings.Replace(n.Body, matches[0], "", 1)
	n.Meta = meta

	return nil
}